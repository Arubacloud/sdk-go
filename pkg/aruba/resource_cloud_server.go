package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/compute"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// CloudServer is the wrapper for an Aruba Cloud Compute server (a direct child of a Project).
// Construct with aruba.NewCloudServer() and bind it via IntoProject(project), WithVPC(vpc),
// WithBootVolume(volume), etc.
//
// Schema asymmetry: the request side uses FlavorName *string under the "flavorName" wire
// field; the response side returns a full Flavor struct under the "flavor" wire field.
// This wrapper exposes WithFlavor(flavor) for the request leg and Flavor() / FlavorRaw()
// for the response leg.
//
// The response also carries Template ReferenceResource (no request equivalent); this
// wrapper exposes Template() as a read-only getter.
type CloudServer struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	// Request-side scalars.
	zone      *string
	flavor    *string
	userData  *string
	vpcPreset *bool

	// Body-refs (single).
	vpcRef        *string
	bootVolumeRef *string
	keyPairRef    *string
	elasticIPRef  *string

	// Body-refs (multi-slice).
	subnetRefs        []string
	securityGroupRefs []string

	// Hydrated response.
	response *types.CloudServerResponse

	// Action executor — set by the adapter when this wrapper is produced by a real client
	// call. Locally-constructed wrappers (NewCloudServer()) have actions == nil and will
	// return a clear error when PowerOn/PowerOff/SetPassword are called.
	actions cloudServerActions
}

// Setters (chainable).

func (cs *CloudServer) IntoProject(p Ref) *CloudServer        { cs.intoProject(p); return cs }
func (cs *CloudServer) WithName(n string) *CloudServer        { cs.withName(n); return cs }
func (cs *CloudServer) AddTag(t string) *CloudServer          { cs.addTag(t); return cs }
func (cs *CloudServer) RemoveTag(t string) *CloudServer       { cs.removeTag(t); return cs }
func (cs *CloudServer) ReplaceTags(ts ...string) *CloudServer { cs.replaceTags(ts...); return cs }
func (cs *CloudServer) WithLocation(loc string) *CloudServer  { cs.withLocation(loc); return cs }
func (cs *CloudServer) InRegion(region string) *CloudServer   { cs.withLocation(region); return cs }

func (cs *CloudServer) InZone(zone string) *CloudServer       { cs.zone = &zone; return cs }
func (cs *CloudServer) WithFlavor(flavor string) *CloudServer { cs.flavor = &flavor; return cs }
func (cs *CloudServer) WithUserData(b64 string) *CloudServer  { cs.userData = &b64; return cs }
func (cs *CloudServer) WithVPCPreset(b bool) *CloudServer     { cs.vpcPreset = &b; return cs }

// Single body-ref setters.

func (cs *CloudServer) WithVPC(v Ref) *CloudServer { return cs.setSingleRef("WithVPC", v, &cs.vpcRef) }
func (cs *CloudServer) WithBootVolume(vol Ref) *CloudServer {
	return cs.setSingleRef("WithBootVolume", vol, &cs.bootVolumeRef)
}
func (cs *CloudServer) WithKeyPair(kp Ref) *CloudServer {
	return cs.setSingleRef("WithKeyPair", kp, &cs.keyPairRef)
}
func (cs *CloudServer) WithElasticIP(eip Ref) *CloudServer {
	return cs.setSingleRef("WithElasticIP", eip, &cs.elasticIPRef)
}

// Multi-ref slice setters.

func (cs *CloudServer) AddSubnet(s Ref) *CloudServer {
	return cs.appendRef("AddSubnet", s, &cs.subnetRefs)
}
func (cs *CloudServer) AddSecurityGroup(sg Ref) *CloudServer {
	return cs.appendRef("AddSecurityGroup", sg, &cs.securityGroupRefs)
}

// Internal ref helpers.

func (cs *CloudServer) setSingleRef(label string, r Ref, dst **string) *CloudServer {
	uri := r.URI()
	if uri == "" {
		cs.addErr(fmt.Errorf("%s: empty URI", label))
		return cs
	}
	*dst = &uri
	return cs
}

func (cs *CloudServer) appendRef(label string, r Ref, dst *[]string) *CloudServer {
	uri := r.URI()
	if uri == "" {
		cs.addErr(fmt.Errorf("%s: empty URI", label))
		return cs
	}
	*dst = append(*dst, uri)
	return cs
}

// Ref + ID accessors.

func (cs *CloudServer) URI() string           { return cs.RespURI() }
func (cs *CloudServer) CloudServerID() string { return cs.ID() }

// Accessors.

func (cs *CloudServer) Raw() *types.CloudServerResponse      { return cs.response }
func (cs *CloudServer) RawRequest() types.CloudServerRequest { return cs.toRequest() }

func (cs *CloudServer) Zone() string {
	return cloudServerDerefString(cs.zone)
}

// Flavor returns the flavor name. On a hydrated response the value comes from the
// response's Flavor.Name; before hydration it returns what was passed to WithFlavor.
func (cs *CloudServer) Flavor() string {
	if cs.response != nil && cs.response.Properties.Flavor.Name != "" {
		return cs.response.Properties.Flavor.Name
	}
	return cloudServerDerefString(cs.flavor)
}

// FlavorRaw returns the full flavor struct from the last response, or nil.
func (cs *CloudServer) FlavorRaw() *types.CloudServerFlavorResponse {
	if cs.response == nil {
		return nil
	}
	return &cs.response.Properties.Flavor
}

// Template returns the template URI from the last response (read-only; no request equivalent).
func (cs *CloudServer) Template() string {
	if cs.response == nil {
		return ""
	}
	return cs.response.Properties.Template.URI
}

func (cs *CloudServer) VPC() string {
	if cs.response != nil && cs.response.Properties.VPC.URI != "" {
		return cs.response.Properties.VPC.URI
	}
	return cloudServerDerefString(cs.vpcRef)
}

func (cs *CloudServer) BootVolume() string {
	if cs.response != nil && cs.response.Properties.BootVolume.URI != "" {
		return cs.response.Properties.BootVolume.URI
	}
	return cloudServerDerefString(cs.bootVolumeRef)
}

func (cs *CloudServer) KeyPair() string {
	if cs.response != nil && cs.response.Properties.KeyPair.URI != "" {
		return cs.response.Properties.KeyPair.URI
	}
	return cloudServerDerefString(cs.keyPairRef)
}

func (cs *CloudServer) NetworkInterfaces() []types.CloudServerNetworkInterfaceDetails {
	if cs.response == nil {
		return nil
	}
	return cs.response.Properties.NetworkInterfaces
}

// Action methods (require hydration via a client Get/Create/Update/List call).

func (cs *CloudServer) PowerOn(ctx context.Context, opts ...CallOption) error {
	if err := cs.preActionCheck("PowerOn"); err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := cs.actions.powerOn(ctx, cs.ProjectID(), cs.CloudServerID(), rp)
	populateHTTPEnvelope(&cs.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		cs.fromResponse(resp.Data)
	}
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (cs *CloudServer) PowerOff(ctx context.Context, opts ...CallOption) error {
	if err := cs.preActionCheck("PowerOff"); err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := cs.actions.powerOff(ctx, cs.ProjectID(), cs.CloudServerID(), rp)
	populateHTTPEnvelope(&cs.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		cs.fromResponse(resp.Data)
	}
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (cs *CloudServer) SetPassword(ctx context.Context, password string, opts ...CallOption) error {
	if err := cs.preActionCheck("SetPassword"); err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := cs.actions.setPassword(ctx, cs.ProjectID(), cs.CloudServerID(), password, rp)
	populateHTTPEnvelope(&cs.httpEnvelopeMixin, resp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (cs *CloudServer) preActionCheck(label string) error {
	if cs.actions == nil {
		return fmt.Errorf("%s: this *CloudServer was not obtained via a client call (no action executor) — fetch via Get/Create/Update/List first", label)
	}
	if cs.CloudServerID() == "" {
		return fmt.Errorf("%s: missing cloud-server ID", label)
	}
	if cs.ProjectID() == "" {
		return fmt.Errorf("%s: missing project ID", label)
	}
	return nil
}

// Wire conversions.

func (cs *CloudServer) toRequest() types.CloudServerRequest {
	props := types.CloudServerPropertiesRequest{}
	if cs.zone != nil {
		props.Zone = *cs.zone
	}
	if cs.vpcPreset != nil {
		props.VPCPreset = *cs.vpcPreset
	}
	if cs.flavor != nil {
		v := *cs.flavor
		props.FlavorName = &v // wire field is "flavorName"; wrapper field is just "flavor"
	}
	if cs.userData != nil {
		v := *cs.userData
		props.UserData = &v
	}
	if cs.vpcRef != nil {
		props.VPC = types.ReferenceResource{URI: *cs.vpcRef}
	}
	if cs.bootVolumeRef != nil {
		props.BootVolume = types.ReferenceResource{URI: *cs.bootVolumeRef}
	}
	if cs.keyPairRef != nil {
		props.KeyPair = &types.ReferenceResource{URI: *cs.keyPairRef}
	}
	if cs.elasticIPRef != nil {
		props.ElasticIP = &types.ReferenceResource{URI: *cs.elasticIPRef}
	}
	if len(cs.subnetRefs) > 0 {
		props.Subnets = make([]types.ReferenceResource, 0, len(cs.subnetRefs))
		for _, u := range cs.subnetRefs {
			props.Subnets = append(props.Subnets, types.ReferenceResource{URI: u})
		}
	}
	if len(cs.securityGroupRefs) > 0 {
		props.SecurityGroups = make([]types.ReferenceResource, 0, len(cs.securityGroupRefs))
		for _, u := range cs.securityGroupRefs {
			props.SecurityGroups = append(props.SecurityGroups, types.ReferenceResource{URI: u})
		}
	}
	return types.CloudServerRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: cs.toMetadata(),
			Location:                cs.toLocation(),
		},
		Properties: props,
	}
}

func (cs *CloudServer) fromResponse(resp *types.CloudServerResponse) {
	if resp == nil {
		return
	}
	cs.response = resp
	cs.setMeta(&resp.Metadata)
	cs.withName(cloudServerDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		cs.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		cs.withLocation(resp.Metadata.LocationResponse.Value)
	}
	cs.setLinked(resp.Properties.LinkedResources)
	cs.setStatus(&resp.Status)

	if resp.Properties.Zone != "" {
		v := resp.Properties.Zone
		cs.zone = &v
	}
	if resp.Properties.Flavor.Name != "" {
		v := resp.Properties.Flavor.Name
		cs.flavor = &v
	}
	if resp.Properties.VPC.URI != "" {
		v := resp.Properties.VPC.URI
		cs.vpcRef = &v
	}
	if resp.Properties.BootVolume.URI != "" {
		v := resp.Properties.BootVolume.URI
		cs.bootVolumeRef = &v
	}
	if resp.Properties.KeyPair.URI != "" {
		v := resp.Properties.KeyPair.URI
		cs.keyPairRef = &v
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		cs.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if cs.projectID == "" && cs.RespURI() != "" {
		ids := parseURIIDs(cs.RespURI())
		cs.projectID = ids["projects"]
	}
}

func cloudServerDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ---------------------------------------------------------------------------
// Low-level client seam + adapter
// ---------------------------------------------------------------------------

// cloudServerActions is an internal interface satisfied by cloudServersClientAdapter. It
// allows *CloudServer to dispatch PowerOn/PowerOff/SetPassword without leaking the adapter
// into the public API.
type cloudServerActions interface {
	powerOn(ctx context.Context, projectID, cloudServerID string, rp *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	powerOff(ctx context.Context, projectID, cloudServerID string, rp *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	setPassword(ctx context.Context, projectID, cloudServerID, password string, rp *types.RequestParameters) (*types.Response[any], error)
}

type cloudServerLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.CloudServerList], error)
	Get(ctx context.Context, projectID, cloudServerID string, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	Create(ctx context.Context, projectID string, body types.CloudServerRequest, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	Update(ctx context.Context, projectID, cloudServerID string, body types.CloudServerRequest, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	Delete(ctx context.Context, projectID, cloudServerID string, params *types.RequestParameters) (*types.Response[any], error)
	PowerOn(ctx context.Context, projectID, cloudServerID string, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	PowerOff(ctx context.Context, projectID, cloudServerID string, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	SetPassword(ctx context.Context, projectID, cloudServerID string, body types.CloudServerPasswordRequest, params *types.RequestParameters) (*types.Response[any], error)
}

type cloudServersClientAdapter struct{ low cloudServerLowLevelClient }

var _ cloudServerActions = (*cloudServersClientAdapter)(nil)

func newCloudServersClientAdapter(rest *restclient.Client) *cloudServersClientAdapter {
	if rest == nil {
		return &cloudServersClientAdapter{}
	}
	return &cloudServersClientAdapter{low: compute.NewCloudServersClientImpl(rest)}
}

func (a *cloudServersClientAdapter) Create(ctx context.Context, cs *CloudServer, opts ...CallOption) (*CloudServer, error) {
	if err := cs.Err(); err != nil {
		return cs, err
	}
	if cs.ProjectID() == "" {
		return cs, fmt.Errorf("Create: CloudServer has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, cs.ProjectID(), cs.toRequest(), rp)
	populateHTTPEnvelope(&cs.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		cs.fromResponse(resp.Data)
	}
	cs.actions = a
	if err != nil {
		return cs, err
	}
	if resp != nil && !resp.IsSuccess() {
		return cs, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return cs, nil
}

func (a *cloudServersClientAdapter) Update(ctx context.Context, cs *CloudServer, opts ...CallOption) (*CloudServer, error) {
	if err := cs.Err(); err != nil {
		return cs, err
	}
	if cs.CloudServerID() == "" {
		return cs, fmt.Errorf("Update: CloudServer has no ID")
	}
	if cs.ProjectID() == "" {
		return cs, fmt.Errorf("Update: CloudServer has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, cs.ProjectID(), cs.CloudServerID(), cs.toRequest(), rp)
	populateHTTPEnvelope(&cs.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		cs.fromResponse(resp.Data)
	}
	cs.actions = a
	if err != nil {
		return cs, err
	}
	if resp != nil && !resp.IsSuccess() {
		return cs, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return cs, nil
}

func (a *cloudServersClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*CloudServer, error) {
	projectID, cloudServerID, err := cloudServerIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, cloudServerID, rp)
	out := &CloudServer{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	if out.projectID == "" {
		out.projectID = projectID
	}
	out.actions = a
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *cloudServersClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, cloudServerID, err := cloudServerIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, cloudServerID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *cloudServersClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*CloudServer], error) {
	projectID, err := projectIDFromRef(project)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*CloudServer
	if resp != nil && resp.Data != nil {
		items = make([]*CloudServer, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			cs := &CloudServer{}
			cs.fromResponse(&resp.Data.Values[i])
			if cs.projectID == "" {
				cs.projectID = projectID
			}
			cs.actions = a
			items = append(items, cs)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*CloudServer], error) {
		return nil, fmt.Errorf("List pagination by URL not yet wired; re-call List with adjusted CallOptions")
	}
	var total int64
	var self, prev, next, first, last string
	if resp != nil && resp.Data != nil {
		total = resp.Data.Total
		self = resp.Data.Self
		prev = resp.Data.Prev
		next = resp.Data.Next
		first = resp.Data.First
		last = resp.Data.Last
	}
	return newList(items, total, self, prev, next, first, last, resp, opts, refetch), nil
}

// Internal action methods — satisfy cloudServerActions; called by *CloudServer action methods.

func (a *cloudServersClientAdapter) powerOn(ctx context.Context, projectID, cloudServerID string, rp *types.RequestParameters) (*types.Response[types.CloudServerResponse], error) {
	return a.low.PowerOn(ctx, projectID, cloudServerID, rp)
}

func (a *cloudServersClientAdapter) powerOff(ctx context.Context, projectID, cloudServerID string, rp *types.RequestParameters) (*types.Response[types.CloudServerResponse], error) {
	return a.low.PowerOff(ctx, projectID, cloudServerID, rp)
}

func (a *cloudServersClientAdapter) setPassword(ctx context.Context, projectID, cloudServerID, password string, rp *types.RequestParameters) (*types.Response[any], error) {
	return a.low.SetPassword(ctx, projectID, cloudServerID, types.CloudServerPasswordRequest{Password: password}, rp)
}

// cloudServerIDsFromRef extracts (projectID, cloudServerID) from a Ref.
func cloudServerIDsFromRef(ref Ref) (projectID, cloudServerID string, err error) {
	csID, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withCloudServerID); ok {
			return w.CloudServerID(), true
		}
		return "", false
	}, "cloudServers")
	if !ok || csID == "" {
		return "", "", fmt.Errorf("cannot determine CloudServer ID from Ref %q", ref.URI())
	}
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid == "" {
		return "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return pid, csID, nil
}

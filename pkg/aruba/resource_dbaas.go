package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/database"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// DBaaS is the wrapper for an Aruba Cloud Database-as-a-Service instance
// (a direct child of a Project). Construct with aruba.NewDBaaS() and bind
// it via IntoProject(project), OfEngine, WithServerFlavor, WithStorage,
// WithVPC/WithSubnet/WithSecurityGroup/WithElasticIP, etc.
//
// Schema asymmetries:
//   - Engine: request emits Engine.ID; response returns a full
//     DBaaSEngineResponse{Type,Name,Version,...}. Engine() reads .Type
//     (the human-meaningful identifier on the response side).
//   - Flavor: request emits Flavor.Name; response returns the full
//     DBaaSFlavorResponse struct.
//   - Networking: request emits 4 raw URI strings (VPCURI/SubnetURI/
//     SecurityGroupURI/ElasticIPURI); response returns 4 *ReferenceResource
//     objects. Read-back getters (VPC/Subnet/SecurityGroup/ElasticIP)
//     prefer the response side, falling back to the locally-set URI.
//   - Zone: Go field is "Zone" but the wire JSON tag is "dataCenter".
//   - Autoscaling: request emits {Enabled,AvailableSpace,StepSize};
//     response returns {Status,AvailableSpace,StepSize,RuleID}.
//     AutoscalingEnabled() reads only the locally-set value;
//     AutoscalingStatus() / AutoscalingRuleID() read only the response;
//     AutoscalingAvailableSpace() / AutoscalingStepSize() prefer the
//     response and fall back to the locally-set value.
//     fromResponse does NOT back-populate request-side fields, so an
//     Update after Get omits the autoscaling block unless the caller
//     re-asserts intent via WithAutoscaling/WithoutAutoscaling.
type DBaaS struct {
	errMixin
	metadataMixin
	zonalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	// Request-side scalars.
	engine                    *DatabaseEngine // wire: Engine.ID
	flavor                    *DBaaSFlavor    // wire: Flavor.Name
	storageGB                 *int32          // wire: Storage.SizeGB
	autoscalingEnabled        *bool           // wire: Autoscaling.Enabled
	autoscalingAvailableSpace *int32          // wire: Autoscaling.AvailableSpace
	autoscalingStepSize       *int32          // wire: Autoscaling.StepSize
	billingPeriod             *BillingPeriod  // wire: BillingPlan.BillingPeriod

	// Networking refs.
	vpcRef           *string
	subnetRef        *string
	securityGroupRef *string
	elasticIPRef     *string

	// Hydrated response.
	response *types.DBaaSResponse
}

// Setters (chainable).

func (d *DBaaS) IntoProject(p Ref) *DBaaS        { d.intoProject(p); return d }
func (d *DBaaS) WithName(n string) *DBaaS        { d.withName(n); return d }
func (d *DBaaS) AddTag(t string) *DBaaS          { d.addTag(t); return d }
func (d *DBaaS) RemoveTag(t string) *DBaaS       { d.removeTag(t); return d }
func (d *DBaaS) ReplaceTags(ts ...string) *DBaaS { d.replaceTags(ts...); return d }
func (d *DBaaS) InRegion(region Region) *DBaaS   { d.inRegion(region); return d }

func (d *DBaaS) InZone(zone Zone) *DBaaS                    { d.inZone(zone); return d }
func (d *DBaaS) OfEngine(engine DatabaseEngine) *DBaaS      { d.engine = &engine; return d }
func (d *DBaaS) WithServerFlavor(flavor DBaaSFlavor) *DBaaS { d.flavor = &flavor; return d }
func (d *DBaaS) WithStorageGB(gb int) *DBaaS                { v := int32(gb); d.storageGB = &v; return d }

// WithAutoscaling enables autoscaling and pins the available-space threshold and
// step size in GB. Mirrors NodePool.WithAutoscaling(min, max) from resource_kaas_nodepool.go.
func (d *DBaaS) WithAutoscaling(availableSpaceGB, stepSizeGB int) *DBaaS {
	t := true
	av := int32(availableSpaceGB)
	ss := int32(stepSizeGB)
	d.autoscalingEnabled = &t
	d.autoscalingAvailableSpace = &av
	d.autoscalingStepSize = &ss
	return d
}

// WithoutAutoscaling explicitly disables autoscaling and clears the bounds.
func (d *DBaaS) WithoutAutoscaling() *DBaaS {
	f := false
	d.autoscalingEnabled = &f
	d.autoscalingAvailableSpace = nil
	d.autoscalingStepSize = nil
	return d
}
func (d *DBaaS) WithBillingPeriod(period BillingPeriod) *DBaaS { d.billingPeriod = &period; return d }

func (d *DBaaS) WithVPC(v Ref) *DBaaS    { return d.setSingleRef("WithVPC", v, &d.vpcRef) }
func (d *DBaaS) WithSubnet(s Ref) *DBaaS { return d.setSingleRef("WithSubnet", s, &d.subnetRef) }
func (d *DBaaS) WithSecurityGroup(sg Ref) *DBaaS {
	return d.setSingleRef("WithSecurityGroup", sg, &d.securityGroupRef)
}
func (d *DBaaS) WithElasticIP(eip Ref) *DBaaS {
	return d.setSingleRef("WithElasticIP", eip, &d.elasticIPRef)
}

func (d *DBaaS) setSingleRef(label string, r Ref, dst **string) *DBaaS {
	uri := r.URI()
	if uri == "" {
		d.addErr(fmt.Errorf("%s: empty URI", label))
		return d
	}
	*dst = &uri
	return d
}

// Ref + ID accessors.

func (d *DBaaS) URI() string     { return d.RespURI() }
func (d *DBaaS) DBaaSID() string { return d.ID() }

// Accessors.

func (d *DBaaS) Raw() *types.DBaaSResponse      { return d.response }
func (d *DBaaS) RawRequest() types.DBaaSRequest { return d.toRequest() }

// Engine returns the engine identifier. On a hydrated response the value comes
// from Engine.Type; before hydration it returns what was passed to OfEngine.
func (d *DBaaS) Engine() DatabaseEngine {
	if d.response != nil && d.response.Properties.Engine != nil && d.response.Properties.Engine.Type != nil {
		return DatabaseEngine(*d.response.Properties.Engine.Type)
	}
	if d.engine == nil {
		return ""
	}
	return *d.engine
}

// EngineRaw returns the full engine struct from the last response, or nil.
func (d *DBaaS) EngineRaw() *types.DBaaSEngineResponse {
	if d.response == nil {
		return nil
	}
	return d.response.Properties.Engine
}

// Flavor returns the flavor name. On a hydrated response the value comes from
// Flavor.Name; before hydration it returns what was passed to WithServerFlavor.
func (d *DBaaS) Flavor() DBaaSFlavor {
	if d.response != nil && d.response.Properties.Flavor != nil && d.response.Properties.Flavor.Name != nil {
		return DBaaSFlavor(*d.response.Properties.Flavor.Name)
	}
	if d.flavor == nil {
		return ""
	}
	return *d.flavor
}

// FlavorRaw returns the full flavor struct from the last response, or nil.
func (d *DBaaS) FlavorRaw() *types.DBaaSFlavorResponse {
	if d.response == nil {
		return nil
	}
	return d.response.Properties.Flavor
}

// StorageGB returns the storage size in GB. On a hydrated response the value comes
// from Storage.SizeGB; before hydration it returns what was passed to WithStorageGB.
func (d *DBaaS) StorageGB() int {
	if d.response != nil && d.response.Properties.Storage != nil && d.response.Properties.Storage.SizeGB != nil {
		return int(*d.response.Properties.Storage.SizeGB)
	}
	if d.storageGB != nil {
		return int(*d.storageGB)
	}
	return 0
}

// BillingPeriod returns the billing period. On a hydrated response the value comes
// from BillingPlan.BillingPeriod; before hydration it returns what was passed to
// WithBillingPeriod.
func (d *DBaaS) BillingPeriod() BillingPeriod {
	if d.response != nil && d.response.Properties.BillingPlan != nil && d.response.Properties.BillingPlan.BillingPeriod != nil {
		return *d.response.Properties.BillingPlan.BillingPeriod
	}
	if d.billingPeriod == nil {
		return ""
	}
	return *d.billingPeriod
}

// AutoscalingEnabled returns the locally-set Enabled flag (request-side intent).
// The response side carries no Enabled field — see AutoscalingStatus() for the
// platform-reported state.
func (d *DBaaS) AutoscalingEnabled() bool {
	if d.autoscalingEnabled != nil {
		return *d.autoscalingEnabled
	}
	return false
}

// AutoscalingStatus returns the response-side autoscaling status string.
// Empty before hydration.
func (d *DBaaS) AutoscalingStatus() string {
	if d.response != nil && d.response.Properties.Autoscaling != nil &&
		d.response.Properties.Autoscaling.Status != nil {
		return *d.response.Properties.Autoscaling.Status
	}
	return ""
}

// AutoscalingAvailableSpaceGB returns the available-space threshold in GB.
// Hydrated response wins; otherwise returns the locally-set value, else 0.
func (d *DBaaS) AutoscalingAvailableSpaceGB() int {
	if d.response != nil && d.response.Properties.Autoscaling != nil &&
		d.response.Properties.Autoscaling.AvailableSpace != nil {
		return int(*d.response.Properties.Autoscaling.AvailableSpace)
	}
	if d.autoscalingAvailableSpace != nil {
		return int(*d.autoscalingAvailableSpace)
	}
	return 0
}

// AutoscalingStepSizeGB returns the step size in GB.
// Hydrated response wins; otherwise returns the locally-set value, else 0.
func (d *DBaaS) AutoscalingStepSizeGB() int {
	if d.response != nil && d.response.Properties.Autoscaling != nil &&
		d.response.Properties.Autoscaling.StepSize != nil {
		return int(*d.response.Properties.Autoscaling.StepSize)
	}
	if d.autoscalingStepSize != nil {
		return int(*d.autoscalingStepSize)
	}
	return 0
}

// AutoscalingRuleID returns the response-side rule identifier.
// Empty before hydration.
func (d *DBaaS) AutoscalingRuleID() string {
	if d.response != nil && d.response.Properties.Autoscaling != nil &&
		d.response.Properties.Autoscaling.RuleID != nil {
		return *d.response.Properties.Autoscaling.RuleID
	}
	return ""
}

// AutoscalingRaw returns the full autoscaling response struct, or nil before hydration.
func (d *DBaaS) AutoscalingRaw() *types.DBaaSAutoscalingResponse {
	if d.response == nil {
		return nil
	}
	return d.response.Properties.Autoscaling
}

func (d *DBaaS) VPC() string {
	return dbaasNetworkingURI(d.response, func(n *types.DBaaSNetworkingResponse) *types.ReferenceResource { return n.VPC }, d.vpcRef)
}

func (d *DBaaS) Subnet() string {
	return dbaasNetworkingURI(d.response, func(n *types.DBaaSNetworkingResponse) *types.ReferenceResource { return n.Subnet }, d.subnetRef)
}

func (d *DBaaS) SecurityGroup() string {
	return dbaasNetworkingURI(d.response, func(n *types.DBaaSNetworkingResponse) *types.ReferenceResource { return n.SecurityGroup }, d.securityGroupRef)
}

func (d *DBaaS) ElasticIP() string {
	return dbaasNetworkingURI(d.response, func(n *types.DBaaSNetworkingResponse) *types.ReferenceResource { return n.ElasticIP }, d.elasticIPRef)
}

// Wire conversions.

func (d *DBaaS) toRequest() types.DBaaSRequest {
	props := types.DBaaSPropertiesRequest{}
	props.Zone = d.zonePtr()
	if d.engine != nil {
		props.Engine = &types.DBaaSEngine{ID: d.engine}
	}
	if d.flavor != nil {
		props.Flavor = &types.DBaaSFlavorSpec{Name: d.flavor}
	}
	if d.storageGB != nil {
		v := *d.storageGB
		props.Storage = &types.DBaaSStorage{SizeGB: &v}
	}
	if d.autoscalingEnabled != nil || d.autoscalingAvailableSpace != nil || d.autoscalingStepSize != nil {
		a := &types.DBaaSAutoscaling{}
		if d.autoscalingEnabled != nil {
			v := *d.autoscalingEnabled
			a.Enabled = &v
		}
		if d.autoscalingAvailableSpace != nil {
			v := *d.autoscalingAvailableSpace
			a.AvailableSpace = &v
		}
		if d.autoscalingStepSize != nil {
			v := *d.autoscalingStepSize
			a.StepSize = &v
		}
		props.Autoscaling = a
	}
	if d.billingPeriod != nil {
		v := *d.billingPeriod
		props.BillingPlan = &types.DBaaSBillingPlan{BillingPeriod: &v}
	}
	if d.vpcRef != nil || d.subnetRef != nil || d.securityGroupRef != nil || d.elasticIPRef != nil {
		net := &types.DBaaSNetworking{}
		if d.vpcRef != nil {
			net.VPCURI = d.vpcRef
		}
		if d.subnetRef != nil {
			net.SubnetURI = d.subnetRef
		}
		if d.securityGroupRef != nil {
			net.SecurityGroupURI = d.securityGroupRef
		}
		if d.elasticIPRef != nil {
			net.ElasticIPURI = d.elasticIPRef
		}
		props.Networking = net
	}
	return types.DBaaSRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: d.toMetadata(),
			Location:                d.toLocation(),
		},
		Properties: props,
	}
}

func (d *DBaaS) fromResponse(resp *types.DBaaSResponse) {
	if resp == nil {
		return
	}
	d.response = resp
	d.setMeta(&resp.Metadata)
	d.withName(dbaasDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		d.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		d.inRegion(resp.Metadata.LocationResponse.Value)
	}
	d.setLinked(resp.Properties.LinkedResources)
	d.setStatus(&resp.Status)
	d.setTerminalStates(dbaasTerminalStates)

	// Hydrate request-side fields from response for round-trip Update support.
	if resp.Properties.Engine != nil && resp.Properties.Engine.Type != nil {
		e := DatabaseEngine(*resp.Properties.Engine.Type)
		d.engine = &e
	}
	if resp.Properties.Flavor != nil && resp.Properties.Flavor.Name != nil {
		f := DBaaSFlavor(*resp.Properties.Flavor.Name)
		d.flavor = &f
	}
	if resp.Properties.Storage != nil && resp.Properties.Storage.SizeGB != nil {
		v := *resp.Properties.Storage.SizeGB
		d.storageGB = &v
	}
	if resp.Properties.BillingPlan != nil && resp.Properties.BillingPlan.BillingPeriod != nil {
		v := *resp.Properties.BillingPlan.BillingPeriod
		d.billingPeriod = &v
	}
	if resp.Properties.Networking != nil {
		if resp.Properties.Networking.VPC != nil && resp.Properties.Networking.VPC.URI != "" {
			v := resp.Properties.Networking.VPC.URI
			d.vpcRef = &v
		}
		if resp.Properties.Networking.Subnet != nil && resp.Properties.Networking.Subnet.URI != "" {
			v := resp.Properties.Networking.Subnet.URI
			d.subnetRef = &v
		}
		if resp.Properties.Networking.SecurityGroup != nil && resp.Properties.Networking.SecurityGroup.URI != "" {
			v := resp.Properties.Networking.SecurityGroup.URI
			d.securityGroupRef = &v
		}
		if resp.Properties.Networking.ElasticIP != nil && resp.Properties.Networking.ElasticIP.URI != "" {
			v := resp.Properties.Networking.ElasticIP.URI
			d.elasticIPRef = &v
		}
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		d.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if d.projectID == "" && d.RespURI() != "" {
		ids := parseURIIDs(d.RespURI())
		d.projectID = ids["projects"]
	}
}

func dbaasDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func dbaasNetworkingURI(resp *types.DBaaSResponse, pick func(*types.DBaaSNetworkingResponse) *types.ReferenceResource, fallback *string) string {
	if resp != nil && resp.Properties.Networking != nil {
		if r := pick(resp.Properties.Networking); r != nil && r.URI != "" {
			return r.URI
		}
	}
	return dbaasDerefString(fallback)
}

var dbaasTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
	"Failed": false,
}

// ---------------------------------------------------------------------------
// DBaaS low-level client, adapter, and helpers
// ---------------------------------------------------------------------------

type dbaasLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.DBaaSList], error)
	Get(ctx context.Context, projectID, dbaasID string, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	Create(ctx context.Context, projectID string, body types.DBaaSRequest, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	Update(ctx context.Context, projectID, dbaasID string, body types.DBaaSRequest, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	Delete(ctx context.Context, projectID, dbaasID string, params *types.RequestParameters) (*types.Response[any], error)
}

type dbaasClientAdapter struct{ low dbaasLowLevelClient }

func newDBaaSClientAdapter(rest *restclient.Client) *dbaasClientAdapter {
	if rest == nil {
		return &dbaasClientAdapter{}
	}
	return &dbaasClientAdapter{low: database.NewDBaaSClientImpl(rest)}
}

func (a *dbaasClientAdapter) Create(ctx context.Context, d *DBaaS, opts ...CallOption) (*DBaaS, error) {
	if err := d.Err(); err != nil {
		return d, err
	}
	if d.ProjectID() == "" {
		return d, fmt.Errorf("Create: DBaaS has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, d.ProjectID(), d.toRequest(), rp)
	populateHTTPEnvelope(&d.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		d.fromResponse(resp.Data)
		d.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, d)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				d.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return d, err
	}
	if resp != nil && !resp.IsSuccess() {
		return d, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return d, nil
}

func (a *dbaasClientAdapter) Update(ctx context.Context, d *DBaaS, opts ...CallOption) (*DBaaS, error) {
	if err := d.Err(); err != nil {
		return d, err
	}
	if d.DBaaSID() == "" {
		return d, fmt.Errorf("Update: DBaaS has no ID")
	}
	if d.ProjectID() == "" {
		return d, fmt.Errorf("Update: DBaaS has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, d.ProjectID(), d.DBaaSID(), d.toRequest(), rp)
	populateHTTPEnvelope(&d.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		d.fromResponse(resp.Data)
		d.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, d)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				d.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return d, err
	}
	if resp != nil && !resp.IsSuccess() {
		return d, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return d, nil
}

func (a *dbaasClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*DBaaS, error) {
	projectID, dbaasID, err := dbaasIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, dbaasID, rp)
	out := &DBaaS{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
		out.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, out)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				out.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if out.projectID == "" {
		out.projectID = projectID
	}
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *dbaasClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, dbaasID, err := dbaasIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, dbaasID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *dbaasClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*DBaaS], error) {
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
	var items []*DBaaS
	if resp != nil && resp.Data != nil {
		items = make([]*DBaaS, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			d := &DBaaS{}
			d.fromResponse(&resp.Data.Values[i])
			d.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, d)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					d.fromResponse(fresh.Raw())
				}
				return nil
			})
			if d.projectID == "" {
				d.projectID = projectID
			}
			items = append(items, d)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*DBaaS], error) {
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

// dbaasIDsFromRef extracts (projectID, dbaasID) from a Ref.
func dbaasIDsFromRef(ref Ref) (projectID, dbaasID string, err error) {
	did, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withDBaaSID); ok {
			return w.DBaaSID(), true
		}
		return "", false
	}, "dbaas")
	if !ok || did == "" {
		return "", "", fmt.Errorf("cannot determine DBaaS ID from Ref %q", ref.URI())
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
	return pid, did, nil
}

package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/container"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// KaaS is the wrapper for an Aruba Cloud Kubernetes-as-a-Service cluster
// (a direct child of a Project). Construct with aruba.NewKaaS() and bind it
// via IntoProject(project), WithVPC(vpc), WithSubnet(subnet), etc.
//
// Family A: regional, Metadata/Properties envelope, location-aware.
// Supports full CRUD. Update emits KaaSUpdateRequest (narrower than KaaSRequest):
// only KubernetesVersion, NodePools, HA, Storage, and BillingPlan are mutable.
//
// Schema asymmetry (request vs. response):
//   - "nodePools" (request) vs. "nodesPool" (response)
//   - "nodeCidr" (request) vs. "nodecidr" (response)
//   - "podCidr" (request) vs. "podcidr" (response)
//   - NodePoolProperties.Instance string (request) vs. InstanceResponse{ID,Name} (response)
//   - NodePoolProperties.Zone string (request JSON "dataCenter") vs. DataCenterResponse{Code,Name} (response)
//
// Path: /projects/{projectID}/providers/Aruba.Container/kaas[/{kaasID}]
type KaaS struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	// Body-refs (single).
	vpcRef    *string
	subnetRef *string

	// Plain-string / scalar body fields.
	securityGroupName    *string
	nodeCIDRAddress      *string
	nodeCIDRName         *string
	podCIDR              *string
	kubernetesVersion    *string
	ha                   *bool
	storageMaxCumulative *int32 // wire: storage.maxCumulativeVolumeSize
	billingPeriod        *BillingPeriod
	identityClientID     *string
	identityClientSecret *string
	apiServerProfile     *types.APIServerAccessProfileProperties // structural pass-through

	// Sub-builders.
	nodePools []*NodePool

	// Action executor — set by adapter; nil on locally-built wrappers.
	actions kaasActions

	response *types.KaaSResponse
}

// ---------------------------------------------------------------------------
// Standard setters
// ---------------------------------------------------------------------------

func (k *KaaS) IntoProject(p Ref) *KaaS                 { k.intoProject(p); return k }
func (k *KaaS) WithName(n string) *KaaS                 { k.withName(n); return k }
func (k *KaaS) AddTag(t string) *KaaS                   { k.addTag(t); return k }
func (k *KaaS) RemoveTag(t string) *KaaS                { k.removeTag(t); return k }
func (k *KaaS) ReplaceTags(ts ...string) *KaaS          { k.replaceTags(ts...); return k }
func (k *KaaS) WithLocation(loc Region) *KaaS           { k.withLocation(loc); return k }
func (k *KaaS) InRegion(region Region) *KaaS            { k.withLocation(region); return k }
func (k *KaaS) WithKubernetesVersion(v string) *KaaS    { k.kubernetesVersion = &v; return k }
func (k *KaaS) WithPodCIDR(cidr string) *KaaS           { k.podCIDR = &cidr; return k }
func (k *KaaS) WithHA(enabled bool) *KaaS               { k.ha = &enabled; return k }
func (k *KaaS) WithBillingPeriod(period BillingPeriod) *KaaS { k.billingPeriod = &period; return k }
func (k *KaaS) WithSecurityGroupName(name string) *KaaS { k.securityGroupName = &name; return k }

// WithNodeCIDR sets the node CIDR block (address and name).
// The wire type is NodeCIDRProperties{Address, Name}.
func (k *KaaS) WithNodeCIDR(address, name string) *KaaS {
	k.nodeCIDRAddress = &address
	k.nodeCIDRName = &name
	return k
}

// WithStorageGB sets the maximum cumulative volume size in GB.
func (k *KaaS) WithStorageGB(gb int) *KaaS {
	v := int32(gb)
	k.storageMaxCumulative = &v
	return k
}

// WithIdentity sets the managed identity credentials.
func (k *KaaS) WithIdentity(clientID, clientSecret string) *KaaS {
	k.identityClientID = &clientID
	k.identityClientSecret = &clientSecret
	return k
}

// WithAPIServerAccessProfile sets the API server access profile.
func (k *KaaS) WithAPIServerAccessProfile(p *types.APIServerAccessProfileProperties) *KaaS {
	k.apiServerProfile = p
	return k
}

// ---------------------------------------------------------------------------
// Body-ref setters
// ---------------------------------------------------------------------------

func (k *KaaS) WithVPC(v Ref) *KaaS    { return k.setSingleRef("WithVPC", v, &k.vpcRef) }
func (k *KaaS) WithSubnet(s Ref) *KaaS { return k.setSingleRef("WithSubnet", s, &k.subnetRef) }

func (k *KaaS) setSingleRef(label string, ref Ref, dst **string) *KaaS {
	uri := ref.URI()
	if uri == "" {
		k.addErr(fmt.Errorf("%s: empty URI", label))
		return k
	}
	*dst = &uri
	return k
}

// ---------------------------------------------------------------------------
// Node pool sub-builder
// ---------------------------------------------------------------------------

// AddNodePool appends np to the cluster's node pool list.
// Errors accumulated on np are drained into k at attachment time.
func (k *KaaS) AddNodePool(np *NodePool) *KaaS {
	if np == nil {
		return k
	}
	for _, e := range np.errs {
		k.addErr(e)
	}
	k.nodePools = append(k.nodePools, np)
	return k
}

// ---------------------------------------------------------------------------
// Action method
// ---------------------------------------------------------------------------

// DownloadKubeconfig downloads the kubeconfig for this cluster and returns it
// as raw bytes (the YAML content). Requires the wrapper to have been obtained
// via a client call (Get/Create/Update/List); locally-built wrappers return a
// clear error.
func (k *KaaS) DownloadKubeconfig(ctx context.Context, opts ...CallOption) ([]byte, error) {
	if err := k.preActionCheck("DownloadKubeconfig"); err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := k.actions.downloadKubeconfig(ctx, k.ProjectID(), k.KaaSID(), rp)
	populateHTTPEnvelope(&k.httpEnvelopeMixin, resp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	if resp == nil || resp.Data == nil {
		return nil, nil
	}
	return []byte(resp.Data.Content), nil
}

func (k *KaaS) preActionCheck(label string) error {
	if k.actions == nil {
		return fmt.Errorf("%s: this *KaaS was not obtained via a client call (no action executor) — fetch via Get/Create/Update/List first", label)
	}
	if k.KaaSID() == "" {
		return fmt.Errorf("%s: missing KaaS ID", label)
	}
	if k.ProjectID() == "" {
		return fmt.Errorf("%s: missing project ID", label)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Ref + ID accessors
// ---------------------------------------------------------------------------

func (k *KaaS) URI() string    { return k.RespURI() }
func (k *KaaS) KaaSID() string { return k.ID() }

// ---------------------------------------------------------------------------
// Raw accessors
// ---------------------------------------------------------------------------

func (k *KaaS) Raw() *types.KaaSResponse      { return k.response }
func (k *KaaS) RawRequest() types.KaaSRequest { return k.toRequest() }

// ---------------------------------------------------------------------------
// Response-preferring accessors
// ---------------------------------------------------------------------------

func (k *KaaS) VPC() string {
	if k.response != nil && k.response.Properties.VPC.URI != nil {
		return *k.response.Properties.VPC.URI
	}
	return kaasDeref(k.vpcRef)
}

func (k *KaaS) Subnet() string {
	if k.response != nil && k.response.Properties.Subnet.URI != nil {
		return *k.response.Properties.Subnet.URI
	}
	return kaasDeref(k.subnetRef)
}

func (k *KaaS) SecurityGroupName() string {
	if k.response != nil && k.response.Properties.SecurityGroup.Name != nil {
		return *k.response.Properties.SecurityGroup.Name
	}
	return kaasDeref(k.securityGroupName)
}

func (k *KaaS) KubernetesVersion() string {
	if k.response != nil && k.response.Properties.KubernetesVersion.Value != nil {
		return *k.response.Properties.KubernetesVersion.Value
	}
	return kaasDeref(k.kubernetesVersion)
}

func (k *KaaS) BillingPeriod() BillingPeriod {
	if k.response != nil && k.response.Properties.BillingPlan != nil && k.response.Properties.BillingPlan.BillingPeriod != nil {
		return *k.response.Properties.BillingPlan.BillingPeriod
	}
	if k.billingPeriod == nil {
		return ""
	}
	return *k.billingPeriod
}

// ---------------------------------------------------------------------------
// Wire conversions
// ---------------------------------------------------------------------------

func (k *KaaS) toRequest() types.KaaSRequest {
	props := types.KaaSPropertiesRequest{
		VPC:    types.ReferenceResource{URI: kaasDeref(k.vpcRef)},
		Subnet: types.ReferenceResource{URI: kaasDeref(k.subnetRef)},
		SecurityGroup: types.SecurityGroupProperties{
			Name: kaasDeref(k.securityGroupName),
		},
		NodeCIDR: types.NodeCIDRProperties{
			Address: kaasDeref(k.nodeCIDRAddress),
			Name:    kaasDeref(k.nodeCIDRName),
		},
		PodCIDR:           k.podCIDR,
		KubernetesVersion: types.KubernetesVersionInfo{Value: kaasDeref(k.kubernetesVersion)},
		HA:                k.ha,
		BillingPlan: func() types.BillingPeriodResource {
			var bp BillingPeriod
			if k.billingPeriod != nil {
				bp = *k.billingPeriod
			}
			return types.BillingPeriodResource{BillingPeriod: bp}
		}(),
	}
	if k.storageMaxCumulative != nil {
		props.Storage = types.StorageKubernetes{MaxCumulativeVolumeSize: k.storageMaxCumulative}
	}
	if k.identityClientID != nil || k.identityClientSecret != nil {
		props.Identity = &types.IdentityProperties{
			ClientID:     k.identityClientID,
			ClientSecret: k.identityClientSecret,
		}
	}
	if k.apiServerProfile != nil {
		props.APIServerAccessProfile = k.apiServerProfile
	}
	if len(k.nodePools) > 0 {
		props.NodePools = make([]types.NodePoolProperties, 0, len(k.nodePools))
		for _, np := range k.nodePools {
			props.NodePools = append(props.NodePools, np.build())
		}
	}
	return types.KaaSRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: k.toMetadata(),
			Location:                k.toLocation(),
		},
		Properties: props,
	}
}

// toUpdateRequest emits KaaSUpdateRequest, which exposes only the mutable
// fields (KubernetesVersion, NodePools, HA, Storage, BillingPlan).
// VPC, Subnet, SecurityGroup, and CIDRs are immutable after creation.
func (k *KaaS) toUpdateRequest() types.KaaSUpdateRequest {
	props := types.KaaSPropertiesUpdateRequest{
		KubernetesVersion: types.KubernetesVersionInfoUpdate{
			Value: kaasDeref(k.kubernetesVersion),
		},
		HA: k.ha,
	}
	if len(k.nodePools) > 0 {
		props.NodePools = make([]types.NodePoolProperties, 0, len(k.nodePools))
		for _, np := range k.nodePools {
			props.NodePools = append(props.NodePools, np.build())
		}
	}
	if k.storageMaxCumulative != nil {
		props.Storage = &types.StorageKubernetes{MaxCumulativeVolumeSize: k.storageMaxCumulative}
	}
	if k.billingPeriod != nil {
		v := *k.billingPeriod
		props.BillingPlan = &types.BillingPeriodResource{BillingPeriod: v}
	}
	return types.KaaSUpdateRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: k.toMetadata(),
			Location:                k.toLocation(),
		},
		Properties: props,
	}
}

func (k *KaaS) fromResponse(resp *types.KaaSResponse) {
	if resp == nil {
		return
	}
	k.response = resp
	k.setMeta(&resp.Metadata)
	k.withName(kaasDeref(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		k.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		k.withLocation(resp.Metadata.LocationResponse.Value)
	}
	k.setStatus(&resp.Status)
	k.setTerminalStates(kaasTerminalStates)
	k.setLinked(resp.Properties.LinkedResources)
	k.kaasHydrateCacheFromProps(resp.Properties)
	k.nodePools = kaasRebuildNodePools(resp.Properties.NodePools)
	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		k.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if k.projectID == "" && k.RespURI() != "" {
		if pid := parseURIIDs(k.RespURI())["projects"]; pid != "" {
			k.projectID = pid
		}
	}
}

func (k *KaaS) kaasHydrateCacheFromProps(props types.KaaSPropertiesResponse) {
	if props.VPC.URI != nil && *props.VPC.URI != "" {
		v := *props.VPC.URI
		k.vpcRef = &v
	}
	if props.Subnet.URI != nil && *props.Subnet.URI != "" {
		v := *props.Subnet.URI
		k.subnetRef = &v
	}
	if props.SecurityGroup.Name != nil && *props.SecurityGroup.Name != "" {
		v := *props.SecurityGroup.Name
		k.securityGroupName = &v
	}
	if props.NodeCIDR.Address != nil && *props.NodeCIDR.Address != "" {
		v := *props.NodeCIDR.Address
		k.nodeCIDRAddress = &v
	}
	if props.NodeCIDR.Name != nil && *props.NodeCIDR.Name != "" {
		v := *props.NodeCIDR.Name
		k.nodeCIDRName = &v
	}
	if props.PodCIDR != nil && props.PodCIDR.Address != nil {
		v := *props.PodCIDR.Address
		k.podCIDR = &v
	}
	if props.KubernetesVersion.Value != nil && *props.KubernetesVersion.Value != "" {
		v := *props.KubernetesVersion.Value
		k.kubernetesVersion = &v
	}
	k.ha = props.HA
	if props.Storage != nil && props.Storage.MaxCumulativeVolumeSize != nil {
		v := *props.Storage.MaxCumulativeVolumeSize
		k.storageMaxCumulative = &v
	}
	if props.BillingPlan != nil && props.BillingPlan.BillingPeriod != nil &&
		*props.BillingPlan.BillingPeriod != "" {
		v := *props.BillingPlan.BillingPeriod
		k.billingPeriod = &v
	}
	if props.Identity != nil && props.Identity.ClientID != nil {
		v := *props.Identity.ClientID
		k.identityClientID = &v
		// ClientSecret is not returned in the response — caller must re-set on Update.
	}
}

// kaasRebuildNodePools flattens response-side object types (InstanceResponse,
// DataCenterResponse) back to plain strings so toUpdateRequest() round-trips correctly.
func kaasRebuildNodePools(pools *[]types.NodePoolPropertiesResponse) []*NodePool {
	if pools == nil {
		return nil
	}
	result := make([]*NodePool, 0, len(*pools))
	for _, rp := range *pools {
		np := &NodePool{}
		if rp.Name != nil {
			v := *rp.Name
			np.name = &v
		}
		if rp.Nodes != nil {
			v := *rp.Nodes
			np.nodes = &v
		}
		if rp.Instance != nil && rp.Instance.Name != nil {
			v := *rp.Instance.Name
			np.instance = &v
		}
		if rp.DataCenter != nil && rp.DataCenter.Code != nil {
			v := Zone(*rp.DataCenter.Code)
			np.zone = &v
		}
		if rp.MinCount != nil {
			v := *rp.MinCount
			np.minCount = &v
		}
		if rp.MaxCount != nil {
			v := *rp.MaxCount
			np.maxCount = &v
		}
		b := rp.Autoscaling
		np.autoscaling = &b
		result = append(result, np)
	}
	return result
}

func kaasDeref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ---------------------------------------------------------------------------
// kaasIDsFromRef
// ---------------------------------------------------------------------------

func kaasIDsFromRef(ref Ref) (projectID, kaasID string, err error) {
	kid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withKaaSID); ok {
			return w.KaaSID(), true
		}
		return "", false
	}, "kaas")
	if !ok || kid == "" {
		return "", "", fmt.Errorf("cannot determine KaaS ID from Ref %q", ref.URI())
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
	return pid, kid, nil
}

// ---------------------------------------------------------------------------
// kaasActions — internal interface for action dispatch
// ---------------------------------------------------------------------------

type kaasActions interface {
	downloadKubeconfig(ctx context.Context, projectID, kaasID string, rp *types.RequestParameters) (*types.Response[types.KaaSKubeconfigResponse], error)
}

var kaasTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
	"Failed": false,
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type kaasLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.KaaSList], error)
	Get(ctx context.Context, projectID, kaasID string, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	Create(ctx context.Context, projectID string, body types.KaaSRequest, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	Update(ctx context.Context, projectID, kaasID string, body types.KaaSUpdateRequest, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	Delete(ctx context.Context, projectID, kaasID string, params *types.RequestParameters) (*types.Response[any], error)
	DownloadKubeconfig(ctx context.Context, projectID, kaasID string, params *types.RequestParameters) (*types.Response[types.KaaSKubeconfigResponse], error)
}

type kaasClientAdapter struct {
	low kaasLowLevelClient
}

var _ kaasActions = (*kaasClientAdapter)(nil)

func newKaaSClientAdapter(rest *restclient.Client) *kaasClientAdapter {
	if rest == nil {
		return &kaasClientAdapter{}
	}
	return &kaasClientAdapter{low: container.NewKaaSClientImpl(rest)}
}

func (a *kaasClientAdapter) Create(ctx context.Context, k *KaaS, opts ...CallOption) (*KaaS, error) {
	if err := k.Err(); err != nil {
		return k, err
	}
	if k.ProjectID() == "" {
		return k, fmt.Errorf("Create: KaaS has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, k.ProjectID(), k.toRequest(), rp)
	populateHTTPEnvelope(&k.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		k.fromResponse(resp.Data)
		k.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, k)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				k.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	k.actions = a
	if err != nil {
		return k, err
	}
	if resp != nil && !resp.IsSuccess() {
		return k, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return k, nil
}

func (a *kaasClientAdapter) Update(ctx context.Context, k *KaaS, opts ...CallOption) (*KaaS, error) {
	if err := k.Err(); err != nil {
		return k, err
	}
	if k.KaaSID() == "" {
		return k, fmt.Errorf("Update: KaaS has no ID")
	}
	if k.ProjectID() == "" {
		return k, fmt.Errorf("Update: KaaS has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, k.ProjectID(), k.KaaSID(), k.toUpdateRequest(), rp)
	populateHTTPEnvelope(&k.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		k.fromResponse(resp.Data)
		k.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, k)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				k.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	k.actions = a
	if err != nil {
		return k, err
	}
	if resp != nil && !resp.IsSuccess() {
		return k, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return k, nil
}

func (a *kaasClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*KaaS, error) {
	projectID, kaasID, err := kaasIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, kaasID, rp)
	out := &KaaS{}
	out.projectID = projectID
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
	out.actions = a
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *kaasClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, kaasID, err := kaasIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, kaasID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *kaasClientAdapter) List(ctx context.Context, parent Ref, opts ...CallOption) (*List[*KaaS], error) {
	projectID, err := projectIDFromRef(parent)
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
	var items []*KaaS
	if resp != nil && resp.Data != nil {
		items = make([]*KaaS, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			k := &KaaS{}
			k.projectID = projectID
			k.fromResponse(&resp.Data.Values[i])
			k.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, k)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					k.fromResponse(fresh.Raw())
				}
				return nil
			})
			if k.projectID == "" {
				k.projectID = projectID
			}
			k.actions = a
			items = append(items, k)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*KaaS], error) {
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

// downloadKubeconfig satisfies kaasActions (lowercase, internal interface).
func (a *kaasClientAdapter) downloadKubeconfig(ctx context.Context, projectID, kaasID string, rp *types.RequestParameters) (*types.Response[types.KaaSKubeconfigResponse], error) {
	return a.low.DownloadKubeconfig(ctx, projectID, kaasID, rp)
}

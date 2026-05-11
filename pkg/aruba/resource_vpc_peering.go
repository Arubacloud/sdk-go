package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// ---- Wrapper ----

// VPCPeering is the wrapper for an Aruba Cloud VPC Peering (a child of a VPC).
// Construct with aruba.NewVPCPeering() and bind it via IntoVPC(vpc).
//
// Wraps types.VPCPeeringResponse / types.VPCPeeringRequest. The wrapper carries
// pointer-typed private fields so unset values round-trip through
// the JSON layer correctly.
type VPCPeering struct {
	errMixin
	metadataMixin
	regionalMixin // VPCPeeringRequest.Metadata is RegionalResourceMetadataRequest
	vpcScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	remoteVPC *types.ReferenceResource
	response  *types.VPCPeeringResponse
}

// Setters — chainable, general → specific

// IntoVPC binds this VPCPeering to its parent VPC. Required before Create.
func (p *VPCPeering) IntoVPC(v Ref) *VPCPeering { p.intoVPC(v); return p }

// WithName sets the resource name. Required by the API.
func (p *VPCPeering) WithName(n string) *VPCPeering { p.withName(n); return p }

// AddTag appends a tag for filtering and accounting.
func (p *VPCPeering) AddTag(t string) *VPCPeering { p.addTag(t); return p }

// RemoveTag removes a previously-added tag. No-op if absent.
func (p *VPCPeering) RemoveTag(t string) *VPCPeering { p.removeTag(t); return p }

// ReplaceTags replaces the entire tag set with the given values.
func (p *VPCPeering) ReplaceTags(ts ...string) *VPCPeering { p.replaceTags(ts...); return p }

// InRegion sets the region for this resource.
func (p *VPCPeering) InRegion(region Region) *VPCPeering { p.inRegion(region); return p }

// WithRemoteVPC stores the remote VPC URI in the request body Properties.
// Records a setter-time error if the supplied Ref's URI is empty.
func (p *VPCPeering) WithRemoteVPC(v Ref) *VPCPeering {
	uri := v.URI()
	if uri == "" {
		p.addErr(fmt.Errorf("WithRemoteVPC: remote VPC Ref has empty URI"))
		return p
	}
	p.remoteVPC = &types.ReferenceResource{URI: uri}
	return p
}

// Getters — general → specific

// URI satisfies Ref.
func (p *VPCPeering) URI() string { return p.RespURI() }

// VPCPeeringID satisfies withVPCPeeringID so child wrappers (VPCPeeringRoute)
// can extract this ID via typed assertion.
func (p *VPCPeering) VPCPeeringID() string { return p.ID() }

// Raw shadows responseMetadataMixin.Raw() with the full VPCPeering response.
func (p *VPCPeering) Raw() *types.VPCPeeringResponse { return p.response }

// RawRequest returns what toRequest() would emit right now.
func (p *VPCPeering) RawRequest() types.VPCPeeringRequest { return p.toRequest() }

// RemoteVPCURI returns the configured remote VPC URI ("" if unset).
func (p *VPCPeering) RemoteVPCURI() string {
	if p.remoteVPC == nil {
		return ""
	}
	return p.remoteVPC.URI
}

// Wire converters

// toRequest assembles the Create/Update body from current setter state. Defaults are applied at the wire boundary.
func (p *VPCPeering) toRequest() types.VPCPeeringRequest {
	props := types.VPCPeeringPropertiesRequest{}
	if p.remoteVPC != nil {
		props.RemoteVPC = p.remoteVPC
	}
	return types.VPCPeeringRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: p.toMetadata(),
			Location:                p.toLocation(),
		},
		Properties: props,
	}
}

// fromResponse hydrates the wrapper from a server reply. Nil-safe.
func (p *VPCPeering) fromResponse(resp *types.VPCPeeringResponse) {
	if resp == nil {
		return
	}
	p.response = resp
	p.setMeta(&resp.Metadata)
	p.withName(vpcPeeringDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		p.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		p.inRegion(resp.Metadata.LocationResponse.Value)
	}
	p.setStatus(&resp.Status)
	p.setTerminalStates(vpcPeeringTerminalStates)
	p.setLinked(resp.Properties.LinkedResources)

	if resp.Properties.RemoteVPC != nil {
		rv := *resp.Properties.RemoteVPC
		p.remoteVPC = &rv
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		p.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if (p.vpcID == "" || p.projectID == "") && p.RespURI() != "" {
		ids := parseURIIDs(p.RespURI())
		if p.vpcID == "" {
			p.vpcID = ids["vpcs"]
		}
		if p.projectID == "" {
			p.projectID = ids["projects"]
		}
	}
}

func vpcPeeringDerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

var vpcPeeringTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
	"Failed": false,
}

// ---- Low-level client interface ----

// vpcPeeringLowLevelClient is the contract the wrapper depends on. Returning
// *types.Response[T] preserves HTTP envelope details (status code, headers,
// raw body) for the wrapper's diagnostics.
type vpcPeeringLowLevelClient interface {
	List(ctx context.Context, projectID, vpcID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringList], error)
	Get(ctx context.Context, projectID, vpcID, vpcPeeringID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Create(ctx context.Context, projectID, vpcID string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Update(ctx context.Context, projectID, vpcID, vpcPeeringID string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Delete(ctx context.Context, projectID, vpcID, vpcPeeringID string, params *types.RequestParameters) (*types.Response[any], error)
}

// ---- Adapter ----

// vpcPeeringsClientAdapter bridges the wrapper API (chainable, error-accumulating,
// wire-shape-hidden) to the low-level client (parameter-explicit, returning
// typed wire structs). Translates VPCPeering ↔ types.VPCPeeringRequest/Response and
// surfaces HTTP errors as *aruba.HTTPError.
type vpcPeeringsClientAdapter struct{ low vpcPeeringLowLevelClient }

var _ VPCPeeringsClient = (*vpcPeeringsClientAdapter)(nil)

func newVPCPeeringsClientAdapter(rest *restclient.Client) *vpcPeeringsClientAdapter {
	if rest == nil {
		return &vpcPeeringsClientAdapter{}
	}
	return &vpcPeeringsClientAdapter{low: network.NewVPCPeeringsClientImpl(rest)}
}

// Create posts a new VPCPeering to the API and hydrates the wrapper from the response.
func (a *vpcPeeringsClientAdapter) Create(ctx context.Context, peering *VPCPeering, opts ...CallOption) (*VPCPeering, error) {
	if err := peering.Err(); err != nil {
		return peering, err
	}
	if peering.VPCID() == "" || peering.ProjectID() == "" {
		return peering, fmt.Errorf("Create: VPC peering has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, peering.ProjectID(), peering.VPCID(), peering.toRequest(), rp)
	populateHTTPEnvelope(&peering.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		peering.fromResponse(resp.Data)
		peering.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, peering)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				peering.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return peering, err
	}
	if resp != nil && !resp.IsSuccess() {
		return peering, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return peering, nil
}

// Get fetches a VPCPeering by Ref and returns a freshly hydrated wrapper.
func (a *vpcPeeringsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPCPeering, error) {
	projectID, vpcID, vpcPeeringID, err := vpcPeeringIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, vpcPeeringID, rp)
	out := &VPCPeering{}
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
	if out.vpcID == "" {
		out.vpcID = vpcID
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

// Update sends a PUT for the current wrapper state. Requires ID and parent.
func (a *vpcPeeringsClientAdapter) Update(ctx context.Context, peering *VPCPeering, opts ...CallOption) (*VPCPeering, error) {
	if err := peering.Err(); err != nil {
		return peering, err
	}
	if peering.ID() == "" {
		return peering, fmt.Errorf("Update: VPC peering has no ID — call Get first or seed from response metadata")
	}
	if peering.VPCID() == "" || peering.ProjectID() == "" {
		return peering, fmt.Errorf("Update: VPC peering has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, peering.ProjectID(), peering.VPCID(), peering.ID(), peering.toRequest(), rp)
	populateHTTPEnvelope(&peering.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		peering.fromResponse(resp.Data)
		peering.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, peering)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				peering.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return peering, err
	}
	if resp != nil && !resp.IsSuccess() {
		return peering, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return peering, nil
}

// Delete removes the VPCPeering identified by Ref.
func (a *vpcPeeringsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, vpcPeeringID, err := vpcPeeringIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, vpcPeeringID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

// List returns a paginated list of VPCPeering in the given parent scope.
func (a *vpcPeeringsClientAdapter) List(ctx context.Context, vpc Ref, opts ...CallOption) (*List[*VPCPeering], error) {
	projectID, vpcID, err := vpcIDsFromRef(vpc)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, vpcID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*VPCPeering
	if resp != nil && resp.Data != nil {
		items = make([]*VPCPeering, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			p := &VPCPeering{}
			p.fromResponse(&resp.Data.Values[i])
			p.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, p)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					p.fromResponse(fresh.Raw())
				}
				return nil
			})
			if p.vpcID == "" {
				p.vpcID = vpcID
			}
			if p.projectID == "" {
				p.projectID = projectID
			}
			items = append(items, p)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*VPCPeering], error) {
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

// vpcPeeringIDsFromRef extracts (projectID, vpcID, vpcPeeringID) from a Ref.
// Accepts the production camelCase segment "vpcPeerings" and the mixin/test form "peerings".
func vpcPeeringIDsFromRef(ref Ref) (projectID, vpcID, vpcPeeringID string, err error) {
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCPeeringID); ok {
			return w.VPCPeeringID(), true
		}
		return "", false
	}, "vpc-peerings")
	if !ok || pid == "" {
		m := parseURIIDs(ref.URI())
		if v := m["vpcPeerings"]; v != "" {
			pid = v
			ok = true
		}
		if pid == "" {
			if v := m["peerings"]; v != "" {
				pid = v
				ok = true
			}
		}
	}
	if !ok || pid == "" {
		return "", "", "", fmt.Errorf("cannot determine VPC peering ID from Ref %q", ref.URI())
	}
	vid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid == "" {
		return "", "", "", fmt.Errorf("cannot determine VPC ID from Ref %q", ref.URI())
	}
	projID, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || projID == "" {
		return "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return projID, vid, pid, nil
}

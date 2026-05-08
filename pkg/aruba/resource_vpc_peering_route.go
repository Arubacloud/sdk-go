package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// VPCPeeringRoute is the wrapper for an Aruba Cloud VPC Peering Route
// (a direct child of a VPCPeering, grandchild of a VPC).
// Construct with aruba.NewVPCPeeringRoute() and bind it via IntoVPCPeering(peering).
type VPCPeeringRoute struct {
	errMixin
	metadataMixin
	regionalMixin
	vpcPeeringScopedMixin
	responseMetadataMixin
	statusMixin
	httpEnvelopeMixin

	localCIDR     *string
	remoteCIDR    *string
	billingPeriod *BillingPeriod
	response      *types.VPCPeeringRouteResponse
}

// Setters (chainable).

func (r *VPCPeeringRoute) IntoVPCPeering(p Ref) *VPCPeeringRoute     { r.intoVPCPeering(p); return r }
func (r *VPCPeeringRoute) WithName(n string) *VPCPeeringRoute        { r.withName(n); return r }
func (r *VPCPeeringRoute) AddTag(t string) *VPCPeeringRoute          { r.addTag(t); return r }
func (r *VPCPeeringRoute) RemoveTag(t string) *VPCPeeringRoute       { r.removeTag(t); return r }
func (r *VPCPeeringRoute) ReplaceTags(ts ...string) *VPCPeeringRoute { r.replaceTags(ts...); return r }
func (r *VPCPeeringRoute) WithLocation(loc Region) *VPCPeeringRoute  { r.withLocation(loc); return r }
func (r *VPCPeeringRoute) InRegion(region Region) *VPCPeeringRoute   { r.inRegion(region); return r }

func (r *VPCPeeringRoute) WithLocalCIDR(cidr string) *VPCPeeringRoute { r.localCIDR = &cidr; return r }
func (r *VPCPeeringRoute) WithRemoteCIDR(cidr string) *VPCPeeringRoute {
	r.remoteCIDR = &cidr
	return r
}
func (r *VPCPeeringRoute) WithBillingPeriod(p BillingPeriod) *VPCPeeringRoute {
	r.billingPeriod = &p
	return r
}

// URI satisfies Ref.
func (r *VPCPeeringRoute) URI() string { return r.RespURI() }

// VPCPeeringRouteID satisfies withVPCPeeringRouteID.
func (r *VPCPeeringRoute) VPCPeeringRouteID() string { return r.ID() }

// Raw shadows responseMetadataMixin.Raw() with the typed VPC peering route response.
func (r *VPCPeeringRoute) Raw() *types.VPCPeeringRouteResponse { return r.response }

// RawRequest returns what toRequest() would emit right now.
func (r *VPCPeeringRoute) RawRequest() types.VPCPeeringRouteRequest { return r.toRequest() }

// LocalCIDR returns the configured local network CIDR ("" if unset).
func (r *VPCPeeringRoute) LocalCIDR() string {
	if r.localCIDR == nil {
		return ""
	}
	return *r.localCIDR
}

// RemoteCIDR returns the configured remote network CIDR ("" if unset).
func (r *VPCPeeringRoute) RemoteCIDR() string {
	if r.remoteCIDR == nil {
		return ""
	}
	return *r.remoteCIDR
}

// BillingPeriod returns the configured billing period ("" if unset).
func (r *VPCPeeringRoute) BillingPeriod() BillingPeriod {
	if r.billingPeriod == nil {
		return ""
	}
	return *r.billingPeriod
}

func (r *VPCPeeringRoute) toRequest() types.VPCPeeringRouteRequest {
	var bp BillingPeriod
	if r.billingPeriod != nil {
		bp = *r.billingPeriod
	}
	props := types.VPCPeeringRoutePropertiesRequest{
		BillingPlan: types.BillingPeriodResource{BillingPeriod: bp},
	}
	if r.localCIDR != nil {
		props.LocalNetworkAddress = *r.localCIDR
	}
	if r.remoteCIDR != nil {
		props.RemoteNetworkAddress = *r.remoteCIDR
	}
	return types.VPCPeeringRouteRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: r.toMetadata(),
			Location:                r.toLocation(),
		},
		Properties: props,
	}
}

func (r *VPCPeeringRoute) fromResponse(resp *types.VPCPeeringRouteResponse) {
	if resp == nil {
		return
	}
	r.response = resp
	r.setMeta(&resp.Metadata)
	r.withName(vpcPeeringRouteDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		r.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		r.withLocation(resp.Metadata.LocationResponse.Value)
	}
	r.setStatus(&resp.Status)
	r.setTerminalStates(vpcPeeringRouteTerminalStates)

	if resp.Properties.LocalNetworkAddress != "" {
		v := resp.Properties.LocalNetworkAddress
		r.localCIDR = &v
	}
	if resp.Properties.RemoteNetworkAddress != "" {
		v := resp.Properties.RemoteNetworkAddress
		r.remoteCIDR = &v
	}
	if resp.Properties.BillingPlan.BillingPeriod != "" {
		v := resp.Properties.BillingPlan.BillingPeriod
		r.billingPeriod = &v
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		r.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if (r.vpcID == "" || r.projectID == "" || r.vpcPeeringID == "") && r.RespURI() != "" {
		ids := parseURIIDs(r.RespURI())
		if r.vpcID == "" {
			r.vpcID = ids["vpcs"]
		}
		if r.projectID == "" {
			r.projectID = ids["projects"]
		}
		if r.vpcPeeringID == "" {
			// Production URI uses "vpcPeerings"; mixin/test URIs use "peerings".
			if v := ids["vpcPeerings"]; v != "" {
				r.vpcPeeringID = v
			}
			if r.vpcPeeringID == "" {
				r.vpcPeeringID = ids["peerings"]
			}
		}
	}
}

func vpcPeeringRouteDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

var vpcPeeringRouteTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
	"Failed": false,
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type vpcPeeringRouteLowLevelClient interface {
	List(ctx context.Context, projectID, vpcID, vpcPeeringID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteList], error)
	Get(ctx context.Context, projectID, vpcID, vpcPeeringID, vpcPeeringRouteID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Create(ctx context.Context, projectID, vpcID, vpcPeeringID string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Update(ctx context.Context, projectID, vpcID, vpcPeeringID, vpcPeeringRouteID string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Delete(ctx context.Context, projectID, vpcID, vpcPeeringID, vpcPeeringRouteID string, params *types.RequestParameters) (*types.Response[any], error)
}

type vpcPeeringRoutesClientAdapter struct{ low vpcPeeringRouteLowLevelClient }

func newVPCPeeringRoutesClientAdapter(rest *restclient.Client) *vpcPeeringRoutesClientAdapter {
	if rest == nil {
		return &vpcPeeringRoutesClientAdapter{}
	}
	return &vpcPeeringRoutesClientAdapter{low: network.NewVPCPeeringRoutesClientImpl(rest)}
}

func (a *vpcPeeringRoutesClientAdapter) Create(ctx context.Context, route *VPCPeeringRoute, opts ...CallOption) (*VPCPeeringRoute, error) {
	if err := route.Err(); err != nil {
		return route, err
	}
	if route.VPCPeeringID() == "" || route.VPCID() == "" || route.ProjectID() == "" {
		return route, fmt.Errorf("Create: VPC peering route has no parent peering — call IntoVPCPeering first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, route.ProjectID(), route.VPCID(), route.VPCPeeringID(), route.toRequest(), rp)
	populateHTTPEnvelope(&route.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		route.fromResponse(resp.Data)
		route.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, route)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				route.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return route, err
	}
	if resp != nil && !resp.IsSuccess() {
		return route, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return route, nil
}

func (a *vpcPeeringRoutesClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPCPeeringRoute, error) {
	projectID, vpcID, vpcPeeringID, routeID, err := vpcPeeringRouteIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, vpcPeeringID, routeID, rp)
	out := &VPCPeeringRoute{}
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
	if out.vpcPeeringID == "" {
		out.vpcPeeringID = vpcPeeringID
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

func (a *vpcPeeringRoutesClientAdapter) Update(ctx context.Context, route *VPCPeeringRoute, opts ...CallOption) (*VPCPeeringRoute, error) {
	if err := route.Err(); err != nil {
		return route, err
	}
	if route.ID() == "" {
		return route, fmt.Errorf("Update: VPC peering route has no ID — call Get first or seed from response metadata")
	}
	if route.VPCPeeringID() == "" || route.VPCID() == "" || route.ProjectID() == "" {
		return route, fmt.Errorf("Update: VPC peering route has no parent peering — call IntoVPCPeering first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, route.ProjectID(), route.VPCID(), route.VPCPeeringID(), route.ID(), route.toRequest(), rp)
	populateHTTPEnvelope(&route.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		route.fromResponse(resp.Data)
		route.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, route)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				route.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return route, err
	}
	if resp != nil && !resp.IsSuccess() {
		return route, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return route, nil
}

func (a *vpcPeeringRoutesClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, vpcPeeringID, routeID, err := vpcPeeringRouteIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, vpcPeeringID, routeID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *vpcPeeringRoutesClientAdapter) List(ctx context.Context, peering Ref, opts ...CallOption) (*List[*VPCPeeringRoute], error) {
	projectID, vpcID, vpcPeeringID, err := vpcPeeringIDsFromRef(peering)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, vpcID, vpcPeeringID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*VPCPeeringRoute
	if resp != nil && resp.Data != nil {
		items = make([]*VPCPeeringRoute, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			r := &VPCPeeringRoute{}
			r.fromResponse(&resp.Data.Values[i])
			r.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, r)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					r.fromResponse(fresh.Raw())
				}
				return nil
			})
			if r.vpcPeeringID == "" {
				r.vpcPeeringID = vpcPeeringID
			}
			if r.vpcID == "" {
				r.vpcID = vpcID
			}
			if r.projectID == "" {
				r.projectID = projectID
			}
			items = append(items, r)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*VPCPeeringRoute], error) {
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

// vpcPeeringRouteIDsFromRef extracts (projectID, vpcID, vpcPeeringID, vpcPeeringRouteID) from a Ref.
// Accepts the production camelCase segment "vpcPeeringRoutes" and the test form "vpc-peering-routes".
// For the peering parent, accepts both "vpcPeerings" and "peerings".
func vpcPeeringRouteIDsFromRef(ref Ref) (projectID, vpcID, vpcPeeringID, vpcPeeringRouteID string, err error) {
	rid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCPeeringRouteID); ok {
			return w.VPCPeeringRouteID(), true
		}
		return "", false
	}, "vpc-peering-routes")
	if !ok {
		if v := parseURIIDs(ref.URI())["vpcPeeringRoutes"]; v != "" {
			rid = v
			ok = true
		}
	}
	if !ok {
		return "", "", "", "", fmt.Errorf("cannot determine VPC peering route ID from Ref %q", ref.URI())
	}
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCPeeringID); ok {
			return w.VPCPeeringID(), true
		}
		return "", false
	}, "vpc-peerings")
	if !ok {
		m := parseURIIDs(ref.URI())
		if v := m["vpcPeerings"]; v != "" {
			pid = v
			ok = true
		} else if v := m["peerings"]; v != "" {
			pid = v
			ok = true
		}
	}
	if !ok {
		return "", "", "", "", fmt.Errorf("cannot determine VPC peering ID from Ref %q", ref.URI())
	}
	vid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid == "" {
		return "", "", "", "", fmt.Errorf("cannot determine VPC ID from Ref %q", ref.URI())
	}
	projID, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || projID == "" {
		return "", "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return projID, vid, pid, rid, nil
}

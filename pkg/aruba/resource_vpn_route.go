package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// VPNRoute is the wrapper for an Aruba Cloud VPN Route (a direct child of a VPNTunnel).
// Construct with NewVPNRoute() and bind it via IntoVPNTunnel(tunnel).
type VPNRoute struct {
	errMixin
	metadataMixin
	regionalMixin
	vpnTunnelScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	cloudSubnet  *string
	onPremSubnet *string

	response *types.VPNRouteResponse
}

// Setters (chainable).

func (r *VPNRoute) IntoVPNTunnel(t Ref) *VPNRoute        { r.intoVPNTunnel(t); return r }
func (r *VPNRoute) WithName(n string) *VPNRoute          { r.withName(n); return r }
func (r *VPNRoute) AddTag(tag string) *VPNRoute          { r.addTag(tag); return r }
func (r *VPNRoute) RemoveTag(tag string) *VPNRoute       { r.removeTag(tag); return r }
func (r *VPNRoute) ReplaceTags(tags ...string) *VPNRoute { r.replaceTags(tags...); return r }
func (r *VPNRoute) WithLocation(loc Region) *VPNRoute    { r.withLocation(loc); return r }
func (r *VPNRoute) InRegion(region Region) *VPNRoute     { r.inRegion(region); return r }

func (r *VPNRoute) WithCloudSubnet(cidr string) *VPNRoute  { r.cloudSubnet = &cidr; return r }
func (r *VPNRoute) WithOnPremSubnet(cidr string) *VPNRoute { r.onPremSubnet = &cidr; return r }

// URI satisfies Ref.
func (r *VPNRoute) URI() string { return r.RespURI() }

// VPNRouteID satisfies withVPNRouteID.
func (r *VPNRoute) VPNRouteID() string { return r.ID() }

// Raw shadows responseMetadataMixin.Raw() with the typed VPN route response.
func (r *VPNRoute) Raw() *types.VPNRouteResponse { return r.response }

// RawRequest returns what toRequest() would emit right now.
func (r *VPNRoute) RawRequest() types.VPNRouteRequest { return r.toRequest() }

// CloudSubnet returns the configured cloud subnet CIDR ("" if unset).
func (r *VPNRoute) CloudSubnet() string { return vpnRouteDerefString(r.cloudSubnet) }

// OnPremSubnet returns the configured on-premises subnet CIDR ("" if unset).
func (r *VPNRoute) OnPremSubnet() string { return vpnRouteDerefString(r.onPremSubnet) }

func (r *VPNRoute) toRequest() types.VPNRouteRequest {
	props := types.VPNRoutePropertiesRequest{
		CloudSubnet:  vpnRouteDerefString(r.cloudSubnet),
		OnPremSubnet: vpnRouteDerefString(r.onPremSubnet),
	}
	return types.VPNRouteRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: r.toMetadata(),
			Location:                r.toLocation(),
		},
		Properties: props,
	}
}

func (r *VPNRoute) fromResponse(resp *types.VPNRouteResponse) {
	if resp == nil {
		return
	}
	r.response = resp
	r.setMeta(&resp.Metadata)
	r.withName(vpnRouteDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		r.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		r.withLocation(resp.Metadata.LocationResponse.Value)
	}
	r.setStatus(&resp.Status)
	r.setTerminalStates(vpnRouteTerminalStates)
	if len(resp.Properties.LinkedResources) > 0 {
		r.setLinked(resp.Properties.LinkedResources)
	}
	if resp.Properties.CloudSubnet != "" {
		v := resp.Properties.CloudSubnet
		r.cloudSubnet = &v
	}
	if resp.Properties.OnPremSubnet != "" {
		v := resp.Properties.OnPremSubnet
		r.onPremSubnet = &v
	}
	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		r.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if (r.projectID == "" || r.vpnTunnelID == "") && r.RespURI() != "" {
		ids := parseURIIDs(r.RespURI())
		if r.projectID == "" {
			r.projectID = ids["projects"]
		}
		if r.vpnTunnelID == "" {
			// Production URI uses "vpnTunnels"; mixin/test form uses "vpn-tunnels".
			if v := ids["vpnTunnels"]; v != "" {
				r.vpnTunnelID = v
			} else {
				r.vpnTunnelID = ids["vpn-tunnels"]
			}
		}
	}
}

func vpnRouteDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

var vpnRouteTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
	"Failed": false,
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type vpnRouteLowLevelClient interface {
	List(ctx context.Context, projectID, vpnTunnelID string, params *types.RequestParameters) (*types.Response[types.VPNRouteList], error)
	Get(ctx context.Context, projectID, vpnTunnelID, vpnRouteID string, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	Create(ctx context.Context, projectID, vpnTunnelID string, body types.VPNRouteRequest, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	Update(ctx context.Context, projectID, vpnTunnelID, vpnRouteID string, body types.VPNRouteRequest, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	Delete(ctx context.Context, projectID, vpnTunnelID, vpnRouteID string, params *types.RequestParameters) (*types.Response[any], error)
}

type vpnRoutesClientAdapter struct{ low vpnRouteLowLevelClient }

func newVPNRoutesClientAdapter(rest *restclient.Client) *vpnRoutesClientAdapter {
	if rest == nil {
		return &vpnRoutesClientAdapter{}
	}
	return &vpnRoutesClientAdapter{low: network.NewVPNRoutesClientImpl(rest)}
}

func (a *vpnRoutesClientAdapter) Create(ctx context.Context, r *VPNRoute, opts ...CallOption) (*VPNRoute, error) {
	if err := r.Err(); err != nil {
		return r, err
	}
	if r.VPNTunnelID() == "" || r.ProjectID() == "" {
		return r, fmt.Errorf("Create: VPN route has no parent tunnel — call IntoVPNTunnel first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, r.ProjectID(), r.VPNTunnelID(), r.toRequest(), rp)
	populateHTTPEnvelope(&r.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		r.fromResponse(resp.Data)
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
	}
	if err != nil {
		return r, err
	}
	if resp != nil && !resp.IsSuccess() {
		return r, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return r, nil
}

func (a *vpnRoutesClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPNRoute, error) {
	projectID, vpnTunnelID, vpnRouteID, err := vpnRouteIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpnTunnelID, vpnRouteID, rp)
	out := &VPNRoute{}
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
	if out.vpnTunnelID == "" {
		out.vpnTunnelID = vpnTunnelID
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

func (a *vpnRoutesClientAdapter) Update(ctx context.Context, r *VPNRoute, opts ...CallOption) (*VPNRoute, error) {
	if err := r.Err(); err != nil {
		return r, err
	}
	if r.ID() == "" {
		return r, fmt.Errorf("Update: VPN route has no ID — call Get first or seed from response metadata")
	}
	if r.VPNTunnelID() == "" || r.ProjectID() == "" {
		return r, fmt.Errorf("Update: VPN route has no parent tunnel — call IntoVPNTunnel first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, r.ProjectID(), r.VPNTunnelID(), r.ID(), r.toRequest(), rp)
	populateHTTPEnvelope(&r.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		r.fromResponse(resp.Data)
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
	}
	if err != nil {
		return r, err
	}
	if resp != nil && !resp.IsSuccess() {
		return r, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return r, nil
}

func (a *vpnRoutesClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpnTunnelID, vpnRouteID, err := vpnRouteIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpnTunnelID, vpnRouteID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *vpnRoutesClientAdapter) List(ctx context.Context, tunnel Ref, opts ...CallOption) (*List[*VPNRoute], error) {
	projectID, vpnTunnelID, err := vpnTunnelIDsFromRef(tunnel)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, vpnTunnelID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*VPNRoute
	if resp != nil && resp.Data != nil {
		items = make([]*VPNRoute, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			v := &VPNRoute{}
			v.fromResponse(&resp.Data.Values[i])
			v.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, v)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					v.fromResponse(fresh.Raw())
				}
				return nil
			})
			if v.vpnTunnelID == "" {
				v.vpnTunnelID = vpnTunnelID
			}
			if v.projectID == "" {
				v.projectID = projectID
			}
			items = append(items, v)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*VPNRoute], error) {
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

// vpnRouteIDsFromRef extracts (projectID, vpnTunnelID, vpnRouteID) from a Ref.
// Accepts camelCase segments ("vpnTunnels", "vpnRoutes") and kebab-case forms.
func vpnRouteIDsFromRef(ref Ref) (projectID, vpnTunnelID, vpnRouteID string, err error) {
	rid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPNRouteID); ok {
			return w.VPNRouteID(), true
		}
		return "", false
	}, "vpn-routes")
	if !ok {
		if v := parseURIIDs(ref.URI())["vpnRoutes"]; v != "" {
			rid = v
			ok = true
		}
	}
	if !ok || rid == "" {
		return "", "", "", fmt.Errorf("cannot determine VPN route ID from Ref %q", ref.URI())
	}
	tid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPNTunnelID); ok {
			return w.VPNTunnelID(), true
		}
		return "", false
	}, "vpn-tunnels")
	if !ok {
		if v := parseURIIDs(ref.URI())["vpnTunnels"]; v != "" {
			tid = v
			ok = true
		}
	}
	if !ok || tid == "" {
		return "", "", "", fmt.Errorf("cannot determine VPN tunnel ID from Ref %q", ref.URI())
	}
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid == "" {
		return "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return pid, tid, rid, nil
}

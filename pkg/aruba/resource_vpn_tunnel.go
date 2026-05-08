package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// VPNTunnel is the wrapper for an Aruba Cloud VPN Tunnel (direct child of a project).
// Construct with NewVPNTunnel() and bind it via IntoProject(project).
type VPNTunnel struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	vpnType            *string
	vpnClientProtocol  *string
	billingPeriod      *string
	peerClientPublicIP *string

	ipConfig *VPNIPConfig
	ike      *VPNIKE
	esp      *VPNESP
	psk      *VPNPSK

	response *types.VPNTunnelResponse
}

// Setters (chainable).

func (t *VPNTunnel) IntoProject(p Ref) *VPNTunnel               { t.intoProject(p); return t }
func (t *VPNTunnel) WithName(n string) *VPNTunnel               { t.withName(n); return t }
func (t *VPNTunnel) AddTag(tag string) *VPNTunnel               { t.addTag(tag); return t }
func (t *VPNTunnel) RemoveTag(tag string) *VPNTunnel            { t.removeTag(tag); return t }
func (t *VPNTunnel) ReplaceTags(tags ...string) *VPNTunnel      { t.replaceTags(tags...); return t }
func (t *VPNTunnel) WithLocation(loc string) *VPNTunnel         { t.withLocation(loc); return t }
func (t *VPNTunnel) InRegion(r string) *VPNTunnel               { t.inRegion(r); return t }
func (t *VPNTunnel) WithVPNType(s string) *VPNTunnel            { t.vpnType = &s; return t }
func (t *VPNTunnel) WithVPNClientProtocol(s string) *VPNTunnel  { t.vpnClientProtocol = &s; return t }
func (t *VPNTunnel) WithBillingPeriod(s string) *VPNTunnel      { t.billingPeriod = &s; return t }
func (t *VPNTunnel) WithPeerClientPublicIP(s string) *VPNTunnel { t.peerClientPublicIP = &s; return t }

func (t *VPNTunnel) WithIPConfig(c *VPNIPConfig) *VPNTunnel {
	t.ipConfig = c
	if c != nil {
		for _, e := range c.errs {
			t.addErr(e)
		}
	}
	return t
}

func (t *VPNTunnel) WithIKESettings(k *VPNIKE) *VPNTunnel {
	t.ike = k
	if k != nil {
		for _, e := range k.errs {
			t.addErr(e)
		}
	}
	return t
}

func (t *VPNTunnel) WithESPSettings(e *VPNESP) *VPNTunnel {
	t.esp = e
	if e != nil {
		for _, err := range e.errs {
			t.addErr(err)
		}
	}
	return t
}

func (t *VPNTunnel) WithPSKSettings(p *VPNPSK) *VPNTunnel {
	t.psk = p
	if p != nil {
		for _, e := range p.errs {
			t.addErr(e)
		}
	}
	return t
}

// URI satisfies Ref.
func (t *VPNTunnel) URI() string { return t.RespURI() }

// VPNTunnelID satisfies withVPNTunnelID.
func (t *VPNTunnel) VPNTunnelID() string { return t.ID() }

// Raw shadows responseMetadataMixin.Raw() with the typed VPN tunnel response.
func (t *VPNTunnel) Raw() *types.VPNTunnelResponse { return t.response }

// RawRequest returns what toRequest() would emit right now.
func (t *VPNTunnel) RawRequest() types.VPNTunnelRequest { return t.toRequest() }

// Read accessors.

func (t *VPNTunnel) IPConfig() *VPNIPConfig     { return t.ipConfig }
func (t *VPNTunnel) IKE() *VPNIKE               { return t.ike }
func (t *VPNTunnel) ESP() *VPNESP               { return t.esp }
func (t *VPNTunnel) PSK() *VPNPSK               { return t.psk }
func (t *VPNTunnel) VPNType() string            { return vpnTunnelDerefString(t.vpnType) }
func (t *VPNTunnel) VPNClientProtocol() string  { return vpnTunnelDerefString(t.vpnClientProtocol) }
func (t *VPNTunnel) BillingPeriod() string      { return vpnTunnelDerefString(t.billingPeriod) }
func (t *VPNTunnel) PeerClientPublicIP() string { return vpnTunnelDerefString(t.peerClientPublicIP) }

func (t *VPNTunnel) toRequest() types.VPNTunnelRequest {
	props := types.VPNTunnelPropertiesRequest{
		VPNType:           t.vpnType,
		VPNClientProtocol: t.vpnClientProtocol,
	}
	if t.ipConfig != nil {
		props.IPConfigurations = t.ipConfig.build()
	}
	if t.ike != nil || t.esp != nil || t.psk != nil || t.peerClientPublicIP != nil {
		cs := &types.VPNClientSettings{PeerClientPublicIP: t.peerClientPublicIP}
		if t.ike != nil {
			cs.IKE = t.ike.build()
		}
		if t.esp != nil {
			cs.ESP = t.esp.build()
		}
		if t.psk != nil {
			cs.PSK = t.psk.build()
		}
		props.VPNClientSettings = cs
	}
	if t.billingPeriod != nil {
		props.BillingPlan = &types.BillingPeriodResource{BillingPeriod: *t.billingPeriod}
	}
	return types.VPNTunnelRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: t.toMetadata(),
			Location:                t.toLocation(),
		},
		Properties: props,
	}
}

func (t *VPNTunnel) fromResponse(resp *types.VPNTunnelResponse) {
	if resp == nil {
		return
	}
	t.response = resp
	t.setMeta(&resp.Metadata)
	t.withName(vpnTunnelDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		t.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		t.withLocation(resp.Metadata.LocationResponse.Value)
	}
	t.setStatus(&resp.Status)
	t.setTerminalStates(vpnTunnelTerminalStates)

	if resp.Properties.VPNType != nil {
		v := *resp.Properties.VPNType
		t.vpnType = &v
	}
	if resp.Properties.VPNClientProtocol != nil {
		v := *resp.Properties.VPNClientProtocol
		t.vpnClientProtocol = &v
	}
	if resp.Properties.BillingPlan != nil && resp.Properties.BillingPlan.BillingPeriod != "" {
		bp := resp.Properties.BillingPlan.BillingPeriod
		t.billingPeriod = &bp
	}
	if cs := resp.Properties.VPNClientSettings; cs != nil && cs.PeerClientPublicIP != nil {
		v := *cs.PeerClientPublicIP
		t.peerClientPublicIP = &v
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		t.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if t.projectID == "" && t.RespURI() != "" {
		if id := parseURIIDs(t.RespURI())["projects"]; id != "" {
			t.projectID = id
		}
	}
}

func vpnTunnelDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

var vpnTunnelTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
	"Failed": false,
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type vpnTunnelLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.VPNTunnelList], error)
	Get(ctx context.Context, projectID, vpnTunnelID string, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	Create(ctx context.Context, projectID string, body types.VPNTunnelRequest, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	Update(ctx context.Context, projectID, vpnTunnelID string, body types.VPNTunnelRequest, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	Delete(ctx context.Context, projectID, vpnTunnelID string, params *types.RequestParameters) (*types.Response[any], error)
}

type vpnTunnelsClientAdapter struct{ low vpnTunnelLowLevelClient }

func newVPNTunnelsClientAdapter(rest *restclient.Client) *vpnTunnelsClientAdapter {
	if rest == nil {
		return &vpnTunnelsClientAdapter{}
	}
	return &vpnTunnelsClientAdapter{low: network.NewVPNTunnelsClientImpl(rest)}
}

func (a *vpnTunnelsClientAdapter) Create(ctx context.Context, t *VPNTunnel, opts ...CallOption) (*VPNTunnel, error) {
	if err := t.Err(); err != nil {
		return t, err
	}
	if t.ProjectID() == "" {
		return t, fmt.Errorf("Create: VPN tunnel has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, t.ProjectID(), t.toRequest(), rp)
	populateHTTPEnvelope(&t.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		t.fromResponse(resp.Data)
		t.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, t)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				t.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return t, err
	}
	if resp != nil && !resp.IsSuccess() {
		return t, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return t, nil
}

func (a *vpnTunnelsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPNTunnel, error) {
	projectID, vpnTunnelID, err := vpnTunnelIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpnTunnelID, rp)
	out := &VPNTunnel{}
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

func (a *vpnTunnelsClientAdapter) Update(ctx context.Context, t *VPNTunnel, opts ...CallOption) (*VPNTunnel, error) {
	if err := t.Err(); err != nil {
		return t, err
	}
	if t.ID() == "" {
		return t, fmt.Errorf("Update: VPN tunnel has no ID — call Get first or seed from response metadata")
	}
	if t.ProjectID() == "" {
		return t, fmt.Errorf("Update: VPN tunnel has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, t.ProjectID(), t.ID(), t.toRequest(), rp)
	populateHTTPEnvelope(&t.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		t.fromResponse(resp.Data)
		t.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, t)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				t.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return t, err
	}
	if resp != nil && !resp.IsSuccess() {
		return t, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return t, nil
}

func (a *vpnTunnelsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpnTunnelID, err := vpnTunnelIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpnTunnelID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *vpnTunnelsClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*VPNTunnel], error) {
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
	var items []*VPNTunnel
	if resp != nil && resp.Data != nil {
		items = make([]*VPNTunnel, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			v := &VPNTunnel{}
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
			if v.projectID == "" {
				v.projectID = projectID
			}
			items = append(items, v)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*VPNTunnel], error) {
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

// vpnTunnelIDsFromRef extracts (projectID, vpnTunnelID) from a Ref.
// Accepts the production camelCase segment "vpnTunnels" and the mixin/test form "vpn-tunnels".
func vpnTunnelIDsFromRef(ref Ref) (projectID, vpnTunnelID string, err error) {
	tid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPNTunnelID); ok {
			return w.VPNTunnelID(), true
		}
		return "", false
	}, "vpnTunnels")
	if !ok || tid == "" {
		if v := parseURIIDs(ref.URI())["vpn-tunnels"]; v != "" {
			tid = v
			ok = true
		}
	}
	if !ok || tid == "" {
		return "", "", fmt.Errorf("cannot determine VPN tunnel ID from Ref %q", ref.URI())
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
	return pid, tid, nil
}

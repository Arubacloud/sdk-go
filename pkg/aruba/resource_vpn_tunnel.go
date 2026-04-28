package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

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

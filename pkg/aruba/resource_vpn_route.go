package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

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
func (r *VPNRoute) WithLocation(loc string) *VPNRoute    { r.withLocation(loc); return r }
func (r *VPNRoute) InRegion(region string) *VPNRoute     { r.inRegion(region); return r }

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

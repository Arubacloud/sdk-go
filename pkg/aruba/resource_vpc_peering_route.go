package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

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
	billingPeriod *string
	response      *types.VPCPeeringRouteResponse
}

// Setters (chainable).

func (r *VPCPeeringRoute) IntoVPCPeering(p Ref) *VPCPeeringRoute     { r.intoVPCPeering(p); return r }
func (r *VPCPeeringRoute) WithName(n string) *VPCPeeringRoute        { r.withName(n); return r }
func (r *VPCPeeringRoute) AddTag(t string) *VPCPeeringRoute          { r.addTag(t); return r }
func (r *VPCPeeringRoute) RemoveTag(t string) *VPCPeeringRoute       { r.removeTag(t); return r }
func (r *VPCPeeringRoute) ReplaceTags(ts ...string) *VPCPeeringRoute { r.replaceTags(ts...); return r }
func (r *VPCPeeringRoute) WithLocation(loc string) *VPCPeeringRoute  { r.withLocation(loc); return r }
func (r *VPCPeeringRoute) InRegion(region string) *VPCPeeringRoute   { r.inRegion(region); return r }

func (r *VPCPeeringRoute) WithLocalCIDR(cidr string) *VPCPeeringRoute { r.localCIDR = &cidr; return r }
func (r *VPCPeeringRoute) WithRemoteCIDR(cidr string) *VPCPeeringRoute {
	r.remoteCIDR = &cidr
	return r
}
func (r *VPCPeeringRoute) WithBillingPeriod(p string) *VPCPeeringRoute {
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
func (r *VPCPeeringRoute) BillingPeriod() string {
	if r.billingPeriod == nil {
		return ""
	}
	return *r.billingPeriod
}

func (r *VPCPeeringRoute) toRequest() types.VPCPeeringRouteRequest {
	var bp string
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

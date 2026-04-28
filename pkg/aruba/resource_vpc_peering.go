package aruba

import (
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

// VPCPeering wraps an Aruba Cloud VPC Peering (a direct child of a VPC, with regional metadata).
// Construct with aruba.NewVPCPeering() and bind it to a parent VPC via IntoVPC(vpc).
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

// Setters (chainable).

func (p *VPCPeering) IntoVPC(v Ref) *VPCPeering            { p.intoVPC(v); return p }
func (p *VPCPeering) WithName(n string) *VPCPeering        { p.withName(n); return p }
func (p *VPCPeering) AddTag(t string) *VPCPeering          { p.addTag(t); return p }
func (p *VPCPeering) RemoveTag(t string) *VPCPeering       { p.removeTag(t); return p }
func (p *VPCPeering) ReplaceTags(ts ...string) *VPCPeering { p.replaceTags(ts...); return p }
func (p *VPCPeering) WithLocation(loc string) *VPCPeering  { p.withLocation(loc); return p }
func (p *VPCPeering) InRegion(region string) *VPCPeering   { p.inRegion(region); return p }

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
		p.withLocation(resp.Metadata.LocationResponse.Value)
	}
	p.setStatus(&resp.Status)
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

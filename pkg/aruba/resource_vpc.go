package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// VPC is the wrapper for an Aruba Cloud VPC. Construct with aruba.NewVPC()
// and bind it to a project via IntoProject(parent). Pass to VPCsClient.Create
// / .Update or receive from .Get / .List.
type VPC struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	defaultVPC *bool
	preset     *bool
	response   *types.VPCResponse
}

func (v *VPC) IntoProject(p Ref) *VPC        { v.intoProject(p); return v }
func (v *VPC) WithName(n string) *VPC        { v.withName(n); return v }
func (v *VPC) AddTag(t string) *VPC          { v.addTag(t); return v }
func (v *VPC) RemoveTag(t string) *VPC       { v.removeTag(t); return v }
func (v *VPC) ReplaceTags(ts ...string) *VPC { v.replaceTags(ts...); return v }
func (v *VPC) WithLocation(loc string) *VPC  { v.withLocation(loc); return v }
func (v *VPC) InRegion(region string) *VPC   { v.withLocation(region); return v }
func (v *VPC) WithDefault(b bool) *VPC       { v.defaultVPC = &b; return v }
func (v *VPC) WithPreset(b bool) *VPC        { v.preset = &b; return v }

// URI satisfies Ref.
func (v *VPC) URI() string { return v.RespURI() }

// VPCID satisfies withVPCID so children's IntoVPC can extract the parent ID.
func (v *VPC) VPCID() string { return v.ID() }

// Raw shadows the promoted responseMetadataMixin.Raw() returning the full response.
func (v *VPC) Raw() *types.VPCResponse { return v.response }

// RawRequest returns the wire-level request that toRequest() would emit.
func (v *VPC) RawRequest() types.VPCRequest { return v.toRequest() }

// IsDefault returns true if this VPC is the account-region default.
func (v *VPC) IsDefault() bool {
	if v.defaultVPC == nil {
		return false
	}
	return *v.defaultVPC
}

// IsPreset returns true if the VPC was created with a preset subnet/SG.
func (v *VPC) IsPreset() bool {
	if v.preset == nil {
		return false
	}
	return *v.preset
}

func (v *VPC) toRequest() types.VPCRequest {
	var props *types.VPCProperties
	if v.defaultVPC != nil || v.preset != nil {
		props = &types.VPCProperties{Default: v.defaultVPC, Preset: v.preset}
	}
	return types.VPCRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: v.toMetadata(),
			Location:                v.toLocation(),
		},
		Properties: types.VPCPropertiesRequest{Properties: props},
	}
}

func (v *VPC) fromResponse(resp *types.VPCResponse) {
	if resp == nil {
		return
	}
	v.response = resp
	v.setMeta(&resp.Metadata)
	v.withName(vpcDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		v.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		v.withLocation(resp.Metadata.LocationResponse.Value)
	}
	v.setStatus(&resp.Status)
	v.setLinked(resp.Properties.LinkedResources)
	d := resp.Properties.Default
	v.defaultVPC = &d
	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		v.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
}

func vpcDerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

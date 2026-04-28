package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// SecurityGroup is the wrapper for an Aruba Cloud Security Group (a direct child of a VPC).
// Construct with aruba.NewSecurityGroup() and bind it to a parent VPC via IntoVPC(vpc).
type SecurityGroup struct {
	errMixin
	metadataMixin
	vpcScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	defaultSG *bool                        // Properties.Default (request: *bool for omitempty; response: plain bool)
	response  *types.SecurityGroupResponse // backs Raw()
}

// Setters (chainable; promoted methods re-exposed at *SecurityGroup level).
func (sg *SecurityGroup) IntoVPC(v Ref) *SecurityGroup            { sg.intoVPC(v); return sg }
func (sg *SecurityGroup) WithName(n string) *SecurityGroup        { sg.withName(n); return sg }
func (sg *SecurityGroup) AddTag(t string) *SecurityGroup          { sg.addTag(t); return sg }
func (sg *SecurityGroup) RemoveTag(t string) *SecurityGroup       { sg.removeTag(t); return sg }
func (sg *SecurityGroup) ReplaceTags(ts ...string) *SecurityGroup { sg.replaceTags(ts...); return sg }
func (sg *SecurityGroup) WithDefault(b bool) *SecurityGroup       { sg.defaultSG = &b; return sg }

// URI satisfies Ref.
func (sg *SecurityGroup) URI() string { return sg.RespURI() }

// SecurityGroupID satisfies withSecurityGroupID so child wrappers (SecurityGroupRule)
// can extract this ID via typed assertion.
func (sg *SecurityGroup) SecurityGroupID() string { return sg.ID() }

// Raw shadows responseMetadataMixin.Raw() with the full SecurityGroup response.
func (sg *SecurityGroup) Raw() *types.SecurityGroupResponse { return sg.response }

// RawRequest returns what toRequest() would emit right now.
func (sg *SecurityGroup) RawRequest() types.SecurityGroupRequest { return sg.toRequest() }

// Default returns the security group's default flag, or false if unset.
func (sg *SecurityGroup) Default() bool {
	if sg.defaultSG == nil {
		return false
	}
	return *sg.defaultSG
}

func (sg *SecurityGroup) toRequest() types.SecurityGroupRequest {
	return types.SecurityGroupRequest{
		Metadata: sg.toMetadata(),
		Properties: types.SecurityGroupPropertiesRequest{
			Default: sg.defaultSG,
		},
	}
}

func (sg *SecurityGroup) fromResponse(resp *types.SecurityGroupResponse) {
	if resp == nil {
		return
	}
	sg.response = resp
	sg.setMeta(&resp.Metadata)
	sg.withName(securityGroupDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		sg.replaceTags(resp.Metadata.Tags...)
	}
	sg.setStatus(&resp.Status)
	sg.setLinked(resp.Properties.LinkedResources)

	// Properties.Default is plain bool on the response — backfill into our *bool field.
	d := resp.Properties.Default
	sg.defaultSG = &d

	// Backfill ancestor IDs: prefer ProjectResponseMetadata, then URI parse.
	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		sg.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if (sg.vpcID == "" || sg.projectID == "") && sg.RespURI() != "" {
		ids := parseURIIDs(sg.RespURI())
		if sg.vpcID == "" {
			sg.vpcID = ids["vpcs"]
		}
		if sg.projectID == "" {
			sg.projectID = ids["projects"]
		}
	}
}

func securityGroupDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

package aruba

import (
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

// SecurityRule is the wrapper for an Aruba Cloud Security Rule (a direct child of a SecurityGroup).
// Construct with aruba.NewSecurityRule() and bind it via IntoSecurityGroup(sg).
type SecurityRule struct {
	errMixin
	metadataMixin
	regionalMixin
	securityGroupScopedMixin
	responseMetadataMixin
	statusMixin
	httpEnvelopeMixin

	direction *types.RuleDirection
	protocol  *string
	port      *string
	target    *types.RuleTarget
	response  *types.SecurityRuleResponse
}

// Setters (chainable; promoted methods re-exposed at *SecurityRule level).

func (r *SecurityRule) IntoSecurityGroup(sg Ref) *SecurityRule { r.intoSecurityGroup(sg); return r }
func (r *SecurityRule) WithName(n string) *SecurityRule        { r.withName(n); return r }
func (r *SecurityRule) AddTag(t string) *SecurityRule          { r.addTag(t); return r }
func (r *SecurityRule) RemoveTag(t string) *SecurityRule       { r.removeTag(t); return r }
func (r *SecurityRule) ReplaceTags(ts ...string) *SecurityRule { r.replaceTags(ts...); return r }
func (r *SecurityRule) WithLocation(loc string) *SecurityRule  { r.withLocation(loc); return r }
func (r *SecurityRule) InRegion(region string) *SecurityRule   { r.inRegion(region); return r }

// WithDirection sets the rule direction. Accepts "Ingress" or "Egress".
func (r *SecurityRule) WithDirection(dir string) *SecurityRule {
	d := types.RuleDirection(dir)
	r.direction = &d
	return r
}

// WithProtocol sets the L4 protocol. Accepts "ANY", "TCP", "UDP", "ICMP".
func (r *SecurityRule) WithProtocol(proto string) *SecurityRule {
	r.protocol = &proto
	return r
}

// WithPort sets the port specifier — single (e.g., "22"), range ("80-100"), or wildcard ("*").
func (r *SecurityRule) WithPort(port string) *SecurityRule {
	r.port = &port
	return r
}

// WithTargetCIDR sets the target as an IP/CIDR endpoint.
// Mutually exclusive with WithTargetSecurityGroup — setting both records a setter-time error.
func (r *SecurityRule) WithTargetCIDR(cidr string) *SecurityRule {
	if r.target != nil && r.target.Kind == types.EndpointTypeSecurityGroup {
		r.addErr(fmt.Errorf("WithTargetCIDR: target already set to SecurityGroup; pick one"))
		return r
	}
	r.target = &types.RuleTarget{Kind: types.EndpointTypeIP, Value: cidr}
	return r
}

// WithTargetSecurityGroup sets the target as another SecurityGroup endpoint.
// Mutually exclusive with WithTargetCIDR — setting both records a setter-time error.
func (r *SecurityRule) WithTargetSecurityGroup(sg Ref) *SecurityRule {
	if r.target != nil && r.target.Kind == types.EndpointTypeIP {
		r.addErr(fmt.Errorf("WithTargetSecurityGroup: target already set to CIDR; pick one"))
		return r
	}
	uri := sg.URI()
	if uri == "" {
		r.addErr(fmt.Errorf("WithTargetSecurityGroup: target SecurityGroup Ref has empty URI"))
		return r
	}
	r.target = &types.RuleTarget{Kind: types.EndpointTypeSecurityGroup, Value: uri}
	return r
}

// URI satisfies Ref.
func (r *SecurityRule) URI() string { return r.RespURI() }

// SecurityRuleID satisfies withSecurityRuleID.
func (r *SecurityRule) SecurityRuleID() string { return r.ID() }

// Raw shadows responseMetadataMixin.Raw() with the full SecurityRule response.
func (r *SecurityRule) Raw() *types.SecurityRuleResponse { return r.response }

// RawRequest returns what toRequest() would emit right now.
func (r *SecurityRule) RawRequest() types.SecurityRuleRequest { return r.toRequest() }

// Direction returns the configured rule direction (zero value if unset).
func (r *SecurityRule) Direction() types.RuleDirection {
	if r.direction == nil {
		return ""
	}
	return *r.direction
}

// Protocol returns the configured protocol ("" if unset).
func (r *SecurityRule) Protocol() string {
	if r.protocol == nil {
		return ""
	}
	return *r.protocol
}

// Port returns the configured port specifier ("" if unset).
func (r *SecurityRule) Port() string {
	if r.port == nil {
		return ""
	}
	return *r.port
}

// TargetKind returns the configured target endpoint kind ("" if unset).
func (r *SecurityRule) TargetKind() types.EndpointTypeDto {
	if r.target == nil {
		return ""
	}
	return r.target.Kind
}

// TargetValue returns the configured target endpoint value ("" if unset).
func (r *SecurityRule) TargetValue() string {
	if r.target == nil {
		return ""
	}
	return r.target.Value
}

func (r *SecurityRule) toRequest() types.SecurityRuleRequest {
	props := types.SecurityRulePropertiesRequest{}
	if r.direction != nil {
		props.Direction = *r.direction
	}
	if r.protocol != nil {
		props.Protocol = *r.protocol
	}
	if r.port != nil {
		props.Port = *r.port
	}
	if r.target != nil {
		props.Target = r.target
	}
	return types.SecurityRuleRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: r.toMetadata(),
			Location:                r.toLocation(),
		},
		Properties: props,
	}
}

func (r *SecurityRule) fromResponse(resp *types.SecurityRuleResponse) {
	if resp == nil {
		return
	}
	r.response = resp
	r.setMeta(&resp.Metadata)
	r.withName(securityRuleDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		r.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		r.withLocation(resp.Metadata.LocationResponse.Value)
	}
	r.setStatus(&resp.Status)

	if resp.Properties.Direction != "" {
		d := resp.Properties.Direction
		r.direction = &d
	}
	if resp.Properties.Protocol != "" {
		p := resp.Properties.Protocol
		r.protocol = &p
	}
	if resp.Properties.Port != "" {
		p := resp.Properties.Port
		r.port = &p
	}
	if resp.Properties.Target != nil {
		t := *resp.Properties.Target
		r.target = &t
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		r.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if (r.vpcID == "" || r.projectID == "" || r.securityGroupID == "") && r.RespURI() != "" {
		ids := parseURIIDs(r.RespURI())
		if r.vpcID == "" {
			r.vpcID = ids["vpcs"]
		}
		if r.projectID == "" {
			r.projectID = ids["projects"]
		}
		if r.securityGroupID == "" {
			// Production URI uses "securitygroups"; wrapper-test URIs use "security-groups".
			if v := ids["securitygroups"]; v != "" {
				r.securityGroupID = v
			}
			if r.securityGroupID == "" {
				r.securityGroupID = ids["security-groups"]
			}
		}
	}
}

func securityRuleDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
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
	r.setTerminalStates(securityRuleTerminalStates)

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

var securityRuleTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type securityRuleLowLevelClient interface {
	List(ctx context.Context, projectID, vpcID, securityGroupID string, params *types.RequestParameters) (*types.Response[types.SecurityRuleList], error)
	Get(ctx context.Context, projectID, vpcID, securityGroupID, securityRuleID string, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	Create(ctx context.Context, projectID, vpcID, securityGroupID string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	Update(ctx context.Context, projectID, vpcID, securityGroupID, securityRuleID string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	Delete(ctx context.Context, projectID, vpcID, securityGroupID, securityRuleID string, params *types.RequestParameters) (*types.Response[any], error)
}

type securityGroupRulesClientAdapter struct{ low securityRuleLowLevelClient }

func newSecurityGroupRulesClientAdapter(rest *restclient.Client) *securityGroupRulesClientAdapter {
	if rest == nil {
		return &securityGroupRulesClientAdapter{}
	}
	return &securityGroupRulesClientAdapter{
		low: network.NewSecurityGroupRulesClientImpl(
			rest,
			network.NewSecurityGroupsClientImpl(rest, network.NewVPCsClientImpl(rest)),
		),
	}
}

func (a *securityGroupRulesClientAdapter) Create(ctx context.Context, rule *SecurityRule, opts ...CallOption) (*SecurityRule, error) {
	if err := rule.Err(); err != nil {
		return rule, err
	}
	if rule.SecurityGroupID() == "" || rule.VPCID() == "" || rule.ProjectID() == "" {
		return rule, fmt.Errorf("Create: security rule has no SecurityGroup — call IntoSecurityGroup first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, rule.ProjectID(), rule.VPCID(), rule.SecurityGroupID(), rule.toRequest(), rp)
	populateHTTPEnvelope(&rule.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		rule.fromResponse(resp.Data)
		rule.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, rule)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				rule.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return rule, err
	}
	if resp != nil && !resp.IsSuccess() {
		return rule, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return rule, nil
}

func (a *securityGroupRulesClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*SecurityRule, error) {
	projectID, vpcID, securityGroupID, securityRuleID, err := securityRuleIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, securityGroupID, securityRuleID, rp)
	out := &SecurityRule{}
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
	if out.securityGroupID == "" {
		out.securityGroupID = securityGroupID
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

func (a *securityGroupRulesClientAdapter) Update(ctx context.Context, rule *SecurityRule, opts ...CallOption) (*SecurityRule, error) {
	if err := rule.Err(); err != nil {
		return rule, err
	}
	if rule.ID() == "" {
		return rule, fmt.Errorf("Update: security rule has no ID — call Get first or seed from response metadata")
	}
	if rule.SecurityGroupID() == "" || rule.VPCID() == "" || rule.ProjectID() == "" {
		return rule, fmt.Errorf("Update: security rule has no SecurityGroup — call IntoSecurityGroup first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, rule.ProjectID(), rule.VPCID(), rule.SecurityGroupID(), rule.ID(), rule.toRequest(), rp)
	populateHTTPEnvelope(&rule.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		rule.fromResponse(resp.Data)
		rule.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, rule)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				rule.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return rule, err
	}
	if resp != nil && !resp.IsSuccess() {
		return rule, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return rule, nil
}

func (a *securityGroupRulesClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, securityGroupID, securityRuleID, err := securityRuleIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, securityGroupID, securityRuleID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *securityGroupRulesClientAdapter) List(ctx context.Context, sg Ref, opts ...CallOption) (*List[*SecurityRule], error) {
	projectID, vpcID, securityGroupID, err := securityGroupIDsFromRef(sg)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, vpcID, securityGroupID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*SecurityRule
	if resp != nil && resp.Data != nil {
		items = make([]*SecurityRule, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			rule := &SecurityRule{}
			rule.fromResponse(&resp.Data.Values[i])
			rule.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, rule)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					rule.fromResponse(fresh.Raw())
				}
				return nil
			})
			if rule.securityGroupID == "" {
				rule.securityGroupID = securityGroupID
			}
			if rule.vpcID == "" {
				rule.vpcID = vpcID
			}
			if rule.projectID == "" {
				rule.projectID = projectID
			}
			items = append(items, rule)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*SecurityRule], error) {
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

// securityRuleIDsFromRef extracts (projectID, vpcID, securityGroupID, securityRuleID) from a Ref.
// Tries typed assertions first, then falls back to URI path parsing (both hyphenated and no-hyphen forms).
func securityRuleIDsFromRef(ref Ref) (projectID, vpcID, securityGroupID, securityRuleID string, err error) {
	rid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withSecurityRuleID); ok {
			return w.SecurityRuleID(), true
		}
		return "", false
	}, "security-rules")
	if !ok || rid == "" {
		m := parseURIIDs(ref.URI())
		if v := m["securityrules"]; v != "" {
			rid = v
			ok = true
		}
	}
	if !ok || rid == "" {
		return "", "", "", "", fmt.Errorf("cannot determine security rule ID from Ref %q", ref.URI())
	}
	sgid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withSecurityGroupID); ok {
			return w.SecurityGroupID(), true
		}
		return "", false
	}, "security-groups")
	if !ok || sgid == "" {
		m := parseURIIDs(ref.URI())
		if v := m["securitygroups"]; v != "" {
			sgid = v
			ok = true
		}
	}
	if !ok || sgid == "" {
		return "", "", "", "", fmt.Errorf("cannot determine security group ID from Ref %q", ref.URI())
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
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid == "" {
		return "", "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return pid, vid, sgid, rid, nil
}

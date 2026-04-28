package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

type NetworkClient interface {
	ElasticIPs() ElasticIPsClient
	LoadBalancers() LoadBalancersClient
	SecurityGroupRules() SecurityGroupRulesClient
	SecurityGroups() SecurityGroupsClient
	Subnets() SubnetsClient
	VPCPeeringRoutes() VPCPeeringRoutesClient
	VPCPeerings() VPCPeeringsClient
	VPCs() VPCsClient
	VPNRoutes() VPNRoutesClient
	VPNTunnels() VPNTunnelsClient
}

type networkClientImpl struct {
	elasticIPsClient         ElasticIPsClient
	loadBalancersClient      LoadBalancersClient
	securityGroupRulesClient SecurityGroupRulesClient
	securityGroupsClient     SecurityGroupsClient
	subnetsClient            SubnetsClient
	vpcPeeringRoutesClient   VPCPeeringRoutesClient
	vpcPeeringsClient        VPCPeeringsClient
	vpcsClient               VPCsClient
	vpnRoutesClient          VPNRoutesClient
	vpnTunnelsClient         VPNTunnelsClient
}

var _ NetworkClient = (*networkClientImpl)(nil)

func (c *networkClientImpl) ElasticIPs() ElasticIPsClient {
	return c.elasticIPsClient
}
func (c *networkClientImpl) LoadBalancers() LoadBalancersClient {
	return c.loadBalancersClient
}
func (c *networkClientImpl) SecurityGroupRules() SecurityGroupRulesClient {
	return c.securityGroupRulesClient
}
func (c *networkClientImpl) SecurityGroups() SecurityGroupsClient {
	return c.securityGroupsClient
}
func (c *networkClientImpl) Subnets() SubnetsClient {
	return c.subnetsClient
}
func (c *networkClientImpl) VPCPeeringRoutes() VPCPeeringRoutesClient {
	return c.vpcPeeringRoutesClient
}
func (c *networkClientImpl) VPCPeerings() VPCPeeringsClient {
	return c.vpcPeeringsClient
}
func (c *networkClientImpl) VPCs() VPCsClient {
	return c.vpcsClient
}
func (c *networkClientImpl) VPNRoutes() VPNRoutesClient {
	return c.vpnRoutesClient
}
func (c *networkClientImpl) VPNTunnels() VPNTunnelsClient {
	return c.vpnTunnelsClient
}

type ElasticIPsClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*ElasticIP], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*ElasticIP, error)
	Create(ctx context.Context, eip *ElasticIP, opts ...CallOption) (*ElasticIP, error)
	Update(ctx context.Context, eip *ElasticIP, opts ...CallOption) (*ElasticIP, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type elasticIPLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.ElasticList], error)
	Get(ctx context.Context, projectID, elasticIPID string, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Create(ctx context.Context, projectID string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Update(ctx context.Context, projectID, elasticIPID string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Delete(ctx context.Context, projectID, elasticIPID string, params *types.RequestParameters) (*types.Response[any], error)
}

type elasticIPsClientAdapter struct{ low elasticIPLowLevelClient }

func newElasticIPsClientAdapter(rest *restclient.Client) *elasticIPsClientAdapter {
	if rest == nil {
		return &elasticIPsClientAdapter{}
	}
	return &elasticIPsClientAdapter{low: network.NewElasticIPsClientImpl(rest)}
}

func (a *elasticIPsClientAdapter) Create(ctx context.Context, e *ElasticIP, opts ...CallOption) (*ElasticIP, error) {
	if err := e.Err(); err != nil {
		return e, err
	}
	if e.ProjectID() == "" {
		return e, fmt.Errorf("Create: elastic IP has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, e.ProjectID(), e.toRequest(), rp)
	populateHTTPEnvelope(&e.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		e.fromResponse(resp.Data)
	}
	if err != nil {
		return e, err
	}
	if resp != nil && !resp.IsSuccess() {
		return e, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return e, nil
}

func (a *elasticIPsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*ElasticIP, error) {
	projectID, elasticIPID, err := elasticIPIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, elasticIPID, rp)
	out := &ElasticIP{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	out.projectID = projectID
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *elasticIPsClientAdapter) Update(ctx context.Context, e *ElasticIP, opts ...CallOption) (*ElasticIP, error) {
	if err := e.Err(); err != nil {
		return e, err
	}
	if e.ID() == "" {
		return e, fmt.Errorf("Update: elastic IP has no ID — call Get first or seed from response metadata")
	}
	if e.ProjectID() == "" {
		return e, fmt.Errorf("Update: elastic IP has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, e.ProjectID(), e.ID(), e.toRequest(), rp)
	populateHTTPEnvelope(&e.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		e.fromResponse(resp.Data)
	}
	if err != nil {
		return e, err
	}
	if resp != nil && !resp.IsSuccess() {
		return e, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return e, nil
}

func (a *elasticIPsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, elasticIPID, err := elasticIPIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, elasticIPID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *elasticIPsClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*ElasticIP], error) {
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
	var items []*ElasticIP
	if resp != nil && resp.Data != nil {
		items = make([]*ElasticIP, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			e := &ElasticIP{}
			e.fromResponse(&resp.Data.Values[i])
			if e.projectID == "" {
				e.projectID = projectID
			}
			items = append(items, e)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*ElasticIP], error) {
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

// elasticIPIDsFromRef extracts (projectID, elasticIPID) from a Ref. Tries typed
// assertions first, then falls back to URI path parsing.
func elasticIPIDsFromRef(ref Ref) (projectID, elasticIPID string, err error) {
	eid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withElasticIPID); ok {
			return w.ElasticIPID(), true
		}
		return "", false
	}, "elasticIps")
	if !ok || eid == "" {
		return "", "", fmt.Errorf("cannot determine elastic IP ID from Ref %q", ref.URI())
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
	return pid, eid, nil
}

type LoadBalancersClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*LoadBalancer], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*LoadBalancer, error)
}

type loadBalancerLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.LoadBalancerList], error)
	Get(ctx context.Context, projectID, loadBalancerID string, params *types.RequestParameters) (*types.Response[types.LoadBalancerResponse], error)
}

type loadBalancersClientAdapter struct{ low loadBalancerLowLevelClient }

func newLoadBalancersClientAdapter(rest *restclient.Client) *loadBalancersClientAdapter {
	if rest == nil {
		return &loadBalancersClientAdapter{}
	}
	return &loadBalancersClientAdapter{low: network.NewLoadBalancersClientImpl(rest)}
}

func (a *loadBalancersClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*LoadBalancer, error) {
	projectID, loadBalancerID, err := loadBalancerIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, loadBalancerID, rp)
	out := &LoadBalancer{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	out.projectID = projectID
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *loadBalancersClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*LoadBalancer], error) {
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
	var items []*LoadBalancer
	if resp != nil && resp.Data != nil {
		items = make([]*LoadBalancer, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			lb := &LoadBalancer{}
			lb.fromResponse(&resp.Data.Values[i])
			if lb.projectID == "" {
				lb.projectID = projectID
			}
			items = append(items, lb)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*LoadBalancer], error) {
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

// loadBalancerIDsFromRef extracts (projectID, loadBalancerID) from a Ref. Tries typed
// assertions first, then falls back to URI path parsing.
func loadBalancerIDsFromRef(ref Ref) (projectID, loadBalancerID string, err error) {
	lid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withLoadBalancerID); ok {
			return w.LoadBalancerID(), true
		}
		return "", false
	}, "loadbalancers")
	if !ok || lid == "" {
		return "", "", fmt.Errorf("cannot determine load balancer ID from Ref %q", ref.URI())
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
	return pid, lid, nil
}

type SecurityGroupRulesClient interface {
	List(ctx context.Context, securityGroup Ref, opts ...CallOption) (*List[*SecurityRule], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*SecurityRule, error)
	Create(ctx context.Context, rule *SecurityRule, opts ...CallOption) (*SecurityRule, error)
	Update(ctx context.Context, rule *SecurityRule, opts ...CallOption) (*SecurityRule, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

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

type SecurityGroupsClient interface {
	List(ctx context.Context, vpc Ref, opts ...CallOption) (*List[*SecurityGroup], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*SecurityGroup, error)
	Create(ctx context.Context, sg *SecurityGroup, opts ...CallOption) (*SecurityGroup, error)
	Update(ctx context.Context, sg *SecurityGroup, opts ...CallOption) (*SecurityGroup, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type securityGroupLowLevelClient interface {
	List(ctx context.Context, projectID, vpcID string, params *types.RequestParameters) (*types.Response[types.SecurityGroupList], error)
	Get(ctx context.Context, projectID, vpcID, securityGroupID string, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Create(ctx context.Context, projectID, vpcID string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Update(ctx context.Context, projectID, vpcID, securityGroupID string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Delete(ctx context.Context, projectID, vpcID, securityGroupID string, params *types.RequestParameters) (*types.Response[any], error)
}

type securityGroupsClientAdapter struct{ low securityGroupLowLevelClient }

func newSecurityGroupsClientAdapter(rest *restclient.Client) *securityGroupsClientAdapter {
	if rest == nil {
		return &securityGroupsClientAdapter{}
	}
	return &securityGroupsClientAdapter{
		low: network.NewSecurityGroupsClientImpl(rest, network.NewVPCsClientImpl(rest)),
	}
}

func (a *securityGroupsClientAdapter) Create(ctx context.Context, sg *SecurityGroup, opts ...CallOption) (*SecurityGroup, error) {
	if err := sg.Err(); err != nil {
		return sg, err
	}
	if sg.VPCID() == "" || sg.ProjectID() == "" {
		return sg, fmt.Errorf("Create: security group has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, sg.ProjectID(), sg.VPCID(), sg.toRequest(), rp)
	populateHTTPEnvelope(&sg.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		sg.fromResponse(resp.Data)
	}
	if err != nil {
		return sg, err
	}
	if resp != nil && !resp.IsSuccess() {
		return sg, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return sg, nil
}

func (a *securityGroupsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*SecurityGroup, error) {
	projectID, vpcID, securityGroupID, err := securityGroupIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, securityGroupID, rp)
	out := &SecurityGroup{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	out.vpcID = vpcID
	out.projectID = projectID
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *securityGroupsClientAdapter) Update(ctx context.Context, sg *SecurityGroup, opts ...CallOption) (*SecurityGroup, error) {
	if err := sg.Err(); err != nil {
		return sg, err
	}
	if sg.ID() == "" {
		return sg, fmt.Errorf("Update: security group has no ID — call Get first or seed from response metadata")
	}
	if sg.VPCID() == "" || sg.ProjectID() == "" {
		return sg, fmt.Errorf("Update: security group has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, sg.ProjectID(), sg.VPCID(), sg.ID(), sg.toRequest(), rp)
	populateHTTPEnvelope(&sg.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		sg.fromResponse(resp.Data)
	}
	if err != nil {
		return sg, err
	}
	if resp != nil && !resp.IsSuccess() {
		return sg, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return sg, nil
}

func (a *securityGroupsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, securityGroupID, err := securityGroupIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, securityGroupID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *securityGroupsClientAdapter) List(ctx context.Context, vpc Ref, opts ...CallOption) (*List[*SecurityGroup], error) {
	projectID, vpcID, err := vpcIDsFromRef(vpc)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, vpcID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*SecurityGroup
	if resp != nil && resp.Data != nil {
		items = make([]*SecurityGroup, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			sg := &SecurityGroup{}
			sg.fromResponse(&resp.Data.Values[i])
			if sg.vpcID == "" {
				sg.vpcID = vpcID
			}
			if sg.projectID == "" {
				sg.projectID = projectID
			}
			items = append(items, sg)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*SecurityGroup], error) {
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

// securityGroupIDsFromRef extracts (projectID, vpcID, securityGroupID) from a Ref.
// Tries typed assertions first, then falls back to URI path parsing.
func securityGroupIDsFromRef(ref Ref) (projectID, vpcID, securityGroupID string, err error) {
	sgid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withSecurityGroupID); ok {
			return w.SecurityGroupID(), true
		}
		return "", false
	}, "security-groups")
	if !ok || sgid == "" {
		return "", "", "", fmt.Errorf("cannot determine security group ID from Ref %q", ref.URI())
	}
	vid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid == "" {
		return "", "", "", fmt.Errorf("cannot determine VPC ID from Ref %q", ref.URI())
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
	return pid, vid, sgid, nil
}

type SubnetsClient interface {
	List(ctx context.Context, vpc Ref, opts ...CallOption) (*List[*Subnet], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*Subnet, error)
	Create(ctx context.Context, subnet *Subnet, opts ...CallOption) (*Subnet, error)
	Update(ctx context.Context, subnet *Subnet, opts ...CallOption) (*Subnet, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type subnetLowLevelClient interface {
	List(ctx context.Context, projectID, vpcID string, params *types.RequestParameters) (*types.Response[types.SubnetList], error)
	Get(ctx context.Context, projectID, vpcID, subnetID string, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	Create(ctx context.Context, projectID, vpcID string, body types.SubnetRequest, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	Update(ctx context.Context, projectID, vpcID, subnetID string, body types.SubnetRequest, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	Delete(ctx context.Context, projectID, vpcID, subnetID string, params *types.RequestParameters) (*types.Response[any], error)
}

type subnetsClientAdapter struct{ low subnetLowLevelClient }

func newSubnetsClientAdapter(rest *restclient.Client) *subnetsClientAdapter {
	if rest == nil {
		return &subnetsClientAdapter{}
	}
	return &subnetsClientAdapter{
		low: network.NewSubnetsClientImpl(rest, network.NewVPCsClientImpl(rest)),
	}
}

func (a *subnetsClientAdapter) Create(ctx context.Context, s *Subnet, opts ...CallOption) (*Subnet, error) {
	if err := s.Err(); err != nil {
		return s, err
	}
	if s.VPCID() == "" || s.ProjectID() == "" {
		return s, fmt.Errorf("Create: subnet has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, s.ProjectID(), s.VPCID(), s.toRequest(), rp)
	populateHTTPEnvelope(&s.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		s.fromResponse(resp.Data)
	}
	if err != nil {
		return s, err
	}
	if resp != nil && !resp.IsSuccess() {
		return s, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return s, nil
}

func (a *subnetsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*Subnet, error) {
	projectID, vpcID, subnetID, err := subnetIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, subnetID, rp)
	out := &Subnet{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	out.vpcID = vpcID
	out.projectID = projectID
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *subnetsClientAdapter) Update(ctx context.Context, s *Subnet, opts ...CallOption) (*Subnet, error) {
	if err := s.Err(); err != nil {
		return s, err
	}
	if s.ID() == "" {
		return s, fmt.Errorf("Update: subnet has no ID — call Get first or seed from response metadata")
	}
	if s.VPCID() == "" || s.ProjectID() == "" {
		return s, fmt.Errorf("Update: subnet has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, s.ProjectID(), s.VPCID(), s.ID(), s.toRequest(), rp)
	populateHTTPEnvelope(&s.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		s.fromResponse(resp.Data)
	}
	if err != nil {
		return s, err
	}
	if resp != nil && !resp.IsSuccess() {
		return s, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return s, nil
}

func (a *subnetsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, subnetID, err := subnetIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, subnetID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *subnetsClientAdapter) List(ctx context.Context, vpc Ref, opts ...CallOption) (*List[*Subnet], error) {
	projectID, vpcID, err := vpcIDsFromRef(vpc)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, vpcID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*Subnet
	if resp != nil && resp.Data != nil {
		items = make([]*Subnet, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			s := &Subnet{}
			s.fromResponse(&resp.Data.Values[i])
			s.vpcID = vpcID
			if s.projectID == "" {
				s.projectID = projectID
			}
			items = append(items, s)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*Subnet], error) {
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

// subnetIDsFromRef extracts (projectID, vpcID, subnetID) from a Ref. Tries typed
// assertions first, then falls back to URI path parsing.
func subnetIDsFromRef(ref Ref) (projectID, vpcID, subnetID string, err error) {
	sid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withSubnetID); ok {
			return w.SubnetID(), true
		}
		return "", false
	}, "subnets")
	if !ok || sid == "" {
		return "", "", "", fmt.Errorf("cannot determine subnet ID from Ref %q", ref.URI())
	}
	vid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid == "" {
		return "", "", "", fmt.Errorf("cannot determine VPC ID from Ref %q", ref.URI())
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
	return pid, vid, sid, nil
}

type VPCPeeringRoutesClient interface {
	List(ctx context.Context, peering Ref, opts ...CallOption) (*List[*VPCPeeringRoute], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPCPeeringRoute, error)
	Create(ctx context.Context, route *VPCPeeringRoute, opts ...CallOption) (*VPCPeeringRoute, error)
	Update(ctx context.Context, route *VPCPeeringRoute, opts ...CallOption) (*VPCPeeringRoute, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type vpcPeeringRouteLowLevelClient interface {
	List(ctx context.Context, projectID, vpcID, vpcPeeringID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteList], error)
	Get(ctx context.Context, projectID, vpcID, vpcPeeringID, vpcPeeringRouteID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Create(ctx context.Context, projectID, vpcID, vpcPeeringID string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Update(ctx context.Context, projectID, vpcID, vpcPeeringID, vpcPeeringRouteID string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Delete(ctx context.Context, projectID, vpcID, vpcPeeringID, vpcPeeringRouteID string, params *types.RequestParameters) (*types.Response[any], error)
}

type vpcPeeringRoutesClientAdapter struct{ low vpcPeeringRouteLowLevelClient }

func newVPCPeeringRoutesClientAdapter(rest *restclient.Client) *vpcPeeringRoutesClientAdapter {
	if rest == nil {
		return &vpcPeeringRoutesClientAdapter{}
	}
	return &vpcPeeringRoutesClientAdapter{low: network.NewVPCPeeringRoutesClientImpl(rest)}
}

func (a *vpcPeeringRoutesClientAdapter) Create(ctx context.Context, route *VPCPeeringRoute, opts ...CallOption) (*VPCPeeringRoute, error) {
	if err := route.Err(); err != nil {
		return route, err
	}
	if route.VPCPeeringID() == "" || route.VPCID() == "" || route.ProjectID() == "" {
		return route, fmt.Errorf("Create: VPC peering route has no parent peering — call IntoVPCPeering first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, route.ProjectID(), route.VPCID(), route.VPCPeeringID(), route.toRequest(), rp)
	populateHTTPEnvelope(&route.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		route.fromResponse(resp.Data)
	}
	if err != nil {
		return route, err
	}
	if resp != nil && !resp.IsSuccess() {
		return route, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return route, nil
}

func (a *vpcPeeringRoutesClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPCPeeringRoute, error) {
	projectID, vpcID, vpcPeeringID, routeID, err := vpcPeeringRouteIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, vpcPeeringID, routeID, rp)
	out := &VPCPeeringRoute{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	if out.vpcPeeringID == "" {
		out.vpcPeeringID = vpcPeeringID
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

func (a *vpcPeeringRoutesClientAdapter) Update(ctx context.Context, route *VPCPeeringRoute, opts ...CallOption) (*VPCPeeringRoute, error) {
	if err := route.Err(); err != nil {
		return route, err
	}
	if route.ID() == "" {
		return route, fmt.Errorf("Update: VPC peering route has no ID — call Get first or seed from response metadata")
	}
	if route.VPCPeeringID() == "" || route.VPCID() == "" || route.ProjectID() == "" {
		return route, fmt.Errorf("Update: VPC peering route has no parent peering — call IntoVPCPeering first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, route.ProjectID(), route.VPCID(), route.VPCPeeringID(), route.ID(), route.toRequest(), rp)
	populateHTTPEnvelope(&route.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		route.fromResponse(resp.Data)
	}
	if err != nil {
		return route, err
	}
	if resp != nil && !resp.IsSuccess() {
		return route, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return route, nil
}

func (a *vpcPeeringRoutesClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, vpcPeeringID, routeID, err := vpcPeeringRouteIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, vpcPeeringID, routeID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *vpcPeeringRoutesClientAdapter) List(ctx context.Context, peering Ref, opts ...CallOption) (*List[*VPCPeeringRoute], error) {
	projectID, vpcID, vpcPeeringID, err := vpcPeeringIDsFromRef(peering)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, vpcID, vpcPeeringID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*VPCPeeringRoute
	if resp != nil && resp.Data != nil {
		items = make([]*VPCPeeringRoute, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			r := &VPCPeeringRoute{}
			r.fromResponse(&resp.Data.Values[i])
			if r.vpcPeeringID == "" {
				r.vpcPeeringID = vpcPeeringID
			}
			if r.vpcID == "" {
				r.vpcID = vpcID
			}
			if r.projectID == "" {
				r.projectID = projectID
			}
			items = append(items, r)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*VPCPeeringRoute], error) {
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

// vpcPeeringRouteIDsFromRef extracts (projectID, vpcID, vpcPeeringID, vpcPeeringRouteID) from a Ref.
// Accepts the production camelCase segment "vpcPeeringRoutes" and the test form "vpc-peering-routes".
// For the peering parent, accepts both "vpcPeerings" and "peerings".
func vpcPeeringRouteIDsFromRef(ref Ref) (projectID, vpcID, vpcPeeringID, vpcPeeringRouteID string, err error) {
	rid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCPeeringRouteID); ok {
			return w.VPCPeeringRouteID(), true
		}
		return "", false
	}, "vpc-peering-routes")
	if !ok {
		if v := parseURIIDs(ref.URI())["vpcPeeringRoutes"]; v != "" {
			rid = v
			ok = true
		}
	}
	if !ok {
		return "", "", "", "", fmt.Errorf("cannot determine VPC peering route ID from Ref %q", ref.URI())
	}
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCPeeringID); ok {
			return w.VPCPeeringID(), true
		}
		return "", false
	}, "vpc-peerings")
	if !ok {
		m := parseURIIDs(ref.URI())
		if v := m["vpcPeerings"]; v != "" {
			pid = v
			ok = true
		} else if v := m["peerings"]; v != "" {
			pid = v
			ok = true
		}
	}
	if !ok {
		return "", "", "", "", fmt.Errorf("cannot determine VPC peering ID from Ref %q", ref.URI())
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
	projID, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || projID == "" {
		return "", "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return projID, vid, pid, rid, nil
}

type VPCPeeringsClient interface {
	List(ctx context.Context, vpc Ref, opts ...CallOption) (*List[*VPCPeering], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPCPeering, error)
	Create(ctx context.Context, peering *VPCPeering, opts ...CallOption) (*VPCPeering, error)
	Update(ctx context.Context, peering *VPCPeering, opts ...CallOption) (*VPCPeering, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type vpcPeeringLowLevelClient interface {
	List(ctx context.Context, projectID, vpcID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringList], error)
	Get(ctx context.Context, projectID, vpcID, vpcPeeringID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Create(ctx context.Context, projectID, vpcID string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Update(ctx context.Context, projectID, vpcID, vpcPeeringID string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Delete(ctx context.Context, projectID, vpcID, vpcPeeringID string, params *types.RequestParameters) (*types.Response[any], error)
}

type vpcPeeringsClientAdapter struct{ low vpcPeeringLowLevelClient }

func newVPCPeeringsClientAdapter(rest *restclient.Client) *vpcPeeringsClientAdapter {
	if rest == nil {
		return &vpcPeeringsClientAdapter{}
	}
	return &vpcPeeringsClientAdapter{low: network.NewVPCPeeringsClientImpl(rest)}
}

func (a *vpcPeeringsClientAdapter) Create(ctx context.Context, peering *VPCPeering, opts ...CallOption) (*VPCPeering, error) {
	if err := peering.Err(); err != nil {
		return peering, err
	}
	if peering.VPCID() == "" || peering.ProjectID() == "" {
		return peering, fmt.Errorf("Create: VPC peering has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, peering.ProjectID(), peering.VPCID(), peering.toRequest(), rp)
	populateHTTPEnvelope(&peering.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		peering.fromResponse(resp.Data)
	}
	if err != nil {
		return peering, err
	}
	if resp != nil && !resp.IsSuccess() {
		return peering, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return peering, nil
}

func (a *vpcPeeringsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPCPeering, error) {
	projectID, vpcID, vpcPeeringID, err := vpcPeeringIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, vpcPeeringID, rp)
	out := &VPCPeering{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
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

func (a *vpcPeeringsClientAdapter) Update(ctx context.Context, peering *VPCPeering, opts ...CallOption) (*VPCPeering, error) {
	if err := peering.Err(); err != nil {
		return peering, err
	}
	if peering.ID() == "" {
		return peering, fmt.Errorf("Update: VPC peering has no ID — call Get first or seed from response metadata")
	}
	if peering.VPCID() == "" || peering.ProjectID() == "" {
		return peering, fmt.Errorf("Update: VPC peering has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, peering.ProjectID(), peering.VPCID(), peering.ID(), peering.toRequest(), rp)
	populateHTTPEnvelope(&peering.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		peering.fromResponse(resp.Data)
	}
	if err != nil {
		return peering, err
	}
	if resp != nil && !resp.IsSuccess() {
		return peering, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return peering, nil
}

func (a *vpcPeeringsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, vpcPeeringID, err := vpcPeeringIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, vpcPeeringID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *vpcPeeringsClientAdapter) List(ctx context.Context, vpc Ref, opts ...CallOption) (*List[*VPCPeering], error) {
	projectID, vpcID, err := vpcIDsFromRef(vpc)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, vpcID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*VPCPeering
	if resp != nil && resp.Data != nil {
		items = make([]*VPCPeering, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			p := &VPCPeering{}
			p.fromResponse(&resp.Data.Values[i])
			if p.vpcID == "" {
				p.vpcID = vpcID
			}
			if p.projectID == "" {
				p.projectID = projectID
			}
			items = append(items, p)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*VPCPeering], error) {
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

// vpcPeeringIDsFromRef extracts (projectID, vpcID, vpcPeeringID) from a Ref.
// Accepts the production camelCase segment "vpcPeerings" and the mixin/test form "peerings".
func vpcPeeringIDsFromRef(ref Ref) (projectID, vpcID, vpcPeeringID string, err error) {
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCPeeringID); ok {
			return w.VPCPeeringID(), true
		}
		return "", false
	}, "vpc-peerings")
	if !ok || pid == "" {
		m := parseURIIDs(ref.URI())
		if v := m["vpcPeerings"]; v != "" {
			pid = v
			ok = true
		}
		if pid == "" {
			if v := m["peerings"]; v != "" {
				pid = v
				ok = true
			}
		}
	}
	if !ok || pid == "" {
		return "", "", "", fmt.Errorf("cannot determine VPC peering ID from Ref %q", ref.URI())
	}
	vid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid == "" {
		return "", "", "", fmt.Errorf("cannot determine VPC ID from Ref %q", ref.URI())
	}
	projID, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || projID == "" {
		return "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return projID, vid, pid, nil
}

type VPCsClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*VPC], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPC, error)
	Create(ctx context.Context, vpc *VPC, opts ...CallOption) (*VPC, error)
	Update(ctx context.Context, vpc *VPC, opts ...CallOption) (*VPC, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type vpcLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.VPCList], error)
	Get(ctx context.Context, projectID, vpcID string, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	Create(ctx context.Context, projectID string, body types.VPCRequest, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	Update(ctx context.Context, projectID, vpcID string, body types.VPCRequest, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	Delete(ctx context.Context, projectID, vpcID string, params *types.RequestParameters) (*types.Response[any], error)
}

type vpcsClientAdapter struct{ low vpcLowLevelClient }

func newVPCsClientAdapter(rest *restclient.Client) *vpcsClientAdapter {
	if rest == nil {
		return &vpcsClientAdapter{}
	}
	return &vpcsClientAdapter{low: network.NewVPCsClientImpl(rest)}
}

func (a *vpcsClientAdapter) Create(ctx context.Context, v *VPC, opts ...CallOption) (*VPC, error) {
	if err := v.Err(); err != nil {
		return v, err
	}
	if v.ProjectID() == "" {
		return v, fmt.Errorf("Create: VPC has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, v.ProjectID(), v.toRequest(), rp)
	populateHTTPEnvelope(&v.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		v.fromResponse(resp.Data)
	}
	if err != nil {
		return v, err
	}
	if resp != nil && !resp.IsSuccess() {
		return v, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return v, nil
}

func (a *vpcsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPC, error) {
	projectID, vpcID, err := vpcIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, rp)
	out := &VPC{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *vpcsClientAdapter) Update(ctx context.Context, v *VPC, opts ...CallOption) (*VPC, error) {
	if err := v.Err(); err != nil {
		return v, err
	}
	if v.ID() == "" {
		return v, fmt.Errorf("Update: VPC has no ID — call Get first or seed from response metadata")
	}
	if v.ProjectID() == "" {
		return v, fmt.Errorf("Update: VPC has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, v.ProjectID(), v.ID(), v.toRequest(), rp)
	populateHTTPEnvelope(&v.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		v.fromResponse(resp.Data)
	}
	if err != nil {
		return v, err
	}
	if resp != nil && !resp.IsSuccess() {
		return v, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return v, nil
}

func (a *vpcsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, err := vpcIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *vpcsClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*VPC], error) {
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
	var items []*VPC
	if resp != nil && resp.Data != nil {
		items = make([]*VPC, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			v := &VPC{}
			v.fromResponse(&resp.Data.Values[i])
			items = append(items, v)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*VPC], error) {
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

// vpcIDsFromRef extracts (projectID, vpcID) from a Ref. Tries typed assertions
// first, then falls back to URI path parsing.
func vpcIDsFromRef(ref Ref) (projectID, vpcID string, err error) {
	vid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid == "" {
		return "", "", fmt.Errorf("cannot determine VPC ID from Ref %q", ref.URI())
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
	return pid, vid, nil
}

type VPNRoutesClient interface {
	List(ctx context.Context, tunnel Ref, opts ...CallOption) (*List[*VPNRoute], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPNRoute, error)
	Create(ctx context.Context, r *VPNRoute, opts ...CallOption) (*VPNRoute, error)
	Update(ctx context.Context, r *VPNRoute, opts ...CallOption) (*VPNRoute, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

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

type VPNTunnelsClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*VPNTunnel], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPNTunnel, error)
	Create(ctx context.Context, t *VPNTunnel, opts ...CallOption) (*VPNTunnel, error)
	Update(ctx context.Context, t *VPNTunnel, opts ...CallOption) (*VPNTunnel, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

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

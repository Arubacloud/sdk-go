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
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.ElasticList], error)
	Get(ctx context.Context, projectID string, elasticIPID string, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Create(ctx context.Context, projectID string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Update(ctx context.Context, projectID string, elasticIPID string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Delete(ctx context.Context, projectID string, elasticIPID string, params *types.RequestParameters) (*types.Response[any], error)
}

type LoadBalancersClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.LoadBalancerList], error)
	Get(ctx context.Context, projectID string, loadBalancerID string, params *types.RequestParameters) (*types.Response[types.LoadBalancerResponse], error)
}

type SecurityGroupRulesClient interface {
	List(ctx context.Context, projectID string, vpcID string, securityGroupID string, params *types.RequestParameters) (*types.Response[types.SecurityRuleList], error)
	Get(ctx context.Context, projectID string, vpcID string, securityGroupID string, securityGroupRuleID string, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	Create(ctx context.Context, projectID string, vpcID string, securityGroupID string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	Update(ctx context.Context, projectID string, vpcID string, securityGroupID string, securityGroupRuleID string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	Delete(ctx context.Context, projectID string, vpcID string, securityGroupID string, securityGroupRuleID string, params *types.RequestParameters) (*types.Response[any], error)
}

type SecurityGroupsClient interface {
	List(ctx context.Context, projectID string, vpcID string, params *types.RequestParameters) (*types.Response[types.SecurityGroupList], error)
	Get(ctx context.Context, projectID string, vpcID string, securityGroupID string, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Create(ctx context.Context, projectID string, vpcID string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Update(ctx context.Context, projectID string, vpcID string, securityGroupID string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Delete(ctx context.Context, projectID string, vpcID string, securityGroupID string, params *types.RequestParameters) (*types.Response[any], error)
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
	List(ctx context.Context, projectID string, vpcID string, vpcPeeringID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteList], error)
	Get(ctx context.Context, projectID string, vpcID string, vpcPeeringID string, vpcPeeringRouteID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Create(ctx context.Context, projectID string, vpcID string, vpcPeeringID string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Update(ctx context.Context, projectID string, vpcID string, vpcPeeringID string, vpcPeeringRouteID string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Delete(ctx context.Context, projectID string, vpcID string, vpcPeeringID string, vpcPeeringRouteID string, params *types.RequestParameters) (*types.Response[any], error)
}

type VPCPeeringsClient interface {
	List(ctx context.Context, projectID string, vpcID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringList], error)
	Get(ctx context.Context, projectID string, vpcID string, vpcPeeringID string, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Create(ctx context.Context, projectID string, vpcID string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Update(ctx context.Context, projectID string, vpcID string, vpcPeeringID string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Delete(ctx context.Context, projectID string, vpcID string, vpcPeeringID string, params *types.RequestParameters) (*types.Response[any], error)
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
	List(ctx context.Context, projectID string, vpnTunnelID string, params *types.RequestParameters) (*types.Response[types.VPNRouteList], error)
	Get(ctx context.Context, projectID string, vpnTunnelID string, vpnRouteID string, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	Create(ctx context.Context, projectID string, vpnTunnelID string, body types.VPNRouteRequest, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	Update(ctx context.Context, projectID string, vpnTunnelID string, vpnRouteID string, body types.VPNRouteRequest, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	Delete(ctx context.Context, projectID string, vpnTunnelID string, vpnRouteID string, params *types.RequestParameters) (*types.Response[any], error)
}

type VPNTunnelsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.VPNTunnelList], error)
	Get(ctx context.Context, projectID string, vpnTunnelID string, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	Create(ctx context.Context, projectID string, body types.VPNTunnelRequest, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	Update(ctx context.Context, projectID string, vpnTunnelID string, body types.VPNTunnelRequest, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	Delete(ctx context.Context, projectID string, vpnTunnelID string, params *types.RequestParameters) (*types.Response[any], error)
}

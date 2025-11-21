package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
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
	Get(ctx context.Context, projectID string, elasticIPId string, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Create(ctx context.Context, projectID string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Update(ctx context.Context, projectID string, elasticIPId string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Delete(ctx context.Context, projectID string, elasticIPId string, params *types.RequestParameters) (*types.Response[any], error)
}

type LoadBalancersClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.LoadBalancerList], error)
	Get(ctx context.Context, projectID string, loadBalancerId string, params *types.RequestParameters) (*types.Response[types.LoadBalancerResponse], error)
}

type SecurityGroupRulesClient interface {
	List(ctx context.Context, projectID string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[types.SecurityRuleList], error)
	Get(ctx context.Context, projectID string, vpcId string, securityGroupId string, securityGroupRuleId string, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	Create(ctx context.Context, projectID string, vpcId string, securityGroupId string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	Update(ctx context.Context, projectID string, vpcId string, securityGroupId string, securityGroupRuleId string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	Delete(ctx context.Context, projectID string, vpcId string, securityGroupId string, securityGroupRuleId string, params *types.RequestParameters) (*types.Response[any], error)
}

type SecurityGroupsClient interface {
	List(ctx context.Context, projectID string, vpcId string, params *types.RequestParameters) (*types.Response[types.SecurityGroupList], error)
	Get(ctx context.Context, projectID string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Create(ctx context.Context, projectID string, vpcId string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Update(ctx context.Context, projectID string, vpcId string, securityGroupId string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Delete(ctx context.Context, projectID string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[any], error)
}

type SubnetsClient interface {
	List(ctx context.Context, projectID string, vpcId string, params *types.RequestParameters) (*types.Response[types.SubnetList], error)
	Get(ctx context.Context, projectID string, vpcId string, subnetId string, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	Create(ctx context.Context, projectID string, vpcId string, body types.SubnetRequest, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	Update(ctx context.Context, projectID string, vpcId string, subnetId string, body types.SubnetRequest, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	Delete(ctx context.Context, projectID string, vpcId string, subnetId string, params *types.RequestParameters) (*types.Response[any], error)
}

type VPCPeeringRoutesClient interface {
	List(ctx context.Context, projectID string, vpcId string, vpcPeeringId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteList], error)
	Get(ctx context.Context, projectID string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Create(ctx context.Context, projectID string, vpcId string, vpcPeeringId string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Update(ctx context.Context, projectID string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	Delete(ctx context.Context, projectID string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *types.RequestParameters) (*types.Response[any], error)
}

type VPCPeeringsClient interface {
	List(ctx context.Context, projectID string, vpcId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringList], error)
	Get(ctx context.Context, projectID string, vpcId string, vpcPeeringId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Create(ctx context.Context, projectID string, vpcId string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Update(ctx context.Context, projectID string, vpcId string, vpcPeeringId string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	Delete(ctx context.Context, projectID string, vpcId string, vpcPeeringId string, params *types.RequestParameters) (*types.Response[any], error)
}

type VPCsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.VPCList], error)
	Get(ctx context.Context, projectID string, vpcId string, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	Create(ctx context.Context, projectID string, body types.VPCRequest, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	Update(ctx context.Context, projectID string, vpcId string, body types.VPCRequest, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	Delete(ctx context.Context, projectID string, vpcId string, params *types.RequestParameters) (*types.Response[any], error)
}

type VPNRoutesClient interface {
	List(ctx context.Context, projectID string, vpnTunnelId string, params *types.RequestParameters) (*types.Response[types.VPNRouteList], error)
	Get(ctx context.Context, projectID string, vpnTunnelId string, vpnRouteId string, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	Create(ctx context.Context, projectID string, vpnTunnelId string, body types.VPNRouteRequest, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	Update(ctx context.Context, projectID string, vpnTunnelId string, vpnRouteId string, body types.VPNRouteRequest, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	Delete(ctx context.Context, projectID string, vpnTunnelId string, vpnRouteId string, params *types.RequestParameters) (*types.Response[any], error)
}

type VPNTunnelsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.VPNTunnelList], error)
	Get(ctx context.Context, projectID string, vpnTunnelId string, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	Create(ctx context.Context, projectID string, body types.VPNTunnelRequest, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	Update(ctx context.Context, projectID string, vpnTunnelId string, body types.VPNTunnelRequest, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	Delete(ctx context.Context, projectID string, vpnTunnelId string, params *types.RequestParameters) (*types.Response[any], error)
}

package network

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

// NetworkAPI defines the unified interface for all Network operations
type NetworkAPI interface {
	// ElasticIP operations
	ListElasticIPs(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.ElasticList], error)
	GetElasticIP(ctx context.Context, project string, elasticIPId string, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	CreateElasticIP(ctx context.Context, project string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	UpdateElasticIP(ctx context.Context, project string, elasticIPId string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	DeleteElasticIP(ctx context.Context, projectId string, elasticIPId string, params *types.RequestParameters) (*types.Response[any], error)

	// LoadBalancer operations
	ListLoadBalancers(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.LoadBalancerList], error)
	GetLoadBalancer(ctx context.Context, project string, loadBalancerId string, params *types.RequestParameters) (*types.Response[types.LoadBalancerResponse], error)

	// VPC operations
	ListVPCs(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.VPCList], error)
	GetVPC(ctx context.Context, project string, vpcId string, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	CreateVPC(ctx context.Context, project string, body types.VPCRequest, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	UpdateVPC(ctx context.Context, project string, vpcId string, body types.VPCRequest, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	DeleteVPC(ctx context.Context, projectId string, vpcId string, params *types.RequestParameters) (*types.Response[any], error)

	// Subnet operations
	ListSubnets(ctx context.Context, project string, vpcId string, params *types.RequestParameters) (*types.Response[types.SubnetList], error)
	GetSubnet(ctx context.Context, project string, vpcId string, subnetId string, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	CreateSubnet(ctx context.Context, project string, vpcId string, body types.SubnetRequest, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	UpdateSubnet(ctx context.Context, project string, vpcId string, subnetId string, body types.SubnetRequest, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	DeleteSubnet(ctx context.Context, projectId string, vpcId string, subnetId string, params *types.RequestParameters) (*types.Response[any], error)

	// SecurityGroup operations
	ListSecurityGroups(ctx context.Context, project string, vpcId string, params *types.RequestParameters) (*types.Response[types.SecurityGroupList], error)
	GetSecurityGroup(ctx context.Context, project string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	CreateSecurityGroup(ctx context.Context, project string, vpcId string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	UpdateSecurityGroup(ctx context.Context, project string, vpcId string, securityGroupId string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	DeleteSecurityGroup(ctx context.Context, projectId string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[any], error)

	// SecurityGroupRule operations
	ListSecurityGroupRules(ctx context.Context, project string, vpcId string, securityGroupId string, params *types.RequestParameters) (*types.Response[types.SecurityRuleList], error)
	GetSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	CreateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	UpdateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, body types.SecurityRuleRequest, params *types.RequestParameters) (*types.Response[types.SecurityRuleResponse], error)
	DeleteSecurityGroupRule(ctx context.Context, projectId string, vpcId string, securityGroupId string, securityGroupRuleId string, params *types.RequestParameters) (*types.Response[any], error)

	// VpcPeering operations
	ListVpcPeerings(ctx context.Context, project string, vpcId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringList], error)
	GetVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	CreateVpcPeering(ctx context.Context, project string, vpcId string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	UpdateVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, body types.VPCPeeringRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringResponse], error)
	DeleteVpcPeering(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, params *types.RequestParameters) (*types.Response[any], error)

	// VpcPeeringRoute operations
	ListVpcPeeringRoutes(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteList], error)
	GetVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	CreateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	UpdateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, body types.VPCPeeringRouteRequest, params *types.RequestParameters) (*types.Response[types.VPCPeeringRouteResponse], error)
	DeleteVpcPeeringRoute(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *types.RequestParameters) (*types.Response[any], error)

	// VpnTunnel operations
	ListVpnTunnels(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.VPNTunnelList], error)
	GetVpnTunnel(ctx context.Context, project string, vpnTunnelId string, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	CreateVpnTunnel(ctx context.Context, project string, body types.VPNTunnelRequest, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	UpdateVpnTunnel(ctx context.Context, project string, vpnTunnelId string, body types.VPNTunnelRequest, params *types.RequestParameters) (*types.Response[types.VPNTunnelResponse], error)
	DeleteVpnTunnel(ctx context.Context, projectId string, vpnTunnelId string, params *types.RequestParameters) (*types.Response[any], error)

	// VpnRoute operations
	ListVpnRoutes(ctx context.Context, project string, vpnTunnelId string, params *types.RequestParameters) (*types.Response[types.VPNRouteList], error)
	GetVpnRoute(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	CreateVpnRoute(ctx context.Context, project string, vpnTunnelId string, body types.VPNRouteRequest, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	UpdateVpnRoute(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, body types.VPNRouteRequest, params *types.RequestParameters) (*types.Response[types.VPNRouteResponse], error)
	DeleteVpnRoute(ctx context.Context, projectId string, vpnTunnelId string, vpnRouteId string, params *types.RequestParameters) (*types.Response[any], error)
}

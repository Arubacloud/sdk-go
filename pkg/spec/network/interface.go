package network

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// NetworkAPI defines the unified interface for all Network operations
type NetworkAPI interface {
	// ElasticIP operations
	ListElasticIPs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.ElasticList], error)
	GetElasticIP(ctx context.Context, project string, elasticIPId string, params *schema.RequestParameters) (*schema.Response[schema.ElasticIPResponse], error)
	CreateElasticIP(ctx context.Context, project string, body schema.ElasticIPRequest, params *schema.RequestParameters) (*schema.Response[schema.ElasticIPResponse], error)
	UpdateElasticIP(ctx context.Context, project string, elasticIPId string, body schema.ElasticIPRequest, params *schema.RequestParameters) (*schema.Response[schema.ElasticIPResponse], error)
	DeleteElasticIP(ctx context.Context, projectId string, elasticIPId string, params *schema.RequestParameters) (*schema.Response[any], error)

	// LoadBalancer operations
	ListLoadBalancers(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.LoadBalancerList], error)
	GetLoadBalancer(ctx context.Context, project string, loadBalancerId string, params *schema.RequestParameters) (*schema.Response[schema.LoadBalancerResponse], error)

	// VPC operations
	ListVPCs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.VPCList], error)
	GetVPC(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.VPCResponse], error)
	CreateVPC(ctx context.Context, project string, body schema.VPCRequest, params *schema.RequestParameters) (*schema.Response[schema.VPCResponse], error)
	UpdateVPC(ctx context.Context, project string, vpcId string, body schema.VPCRequest, params *schema.RequestParameters) (*schema.Response[schema.VPCResponse], error)
	DeleteVPC(ctx context.Context, projectId string, vpcId string, params *schema.RequestParameters) (*schema.Response[any], error)

	// Subnet operations
	ListSubnets(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.SubnetList], error)
	GetSubnet(ctx context.Context, project string, vpcId string, subnetId string, params *schema.RequestParameters) (*schema.Response[schema.SubnetResponse], error)
	CreateSubnet(ctx context.Context, project string, vpcId string, body schema.SubnetRequest, params *schema.RequestParameters) (*schema.Response[schema.SubnetResponse], error)
	UpdateSubnet(ctx context.Context, project string, vpcId string, subnetId string, body schema.SubnetRequest, params *schema.RequestParameters) (*schema.Response[schema.SubnetResponse], error)
	DeleteSubnet(ctx context.Context, projectId string, vpcId string, subnetId string, params *schema.RequestParameters) (*schema.Response[any], error)

	// SecurityGroup operations
	ListSecurityGroups(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.SecurityGroupList], error)
	GetSecurityGroup(ctx context.Context, project string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*schema.Response[schema.SecurityGroupResponse], error)
	CreateSecurityGroup(ctx context.Context, project string, vpcId string, body schema.SecurityGroupRequest, params *schema.RequestParameters) (*schema.Response[schema.SecurityGroupResponse], error)
	UpdateSecurityGroup(ctx context.Context, project string, vpcId string, securityGroupId string, body schema.SecurityGroupRequest, params *schema.RequestParameters) (*schema.Response[schema.SecurityGroupResponse], error)
	DeleteSecurityGroup(ctx context.Context, projectId string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*schema.Response[any], error)

	// SecurityGroupRule operations
	ListSecurityGroupRules(ctx context.Context, project string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleList], error)
	GetSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleResponse], error)
	CreateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, body schema.SecurityRuleRequest, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleResponse], error)
	UpdateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, body schema.SecurityRuleRequest, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleResponse], error)
	DeleteSecurityGroupRule(ctx context.Context, projectId string, vpcId string, securityGroupId string, securityGroupRuleId string, params *schema.RequestParameters) (*schema.Response[any], error)

	// VpcPeering operations
	ListVpcPeerings(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringList], error)
	GetVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringResponse], error)
	CreateVpcPeering(ctx context.Context, project string, vpcId string, body schema.VPCPeeringRequest, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringResponse], error)
	UpdateVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, body schema.VPCPeeringRequest, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringResponse], error)
	DeleteVpcPeering(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*schema.Response[any], error)

	// VpcPeeringRoute operations
	ListVpcPeeringRoutes(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringRouteList], error)
	GetVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringRouteResponse], error)
	CreateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, body schema.VPCPeeringRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringRouteResponse], error)
	UpdateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, body schema.VPCPeeringRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VPCPeeringRouteResponse], error)
	DeleteVpcPeeringRoute(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*schema.Response[any], error)

	// VpnTunnel operations
	ListVpnTunnels(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.VPNTunnelList], error)
	GetVpnTunnel(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[schema.VPNTunnelResponse], error)
	CreateVpnTunnel(ctx context.Context, project string, body schema.VPNTunnelRequest, params *schema.RequestParameters) (*schema.Response[schema.VPNTunnelResponse], error)
	UpdateVpnTunnel(ctx context.Context, project string, vpnTunnelId string, body schema.VPNTunnelRequest, params *schema.RequestParameters) (*schema.Response[schema.VPNTunnelResponse], error)
	DeleteVpnTunnel(ctx context.Context, projectId string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[any], error)

	// VpnRoute operations
	ListVpnRoutes(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[schema.VPNRouteList], error)
	GetVpnRoute(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, params *schema.RequestParameters) (*schema.Response[schema.VPNRouteResponse], error)
	CreateVpnRoute(ctx context.Context, project string, vpnTunnelId string, body schema.VPNRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VPNRouteResponse], error)
	UpdateVpnRoute(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, body schema.VPNRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VPNRouteResponse], error)
	DeleteVpnRoute(ctx context.Context, projectId string, vpnTunnelId string, vpnRouteId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

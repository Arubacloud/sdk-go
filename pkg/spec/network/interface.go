package network

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ElasticIPAPI defines the interface for Elastic IP operations
type ElasticIPAPI interface {
	ListElasticIPs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.ElasticList], error)
	GetElasticIP(ctx context.Context, project string, elasticIPId string, params *schema.RequestParameters) (*schema.Response[schema.ElasticIpResponse], error)
	CreateElasticIP(ctx context.Context, project string, body schema.ElasticIpRequest, params *schema.RequestParameters) (*schema.Response[schema.ElasticIpResponse], error)
	UpdateElasticIP(ctx context.Context, project string, elasticIPId string, body schema.ElasticIpRequest, params *schema.RequestParameters) (*schema.Response[schema.ElasticIpResponse], error)
	DeleteElasticIP(ctx context.Context, projectId string, elasticIPId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// LoadBalancerAPI defines the interface for Load Balancer operations
type LoadBalancerAPI interface {
	ListLoadBalancers(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.LoadBalancerList], error)
	GetLoadBalancer(ctx context.Context, project string, loadBalancerId string, params *schema.RequestParameters) (*schema.Response[schema.LoadBalancerResponse], error)
}

// VpcAPI defines the interface for VPC operations
type VpcAPI interface {
	ListVPCs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.VpcList], error)
	GetVPC(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.VpcResponse], error)
	CreateVPC(ctx context.Context, project string, body schema.VpcRequest, params *schema.RequestParameters) (*schema.Response[schema.VpcResponse], error)
	UpdateVPC(ctx context.Context, project string, vpcId string, body schema.VpcRequest, params *schema.RequestParameters) (*schema.Response[schema.VpcResponse], error)
	DeleteVPC(ctx context.Context, projectId string, vpcId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// SubnetAPI defines the interface for Subnet operations
type SubnetAPI interface {
	ListSubnets(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.SubnetList], error)
	GetSubnet(ctx context.Context, project string, vpcId string, subnetId string, params *schema.RequestParameters) (*schema.Response[schema.SubnetResponse], error)
	CreateSubnet(ctx context.Context, project string, vpcId string, body schema.SubnetRequest, params *schema.RequestParameters) (*schema.Response[schema.SubnetResponse], error)
	UpdateSubnet(ctx context.Context, project string, vpcId string, subnetId string, body schema.SubnetRequest, params *schema.RequestParameters) (*schema.Response[schema.SubnetResponse], error)
	DeleteSubnet(ctx context.Context, projectId string, vpcId string, subnetId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// SecurityGroupAPI defines the interface for SecurityGroup operations
type SecurityGroupAPI interface {
	ListSecurityGroups(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.SecurityGroupList], error)
	GetSecurityGroup(ctx context.Context, project string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*schema.Response[schema.SecurityGroupResponse], error)
	CreateSecurityGroup(ctx context.Context, project string, vpcId string, body schema.SecurityGroupRequest, params *schema.RequestParameters) (*schema.Response[schema.SecurityGroupResponse], error)
	UpdateSecurityGroup(ctx context.Context, project string, vpcId string, securityGroupId string, body schema.SecurityGroupRequest, params *schema.RequestParameters) (*schema.Response[schema.SecurityGroupResponse], error)
	DeleteSecurityGroup(ctx context.Context, projectId string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// SecurityGroupRuleAPI defines the interface for SecurityGroupRule operations
type SecurityGroupRuleAPI interface {
	ListSecurityGroupRules(ctx context.Context, project string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleList], error)
	GetSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleResponse], error)
	CreateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, body schema.SecurityRuleRequest, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleResponse], error)
	UpdateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, body schema.SecurityRuleRequest, params *schema.RequestParameters) (*schema.Response[schema.SecurityRuleResponse], error)
	DeleteSecurityGroupRule(ctx context.Context, projectId string, vpcId string, securityGroupId string, securityGroupRuleId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// VpcPeeringAPI defines the interface for VPC Peering operations
type VpcPeeringAPI interface {
	ListVpcPeerings(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringList], error)
	GetVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringResponse], error)
	CreateVpcPeering(ctx context.Context, project string, vpcId string, body schema.VpcPeeringRequest, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringResponse], error)
	UpdateVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, body schema.VpcPeeringRequest, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringResponse], error)
	DeleteVpcPeering(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// VpcPeeringRouteAPI defines the interface for VPC Peering Route operations
type VpcPeeringRouteAPI interface {
	ListVpcPeeringRoutes(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringRouteList], error)
	GetVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringRouteResponse], error)
	CreateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, body schema.VpcPeeringRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringRouteResponse], error)
	UpdateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, body schema.VpcPeeringRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringRouteResponse], error)
	DeleteVpcPeeringRoute(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// VpnTunnelAPI defines the interface for VPN Tunnel operations
type VpnTunnelAPI interface {
	ListVpnTunnels(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.VpnTunnelList], error)
	GetVpnTunnel(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[schema.VpnTunnelResponse], error)
	CreateVpnTunnel(ctx context.Context, project string, body schema.VpnTunnelRequest, params *schema.RequestParameters) (*schema.Response[schema.VpnTunnelResponse], error)
	UpdateVpnTunnel(ctx context.Context, project string, vpnTunnelId string, body schema.VpnTunnelRequest, params *schema.RequestParameters) (*schema.Response[schema.VpnTunnelResponse], error)
	DeleteVpnTunnel(ctx context.Context, projectId string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// VpnRouteAPI defines the interface for VPN Route operations
type VpnRouteAPI interface {
	ListVpnRoutes(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*schema.Response[schema.VpnRouteList], error)
	GetVpnRoute(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, params *schema.RequestParameters) (*schema.Response[schema.VpnRouteResponse], error)
	CreateVpnRoute(ctx context.Context, project string, vpnTunnelId string, body schema.VpnRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VpnRouteResponse], error)
	UpdateVpnRoute(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, body schema.VpnRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VpnRouteResponse], error)
	DeleteVpnRoute(ctx context.Context, projectId string, vpnTunnelId string, vpnRouteId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

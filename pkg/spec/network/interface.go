package network

import (
	"context"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ElasticIPAPI defines the interface for Elastic IP operations
type ElasticIPAPI interface {
	ListElasticIPs(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetElasticIP(ctx context.Context, project string, elasticIPId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateElasticIP(ctx context.Context, project string, body schema.Elas, params *schema.RequestParameters) (*http.Response, error)
	DeleteElasticIP(ctx context.Context, projectId string, elasticIPId string, params *schema.RequestParameters) (*http.Response, error)
}

// LoadBalancerAPI defines the interface for Load Balancer operations
type LoadBalancerAPI interface {
	ListLoadBalancers(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetLoadBalancer(ctx context.Context, project string, loadBalancerId string, params *schema.RequestParameters) (*http.Response, error)
}

// VpcAPI defines the interface for VPC operations
type VpcAPI interface {
	ListVPCs(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetVPC(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateVPC(ctx context.Context, project string, body schema.VpcRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteVPC(ctx context.Context, projectId string, vpcId string, params *schema.RequestParameters) (*http.Response, error)
}

// SubnetAPI defines the interface for Subnet operations
type SubnetAPI interface {
	ListSubnets(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*http.Response, error)
	GetSubnet(ctx context.Context, project string, vpcId string, subnetId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateSubnet(ctx context.Context, project string, vpcId string, body schema.SubnetRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteSubnet(ctx context.Context, projectId string, vpcId string, subnetId string, params *schema.RequestParameters) (*http.Response, error)
}

// SecurityGroupAPI defines the interface for SecurityGroup operations
type SecurityGroupAPI interface {
	ListSecurityGroups(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*http.Response, error)
	GetSecurityGroup(ctx context.Context, project string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateSecurityGroup(ctx context.Context, project string, vpcId string, body schema.SecurityGroupRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteSecurityGroup(ctx context.Context, projectId string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*http.Response, error)
}

// SecurityGroupRuleAPI defines the interface for SecurityGroupRule operations
type SecurityGroupRuleAPI interface {
	ListSecurityGroupRules(ctx context.Context, project string, vpcId string, securityGroupId string, params *schema.RequestParameters) (*http.Response, error)
	GetSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, securityGroupRuleId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateSecurityGroupRule(ctx context.Context, project string, vpcId string, securityGroupId string, body schema.SecurityRuleRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteSecurityGroupRule(ctx context.Context, projectId string, vpcId string, securityGroupId string, securityGroupRuleId string, params *schema.RequestParameters) (*http.Response, error)
}

// VpcPeeringAPI defines the interface for VPC Peering operations
type VpcPeeringAPI interface {
	ListVpcPeerings(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*http.Response, error)
	GetVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateVpcPeering(ctx context.Context, project string, vpcId string, body schema.VpcPeeringRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteVpcPeering(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*http.Response, error)
}

// VpcPeeringRouteAPI defines the interface for VPC Peering Route operations
type VpcPeeringRouteAPI interface {
	ListVpcPeeringRoutes(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*http.Response, error)
	GetVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, body schema.VpcPeeringRouteRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteVpcPeeringRoute(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*http.Response, error)
}

// VpnTunnelAPI defines the interface for VPN Tunnel operations
type VpnTunnelAPI interface {
	ListVpnTunnels(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetVpnTunnel(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateVpnTunnel(ctx context.Context, project string, body schema.VpnTunnelRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteVpnTunnel(ctx context.Context, projectId string, vpnTunnelId string, params *schema.RequestParameters) (*http.Response, error)
}

// VpnRouteAPI defines the interface for VPN Route operations
type VpnRouteAPI interface {
	ListVpnRoutes(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*http.Response, error)
	GetVpnRoute(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateVpnRoute(ctx context.Context, project string, vpnTunnelId string, body schema.VpnRouteRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteVpnRoute(ctx context.Context, projectId string, vpnTunnelId string, vpnRouteId string, params *schema.RequestParameters) (*http.Response, error)
}

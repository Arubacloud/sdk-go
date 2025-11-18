package network

// API path constants for network resources
const (
	// VPC Network paths
	VPCNetworksPath = "/projects/%s/providers/Aruba.Network/vpcs"
	VPCNetworkPath  = "/projects/%s/providers/Aruba.Network/vpcs/%s"

	// Subnet paths (nested under VPC)
	SubnetsPath = "/projects/%s/providers/Aruba.Network/vpcs/%s/subnets"
	SubnetPath  = "/projects/%s/providers/Aruba.Network/vpcs/%s/subnets/%s"

	// Security Group paths (nested under VPC)
	SecurityGroupsPath = "/projects/%s/providers/Aruba.Network/vpcs/%s/securitygroups"
	SecurityGroupPath  = "/projects/%s/providers/Aruba.Network/vpcs/%s/securitygroups/%s"

	// Security Group Rule paths (nested under VPC and Security Group)
	SecurityGroupRulesPath = "/projects/%s/providers/Aruba.Network/vpcs/%s/securitygroups/%s/securityrules"
	SecurityGroupRulePath  = "/projects/%s/providers/Aruba.Network/vpcs/%s/securitygroups/%s/securityrules/%s"

	// Elastic IP paths
	ElasticIPsPath = "/projects/%s/providers/Aruba.Network/elasticIps"
	ElasticIPPath  = "/projects/%s/providers/Aruba.Network/elasticIps/%s"

	// Load Balancer paths
	LoadBalancersPath = "/projects/%s/providers/Aruba.Network/loadbalancers"
	LoadBalancerPath  = "/projects/%s/providers/Aruba.Network/loadbalancers/%s"

	// VPC Peering Connection paths
	VPCPeeringsPath = "/projects/%s/providers/Aruba.Network/vpcs/%s/vpcPeerings"
	VPCPeeringPath  = "/projects/%s/providers/Aruba.Network/vpcs/%s/vpcPeerings/%s"

	// VPC Peering Route paths
	VPCPeeringRoutesPath = "/projects/%s/providers/Aruba.Network/vpcs/%s/vpcPeerings/%s/routes"
	VPCPeeringRoutePath  = "/projects/%s/providers/Aruba.Network/vpcs/%s/vpcPeerings/%s/routes/%s"

	// VPN Tunnel paths
	VPNTunnelsPath = "/projects/%s/providers/Aruba.Network/vpntunnels"
	VPNTunnelPath  = "/projects/%s/providers/Aruba.Network/vpntunnels/%s"

	// VPN Route paths
	VPNRoutesPath = "/projects/%s/providers/Aruba.Network/vpntunnels/%s/routes"
	VPNRoutePath  = "/projects/%s/providers/Aruba.Network/vpntunnels/%s/routes/%s"
)

package aruba

import "fmt"

// Typed Ref constructors for network resources. Each constructor builds a Ref
// whose URI matches the canonical Aruba Cloud API path for that resource type.
// Use these instead of hand-building URI strings in downstream consumers.

// VPCRef returns a Ref that points to the VPC with the given IDs.
func VPCRef(projectID, vpcID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/vpcs/%s", projectID, vpcID))
}

// SubnetRef returns a Ref that points to the Subnet nested under a VPC.
func SubnetRef(projectID, vpcID, subnetID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/vpcs/%s/subnets/%s", projectID, vpcID, subnetID))
}

// SecurityGroupRef returns a Ref that points to the SecurityGroup nested under a VPC.
func SecurityGroupRef(projectID, vpcID, sgID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/vpcs/%s/securitygroups/%s", projectID, vpcID, sgID))
}

// SecurityRuleRef returns a Ref that points to the SecurityRule nested under a SecurityGroup.
func SecurityRuleRef(projectID, vpcID, sgID, ruleID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/vpcs/%s/securitygroups/%s/securityrules/%s", projectID, vpcID, sgID, ruleID))
}

// ElasticIPRef returns a Ref that points to the ElasticIP with the given IDs.
func ElasticIPRef(projectID, eipID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/elasticIps/%s", projectID, eipID))
}

// LoadBalancerRef returns a Ref that points to the LoadBalancer with the given IDs.
func LoadBalancerRef(projectID, lbID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/loadbalancers/%s", projectID, lbID))
}

// VPCPeeringRef returns a Ref that points to the VPCPeering nested under a VPC.
func VPCPeeringRef(projectID, vpcID, peeringID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/vpcs/%s/vpcPeerings/%s", projectID, vpcID, peeringID))
}

// VPCPeeringRouteRef returns a Ref that points to the VPCPeeringRoute nested under a VPCPeering.
func VPCPeeringRouteRef(projectID, vpcID, peeringID, routeID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/vpcs/%s/vpcPeerings/%s/vpcPeeringRoutes/%s", projectID, vpcID, peeringID, routeID))
}

// VPNTunnelRef returns a Ref that points to the VPNTunnel with the given IDs.
func VPNTunnelRef(projectID, tunnelID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/vpnTunnels/%s", projectID, tunnelID))
}

// VPNRouteRef returns a Ref that points to the VPNRoute nested under a VPNTunnel.
func VPNRouteRef(projectID, tunnelID, routeID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/vpnTunnels/%s/vpnRoutes/%s", projectID, tunnelID, routeID))
}

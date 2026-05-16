package aruba

import "testing"

func TestVPCRef(t *testing.T) {
	ref := VPCRef("p-1", "vpc-1")
	want := "/projects/p-1/providers/Aruba.Network/vpcs/vpc-1"
	if ref.URI() != want {
		t.Errorf("VPCRef URI = %q, want %q", ref.URI(), want)
	}
	ids := parseURIIDs(ref.URI())
	if ids["projects"] != "p-1" || ids["vpcs"] != "vpc-1" {
		t.Errorf("parseURIIDs = %v", ids)
	}
}

func TestSubnetRef(t *testing.T) {
	ref := SubnetRef("p-1", "vpc-1", "sn-1")
	want := "/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/subnets/sn-1"
	if ref.URI() != want {
		t.Errorf("SubnetRef URI = %q, want %q", ref.URI(), want)
	}
	ids := parseURIIDs(ref.URI())
	if ids["projects"] != "p-1" || ids["vpcs"] != "vpc-1" || ids["subnets"] != "sn-1" {
		t.Errorf("parseURIIDs = %v", ids)
	}
}

func TestSecurityGroupRef(t *testing.T) {
	ref := SecurityGroupRef("p-1", "vpc-1", "sg-1")
	want := "/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1"
	if ref.URI() != want {
		t.Errorf("SecurityGroupRef URI = %q, want %q", ref.URI(), want)
	}
	ids := parseURIIDs(ref.URI())
	if ids["projects"] != "p-1" || ids["vpcs"] != "vpc-1" || ids["securitygroups"] != "sg-1" {
		t.Errorf("parseURIIDs = %v", ids)
	}
}

func TestSecurityRuleRef(t *testing.T) {
	ref := SecurityRuleRef("p-1", "vpc-1", "sg-1", "rule-1")
	want := "/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1/securityrules/rule-1"
	if ref.URI() != want {
		t.Errorf("SecurityRuleRef URI = %q, want %q", ref.URI(), want)
	}
	ids := parseURIIDs(ref.URI())
	if ids["projects"] != "p-1" || ids["securitygroups"] != "sg-1" || ids["securityrules"] != "rule-1" {
		t.Errorf("parseURIIDs = %v", ids)
	}
}

func TestElasticIPRef(t *testing.T) {
	ref := ElasticIPRef("p-1", "eip-1")
	want := "/projects/p-1/providers/Aruba.Network/elasticIps/eip-1"
	if ref.URI() != want {
		t.Errorf("ElasticIPRef URI = %q, want %q", ref.URI(), want)
	}
	ids := parseURIIDs(ref.URI())
	if ids["projects"] != "p-1" || ids["elasticIps"] != "eip-1" {
		t.Errorf("parseURIIDs = %v", ids)
	}
}

func TestLoadBalancerRef(t *testing.T) {
	ref := LoadBalancerRef("p-1", "lb-1")
	want := "/projects/p-1/providers/Aruba.Network/loadbalancers/lb-1"
	if ref.URI() != want {
		t.Errorf("LoadBalancerRef URI = %q, want %q", ref.URI(), want)
	}
	ids := parseURIIDs(ref.URI())
	if ids["projects"] != "p-1" || ids["loadbalancers"] != "lb-1" {
		t.Errorf("parseURIIDs = %v", ids)
	}
}

func TestVPCPeeringRef(t *testing.T) {
	ref := VPCPeeringRef("p-1", "vpc-1", "peer-1")
	want := "/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/vpcPeerings/peer-1"
	if ref.URI() != want {
		t.Errorf("VPCPeeringRef URI = %q, want %q", ref.URI(), want)
	}
	ids := parseURIIDs(ref.URI())
	if ids["projects"] != "p-1" || ids["vpcs"] != "vpc-1" || ids["vpcPeerings"] != "peer-1" {
		t.Errorf("parseURIIDs = %v", ids)
	}
}

func TestVPCPeeringRouteRef(t *testing.T) {
	ref := VPCPeeringRouteRef("p-1", "vpc-1", "peer-1", "rt-1")
	want := "/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/vpcPeerings/peer-1/vpcPeeringRoutes/rt-1"
	if ref.URI() != want {
		t.Errorf("VPCPeeringRouteRef URI = %q, want %q", ref.URI(), want)
	}
	ids := parseURIIDs(ref.URI())
	if ids["projects"] != "p-1" || ids["vpcPeerings"] != "peer-1" || ids["vpcPeeringRoutes"] != "rt-1" {
		t.Errorf("parseURIIDs = %v", ids)
	}
}

func TestVPNTunnelRef(t *testing.T) {
	ref := VPNTunnelRef("p-1", "tun-1")
	want := "/projects/p-1/providers/Aruba.Network/vpnTunnels/tun-1"
	if ref.URI() != want {
		t.Errorf("VPNTunnelRef URI = %q, want %q", ref.URI(), want)
	}
	ids := parseURIIDs(ref.URI())
	if ids["projects"] != "p-1" || ids["vpnTunnels"] != "tun-1" {
		t.Errorf("parseURIIDs = %v", ids)
	}
}

func TestVPNRouteRef(t *testing.T) {
	ref := VPNRouteRef("p-1", "tun-1", "rt-1")
	want := "/projects/p-1/providers/Aruba.Network/vpnTunnels/tun-1/vpnRoutes/rt-1"
	if ref.URI() != want {
		t.Errorf("VPNRouteRef URI = %q, want %q", ref.URI(), want)
	}
	ids := parseURIIDs(ref.URI())
	if ids["projects"] != "p-1" || ids["vpnTunnels"] != "tun-1" || ids["vpnRoutes"] != "rt-1" {
		t.Errorf("parseURIIDs = %v", ids)
	}
}

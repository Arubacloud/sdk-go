package aruba

import (
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

// ---- Sub-builder ----

// VPNIPConfig is a fluent builder for the IPConfigurations block of a VPNTunnel.
// Construct with NewVPNIPConfig() and attach via VPNTunnel.WithIPConfig.
type VPNIPConfig struct {
	errMixin
	vpc        *types.ReferenceResource
	publicIP   *types.ReferenceResource
	subnetName string
	subnetCIDR string
	hasSubnet  bool
}

// WithVPC sets the VPC reference for this IP configuration. Errors if v's URI is empty.
func (c *VPNIPConfig) WithVPC(v Ref) *VPNIPConfig {
	if v == nil || v.URI() == "" {
		c.addErr(fmt.Errorf("WithVPC: VPC Ref has empty URI"))
		return c
	}
	c.vpc = &types.ReferenceResource{URI: v.URI()}
	return c
}

// WithElasticIP sets the public elastic IP reference. Errors if v's URI is empty.
func (c *VPNIPConfig) WithElasticIP(v Ref) *VPNIPConfig {
	if v == nil || v.URI() == "" {
		c.addErr(fmt.Errorf("WithElasticIP: PublicIP Ref has empty URI"))
		return c
	}
	c.publicIP = &types.ReferenceResource{URI: v.URI()}
	return c
}

// WithSubnet sets the subnet name and CIDR block for this IP configuration.
func (c *VPNIPConfig) WithSubnet(name, cidr string) *VPNIPConfig {
	c.subnetName, c.subnetCIDR, c.hasSubnet = name, cidr, true
	return c
}

func (c *VPNIPConfig) build() *types.IPConfigurations {
	if c == nil {
		return nil
	}
	out := &types.IPConfigurations{VPC: c.vpc, PublicIP: c.publicIP}
	if c.hasSubnet {
		out.Subnet = &types.SubnetInfo{Name: c.subnetName, CIDR: c.subnetCIDR}
	}
	return out
}

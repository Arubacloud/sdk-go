package aruba

// URI returns an opaque Ref backed by a raw URI string. Use this when you have
// a URI but not a typed wrapper, for example a resource URI loaded from a config
// file or an environment variable.
//
//	vpc, err := client.FromNetwork().VPCs().Get(ctx, aruba.URI("/projects/p/network/vpcs/v"))
func URI(s string) Ref {
	return uriRef{uri: s}
}

// NewProject returns a fresh *Project ready for fluent setters and a Create call.
func NewProject() *Project { return &Project{} }

// NewVPC returns a fresh *VPC ready for fluent setters and a Create call.
// Binds the projectScopedMixin's error sink to the VPC's errMixin so IntoProject
// failures surface via Err().
func NewVPC() *VPC {
	v := &VPC{}
	v.projectScopedMixin = bindProjectScoped(&v.errMixin)
	return v
}

// NewSubnet returns a fresh *Subnet ready for fluent setters and a Create call.
// Binds vpcScopedMixin's error sink so IntoVPC failures surface via Err().
func NewSubnet() *Subnet {
	s := &Subnet{}
	s.vpcScopedMixin = bindVPCScoped(&s.errMixin)
	return s
}

// NewElasticIP returns a fresh *ElasticIP ready for fluent setters and a Create call.
// Binds projectScopedMixin's error sink so IntoProject failures surface via Err().
func NewElasticIP() *ElasticIP {
	e := &ElasticIP{}
	e.projectScopedMixin = bindProjectScoped(&e.errMixin)
	return e
}

// NewSecurityGroup returns a fresh *SecurityGroup ready for fluent setters and a Create call.
// Binds vpcScopedMixin's error sink so IntoVPC failures surface via Err().
func NewSecurityGroup() *SecurityGroup {
	sg := &SecurityGroup{}
	sg.vpcScopedMixin = bindVPCScoped(&sg.errMixin)
	return sg
}

// NewSecurityRule returns a fresh *SecurityRule ready for fluent setters and a Create call.
// Binds securityGroupScopedMixin's error sink so IntoSecurityGroup failures surface via Err().
func NewSecurityRule() *SecurityRule {
	r := &SecurityRule{}
	r.securityGroupScopedMixin = bindSecurityGroupScoped(&r.errMixin)
	return r
}

// NewVPCPeering returns a fresh *VPCPeering ready for fluent setters and a Create call.
// Binds vpcScopedMixin's error sink so IntoVPC failures surface via Err().
func NewVPCPeering() *VPCPeering {
	p := &VPCPeering{}
	p.vpcScopedMixin = bindVPCScoped(&p.errMixin)
	return p
}

// NewVPCPeeringRoute returns a fresh *VPCPeeringRoute ready for fluent setters and a Create call.
// Binds vpcPeeringScopedMixin's error sink so IntoVPCPeering failures surface via Err().
func NewVPCPeeringRoute() *VPCPeeringRoute {
	r := &VPCPeeringRoute{}
	r.vpcPeeringScopedMixin = bindVPCPeeringScoped(&r.errMixin)
	return r
}

// NewVPNTunnel returns a fresh *VPNTunnel ready for fluent setters and a Create call.
// Binds projectScopedMixin's error sink so IntoProject failures surface via Err().
func NewVPNTunnel() *VPNTunnel {
	t := &VPNTunnel{}
	t.projectScopedMixin = bindProjectScoped(&t.errMixin)
	return t
}

// NewVPNRoute returns a fresh *VPNRoute ready for fluent setters and a Create call.
// Binds vpnTunnelScopedMixin's error sink so IntoVPNTunnel failures surface via Err().
func NewVPNRoute() *VPNRoute {
	r := &VPNRoute{}
	r.vpnTunnelScopedMixin = bindVPNTunnelScoped(&r.errMixin)
	return r
}

// NewBlockStorage returns a fresh *BlockStorage ready for fluent setters and a Create call.
// Binds projectScopedMixin's error sink so IntoProject failures surface via Err().
func NewBlockStorage() *BlockStorage {
	b := &BlockStorage{}
	b.projectScopedMixin = bindProjectScoped(&b.errMixin)
	return b
}

// NewSnapshot returns a fresh *Snapshot ready for fluent setters and a Create call.
// Binds projectScopedMixin's error sink so IntoProject failures surface via Err().
func NewSnapshot() *Snapshot {
	s := &Snapshot{}
	s.projectScopedMixin = bindProjectScoped(&s.errMixin)
	return s
}

// NewStorageBackup returns a fresh *StorageBackup ready for fluent setters and a Create call.
// Binds projectScopedMixin's error sink so IntoProject failures surface via Err().
func NewStorageBackup() *StorageBackup {
	b := &StorageBackup{}
	b.projectScopedMixin = bindProjectScoped(&b.errMixin)
	return b
}

// NewVPNIPConfig returns a fresh *VPNIPConfig sub-builder for configuring IP settings.
func NewVPNIPConfig() *VPNIPConfig { return &VPNIPConfig{} }

// NewVPNIKE returns a fresh *VPNIKE sub-builder for configuring IKE settings.
func NewVPNIKE() *VPNIKE { return &VPNIKE{} }

// NewVPNESP returns a fresh *VPNESP sub-builder for configuring ESP settings.
func NewVPNESP() *VPNESP { return &VPNESP{} }

// NewVPNPSK returns a fresh *VPNPSK sub-builder for configuring PSK settings.
func NewVPNPSK() *VPNPSK { return &VPNPSK{} }

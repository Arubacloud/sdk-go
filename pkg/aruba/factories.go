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

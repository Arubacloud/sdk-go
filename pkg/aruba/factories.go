package aruba

// URI returns an opaque Ref backed by a raw URI string. Use this when you have
// a URI but not a typed wrapper, for example a resource URI loaded from a config
// file or an environment variable.
//
//	vpc, err := client.FromNetwork().VPCs().Get(ctx, aruba.URI("/projects/p/network/vpcs/v"))
func URI(s string) Ref {
	return uriRef{uri: s}
}

package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// Subnet is the wrapper for an Aruba Cloud subnet (a child of a VPC).
// Construct with aruba.NewSubnet() and bind it to a parent VPC via IntoVPC(vpc).
type Subnet struct {
	errMixin
	metadataMixin
	regionalMixin
	vpcScopedMixin // direct parent; populates vpcID + projectID
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	subnetType    *string               // Properties.Type ("Basic" / "Advanced")
	defaultSubnet *bool                 // Properties.Default
	cidr          *string               // Properties.Network.Address
	dhcp          *SubnetDHCP           // Properties.DHCP (sub-builder)
	response      *types.SubnetResponse // backs Raw()
}

func (s *Subnet) IntoVPC(v Ref) *Subnet            { s.intoVPC(v); return s }
func (s *Subnet) WithName(n string) *Subnet        { s.withName(n); return s }
func (s *Subnet) AddTag(t string) *Subnet          { s.addTag(t); return s }
func (s *Subnet) RemoveTag(t string) *Subnet       { s.removeTag(t); return s }
func (s *Subnet) ReplaceTags(ts ...string) *Subnet { s.replaceTags(ts...); return s }
func (s *Subnet) WithLocation(loc string) *Subnet  { s.withLocation(loc); return s }
func (s *Subnet) InRegion(region string) *Subnet   { s.withLocation(region); return s }
func (s *Subnet) WithType(t string) *Subnet        { s.subnetType = &t; return s }
func (s *Subnet) WithDefault(b bool) *Subnet       { s.defaultSubnet = &b; return s }
func (s *Subnet) WithCIDR(cidr string) *Subnet     { s.cidr = &cidr; return s }
func (s *Subnet) WithDHCP(d *SubnetDHCP) *Subnet   { s.dhcp = d; return s }

// URI satisfies Ref.
func (s *Subnet) URI() string { return s.RespURI() }

// SubnetID satisfies withSubnetID so future grandchildren can extract this ID.
func (s *Subnet) SubnetID() string { return s.ID() }

// Raw shadows responseMetadataMixin.Raw() with the full subnet response.
func (s *Subnet) Raw() *types.SubnetResponse { return s.response }

// RawRequest returns what toRequest() would emit right now.
func (s *Subnet) RawRequest() types.SubnetRequest { return s.toRequest() }

func (s *Subnet) Type() string {
	if s.subnetType == nil {
		return ""
	}
	return *s.subnetType
}
func (s *Subnet) IsDefault() bool {
	if s.defaultSubnet == nil {
		return false
	}
	return *s.defaultSubnet
}
func (s *Subnet) CIDR() string {
	if s.cidr == nil {
		return ""
	}
	return *s.cidr
}
func (s *Subnet) DHCP() *SubnetDHCP { return s.dhcp }

func (s *Subnet) toRequest() types.SubnetRequest {
	props := types.SubnetPropertiesRequest{}
	if s.subnetType != nil {
		props.Type = types.SubnetType(*s.subnetType)
	}
	if s.defaultSubnet != nil {
		props.Default = *s.defaultSubnet
	}
	if s.cidr != nil {
		props.Network = &types.SubnetNetwork{Address: *s.cidr}
	}
	if s.dhcp != nil {
		props.DHCP = s.dhcp.toType()
	}
	return types.SubnetRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: s.toMetadata(),
			Location:                s.toLocation(),
		},
		Properties: props,
	}
}

func (s *Subnet) fromResponse(resp *types.SubnetResponse) {
	if resp == nil {
		return
	}
	s.response = resp
	s.setMeta(&resp.Metadata)
	s.withName(subnetDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		s.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		s.withLocation(resp.Metadata.LocationResponse.Value)
	}
	s.setStatus(&resp.Status)
	s.setLinked(resp.Properties.LinkedResources)

	if resp.Properties.Type != "" {
		t := string(resp.Properties.Type)
		s.subnetType = &t
	}
	d := resp.Properties.Default
	s.defaultSubnet = &d
	if resp.Properties.Network != nil && resp.Properties.Network.Address != "" {
		addr := resp.Properties.Network.Address
		s.cidr = &addr
	}
	if resp.Properties.DHCP != nil {
		s.dhcp = dhcpFromType(resp.Properties.DHCP)
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		s.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	// Backfill ancestor IDs from response URI if not already set.
	if (s.vpcID == "" || s.projectID == "") && s.RespURI() != "" {
		ids := parseURIIDs(s.RespURI())
		if s.vpcID == "" {
			s.vpcID = ids["vpcs"]
		}
		if s.projectID == "" {
			s.projectID = ids["projects"]
		}
	}
}

func subnetDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

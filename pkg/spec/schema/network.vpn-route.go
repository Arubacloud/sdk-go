package schema

type VpnRoutePropertiesRequest struct {

	// CloudSubnet CIDR of the cloud subnet
	CloudSubnet string `json:"cloudSubnet"`

	// OnPremSubnet CIDR of the onPrem subnet
	OnPremSubnet string `json:"onPremSubnet"`
}

type VpnRoutePropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	// CloudSubnet CIDR of the cloud subnet
	CloudSubnet string `json:"cloudSubnet"`

	// OnPremSubnet CIDR of the onPrem subnet
	OnPremSubnet string `json:"onPremSubnet"`

	VpnTunnel *ReferenceResource `json:"vpnTunnel,omitempty"`
}

type VpnRouteRequest struct {
	// Metadata of the VPC Route
	Metadata ResourceMetadataRequest `json:"metadata"`

	// Spec contains the VPC Route specification
	Properties VpnRoutePropertiesRequest `json:"properties"`
}

type VpnRouteResponse struct {
	// Metadata of the VPC Route
	Metadata ResourceMetadataResponse `json:"metadata"`
	// Spec contains the VPC Route specification
	Properties VpnRoutePropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type VpnRouteList struct {
	ListResponse
	Values []VpnRouteResponse `json:"values"`
}

package schema

type VPNRoutePropertiesRequest struct {

	// CloudSubnet CIDR of the cloud subnet
	CloudSubnet string `json:"cloudSubnet"`

	// OnPremSubnet CIDR of the onPrem subnet
	OnPremSubnet string `json:"onPremSubnet"`
}

type VPNRoutePropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	// CloudSubnet CIDR of the cloud subnet
	CloudSubnet string `json:"cloudSubnet"`

	// OnPremSubnet CIDR of the onPrem subnet
	OnPremSubnet string `json:"onPremSubnet"`

	VPNTunnel *ReferenceResource `json:"vpnTunnel,omitempty"`
}

type VPNRouteRequest struct {
	// Metadata of the VPC Route
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	// Spec contains the VPC Route specification
	Properties VPNRoutePropertiesRequest `json:"properties"`
}

type VPNRouteResponse struct {
	// Metadata of the VPC Route
	Metadata ResourceMetadataResponse `json:"metadata"`
	// Spec contains the VPC Route specification
	Properties VPNRoutePropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type VPNRouteList struct {
	ListResponse
	Values []VPNRouteResponse `json:"values"`
}

package schema

// VpcPeeringRoutePropertiesRequest contains properties of a VPC peering route to create
type VpcPeeringRoutePropertiesRequest struct {
	// LocalNetworkAddress Local network address in CIDR notation
	LocalNetworkAddress string `json:"localNetworkAddress"`

	// RemoteNetworkAddress Remote network address in CIDR notation
	RemoteNetworkAddress string `json:"remoteNetworkAddress"`

	BillingPlan BillingPeriodResource `json:"billingPlan"`
}

type VpcPeeringRoutePropertiesResponse struct {
	// LocalNetworkAddress Local network address in CIDR notation
	LocalNetworkAddress string `json:"localNetworkAddress"`

	// RemoteNetworkAddress Remote network address in CIDR notation
	RemoteNetworkAddress string `json:"remoteNetworkAddress"`

	BillingPlan BillingPeriodResource `json:"billingPlan"`
}

type VpcPeeringRouteRequest struct {
	// Metadata of the VPC Peering Route
	Metadata ResourceMetadataRequest `json:"metadata"`

	// Spec contains the VPC Peering Route specification
	Properties VpcPeeringRoutePropertiesRequest `json:"properties"`
}

type VpcPeeringRouteResponse struct {
	// Metadata of the VPC Peering Route
	Metadata RegionalResourceMetadataRequest `json:"metadata"`
	// Spec contains the VPC Peering Route specification
	Properties VpcPeeringRoutePropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type VpcPeeringRouteList struct {
	ListResponse
	Values []VpcPeeringRouteResponse `json:"values"`
}

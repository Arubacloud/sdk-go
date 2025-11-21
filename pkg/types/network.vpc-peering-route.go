package types

// VPCPeeringRoutePropertiesRequest contains properties of a VPC peering route to create
type VPCPeeringRoutePropertiesRequest struct {
	// LocalNetworkAddress Local network address in CIDR notation
	LocalNetworkAddress string `json:"localNetworkAddress"`

	// RemoteNetworkAddress Remote network address in CIDR notation
	RemoteNetworkAddress string `json:"remoteNetworkAddress"`

	BillingPlan BillingPeriodResource `json:"billingPlan"`
}

type VPCPeeringRoutePropertiesResponse struct {
	// LocalNetworkAddress Local network address in CIDR notation
	LocalNetworkAddress string `json:"localNetworkAddress"`

	// RemoteNetworkAddress Remote network address in CIDR notation
	RemoteNetworkAddress string `json:"remoteNetworkAddress"`

	BillingPlan BillingPeriodResource `json:"billingPlan"`
}

type VPCPeeringRouteRequest struct {
	// Metadata of the VPC Peering Route
	Metadata ResourceMetadataRequest `json:"metadata"`

	// Spec contains the VPC Peering Route specification
	Properties VPCPeeringRoutePropertiesRequest `json:"properties"`
}

type VPCPeeringRouteResponse struct {
	// Metadata of the VPC Peering Route
	Metadata RegionalResourceMetadataRequest `json:"metadata"`
	// Spec contains the VPC Peering Route specification
	Properties VPCPeeringRoutePropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type VPCPeeringRouteList struct {
	ListResponse
	Values []VPCPeeringRouteResponse `json:"values"`
}

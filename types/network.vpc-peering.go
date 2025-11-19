package types

// VPCPeeringRoutePropertiesRequest contains properties of a VPC peering route to create
type VPCPeeringPropertiesRequest struct {
	RemoteVPC *ReferenceResource `json:"remoteVpc,omitempty"`
}

type VPCPeeringPropertiesResponse struct {
	LinkedResources []LinkedResource   `json:"linkedResources,omitempty"`
	RemoteVPC       *ReferenceResource `json:"remoteVpc,omitempty"`
}

type VPCPeeringRequest struct {
	// Metadata of the VPC Peering
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	// Spec contains the VPC Peering specification
	Properties VPCPeeringPropertiesRequest `json:"properties"`
}

type VPCPeeringResponse struct {
	// Metadata of the VPC Peering
	Metadata ResourceMetadataResponse `json:"metadata"`
	// Spec contains the VPC Peering specification
	Properties VPCPeeringPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type VPCPeeringList struct {
	ListResponse
	Values []VPCPeeringResponse `json:"values"`
}

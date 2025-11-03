package schema

// VpcPeeringRoutePropertiesRequest contains properties of a VPC peering route to create
type VpcPeeringPropertiesRequest struct {
	RemoteVpc *ReferenceResource `json:"remoteVpc,omitempty"`
}

type VpcPeeringPropertiesResponse struct {
	LinkedResources []LinkedResource   `json:"linkedResources,omitempty"`
	RemoteVpc       *ReferenceResource `json:"remoteVpc,omitempty"`
}

type VpcPeeringRequest struct {
	// Metadata of the VPC Peering
	Metadata ResourceMetadataRequest `json:"metadata"`

	// Spec contains the VPC Peering specification
	Properties VpcPeeringPropertiesRequest `json:"properties"`
}

type VpcPeeringResponse struct {
	// Metadata of the VPC Peering
	Metadata ResourceMetadataResponse `json:"metadata"`
	// Spec contains the VPC Peering specification
	Properties VpcPeeringPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type VpcPeeringList struct {
	ListResponse
	Values []VpcPeeringResponse `json:"values"`
}

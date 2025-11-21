package types

type LoadBalancerPropertiesResponse struct {
	// LinkedResources array of resources linked to the Load Balancer (nullable)
	LinkedResources []LinkedResource   `json:"linkedResources,omitempty"`
	Address         *string            `json:"address,omitempty"`
	VPC             *ReferenceResource `json:"vpc,omitempty"`
}

type LoadBalancerResponse struct {
	// Metadata of the Load Balancer
	Metadata ResourceMetadataResponse `json:"metadata"`
	// Spec contains the Load Balancer specification
	Properties LoadBalancerPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type LoadBalancerList struct {
	ListResponse
	Values []LoadBalancerResponse `json:"values"`
}

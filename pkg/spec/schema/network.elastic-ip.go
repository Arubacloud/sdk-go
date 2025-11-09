package schema

type ElasticIpPropertiesRequest struct {
	BillingPlan BillingPeriodResource `json:"billingPlan"`
}

type ElasticIpPropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Address     *string               `json:"address,omitempty"`
	BillingPlan BillingPeriodResource `json:"billingPlan"`
}

type ElasticIpRequest struct {
	Metadata   RegionalResourceMetadataRequest `json:"metadata"`
	Properties ElasticIpPropertiesRequest      `json:"properties"`
}

type ElasticIpResponse struct {
	Metadata   ResourceMetadataResponse    `json:"metadata"`
	Properties ElasticIpPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type ElasticList struct {
	ListResponse
	Values []ElasticIpResponse `json:"values"`
}

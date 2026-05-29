package types

type ElasticIPPropertiesRequest struct {
	BillingPlan *BillingPlan `json:"billingPlan,omitempty"`
}

type ElasticIPPropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Address     *string      `json:"address,omitempty"`
	BillingPlan *BillingPlan `json:"billingPlan,omitempty"`
}

type ElasticIPRequest struct {
	Metadata   RegionalResourceMetadataRequest `json:"metadata"`
	Properties ElasticIPPropertiesRequest      `json:"properties"`
}

type ElasticIPResponse struct {
	Metadata   ResourceMetadataResponse    `json:"metadata"`
	Properties ElasticIPPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type ElasticIPListResponse struct {
	ListResponse
	Values []ElasticIPResponse `json:"values"`
}

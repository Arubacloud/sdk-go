package schema

type ElasticIPPropertiesRequest struct {
	BillingPlan BillingPeriodResource `json:"billingPlan"`
}

type ElasticIPPropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Address     *string               `json:"address,omitempty"`
	BillingPlan BillingPeriodResource `json:"billingPlan"`
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

type ElasticList struct {
	ListResponse
	Values []ElasticIPResponse `json:"values"`
}

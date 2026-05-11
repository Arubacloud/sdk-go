package types

type ElasticIPPropertiesRequest struct {
	BillingPeriod *BillingPeriod `json:"billingPeriod,omitempty"`
}

type ElasticIPPropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Address       *string        `json:"address,omitempty"`
	BillingPeriod *BillingPeriod `json:"billingPeriod,omitempty"`
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

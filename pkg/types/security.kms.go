package types

type KmsPropertiesRequest struct {
	BillingPeriod BillingPeriodResource `json:"billingPeriod"`
}

type KmsPropertiesResponse struct {
	BillingPeriod BillingPeriodResource `json:"billingPeriod"`
}
type KmsRequest struct {
	Metadata   RegionalResourceMetadataRequest `json:"metadata"`
	Properties KmsPropertiesRequest            `json:"properties"`
}

type KmsResponse struct {
	Metadata   ResourceMetadataResponse `json:"metadata"`
	Properties KmsPropertiesResponse    `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type KmsList struct {
	ListResponse
	Values []KmsResponse `json:"values"`
}

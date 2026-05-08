package types

type KeyPairPropertiesRequest struct {
	Value string `json:"value"`
}

type KeyPairPropertiesResult struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Value string `json:"value"`
}

type KeyPairRequest struct {
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	Properties KeyPairPropertiesRequest `json:"properties"`
}

type KeyPairResponse struct {
	Metadata   ResourceMetadataResponse `json:"metadata"`
	Properties KeyPairPropertiesResult  `json:"properties"`
	Status     ResourceStatus           `json:"status"`
}

type KeyPairListResponse struct {
	ListResponse
	Values []KeyPairResponse `json:"values"`
}

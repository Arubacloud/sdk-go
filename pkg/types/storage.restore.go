package types

type RestorePropertiesRequest struct {
	Target ReferenceResource `json:"destinationVolume"`
}

type RestorePropertiesResult struct {
	Destination ReferenceResource `json:"destinationVolume"`
}

type RestoreRequest struct {
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	Properties RestorePropertiesRequest `json:"properties"`
}

type RestoreResponse struct {
	Metadata ResourceMetadataResponse `json:"metadata"`

	Properties RestorePropertiesResult `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type RestoreList struct {
	ListResponse
	Values []RestoreResponse `json:"values"`
}

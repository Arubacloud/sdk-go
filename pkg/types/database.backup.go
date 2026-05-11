package types

type BackupPropertiesRequest struct {
	Zone Zone `json:"datacenter"`

	DBaaS ReferenceResource `json:"dbaas"`

	Database ReferenceResource `json:"database"`

	BillingPeriod *BillingPeriod `json:"billingPeriod,omitempty"`
}

type BackupStorageResponse struct {
	Size int32 `json:"size"`
}

type BackupPropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Zone Zone `json:"datacenter"`

	DBaaS ReferenceResource `json:"dbaas"`

	Database ReferenceResource `json:"database"`

	BillingPeriod *BillingPeriod `json:"billingPeriod,omitempty"`

	Storage BackupStorageResponse `json:"storage"`
}

type BackupRequest struct {
	Metadata   RegionalResourceMetadataRequest `json:"metadata"`
	Properties BackupPropertiesRequest         `json:"properties"`
}

type BackupResponse struct {
	Metadata   ResourceMetadataResponse `json:"metadata"`
	Properties BackupPropertiesResponse `json:"properties"`
	Status     ResourceStatus           `json:"status,omitempty"`
}

type BackupList struct {
	ListResponse
	Values []BackupResponse `json:"values"`
}

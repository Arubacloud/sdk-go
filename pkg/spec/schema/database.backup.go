package schema

type BackupPropertiesRequest struct {
	Zone string `json:"datacenter"`

	DBaaS ReferenceResource `json:"dbaas"`

	Database ReferenceResource `json:"database"`

	BillingPlan BillingPeriodResource `json:"billingPlan"`
}

type BackupStorageResponse struct {
	Size int32 `json:"size"`
}

type BackupPropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Zone string `json:"datacenter"`

	DBaaS ReferenceResource `json:"dbaas"`

	Database ReferenceResource `json:"database"`

	BillingPlan BillingPeriodResource `json:"billingPlan"`

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

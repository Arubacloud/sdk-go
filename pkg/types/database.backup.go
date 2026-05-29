package types

type BackupPropertiesRequest struct {
	Zone Zone `json:"datacenter"`

	DBaaS ReferenceResource `json:"dbaas"`

	Database ReferenceResource `json:"database"`

	BillingPlan *BillingPlan `json:"billingPlan,omitempty"`
}

type BackupStorageResponse struct {
	Size int32 `json:"size"`
}

type BackupPropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Zone Zone `json:"datacenter"`

	DBaaS ReferenceResource `json:"dbaas"`

	Database ReferenceResource `json:"database"`

	BillingPlan *BillingPlan `json:"billingPlan,omitempty"`

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

type DBaaSBackupListResponse struct {
	ListResponse
	Values []BackupResponse `json:"values"`
}

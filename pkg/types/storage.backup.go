package types

type StorageBackupType string

const (
	StorageBackupTypeFull        StorageBackupType = "Full"
	StorageBackupTypeIncremental StorageBackupType = "Incremental"
)

type StorageBackupPropertiesRequest struct {

	// StorageBackupType indicates whether the StorageBackup is full or incremental
	StorageBackupType StorageBackupType `json:"type"`

	// Origin indicates the source volume
	Origin ReferenceResource `json:"sourceVolume"`

	// RetentionDays indicates the number of days to retain the backup
	RetentionDays *int `json:"retentionDays,omitempty"`

	// BillingPeriod indicates the billing period
	BillingPeriod *string `json:"billingPeriod,omitempty"`
}

type StorageBackupPropertiesResult struct {

	// StorageBackupType indicates whether the StorageBackup is full or incremental
	Type StorageBackupType `json:"type"`

	// Origin indicates the source volume
	Origin ReferenceResource `json:"sourceVolume"`

	// RetentionDays indicates the number of days to retain the backup
	RetentionDays *int `json:"retentionDays,omitempty"`

	// BillingPeriod indicates the billing period
	BillingPeriod *string `json:"billingPeriod,omitempty"`
}

type StorageBackupRequest struct {
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	Properties StorageBackupPropertiesRequest `json:"properties"`
}

type StorageBackupResponse struct {
	Metadata ResourceMetadataResponse `json:"metadata"`

	Properties StorageBackupPropertiesResult `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type StorageBackupList struct {
	ListResponse
	Values []StorageBackupResponse `json:"values"`
}

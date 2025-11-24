package types

// BlockStorageType represents the type of block storage
type BlockStorageType string

const (
	BlockStorageTypeStandard    BlockStorageType = "Standard"
	BlockStorageTypePerformance BlockStorageType = "Performance"
)

type BlockStoragePropertiesRequest struct {

	// SizeGB Size of the block storage in GB
	SizeGB int `json:"sizeGb"`

	// BillingPeriod of the block storage
	BillingPeriod string `json:"billingPeriod"`

	//Zone where blockstorage will be created
	Zone string `json:"dataCenter"`

	// Type of block storage. Admissible values: Standard, Performance
	Type BlockStorageType `json:"type"`

	Snapshot *ReferenceResource `json:"snapshot,omitempty"`

	Bootable *bool `json:"bootable,omitempty"`

	Image *string `json:"image,omitempty"`
}

type BlockStoragePropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	// SizeGB Size of the block storage in GB
	SizeGB int `json:"sizeGb"`

	// BillingPeriod Billing plan of the block storage
	BillingPeriod string `json:"billingPeriod"`

	//Zone where blockstorage will be created
	Zone string `json:"dataCenter"`

	// Type of block storage. Admissible values: Standard, Performance
	Type BlockStorageType `json:"type"`

	Snapshot *ReferenceResource `json:"snapshot,omitempty"`

	Bootable *bool `json:"bootable,omitempty"`

	Image *string `json:"image,omitempty"`
}

type BlockStorageRequest struct {
	// Metadata of the Block Storage
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	// Spec contains the Block Storage specification
	Properties BlockStoragePropertiesRequest `json:"properties"`
}

type BlockStorageResponse struct {

	// Metadata of the Block Storage
	Metadata ResourceMetadataResponse `json:"metadata"`

	// Spec contains the Block Storage specification
	Properties BlockStoragePropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type BlockStorageList struct {
	ListResponse
	Values []BlockStorageResponse `json:"values"`
}

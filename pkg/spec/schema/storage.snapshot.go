package schema

type SnapshotPropertiesRequest struct {
	// BillingPeriod The billing period for blockStorage. Only Hour is a valid value (nullable)
	BillingPeriod *string `json:"billingPeriod,omitempty"`

	Volume ReferenceResource `json:"volume,omitempty"`
}

// VolumeInfo contains information about the original volume
type VolumeInfo struct {
	// URI of the volume
	URI *string `json:"uri,omitempty"`

	// Type of the original volume from which the snapshot was created (nullable)
	Name *string `json:"name,omitempty"`

	CompoundResource *ReferenceResource `json:"compoundResource,omitempty"`
}

type SnapshotPropertiesResponse struct {
	// SizeGb The blockStorage's size in gigabyte (nullable)
	SizeGb *int32 `json:"sizeGb,omitempty"`

	// BillingPeriod The billing period for blockStorage. Only Hour is a valid value (nullable)
	BillingPeriod *string `json:"billingPeriod,omitempty"`

	// Volume information about the original volume
	Volume *VolumeInfo `json:"volume,omitempty"`

	// Type of block storage. Admissible values: Standard, Performance
	Type BlockStorageType `json:"type"`

	//Zone where blockstorage will be created
	Zone string `json:"dataCenter"`

	Bootable *bool `json:"bootable,omitempty"`
}

type SnapshotRequest struct {
	// Metadata of the Snapshot
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	// Spec contains the Snapshot specification
	Properties SnapshotPropertiesRequest `json:"properties"`
}

type SnapshotResponse struct {
	// Metadata of the Snapshot
	Metadata ResourceMetadataResponse `json:"metadata"`

	// Spec contains the Snapshot specification
	Properties SnapshotPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type SnapshotList struct {
	ListResponse
	Values []SnapshotResponse `json:"values"`
}

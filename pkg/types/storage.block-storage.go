package types

// VolumeImage identifies a stock OS template (and any bundled software)
// used to provision a bootable BlockStorage volume.
//
// The constants below are derived from Aruba's published catalog at
// https://kb.arubacloud.com/en/computing/template-datasheets/template.aspx
// (plus the OpenClaw and Proxmox VE 8 datasheets on separate pages).
// Aruba may add or retire templates without bumping the API; treat this
// set as a snapshot. Personal/custom templates use a different identifier
// scheme and are not enumerated here.

const (
	// VolumeImageWS22001 — Windows Server 2022 64-bit.
	VolumeImageWS22001 string = "WS22-001"
	// VolumeImageWS19001 — Windows Server 2019 64-bit.
	VolumeImageWS19001 string = "WS19-001"
	// VolumeImageWS16001 — Windows Server 2016 64-bit.
	VolumeImageWS16001 string = "WS16-001"
	// VolumeImageLU24001 — Ubuntu Server 24.04.
	VolumeImageLU24001 string = "LU24-001"
	// VolumeImageLU22001 — Ubuntu Server 22.04 LTS 64-bit.
	VolumeImageLU22001 string = "LU22-001"
	// VolumeImageLU20001 — Ubuntu Server 20.04 LTS 64-bit.
	VolumeImageLU20001 string = "LU20-001"
	// VolumeImageDE12001 — Debian 12.
	VolumeImageDE12001 string = "DE12-001"
	// VolumeImageDE11001 — Debian 11 64-bit.
	VolumeImageDE11001 string = "DE11-001"
	// VolumeImageRO09001 — Rocky Linux 9.
	VolumeImageRO09001 string = "RO09-001"
	// VolumeImageLC09001 — CentOS 9.
	VolumeImageLC09001 string = "LC09-001"
	// VolumeImageAL90001 — AlmaLinux 9.x 64-bit.
	VolumeImageAL90001 string = "AL90-001"
	// VolumeImageAL85001 — AlmaLinux 8.x 64-bit.
	VolumeImageAL85001 string = "AL85-001"
	// VolumeImageLO15001 — openSUSE 15.2 64-bit.
	VolumeImageLO15001 string = "LO15-001"
	// VolumeImageBS13001 — FreeBSD 13 64-bit.
	VolumeImageBS13001 string = "BS13-001"

	// VolumeImageAP85001 — AlmaLinux 8.x 64-bit with Plesk pre-installed.
	VolumeImageAP85001 string = "AP85-001"
	// VolumeImageSQLW22 — Windows Server 2022 64-bit with SQL Server 2022 Web pre-installed.
	VolumeImageSQLW22 string = "SQL-W22"
	// VolumeImageSQLS22 — Windows Server 2022 64-bit with SQL Server 2022 Standard pre-installed.
	VolumeImageSQLS22 string = "SQL-S22"
	// VolumeImageADW132 — Windows Server 2022 64-bit with SQL Server 2019 Web pre-installed.
	VolumeImageADW132 string = "ADW-132"
	// VolumeImageADW131 — Windows Server 2022 64-bit with SQL Server 2019 Standard pre-installed.
	VolumeImageADW131 string = "ADW-131"
	// VolumeImageADW122 — Windows Server 2019 64-bit with SQL Server 2016 Standard pre-installed.
	VolumeImageADW122 string = "ADW-122"
	// VolumeImageADW121 — Windows Server 2019 64-bit with SQL Server 2016 Web pre-installed.
	VolumeImageADW121 string = "ADW-121"
	// VolumeImageARW004 — Windows Server 2022 64-bit with RDS (5/10/15/30 CAL) pre-installed.
	VolumeImageARW004 string = "ARW-004"
	// VolumeImageARW003 — Windows Server 2016 64-bit with RDS (5/10/15/30 CAL) pre-installed.
	VolumeImageARW003 string = "ARW-003"
	// VolumeImageAFE001 — pfSense 2.4.5 64-bit firewall/load balancer appliance.
	VolumeImageAFE001 string = "AFE-001"
	// VolumeImageAFL001 — Endian Firewall Community 3.3.2 firewall/load balancer appliance.
	VolumeImageAFL001 string = "AFL-001"
	// VolumeImageLU20MAI01 — Ubuntu Server 20.04 LTS 64-bit with Mail Server pre-installed.
	VolumeImageLU20MAI01 string = "LU20-MAI01"
	// VolumeImageLU24MN01 — Ubuntu Server 24.04 with MinIO pre-installed.
	VolumeImageLU24MN01 string = "LU24-MN01"
	// VolumeImageLU22MN01 — Ubuntu Server 22.04 LTS 64-bit with MinIO pre-installed.
	VolumeImageLU22MN01 string = "LU22-MN01"
	// VolumeImageLU24VD01 — Ubuntu Server 24.04 Virtual Desktop.
	VolumeImageLU24VD01 string = "LU24-VD01"
	// VolumeImageLU22VD01 — Ubuntu Server 22.04 LTS Virtual Desktop.
	VolumeImageLU22VD01 string = "LU22-VD01"
	// VolumeImageAVL005 — Ubuntu Server 20.04 LTS 64-bit Virtual Desktop.
	VolumeImageAVL005 string = "AVL-005"

	// VolumeImageLU24OC01 — Ubuntu Server 24.04 with OpenClaw AI assistant
	// (NGINX proxy, Fail2ban, UFW firewall, Certbot) pre-installed.
	VolumeImageLU24OC01 string = "LU24-OC01"
	// VolumeImagePX08001 — Debian 12 (Proxmox kernel) with Proxmox VE 8
	// and Fail2ban pre-installed.
	VolumeImagePX08001 string = "PX08-001"
)

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
	BillingPeriod *BillingPeriod `json:"billingPeriod,omitempty"`

	// Zone where blockstorage will be created (optional).
	// If specified, the resource is zonal; otherwise, it is regional.
	Zone *Zone `json:"dataCenter,omitempty"`

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
	BillingPeriod *BillingPeriod `json:"billingPeriod,omitempty"`

	//Zone where blockstorage will be created
	Zone Zone `json:"dataCenter"`

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

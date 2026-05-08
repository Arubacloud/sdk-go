package types

// KubernetesVersion identifies a KaaS Kubernetes version.
//
// The constants below cover versions referenced in SDK fixtures and the KB
// upgrade-versions page. The authoritative list is available via:
//
//	GET /providers/Aruba.Container/kubernetesVersions
type KubernetesVersion string

const (
	KubernetesVersion1282 KubernetesVersion = "1.28.2"
	KubernetesVersion1292 KubernetesVersion = "1.29.2"
	KubernetesVersion1302 KubernetesVersion = "1.30.2"
	KubernetesVersion1332 KubernetesVersion = "1.33.2"
)

// NodePoolInstance identifies a KaaS node pool instance type.
//
// Pattern: K<vCPU>A<RAM>. The constant below covers the SKU referenced in
// the cmd/example reference app. The authoritative list is available via:
//
//	GET /providers/Aruba.Container/instances
type NodePoolInstance string

const (
	NodePoolInstanceK1A2  NodePoolInstance = "K1A2"
	NodePoolInstanceK1A4R NodePoolInstance = "K1A4R"
	NodePoolInstanceK2A4  NodePoolInstance = "K2A4"
	NodePoolInstanceK2A8R NodePoolInstance = "K2A8R"
	NodePoolInstanceK4A8  NodePoolInstance = "K4A8"
	NodePoolInstanceK4A16R NodePoolInstance = "K4A16R"
	NodePoolInstanceK8A16  NodePoolInstance = "K8A16"
	NodePoolInstanceK8A32R NodePoolInstance = "K8A32R"
	NodePoolInstanceK12A24 NodePoolInstance = "K12A24"
	NodePoolInstanceK16A32 NodePoolInstance = "K16A32"
	NodePoolInstanceK24A48 NodePoolInstance = "K24A48"
	NodePoolInstanceK32A64 NodePoolInstance = "K32A64"
)

type NodeCIDRProperties struct {

	// Address in CIDR notation The IP range must be between 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
	Address string `json:"address"`

	// Name of the nodecidr
	Name string `json:"name"`
}

type KubernetesVersionInfo struct {
	Value KubernetesVersion `json:"value"`
}

type KubernetesVersionInfoUpdate struct {
	Value       KubernetesVersion `json:"value"`
	UpgradeDate *string           `json:"upgradeDate,omitempty"`
}

type StorageKubernetes struct {
	MaxCumulativeVolumeSize *int32 `json:"maxCumulativeVolumeSize,omitempty"`
}

type NodePoolProperties struct {

	// Name Nodepool name
	Name string `json:"name"`

	// Nodes Number of nodes
	Nodes int32 `json:"nodes"`

	// Instance Configuration name of the nodes.
	// See metadata section of the API documentation for an updated list of admissible values.
	// For more information, check the documentation.
	Instance NodePoolInstance `json:"instance"`

	// DataCenter Datacenter in which the nodes of the pool will be located.
	// See metadata section of the API documentation for an updated list of admissible values.
	// For more information, check the documentation.
	Zone Zone `json:"dataCenter"`

	// MinCount Minimum number of nodes for autoscaling
	MinCount *int32 `json:"minCount,omitempty"`

	// MaxCount Maximum number of nodes for autoscaling
	MaxCount *int32 `json:"maxCount,omitempty"`

	// Autoscaling Indicates if autoscaling is enabled for this node pool
	Autoscaling bool `json:"autoscaling"`
}

type SecurityGroupProperties struct {
	Name string `json:"name"`
}

type IdentityProperties struct {
	ClientID     *string `json:"clientId,omitempty"`
	ClientSecret *string `json:"clientSecret,omitempty"`
}

type IdentityPropertiesResponse struct {
	ClientID *string `json:"clientId,omitempty"`
}

type APIServerAccessProfileProperties struct {
	AuthorizedIPRanges   *[]string `json:"authorizedIpRanges,omitempty"`
	EnablePrivateCluster bool      `json:"enablePrivateCluster"`
}

type APIServerAccessProfilePropertiesResponse struct {
	AuthorizedIPRanges   *[]string `json:"authorizedIpRanges,omitempty"`
	EnablePrivateCluster bool      `json:"enablePrivateCluster"`
}

type ReferenceResourceResponse struct {
	URI *string `json:"uri,omitempty"`
}

type OpenstackProjectResponse struct {
	ID *string `json:"id,omitempty"`
}

type BillingPeriodResourceResponse struct {
	BillingPeriod *BillingPeriod `json:"billingPeriod,omitempty"`
}

type KaaSPropertiesRequest struct {

	//LinkedResources linked resources to the KaaS cluster
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Preset *bool `json:"preset,omitempty"`

	VPC ReferenceResource `json:"vpc"`

	Subnet ReferenceResource `json:"subnet"`

	NodeCIDR NodeCIDRProperties `json:"nodeCidr"`

	PodCIDR *string `json:"podCidr,omitempty"`

	SecurityGroup SecurityGroupProperties `json:"securityGroup"`

	KubernetesVersion KubernetesVersionInfo `json:"kubernetesVersion"`

	NodePools []NodePoolProperties `json:"nodePools"`

	HA *bool `json:"ha,omitempty"`

	Storage StorageKubernetes `json:"storage,omitempty"`

	BillingPlan BillingPeriodResource `json:"billingPlan"`

	Identity *IdentityProperties `json:"identity,omitempty"`

	APIServerAccessProfile *APIServerAccessProfileProperties `json:"apiServerAccessProfile,omitempty"`
}

type KubernetesVersionInfoUpgradeResponse struct {
	Value *string `json:"value,omitempty"`

	// ScheduledAt Scheduled date and time (nullable)
	ScheduledAt *string `json:"scheduledAt,omitempty"`
}

type InstanceResponse struct {
	// ID Instance identifier (nullable)
	ID *string `json:"id,omitempty"`

	// Name Instance name (nullable)
	Name *string `json:"name,omitempty"`
}

type DataCenterResponse struct {
	// Code Data center code (nullable)
	Code *string `json:"code,omitempty"`

	// Name Data center name (nullable)
	Name *string `json:"name,omitempty"`
}

type NodePoolPropertiesResponse struct {
	// Name Nodepool name (nullable)
	Name *string `json:"name,omitempty"`

	// Nodes Number of nodes (nullable)
	Nodes *int32 `json:"nodes,omitempty"`

	// Instance Configuration of the nodes
	Instance *InstanceResponse `json:"instance,omitempty"`

	// DataCenter Datacenter in which the nodes of the pool will be located
	DataCenter *DataCenterResponse `json:"dataCenter,omitempty"`

	// MinCount Minimum number of nodes for autoscaling (nullable)
	MinCount *int32 `json:"minCount,omitempty"`

	// MaxCount Maximum number of nodes for autoscaling (nullable)
	MaxCount *int32 `json:"maxCount,omitempty"`

	// Autoscaling Indicates if autoscaling is enabled for this node pool
	Autoscaling bool `json:"autoscaling"`

	// CreationDate Creation date and time (nullable)
	CreationDate *string `json:"creationDate,omitempty"`
}

// KubernetesVersionInfoResponse extends KubernetesVersionInfo with additional response fields
type KubernetesVersionInfoResponse struct {
	// Value Value of the version (nullable)
	Value *string `json:"value,omitempty"`

	// EndOfSupportDate End of support date for this version (nullable)
	EndOfSupportDate *string `json:"endOfSupportDate,omitempty"`

	// SellStartDate Start date when this version became available (nullable)
	SellStartDate *string `json:"sellStartDate,omitempty"`

	// SellEndDate End date when this version will no longer be available (nullable)
	SellEndDate *string `json:"sellEndDate,omitempty"`

	// Recommended Indicates if this is the recommended version
	Recommended bool `json:"recommended,omitempty"`

	// UpgradeTo Information about available upgrade
	UpgradeTo *KubernetesVersionInfoUpgradeResponse `json:"upgradeTo,omitempty"`
}

type PodCIDRPropertiesResponse struct {

	// Address in CIDR notation The IP range must be between
	Address *string `json:"address,omitempty"`
}

type NodeCIDRPropertiesResponse struct {

	// Address in CIDR notation The IP range must be between
	Address *string `json:"address,omitempty"`

	Name *string `json:"name,omitempty"`

	URI *string `json:"uri,omitempty"`
}

type KaasSecurityGroupPropertiesResponse struct {
	Name *string `json:"name,omitempty"`

	URI *string `json:"uri,omitempty"`
}

type KaaSPropertiesResponse struct {

	//LinkedResources linked resources to the KaaS cluster
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Preset bool `json:"preset"`

	VPC ReferenceResourceResponse `json:"vpc"`

	Subnet ReferenceResourceResponse `json:"subnet"`

	KubernetesVersion KubernetesVersionInfoResponse `json:"kubernetesVersion"`

	NodePools *[]NodePoolPropertiesResponse `json:"nodesPool,omitempty"`

	PodCIDR *PodCIDRPropertiesResponse `json:"podcidr,omitempty"`

	NodeCIDR NodeCIDRPropertiesResponse `json:"nodecidr"`

	SecurityGroup KaasSecurityGroupPropertiesResponse `json:"securityGroup"`

	HA *bool `json:"ha,omitempty"`

	Storage *StorageKubernetes `json:"storage,omitempty"`

	BillingPlan *BillingPeriodResourceResponse `json:"billingPlan,omitempty"`

	ManagementIP *string `json:"managementIp,omitempty"`

	OpenstackProject *OpenstackProjectResponse `json:"openstackProject,omitempty"`

	Identity *IdentityPropertiesResponse `json:"identity,omitempty"`

	APIServerAccessProfile *APIServerAccessProfilePropertiesResponse `json:"apiServerAccessProfile,omitempty"`
}

type KaaSPropertiesUpdateRequest struct {
	KubernetesVersion KubernetesVersionInfoUpdate `json:"kubernetesVersion"`

	NodePools []NodePoolProperties `json:"nodePools"`

	HA *bool `json:"ha,omitempty"`

	Storage *StorageKubernetes `json:"storage,omitempty"`

	BillingPlan *BillingPeriodResource `json:"billingPlan,omitempty"`
}

type KaaSRequest struct {
	Metadata   RegionalResourceMetadataRequest `json:"metadata"`
	Properties KaaSPropertiesRequest           `json:"properties"`
}

type KaaSUpdateRequest struct {
	Metadata   RegionalResourceMetadataRequest `json:"metadata"`
	Properties KaaSPropertiesUpdateRequest     `json:"properties"`
}

type KaaSResponse struct {
	Metadata   ResourceMetadataResponse `json:"metadata"`
	Properties KaaSPropertiesResponse   `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type KaaSList struct {
	ListResponse
	Values []KaaSResponse `json:"values"`
}

type KaaSKubeconfigResponse struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

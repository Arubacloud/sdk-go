package schema

type NodeCIDRProperties struct {

	// Address in CIDR notation The IP range must be between 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
	Address string `json:"address"`

	// Name of the nodecidr
	Name string `json:"name"`
}

type KubernetesVersionInfo struct {
	Value string `json:"value"`
}

type StorageKubernetes struct {
	MaxCumulativeVolumeSize int32 `json:"maxCumulativeVolumeSize,omitempty"`
}

type NodePoolProperties struct {

	// Name Nodepool name
	Name string `json:"name"`

	// Nodes Number of nodes
	Nodes int32 `json:"nodes"`

	// Instance Configuration name of the nodes.
	// See metadata section of the API documentation for an updated list of admissible values.
	// For more information, check the documentation.
	Instance string `json:"instance"`

	// DataCenter Datacenter in which the nodes of the pool will be located.
	// See metadata section of the API documentation for an updated list of admissible values.
	// For more information, check the documentation.
	Zone string `json:"dataCenter"`
}

type SecurityGroupProperties struct {
	Name string `json:"name"`
}

type KaaSPropertiesRequest struct {

	//LinkedResources linked resources to the KaaS cluster
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Preset bool `json:"preset"`

	VPC ReferenceResource `json:"vpc"`

	Subnet ReferenceResource `json:"subnet"`

	NodeCIDR NodeCIDRProperties `json:"nodeCidr"`

	SecurityGroup SecurityGroupProperties `json:"securityGroup"`

	KubernetesVersion KubernetesVersionInfo `json:"kubernetesVersion"`

	NodePools []NodePoolProperties `json:"nodePools"`

	HA bool `json:"ha"`

	Storage StorageKubernetes `json:"storage,omitempty"`

	BillingPlan BillingPeriodResource `json:"billingPlan"`
}

type KubernetesVersionInfoUpgradeResponse struct {
	Value *string `json:"value,omitempty"`

	// ScheduledAt Scheduled date and time (nullable)
	ScheduledAt *string `json:"scheduledAt,omitempty"`
}

type NodePoolPropertiesResponse struct {
	NodePoolProperties // Embedded struct - inherits the Value field

	// Autoscaling Indicates if autoscaling is enabled for this node pool
	Autoscaling bool `json:"autoscaling,omitempty"`

	// CreationDate Creation date and time (nullable)
	CreationDate *string `json:"creationDate,omitempty"`
}

// KubernetesVersionInfoResponse extends KubernetesVersionInfo with additional response fields
type KubernetesVersionInfoResponse struct {
	KubernetesVersionInfo // Embedded struct - inherits the Value field

	// EndOfSupportDate End of support date for this version (nullable)
	EndOfSupportDate *string `json:"endOfSupportDate,omitempty"`

	// SellStartDate Start date when this version became available (nullable)
	SellStartDate *string `json:"sellStartDate,omitempty"`

	// SellEndDate End date when this version will no longer be available (nullable)
	SellEndDate *string `json:"sellEndDate,omitempty"`

	// Recommended Indicates if this is the recommended version
	Recommended bool `json:"recommended,omitempty"`
}

type PodCIDRPropertiesResponse struct {

	// Address in CIDR notation The IP range must be between
	Address string `json:"address,omitempty"`
}

type NodeCIDRPropertiesResponse struct {

	// Address in CIDR notation The IP range must be between
	Address string `json:"address,omitempty"`

	Name string `json:"name,omitempty"`

	URI string `json:"uri,omitempty"`
}

type KaasSecurityGroupPropertiesResponse struct {
	Name string `json:"name,omitempty"`

	URI string `json:"uri,omitempty"`
}

type KaaSPropertiesResponse struct {

	//LinkedResources linked resources to the KaaS cluster
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Preset bool `json:"preset"`

	VPC ReferenceResource `json:"vpc"`

	Subnet ReferenceResource `json:"subnet"`

	KubernetesVersion KubernetesVersionInfoResponse `json:"kubernetesVersion"`

	NodePools []NodePoolPropertiesResponse `json:"nodesPool"`

	PodCIDR PodCIDRPropertiesResponse `json:"podCidr,omitempty"`

	NodeCIDR NodeCIDRPropertiesResponse `json:"nodeCidr"`

	SecurityGroup KaasSecurityGroupPropertiesResponse `json:"securityGroup"`

	HA bool `json:"ha"`

	Storage StorageKubernetes `json:"storage,omitempty"`

	BillingPlan BillingPeriodResource `json:"billingPlan"`

	ManagementIP *string `json:"managementIp,omitempty"`
}

type KaaSRequest struct {
	Metadata   RegionalResourceMetadataRequest `json:"metadata"`
	Properties KaaSPropertiesRequest           `json:"properties"`
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

package types

// DBaaSEngine contains the database engine configuration
type DBaaSEngine struct {
	// ID Type of DB engine to activate (nullable)
	// For more information, check the documentation.
	ID *string `json:"id,omitempty"`

	// DataCenter Datacenter location (nullable)
	// For more information, check the documentation.
	DataCenter *string `json:"dataCenter,omitempty"`
}

// DBaaSEngineResponse contains the database engine response configuration
type DBaaSEngineResponse struct {
	// ID Engine identifier (nullable)
	ID *string `json:"id,omitempty"`

	// Type Engine type (nullable)
	Type *string `json:"type,omitempty"`

	// Name Engine name (nullable)
	Name *string `json:"name,omitempty"`

	// Version Engine version (nullable)
	Version *string `json:"version,omitempty"`

	// DataCenter Datacenter location (nullable)
	// For more information, check the documentation.
	DataCenter *string `json:"dataCenter,omitempty"`

	// PrivateIPAddress Private IP address (nullable)
	PrivateIPAddress *string `json:"privateIpAddress,omitempty"`
}

// DBaaSFlavor contains the flavor configuration
type DBaaSFlavor struct {
	// Name Type of flavor to use (nullable)
	// For more information, check the documentation.
	Name *string `json:"name,omitempty"`
}

// DBaaSFlavorResponse contains the flavor response configuration
type DBaaSFlavorResponse struct {
	// Name Flavor name (nullable)
	Name *string `json:"name,omitempty"`

	// Category Flavor category (nullable)
	Category *string `json:"category,omitempty"`

	// CPU Number of CPUs (nullable)
	CPU *int32 `json:"cpu,omitempty"`

	// RAM Amount of RAM in MB (nullable)
	RAM *int32 `json:"ram,omitempty"`
}

// DBaaSStorage contains the storage configuration
type DBaaSStorage struct {
	// SizeGB Size in GB to use (nullable)
	SizeGB *int32 `json:"sizeGb,omitempty"`
}

// DBaaSStorageResponse contains the storage response configuration
type DBaaSStorageResponse struct {
	// SizeGB Size in GB (nullable)
	SizeGB *int32 `json:"sizeGb,omitempty"`
}

// DBaaSBillingPlan contains the billing plan configuration
type DBaaSBillingPlan struct {
	// BillingPeriod Type of billing period to use (nullable)
	BillingPeriod *string `json:"billingPeriod,omitempty"`
}

// DBaaSBillingPlanResponse contains the billing plan response configuration
type DBaaSBillingPlanResponse struct {
	// BillingPeriod Billing period (nullable)
	BillingPeriod *string `json:"billingPeriod,omitempty"`
}

// DBaaSNetworking contains the network information to use when creating the new DBaaS
type DBaaSNetworking struct {
	// VPCURI The URI of the VPC resource to bind to this DBaaS instance (nullable)
	// Required when user has at least one VPC (with at least one subnet and a security group).
	VPCURI *string `json:"vpcUri,omitempty"`

	// SubnetURI The URI of the Subnet resource to bind to this DBaaS instance (nullable)
	// It must belong to the VPC defined in VPCURI
	// Required when user has at least one VPC (with at least one subnet and a security group).
	SubnetURI *string `json:"subnetUri,omitempty"`

	// SecurityGroupURI The URI of the SecurityGroup resource to bind to this DBaaS instance (nullable)
	// It must belong to the VPC defined in VPCURI
	// Required when user has at least one VPC (with at least one subnet and a security group).
	SecurityGroupURI *string `json:"securityGroupUri,omitempty"`

	// ElasticIPURI The URI of the ElasticIP resource to bind to this DBaaS instance (nullable)
	ElasticIPURI *string `json:"elasticIpUri,omitempty"`
}

// DBaaSNetworkingResponse contains the network response information
type DBaaSNetworkingResponse struct {
	// VPC VPC resource reference (nullable)
	VPC *ReferenceResource `json:"vpc,omitempty"`

	// Subnet Subnet resource reference (nullable)
	Subnet *ReferenceResource `json:"subnet,omitempty"`

	// SecurityGroup Security group resource reference (nullable)
	SecurityGroup *ReferenceResource `json:"securityGroup,omitempty"`

	// ElasticIP Elastic IP resource reference (nullable)
	ElasticIP *ReferenceResource `json:"elasticIp,omitempty"`
}

// DBaaSAutoscaling contains the autoscaling configuration
type DBaaSAutoscaling struct {
	// Enabled Indicates if autoscaling is enabled (nullable)
	Enabled *bool `json:"enabled,omitempty"`

	// AvailableSpace Available space threshold (nullable)
	AvailableSpace *int32 `json:"availableSpace,omitempty"`

	// StepSize Step size for autoscaling (nullable)
	StepSize *int32 `json:"stepSize,omitempty"`
}

// DBaaSAutoscalingResponse contains the autoscaling response configuration
type DBaaSAutoscalingResponse struct {
	// Status Autoscaling status (nullable)
	Status *string `json:"status,omitempty"`

	// AvailableSpace Available space threshold (nullable)
	AvailableSpace *int32 `json:"availableSpace,omitempty"`

	// StepSize Step size for autoscaling (nullable)
	StepSize *int32 `json:"stepSize,omitempty"`

	// RuleID Rule identifier (nullable)
	RuleID *string `json:"ruleId,omitempty"`
}

// DBaaSPropertiesRequest contains properties required to create a DBaaS instance
type DBaaSPropertiesRequest struct {
	// Engine Database engine configuration
	Engine *DBaaSEngine `json:"engine,omitempty"`

	// Flavor Flavor configuration
	Flavor *DBaaSFlavor `json:"flavor,omitempty"`

	// Storage Storage configuration
	Storage *DBaaSStorage `json:"storage,omitempty"`

	// BillingPlan Billing plan configuration
	BillingPlan *DBaaSBillingPlan `json:"billingPlan,omitempty"`

	// Networking Network information for the DBaaS instance
	Networking *DBaaSNetworking `json:"networking,omitempty"`

	// Autoscaling Autoscaling configuration
	Autoscaling *DBaaSAutoscaling `json:"autoscaling,omitempty"`
}

// DBaaSPropertiesResponse contains the response properties of a DBaaS instance
type DBaaSPropertiesResponse struct {
	// LinkedResources Array of resources linked to the DBaaS instance (nullable)
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	// Engine Database engine response configuration
	Engine *DBaaSEngineResponse `json:"engine,omitempty"`

	// Flavor Flavor response configuration
	Flavor *DBaaSFlavorResponse `json:"flavor,omitempty"`

	// Networking Network response configuration
	Networking *DBaaSNetworkingResponse `json:"networking,omitempty"`

	// Storage Storage response configuration
	Storage *DBaaSStorageResponse `json:"storage,omitempty"`

	// BillingPlan Billing plan response configuration
	BillingPlan *DBaaSBillingPlanResponse `json:"billingPlan,omitempty"`

	// Autoscaling Autoscaling response configuration
	Autoscaling *DBaaSAutoscalingResponse `json:"autoscaling,omitempty"`
}

type DBaaSRequest struct {
	// Metadata of the DBaaS instance
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	// Spec contains the DBaaS instance specification
	Properties DBaaSPropertiesRequest `json:"properties"`
}

type DBaaSResponse struct {
	// Metadata of the DBaaS instance
	Metadata ResourceMetadataResponse `json:"metadata"`

	// Spec contains the DBaaS instance specification
	Properties DBaaSPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type DBaaSList struct {
	ListResponse
	Values []DBaaSResponse `json:"values"`
}

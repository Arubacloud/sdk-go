package types

type UserCredential struct {

	// Username is the administrator username for the container registry
	Username string `json:"username"`
}

type ContainerRegistryPropertiesRequest struct {

	// PublicIp is the public IP associated with the container registry
	PublicIp ReferenceResource `json:"publicIp"`

	VPC ReferenceResource `json:"vpc"`

	// Subnet is the subnet associated with the container registry
	Subnet ReferenceResource `json:"subnet"`

	// SecurityGroup is the security group associated with the container registry
	SecurityGroup ReferenceResource `json:"securityGroup"`

	// BlockStorage is the block storage associated with the container registry
	BlockStorage ReferenceResource `json:"blockStorage"`

	// BillingPlan is the billing plan associated with the container registry
	BillingPlan *BillingPeriodResource `json:"billingPlan,omitempty"`

	// AdminUser is the administrator user for the container registry
	AdminUser *UserCredential `json:"adminUser,omitempty"`

	// Size is the number of concurrent users allowed for the container registry
	ConcurrentUsers *int `json:"size,omitempty"`
}

type ContainerRegistryPropertiesResult struct {

	// PublicIp is the public IP associated with the container registry
	PublicIp ReferenceResource `json:"publicIp"`

	// VPC is the VPC associated with the container registry
	VPC ReferenceResource `json:"vpc"`

	// Subnet is the subnet associated with the container registry
	Subnet ReferenceResource `json:"subnet"`

	// SecurityGroup is the security group associated with the container registry
	SecurityGroup ReferenceResource `json:"securityGroup"`

	// BlockStorage is the block storage associated with the container registry
	BlockStorage ReferenceResource `json:"blockStorage"`

	// BillingPlan is the billing plan associated with the container registry
	BillingPlan *BillingPeriodResource `json:"billingPlan,omitempty"`

	// AdminUser is the administrator user for the container registry
	AdminUser *UserCredential `json:"adminUser,omitempty"`

	// Size is the number of concurrent users allowed for the container registry
	ConcurrentUsers *int `json:"size,omitempty"`
}

type ContainerRegistryRequest struct {
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	Properties ContainerRegistryPropertiesRequest `json:"properties"`
}

type ContainerRegistryResponse struct {
	Metadata   ResourceMetadataResponse          `json:"metadata"`
	Properties ContainerRegistryPropertiesResult `json:"properties"`
	Status     ResourceStatus                    `json:"status,omitempty"`
}

type ContainerRegistryPropertiesResponse struct {
	Metadata   ResourceMetadataResponse          `json:"metadata"`
	Properties ContainerRegistryPropertiesResult `json:"properties"`
}

type ContainerRegistryList struct {
	ListResponse
	Values []ContainerRegistryResponse `json:"values"`
}

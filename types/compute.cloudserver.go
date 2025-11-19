package types

type CloudServerPropertiesRequest struct {
	Zone string `json:"dataCenter"`

	VPC ReferenceResource `json:"vpc"`

	VPCPreset bool `json:"vpcPreset,omitempty"`

	FlavorName *string `json:"flavorName,omitempty"`

	ElastcIP ReferenceResource `json:"elasticIp"`

	BootVolume ReferenceResource `json:"bootVolume"`

	KeyPair ReferenceResource `json:"keyPair"`

	Subnets []ReferenceResource `json:"subnets,omitempty"`

	SecurityGroups []ReferenceResource `json:"securityGroups,omitempty"`
}

type CloudServerFlavorResponse struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Category string `json:"category"`

	CPU int32 `json:"cpu"`

	RAM int32 `json:"ram"`

	HD int32 `json:"hd"`
}

type CloudServerNetworkInterfaceDetails struct {
	Subnet *string `json:"subnet,omitempty"`

	MacAddress *string `json:"macAddress,omitempty"`

	IPs []string `json:"ips,omitempty"`
}

type CloudServerPropertiesResult struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Zone string `json:"dataCenter"`

	VPC ReferenceResource `json:"vpc"`

	Flavor CloudServerFlavorResponse `json:"flavor,omitempty"`

	Template ReferenceResource `json:"template"`

	BootVolume ReferenceResource `json:"bootVolume"`

	KeyPair ReferenceResource `json:"keyPair"`

	NetworkInterfaces []CloudServerNetworkInterfaceDetails `json:"networkInterfaces,omitempty"`
}

type CloudServerRequest struct {
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	Properties CloudServerPropertiesRequest `json:"properties"`
}

type CloudServerResponse struct {
	Metadata   RegionalResourceMetadataRequest `json:"metadata"`
	Properties CloudServerPropertiesResult     `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type CloudServerList struct {
	ListResponse
	Values []CloudServerResponse `json:"values"`
}

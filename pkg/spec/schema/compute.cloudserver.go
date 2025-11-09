package schema

type CloudServerPropertiesRequest struct {
	Zone string `json:"dataCenter"`

	Vpc ReferenceResource `json:"vpc"`

	VpcPreset bool `json:"vpcPreset,omitempty"`

	FlavorName *string `json:"flavorName,omitempty"`

	ElastcIp ReferenceResource `json:"elasticIp"`

	BootVolume ReferenceResource `json:"bootVolume"`

	KeyPair ReferenceResource `json:"keyPair"`

	Subnets []ReferenceResource `json:"subnets,omitempty"`

	SecurityGroups []ReferenceResource `json:"securityGroups,omitempty"`
}

type CloudServerFlavorResponse struct {
	Id string `json:"id"`

	Name string `json:"name"`

	Category string `json:"category"`

	CPU int32 `json:"cpu"`

	RAM int32 `json:"ram"`

	Hd int32 `json:"hd"`
}

type CloudServerNetworkInterfaceDetails struct {
	Subnet *string `json:"subnet, omitempty"`

	MacAddress *string `json:"macAddress, omitempty"`

	IPs []string `json:"ips, omitempty"`
}

type CloudServerPropertiesResult struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Zone string `json:"dataCenter"`

	Vpc ReferenceResource `json:"vpc"`

	Flavor CloudServerFlavorResponse `json:"flavor,omitempty"`

	Template ReferenceResource `json:"template"`

	BootVolume ReferenceResource `json:"bootVolume"`

	KeyPair ReferenceResource `json:"keyPair"`

	NetworkInterfaces []CloudServerNetworkInterfaceDetails `json:"networkInterfaces,omitempty"`
}

type CloudServerRequest struct {
	Metadata ResourceMetadataRequest `json:"metadata"`

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

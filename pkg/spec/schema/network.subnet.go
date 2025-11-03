package schema

// SubnetProperties contains the properties of a Subnet
type SubnetProperties struct {
	// Default indicates if the subnet must be a default subnet. Only one default subnet for region is admissible.
	// Default value: true
	Default *bool `json:"default,omitempty"`

	// Preset indicates if a subnet and a securityGroup with default configuration will be created automatically within the subnet
	// Default value: false
	Preset *bool `json:"preset,omitempty"`
}

// SubnetPropertiesRequest contains the specification for creating a Subnet
type SubnetPropertiesRequest struct {
	// Properties of the subnet (nullable object)
	Properties SubnetProperties `json:"properties"`
}

// SubnetPropertiesResponse contains the specification returned for a Subnet
type SubnetPropertiesResponse struct {
	// LinkedResources array of resources linked to the Subnet (nullable)
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	// Default indicates if the subnet is the default one within the region
	Default bool `json:"default,omitempty"`
}

type SubnetRequest struct {
	// Metadata of the Subnet
	Metadata ResourceMetadataRequest `json:"metadata"`

	// Spec contains the Subnet specification
	Properties SubnetPropertiesRequest `json:"properties"`
}

type SubnetResponse struct {
	// Metadata of the Subnet
	Metadata ResourceMetadataResponse `json:"metadata"`
	// Spec contains the Subnet specification
	Properties SubnetPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type SubnetList struct {
	ListResponse
	Values []SubnetResponse `json:"values"`
}

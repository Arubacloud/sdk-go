package schema

// VpcProperties contains the properties of a VPC
type VpcProperties struct {
	// Default indicates if the vpc must be a default vpc. Only one default vpc for region is admissible.
	// Default value: true
	Default *bool `json:"default,omitempty"`

	// Preset indicates if a subnet and a securityGroup with default configuration will be created automatically within the vpc
	// Default value: false
	Preset *bool `json:"preset,omitempty"`
}

// VpcPropertiesRequest contains the specification for creating a VPC
type VpcPropertiesRequest struct {
	// Properties of the vpc (nullable object)
	Properties *VpcProperties `json:"properties,omitempty"`
}

// VpcPropertiesResponse contains the specification returned for a VPC
type VpcPropertiesResponse struct {
	// LinkedResources array of resources linked to the VPC (nullable)
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	// Default indicates if the vpc is the default one within the region
	Default bool `json:"default,omitempty"`
}

type VpcRequest struct {
	// Metadata of the VPC
	Metadata ResourceMetadataRequest `json:"metadata"`

	// Spec contains the VPC specification
	Properties VpcPropertiesRequest `json:"properties"`
}

type VpcResponse struct {
	// Metadata of the VPC
	Metadata ResourceMetadataResponse `json:"metadata"`
	// Spec contains the VPC specification
	Properties VpcPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type VpcList struct {
	ListResponse
	Values []VpcResponse `json:"values"`
}

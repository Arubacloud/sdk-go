package types

// SubnetType represents the type of subnet
type SubnetType string

const (
	SubnetTypeBasic    SubnetType = "Basic"
	SubnetTypeAdvanced SubnetType = "Advanced"
)

// SubnetNetwork contains the network configuration
type SubnetNetwork struct {
	Address string `json:"address"`
}

// SubnetDHCP contains the DHCP configuration
type SubnetDHCP struct {
	Enabled bool `json:"enabled"`
}

// SubnetPropertiesRequest contains the specification for creating a Subnet
type SubnetPropertiesRequest struct {
	// Type of subnet (Basic or Advanced)
	Type SubnetType `json:"type,omitempty"`

	// Default indicates if the subnet must be a default subnet
	Default bool `json:"default,omitempty"`

	// Network configuration
	Network *SubnetNetwork `json:"network,omitempty"`

	// DHCP configuration
	DHCP *SubnetDHCP `json:"dhcp,omitempty"`
}

// SubnetPropertiesResponse contains the specification returned for a Subnet
type SubnetPropertiesResponse struct {
	// LinkedResources array of resources linked to the Subnet (nullable)
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	// Type of subnet
	Type SubnetType `json:"type,omitempty"`

	// Default indicates if the subnet is the default one within the region
	Default bool `json:"default,omitempty"`

	// Network configuration
	Network *SubnetNetwork `json:"network,omitempty"`

	// DHCP configuration
	DHCP *SubnetDHCP `json:"dhcp,omitempty"`
}

type SubnetRequest struct {
	// Metadata of the Subnet
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

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

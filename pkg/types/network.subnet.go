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

// SubnetDHCPRange contains the DHCP range configuration
type SubnetDHCPRange struct {
	// Start is the starting IP address of the DHCP range
	Start string `json:"start"`
	// Count is the number of IP addresses in the DHCP range
	Count int `json:"count"`
}

// SubnetDHCPRoute contains the DHCP route configuration
type SubnetDHCPRoute struct {
	// Address is the destination network address
	Address string `json:"address"`
	// Gateway is the gateway IP address for the route
	Gateway string `json:"gateway"`
}

// SubnetDHCP contains the DHCP configuration
type SubnetDHCP struct {
	// Enabled indicates if DHCP is enabled
	Enabled bool `json:"enabled"`
	// Range contains the DHCP IP address range
	Range *SubnetDHCPRange `json:"range,omitempty"`
	// Routes contains the DHCP routes configuration
	Routes []SubnetDHCPRoute `json:"routes,omitempty"`
	// DNS contains the DNS server addresses
	DNS []string `json:"dns,omitempty"`
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

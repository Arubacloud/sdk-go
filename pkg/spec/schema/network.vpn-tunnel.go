package schema

// IpConfigurations contains network configuration of the VPN tunnel
type IpConfigurations struct {
	// Vpc reference to the VPC (nullable)
	Vpc *ReferenceResource `json:"vpc,omitempty"`

	// Subnet reference to the subnet (nullable)
	Subnet *ReferenceResource `json:"subnet,omitempty"`

	// PublicIp reference to the public IP (nullable)
	PublicIp *ReferenceResource `json:"publicIp,omitempty"`
}

// IkeSettings contains IKE settings
type IkeSettings struct {
	// Lifetime Lifetime value
	Lifetime int32 `json:"lifetime,omitempty"`

	// Encryption Encryption algorithm (nullable)
	Encryption *string `json:"encryption,omitempty"`

	// Hash Hash algorithm (nullable)
	Hash *string `json:"hash,omitempty"`

	// DhGroup Diffie-Hellman group (nullable)
	DhGroup *string `json:"dhGroup,omitempty"`

	// DpdAction Dead Peer Detection action (nullable)
	DpdAction *string `json:"dpdAction,omitempty"`

	// DpdInterval Dead Peer Detection interval
	DpdInterval int32 `json:"dpdInterval,omitempty"`

	// DpdTimeout Dead Peer Detection timeout
	DpdTimeout int32 `json:"dpdTimeout,omitempty"`
}

// EspSettings contains ESP settings
type EspSettings struct {
	// Lifetime Lifetime value
	Lifetime int32 `json:"lifetime,omitempty"`

	// Encryption Encryption algorithm (nullable)
	Encryption *string `json:"encryption,omitempty"`

	// Hash Hash algorithm (nullable)
	Hash *string `json:"hash,omitempty"`

	// Pfs Perfect Forward Secrecy (nullable)
	Pfs *string `json:"pfs,omitempty"`
}

// PskSettings contains Pre-Shared Key settings
type PskSettings struct {
	// CloudSite Cloud site identifier (nullable)
	CloudSite *string `json:"cloudSite,omitempty"`

	// OnPremSite On-premises site identifier (nullable)
	OnPremSite *string `json:"onPremSite,omitempty"`

	// Secret Pre-shared key secret (nullable)
	Secret *string `json:"secret,omitempty"`
}

// VpnClientSettings contains client settings of the VPN tunnel
type VpnClientSettings struct {
	// Ike IKE settings (nullable)
	Ike *IkeSettings `json:"ike,omitempty"`

	// Esp ESP settings (nullable)
	Esp *EspSettings `json:"esp,omitempty"`

	// Psk Pre-Shared Key settings (nullable)
	Psk *PskSettings `json:"psk,omitempty"`
}

// VpnTunnelPropertiesRequest contains properties of a VPN tunnel
type VpnTunnelPropertiesRequest struct {
	// VpnType Type of VPN tunnel. Admissible values: Site-To-Site (nullable)
	VpnType *string `json:"vpnType,omitempty"`

	// VpnClientProtocol Protocol of the VPN tunnel. Admissible values: ikev2 (nullable)
	VpnClientProtocol *string `json:"vpnClientProtocol,omitempty"`

	// IpConfigurations Network configuration of the VPN tunnel (nullable)
	IpConfigurations *IpConfigurations `json:"ipConfigurations,omitempty"`

	// VpnClientSettings Client settings of the VPN tunnel (nullable)
	VpnClientSettings *VpnClientSettings `json:"vpnClientSettings,omitempty"`

	// PeerClientPublicIp Peer client public IP address (nullable)
	PeerClientPublicIp *string `json:"peerClientPublicIp,omitempty"`

	// BillingPlan Billing plan
	BillingPlan *BillingPeriodResource `json:"billingPlan,omitempty"`
}

// VpnTunnelPropertiesResponse contains the response properties of a VPN tunnel
type VpnTunnelPropertiesResponse struct {
	// VpnType Type of the VPN tunnel (nullable)
	VpnType *string `json:"vpnType,omitempty"`

	// VpnClientProtocol Protocol of the VPN tunnel (nullable)
	VpnClientProtocol *string `json:"vpnClientProtocol,omitempty"`

	// IpConfigurations Network configuration of the VPN tunnel (nullable)
	IpConfigurations *IpConfigurations `json:"ipConfigurations,omitempty"`

	// VpnClientSettings Client settings of the VPN tunnel (nullable)
	VpnClientSettings *VpnClientSettings `json:"vpnClientSettings,omitempty"`

	// PeerClientPublicIp Peer client public IP address (nullable)
	PeerClientPublicIp *string `json:"peerClientPublicIp,omitempty"`

	// RoutesNumber Number of valid VPN routes of the VPN tunnel
	RoutesNumber int32 `json:"routesNumber,omitempty"`

	// BillingPlan Billing plan (nullable)
	BillingPlan *BillingPeriodResource `json:"billingPlan,omitempty"`
}

type VpnTunnelRequest struct {
	// Metadata of the VPN Tunnel
	Metadata ResourceMetadataRequest `json:"metadata"`

	// Spec contains the VPN Tunnel specification
	Properties VpnTunnelPropertiesRequest `json:"properties"`
}

type VpnTunnelResponse struct {
	// Metadata of the VPN Tunnel
	Metadata ResourceMetadataResponse `json:"metadata"`
	// Spec contains the VPN Tunnel specification
	Properties VpnTunnelPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type VpnTunnelList struct {
	ListResponse
	Values []VpnTunnelResponse `json:"values"`
}

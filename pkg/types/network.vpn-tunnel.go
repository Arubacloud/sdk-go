package types

// IPConfigurations contains network configuration of the VPN tunnel
type IPConfigurations struct {
	// VPC reference to the VPC (nullable)
	VPC *ReferenceResource `json:"vpc,omitempty"`

	// Subnet reference to the subnet (nullable)
	Subnet *ReferenceResource `json:"subnet,omitempty"`

	// PublicIP reference to the public IP (nullable)
	PublicIP *ReferenceResource `json:"publicIp,omitempty"`
}

// IKESettings contains IKE settings
type IKESettings struct {
	// Lifetime Lifetime value
	Lifetime int32 `json:"lifetime,omitempty"`

	// Encryption Encryption algorithm (nullable)
	Encryption *string `json:"encryption,omitempty"`

	// Hash Hash algorithm (nullable)
	Hash *string `json:"hash,omitempty"`

	// DHGroup Diffie-Hellman group (nullable)
	DHGroup *string `json:"dhGroup,omitempty"`

	// DPDAction Dead Peer Detection action (nullable)
	DPDAction *string `json:"dpdAction,omitempty"`

	// DPDInterval Dead Peer Detection interval
	DPDInterval int32 `json:"dpdInterval,omitempty"`

	// DPDTimeout Dead Peer Detection timeout
	DPDTimeout int32 `json:"dpdTimeout,omitempty"`
}

// ESPSettings contains ESP settings
type ESPSettings struct {
	// Lifetime Lifetime value
	Lifetime int32 `json:"lifetime,omitempty"`

	// Encryption Encryption algorithm (nullable)
	Encryption *string `json:"encryption,omitempty"`

	// Hash Hash algorithm (nullable)
	Hash *string `json:"hash,omitempty"`

	// PFS Perfect Forward Secrecy (nullable)
	PFS *string `json:"pfs,omitempty"`
}

// PSKSettings contains Pre-Shared Key settings
type PSKSettings struct {
	// CloudSite Cloud site identifier (nullable)
	CloudSite *string `json:"cloudSite,omitempty"`

	// OnPremSite On-premises site identifier (nullable)
	OnPremSite *string `json:"onPremSite,omitempty"`

	// Secret Pre-shared key secret (nullable)
	Secret *string `json:"secret,omitempty"`
}

// VPNClientSettings contains client settings of the VPN tunnel
type VPNClientSettings struct {
	// IKE settings (nullable)
	IKE *IKESettings `json:"ike,omitempty"`

	// ESP settings (nullable)
	ESP *ESPSettings `json:"esp,omitempty"`

	// PSK Pre-Shared Key settings (nullable)
	PSK *PSKSettings `json:"psk,omitempty"`
}

// VPNTunnelPropertiesRequest contains properties of a VPN tunnel
type VPNTunnelPropertiesRequest struct {
	// VPNType Type of VPN tunnel. Admissible values: Site-To-Site (nullable)
	VPNType *string `json:"vpnType,omitempty"`

	// VPNClientProtocol Protocol of the VPN tunnel. Admissible values: ikev2 (nullable)
	VPNClientProtocol *string `json:"vpnClientProtocol,omitempty"`

	// IPConfigurations Network configuration of the VPN tunnel (nullable)
	IPConfigurations *IPConfigurations `json:"ipConfigurations,omitempty"`

	// VPNClientSettings Client settings of the VPN tunnel (nullable)
	VPNClientSettings *VPNClientSettings `json:"vpnClientSettings,omitempty"`

	// PeerClientPublicIP Peer client public IP address (nullable)
	PeerClientPublicIP *string `json:"peerClientPublicIp,omitempty"`

	// BillingPlan Billing plan
	BillingPlan *BillingPeriodResource `json:"billingPlan,omitempty"`
}

// VPNTunnelPropertiesResponse contains the response properties of a VPN tunnel
type VPNTunnelPropertiesResponse struct {
	// VPNType Type of the VPN tunnel (nullable)
	VPNType *string `json:"vpnType,omitempty"`

	// VPNClientProtocol Protocol of the VPN tunnel (nullable)
	VPNClientProtocol *string `json:"vpnClientProtocol,omitempty"`

	// IPConfigurations Network configuration of the VPN tunnel (nullable)
	IPConfigurations *IPConfigurations `json:"ipConfigurations,omitempty"`

	// VPNClientSettings Client settings of the VPN tunnel (nullable)
	VPNClientSettings *VPNClientSettings `json:"vpnClientSettings,omitempty"`

	// PeerClientPublicIP Peer client public IP address (nullable)
	PeerClientPublicIP *string `json:"peerClientPublicIp,omitempty"`

	// RoutesNumber Number of valid VPN routes of the VPN tunnel
	RoutesNumber int32 `json:"routesNumber,omitempty"`

	// BillingPlan Billing plan (nullable)
	BillingPlan *BillingPeriodResource `json:"billingPlan,omitempty"`
}

type VPNTunnelRequest struct {
	// Metadata of the VPN Tunnel
	Metadata RegionalResourceMetadataRequest `json:"metadata"`

	// Spec contains the VPN Tunnel specification
	Properties VPNTunnelPropertiesRequest `json:"properties"`
}

type VPNTunnelResponse struct {
	// Metadata of the VPN Tunnel
	Metadata ResourceMetadataResponse `json:"metadata"`
	// Spec contains the VPN Tunnel specification
	Properties VPNTunnelPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type VPNTunnelList struct {
	ListResponse
	Values []VPNTunnelResponse `json:"values"`
}

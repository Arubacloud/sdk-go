package types

// VPN Encryption algorithm constants for IKESettings.Encryption and ESPSettings.Encryption.
const (
	VPNEncryptionAES128           = "aes128"
	VPNEncryptionAES192           = "aes192"
	VPNEncryptionAES256           = "aes256"
	VPNEncryptionAES128CTR        = "aes128ctr"
	VPNEncryptionAES192CTR        = "aes192ctr"
	VPNEncryptionAES256CTR        = "aes256ctr"
	VPNEncryptionAES128CCM64      = "aes128ccm64"
	VPNEncryptionAES128CCM96      = "aes128ccm96"
	VPNEncryptionAES128CCM128     = "aes128ccm128"
	VPNEncryptionAES192CCM64      = "aes192ccm64"
	VPNEncryptionAES192CCM96      = "aes192ccm96"
	VPNEncryptionAES192CCM128     = "aes192ccm128"
	VPNEncryptionAES256CCM64      = "aes256ccm64"
	VPNEncryptionAES256CCM96      = "aes256ccm96"
	VPNEncryptionAES256CCM128     = "aes256ccm128"
	VPNEncryptionAES128GCM64      = "aes128gcm64"
	VPNEncryptionAES128GCM96      = "aes128gcm96"
	VPNEncryptionAES128GCM128     = "aes128gcm128"
	VPNEncryptionAES192GCM64      = "aes192gcm64"
	VPNEncryptionAES192GCM96      = "aes192gcm96"
	VPNEncryptionAES192GCM128     = "aes192gcm128"
	VPNEncryptionAES256GCM64      = "aes256gcm64"
	VPNEncryptionAES256GCM96      = "aes256gcm96"
	VPNEncryptionAES256GCM128     = "aes256gcm128"
	VPNEncryptionAES128GMAC       = "aes128gmac"
	VPNEncryptionAES192GMAC       = "aes192gmac"
	VPNEncryptionAES256GMAC       = "aes256gmac"
	VPNEncryption3DES             = "3des"
	VPNEncryptionBlowfish128      = "blowfish128"
	VPNEncryptionBlowfish192      = "blowfish192"
	VPNEncryptionBlowfish256      = "blowfish256"
	VPNEncryptionCamellia128      = "camellia128"
	VPNEncryptionCamellia192      = "camellia192"
	VPNEncryptionCamellia256      = "camellia256"
	VPNEncryptionCamellia128CTR   = "camellia128ctr"
	VPNEncryptionCamellia192CTR   = "camellia192ctr"
	VPNEncryptionCamellia256CTR   = "camellia256ctr"
	VPNEncryptionCamellia128CCM64 = "camellia128ccm64"
	VPNEncryptionCamellia128CCM96 = "camellia128ccm96"
	VPNEncryptionCamellia128CCM128 = "camellia128ccm128"
	VPNEncryptionCamellia192CCM64 = "camellia192ccm64"
	VPNEncryptionCamellia192CCM96 = "camellia192ccm96"
	VPNEncryptionCamellia192CCM128 = "camellia192ccm128"
	VPNEncryptionCamellia256CCM64 = "camellia256ccm64"
	VPNEncryptionCamellia256CCM96 = "camellia256ccm96"
	VPNEncryptionCamellia256CCM128 = "camellia256ccm128"
	VPNEncryptionSerpent128       = "serpent128"
	VPNEncryptionSerpent192       = "serpent192"
	VPNEncryptionSerpent256       = "serpent256"
	VPNEncryptionTwofish128       = "twofish128"
	VPNEncryptionTwofish192       = "twofish192"
	VPNEncryptionTwofish256       = "twofish256"
	VPNEncryptionCAST128          = "cast128"
	VPNEncryptionChaCha20Poly1305 = "chacha20poly1305"
)

// VPN Hash algorithm constants for IKESettings.Hash and ESPSettings.Hash.
const (
	VPNHashMD5       = "md5"
	VPNHashMD5128    = "md5_128"
	VPNHashSHA1      = "sha1"
	VPNHashSHA1160   = "sha1_160"
	VPNHashSHA256    = "sha256"
	VPNHashSHA25696  = "sha256_96"
	VPNHashSHA384    = "sha384"
	VPNHashSHA512    = "sha512"
	VPNHashAESXCBC   = "aesxcbc"
	VPNHashAESCMAC   = "aescmac"
	VPNHashAES128GMAC = "aes128gmac"
	VPNHashAES192GMAC = "aes192gmac"
	VPNHashAES256GMAC = "aes256gmac"
)

// VPN DH group constants for IKESettings.DHGroup.
const (
	VPNDHGroup1  = "1"
	VPNDHGroup2  = "2"
	VPNDHGroup5  = "5"
	VPNDHGroup14 = "14"
	VPNDHGroup15 = "15"
	VPNDHGroup16 = "16"
	VPNDHGroup17 = "17"
	VPNDHGroup18 = "18"
	VPNDHGroup19 = "19"
	VPNDHGroup20 = "20"
	VPNDHGroup21 = "21"
	VPNDHGroup22 = "22"
	VPNDHGroup23 = "23"
	VPNDHGroup24 = "24"
	VPNDHGroup25 = "25"
	VPNDHGroup26 = "26"
	VPNDHGroup27 = "27"
	VPNDHGroup28 = "28"
	VPNDHGroup29 = "29"
	VPNDHGroup30 = "30"
	VPNDHGroup31 = "31"
	VPNDHGroup32 = "32"
)

// VPN DPD action constants for IKESettings.DPDAction.
const (
	VPNDPDActionTrap    = "trap"
	VPNDPDActionClear   = "clear"
	VPNDPDActionRestart = "restart"
)

// VPN PFS group constants for ESPSettings.PFS.
const (
	VPNPFSEnable     = "enable"
	VPNPFSDisable    = "disable"
	VPNPFSDHGroup1   = "dh-group1"
	VPNPFSDHGroup2   = "dh-group2"
	VPNPFSDHGroup5   = "dh-group5"
	VPNPFSDHGroup14  = "dh-group14"
	VPNPFSDHGroup15  = "dh-group15"
	VPNPFSDHGroup16  = "dh-group16"
	VPNPFSDHGroup17  = "dh-group17"
	VPNPFSDHGroup18  = "dh-group18"
	VPNPFSDHGroup19  = "dh-group19"
	VPNPFSDHGroup20  = "dh-group20"
	VPNPFSDHGroup21  = "dh-group21"
	VPNPFSDHGroup22  = "dh-group22"
	VPNPFSDHGroup23  = "dh-group23"
	VPNPFSDHGroup24  = "dh-group24"
	VPNPFSDHGroup25  = "dh-group25"
	VPNPFSDHGroup26  = "dh-group26"
	VPNPFSDHGroup27  = "dh-group27"
	VPNPFSDHGroup28  = "dh-group28"
	VPNPFSDHGroup29  = "dh-group29"
	VPNPFSDHGroup30  = "dh-group30"
	VPNPFSDHGroup31  = "dh-group31"
	VPNPFSDHGroup32  = "dh-group32"
)

// IPConfigurations contains network configuration of the VPN tunnel
// SubnetInfo contains subnet CIDR and name for VPN tunnel IP configuration
type SubnetInfo struct {
	CIDR string `json:"cidr,omitempty"`
	Name string `json:"name,omitempty"`
}

// IPConfigurations contains network configuration of the VPN tunnel
type IPConfigurations struct {
	// VPC reference to the VPC (nullable)
	VPC *ReferenceResource `json:"vpc,omitempty"`

	// Subnet info (nullable)
	Subnet *SubnetInfo `json:"subnet,omitempty"`

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

	// PeerClientPublicIP Peer client public IP address (nullable)
	PeerClientPublicIP *string `json:"peerClientPublicIp,omitempty"`
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

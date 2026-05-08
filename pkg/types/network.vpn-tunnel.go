package types

// IKEEncryption is the encryption algorithm for IKESettings.Encryption.
//
// GET /providers/Aruba.Network/vpn-tunnels for the live catalog.
type IKEEncryption string

const (
	IKEEncryptionAES128            IKEEncryption = "aes128"
	IKEEncryptionAES192            IKEEncryption = "aes192"
	IKEEncryptionAES256            IKEEncryption = "aes256"
	IKEEncryptionAES128CTR         IKEEncryption = "aes128ctr"
	IKEEncryptionAES192CTR         IKEEncryption = "aes192ctr"
	IKEEncryptionAES256CTR         IKEEncryption = "aes256ctr"
	IKEEncryptionAES128CCM64       IKEEncryption = "aes128ccm64"
	IKEEncryptionAES128CCM96       IKEEncryption = "aes128ccm96"
	IKEEncryptionAES128CCM128      IKEEncryption = "aes128ccm128"
	IKEEncryptionAES192CCM64       IKEEncryption = "aes192ccm64"
	IKEEncryptionAES192CCM96       IKEEncryption = "aes192ccm96"
	IKEEncryptionAES192CCM128      IKEEncryption = "aes192ccm128"
	IKEEncryptionAES256CCM64       IKEEncryption = "aes256ccm64"
	IKEEncryptionAES256CCM96       IKEEncryption = "aes256ccm96"
	IKEEncryptionAES256CCM128      IKEEncryption = "aes256ccm128"
	IKEEncryptionAES128GCM64       IKEEncryption = "aes128gcm64"
	IKEEncryptionAES128GCM96       IKEEncryption = "aes128gcm96"
	IKEEncryptionAES128GCM128      IKEEncryption = "aes128gcm128"
	IKEEncryptionAES192GCM64       IKEEncryption = "aes192gcm64"
	IKEEncryptionAES192GCM96       IKEEncryption = "aes192gcm96"
	IKEEncryptionAES192GCM128      IKEEncryption = "aes192gcm128"
	IKEEncryptionAES256GCM64       IKEEncryption = "aes256gcm64"
	IKEEncryptionAES256GCM96       IKEEncryption = "aes256gcm96"
	IKEEncryptionAES256GCM128      IKEEncryption = "aes256gcm128"
	IKEEncryptionAES128GMAC        IKEEncryption = "aes128gmac"
	IKEEncryptionAES192GMAC        IKEEncryption = "aes192gmac"
	IKEEncryptionAES256GMAC        IKEEncryption = "aes256gmac"
	IKEEncryption3DES              IKEEncryption = "3des"
	IKEEncryptionBlowfish128       IKEEncryption = "blowfish128"
	IKEEncryptionBlowfish192       IKEEncryption = "blowfish192"
	IKEEncryptionBlowfish256       IKEEncryption = "blowfish256"
	IKEEncryptionCamellia128       IKEEncryption = "camellia128"
	IKEEncryptionCamellia192       IKEEncryption = "camellia192"
	IKEEncryptionCamellia256       IKEEncryption = "camellia256"
	IKEEncryptionCamellia128CTR    IKEEncryption = "camellia128ctr"
	IKEEncryptionCamellia192CTR    IKEEncryption = "camellia192ctr"
	IKEEncryptionCamellia256CTR    IKEEncryption = "camellia256ctr"
	IKEEncryptionCamellia128CCM64  IKEEncryption = "camellia128ccm64"
	IKEEncryptionCamellia128CCM96  IKEEncryption = "camellia128ccm96"
	IKEEncryptionCamellia128CCM128 IKEEncryption = "camellia128ccm128"
	IKEEncryptionCamellia192CCM64  IKEEncryption = "camellia192ccm64"
	IKEEncryptionCamellia192CCM96  IKEEncryption = "camellia192ccm96"
	IKEEncryptionCamellia192CCM128 IKEEncryption = "camellia192ccm128"
	IKEEncryptionCamellia256CCM64  IKEEncryption = "camellia256ccm64"
	IKEEncryptionCamellia256CCM96  IKEEncryption = "camellia256ccm96"
	IKEEncryptionCamellia256CCM128 IKEEncryption = "camellia256ccm128"
	IKEEncryptionSerpent128        IKEEncryption = "serpent128"
	IKEEncryptionSerpent192        IKEEncryption = "serpent192"
	IKEEncryptionSerpent256        IKEEncryption = "serpent256"
	IKEEncryptionTwofish128        IKEEncryption = "twofish128"
	IKEEncryptionTwofish192        IKEEncryption = "twofish192"
	IKEEncryptionTwofish256        IKEEncryption = "twofish256"
	IKEEncryptionCAST128           IKEEncryption = "cast128"
	IKEEncryptionChaCha20Poly1305  IKEEncryption = "chacha20poly1305"
)

// IKEHash is the hash algorithm for IKESettings.Hash.
type IKEHash string

const (
	IKEHashMD5        IKEHash = "md5"
	IKEHashMD5128     IKEHash = "md5_128"
	IKEHashSHA1       IKEHash = "sha1"
	IKEHashSHA1160    IKEHash = "sha1_160"
	IKEHashSHA256     IKEHash = "sha256"
	IKEHashSHA25696   IKEHash = "sha256_96"
	IKEHashSHA384     IKEHash = "sha384"
	IKEHashSHA512     IKEHash = "sha512"
	IKEHashAESXCBC    IKEHash = "aesxcbc"
	IKEHashAESCMAC    IKEHash = "aescmac"
	IKEHashAES128GMAC IKEHash = "aes128gmac"
	IKEHashAES192GMAC IKEHash = "aes192gmac"
	IKEHashAES256GMAC IKEHash = "aes256gmac"
)

// IKEDHGroup is the Diffie-Hellman group for IKESettings.DHGroup.
type IKEDHGroup string

const (
	IKEDHGroup1  IKEDHGroup = "1"
	IKEDHGroup2  IKEDHGroup = "2"
	IKEDHGroup5  IKEDHGroup = "5"
	IKEDHGroup14 IKEDHGroup = "14"
	IKEDHGroup15 IKEDHGroup = "15"
	IKEDHGroup16 IKEDHGroup = "16"
	IKEDHGroup17 IKEDHGroup = "17"
	IKEDHGroup18 IKEDHGroup = "18"
	IKEDHGroup19 IKEDHGroup = "19"
	IKEDHGroup20 IKEDHGroup = "20"
	IKEDHGroup21 IKEDHGroup = "21"
	IKEDHGroup22 IKEDHGroup = "22"
	IKEDHGroup23 IKEDHGroup = "23"
	IKEDHGroup24 IKEDHGroup = "24"
	IKEDHGroup25 IKEDHGroup = "25"
	IKEDHGroup26 IKEDHGroup = "26"
	IKEDHGroup27 IKEDHGroup = "27"
	IKEDHGroup28 IKEDHGroup = "28"
	IKEDHGroup29 IKEDHGroup = "29"
	IKEDHGroup30 IKEDHGroup = "30"
	IKEDHGroup31 IKEDHGroup = "31"
	IKEDHGroup32 IKEDHGroup = "32"
)

// IKEDPDAction is the Dead Peer Detection action for IKESettings.DPDAction.
type IKEDPDAction string

const (
	IKEDPDActionTrap    IKEDPDAction = "trap"
	IKEDPDActionClear   IKEDPDAction = "clear"
	IKEDPDActionRestart IKEDPDAction = "restart"
)

// VPN PFS group constants for ESPSettings.PFS (replaced by ESPPFSGroup in a later commit).
const (
	VPNPFSEnable    = "enable"
	VPNPFSDisable   = "disable"
	VPNPFSDHGroup1  = "dh-group1"
	VPNPFSDHGroup2  = "dh-group2"
	VPNPFSDHGroup5  = "dh-group5"
	VPNPFSDHGroup14 = "dh-group14"
	VPNPFSDHGroup15 = "dh-group15"
	VPNPFSDHGroup16 = "dh-group16"
	VPNPFSDHGroup17 = "dh-group17"
	VPNPFSDHGroup18 = "dh-group18"
	VPNPFSDHGroup19 = "dh-group19"
	VPNPFSDHGroup20 = "dh-group20"
	VPNPFSDHGroup21 = "dh-group21"
	VPNPFSDHGroup22 = "dh-group22"
	VPNPFSDHGroup23 = "dh-group23"
	VPNPFSDHGroup24 = "dh-group24"
	VPNPFSDHGroup25 = "dh-group25"
	VPNPFSDHGroup26 = "dh-group26"
	VPNPFSDHGroup27 = "dh-group27"
	VPNPFSDHGroup28 = "dh-group28"
	VPNPFSDHGroup29 = "dh-group29"
	VPNPFSDHGroup30 = "dh-group30"
	VPNPFSDHGroup31 = "dh-group31"
	VPNPFSDHGroup32 = "dh-group32"
)

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
	Encryption *IKEEncryption `json:"encryption,omitempty"`

	// Hash Hash algorithm (nullable)
	Hash *IKEHash `json:"hash,omitempty"`

	// DHGroup Diffie-Hellman group (nullable)
	DHGroup *IKEDHGroup `json:"dhGroup,omitempty"`

	// DPDAction Dead Peer Detection action (nullable)
	DPDAction *IKEDPDAction `json:"dpdAction,omitempty"`

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

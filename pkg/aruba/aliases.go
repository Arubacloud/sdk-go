package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// ---------------------------------------------------------------------------
// Type aliases — callers need only import pkg/aruba.
// ---------------------------------------------------------------------------

// Enum type aliases (string-based discriminated types from pkg/types).
type (
	// Network
	RuleDirection   = types.RuleDirection
	RuleProtocol    = types.RuleProtocol
	EndpointTypeDto = types.EndpointTypeDto
	SubnetType      = types.SubnetType

	// Storage
	BlockStorageType  = types.BlockStorageType
	StorageBackupType = types.StorageBackupType

	// Container
	ContainerRegistrySizeFlavor = types.ContainerRegistrySizeFlavor

	// Security / KMS
	KeyAlgorithm      = types.KeyAlgorithm
	KeyCreationSource = types.KeyCreationSource
	KeyType           = types.KeyType
	KeyStatus         = types.KeyStatus
	ServiceStatus     = types.ServiceStatus

	// Schedule
	JobType           = types.JobType
	RecurrenceType    = types.RecurrenceType
	DeactiveReasonDto = types.DeactiveReasonDto

	// Metrics / Alerts
	ActionType = types.ActionType

	// Location / zone
	Region = types.Region
	Zone   = types.Zone

	// Billing
	BillingPeriod = types.BillingPeriod

	// Compute
	CloudServerFlavor = types.CloudServerFlavor

	// Database
	DatabaseEngine = types.DatabaseEngine
	DBaaSFlavor    = types.DBaaSFlavor

	// Container / KaaS
	KubernetesVersion = types.KubernetesVersion
	NodePoolInstance  = types.NodePoolInstance

	// Parameters
	AcceptHeader = types.AcceptHeader

	// VPN IKE crypto
	IKEEncryption = types.IKEEncryption
	IKEHash       = types.IKEHash
	IKEDHGroup    = types.IKEDHGroup
	IKEDPDAction  = types.IKEDPDAction

	// VPN ESP crypto
	ESPEncryption = types.ESPEncryption
	ESPHash       = types.ESPHash
	ESPPFSGroup   = types.ESPPFSGroup

	// VPN tunnel types
	VPNType           = types.VPNType
	VPNClientProtocol = types.VPNClientProtocol

	// Schedule
	HTTPVerb = types.HTTPVerb
)

// ---------------------------------------------------------------------------
// Const re-exports — callers use aruba.X without importing pkg/types.
// ---------------------------------------------------------------------------

const (
	// Network — security rule direction
	RuleDirectionIngress = types.RuleDirectionIngress
	RuleDirectionEgress  = types.RuleDirectionEgress

	// Network — security rule protocol
	RuleProtocolANY  = types.RuleProtocolANY
	RuleProtocolTCP  = types.RuleProtocolTCP
	RuleProtocolUDP  = types.RuleProtocolUDP
	RuleProtocolICMP = types.RuleProtocolICMP

	// Network — endpoint type
	EndpointTypeIP            = types.EndpointTypeIP
	EndpointTypeSecurityGroup = types.EndpointTypeSecurityGroup

	// Network — subnet type
	SubnetTypeBasic    = types.SubnetTypeBasic
	SubnetTypeAdvanced = types.SubnetTypeAdvanced

	// Storage — block storage type
	BlockStorageTypeStandard    = types.BlockStorageTypeStandard
	BlockStorageTypePerformance = types.BlockStorageTypePerformance

	// Storage — backup type
	StorageBackupTypeFull        = types.StorageBackupTypeFull
	StorageBackupTypeIncremental = types.StorageBackupTypeIncremental

	// Container — registry size flavor
	ContainerRegistrySizeFlavorSmall    = types.ContainerRegistrySizeFlavorSmall
	ContainerRegistrySizeFlavorMedium   = types.ContainerRegistrySizeFlavorMedium
	ContainerRegistrySizeFlavorHighPerf = types.ContainerRegistrySizeFlavorHighPerf

	// Security — key algorithm
	KeyAlgorithmAes = types.KeyAlgorithmAes
	KeyAlgorithmRsa = types.KeyAlgorithmRsa

	// Security — key creation source
	KeyCreationSourceCmp   = types.KeyCreationSourceCmp
	KeyCreationSourceOther = types.KeyCreationSourceOther

	// Security — key type
	KeyTypeSymmetric  = types.KeyTypeSymmetric
	KeyTypeAsymmetric = types.KeyTypeAsymmetric

	// Security — key status
	KeyStatusActive     = types.KeyStatusActive
	KeyStatusInCreation = types.KeyStatusInCreation
	KeyStatusDeleting   = types.KeyStatusDeleting
	KeyStatusDeleted    = types.KeyStatusDeleted
	KeyStatusFailed     = types.KeyStatusFailed

	// Security — KMIP service status
	ServiceStatusInCreation           = types.ServiceStatusInCreation
	ServiceStatusActive               = types.ServiceStatusActive
	ServiceStatusUpdating             = types.ServiceStatusUpdating
	ServiceStatusDeleting             = types.ServiceStatusDeleting
	ServiceStatusDeleted              = types.ServiceStatusDeleted
	ServiceStatusFailed               = types.ServiceStatusFailed
	ServiceStatusCertificateAvailable = types.ServiceStatusCertificateAvailable

	// Schedule — job type
	JobTypeOneShot   = types.JobTypeOneShot
	JobTypeRecurring = types.JobTypeRecurring

	// Schedule — recurrence type
	RecurrenceTypeHourly  = types.RecurrenceTypeHourly
	RecurrenceTypeDaily   = types.RecurrenceTypeDaily
	RecurrenceTypeWeekly  = types.RecurrenceTypeWeekly
	RecurrenceTypeMonthly = types.RecurrenceTypeMonthly
	RecurrenceTypeCustom  = types.RecurrenceTypeCustom

	// Schedule — deactivation reason
	DeactiveReasonNone            = types.DeactiveReasonNone
	DeactiveReasonManual          = types.DeactiveReasonManual
	DeactiveReasonResourceDeleted = types.DeactiveReasonResourceDeleted

	// Location
	RegionITBGBergamo = types.RegionITBGBergamo

	// Zone
	ZoneITBG1 = types.ZoneITBG1
	ZoneITBG2 = types.ZoneITBG2
	ZoneITBG3 = types.ZoneITBG3

	// Billing period
	BillingPeriodHour = types.BillingPeriodHour

	// Compute — cloud server flavor
	CloudServerFlavorCSO1A2   = types.CloudServerFlavorCSO1A2
	CloudServerFlavorCSO1A4   = types.CloudServerFlavorCSO1A4
	CloudServerFlavorCSO2A4   = types.CloudServerFlavorCSO2A4
	CloudServerFlavorCSO2A8   = types.CloudServerFlavorCSO2A8
	CloudServerFlavorCSO4A8   = types.CloudServerFlavorCSO4A8
	CloudServerFlavorCSO4A16  = types.CloudServerFlavorCSO4A16
	CloudServerFlavorCSO8A16  = types.CloudServerFlavorCSO8A16
	CloudServerFlavorCSO8A32  = types.CloudServerFlavorCSO8A32
	CloudServerFlavorCSO12A24 = types.CloudServerFlavorCSO12A24
	CloudServerFlavorCSO16A32 = types.CloudServerFlavorCSO16A32
	CloudServerFlavorCSO16A64 = types.CloudServerFlavorCSO16A64
	CloudServerFlavorCSO24A48 = types.CloudServerFlavorCSO24A48
	CloudServerFlavorCSO32A64 = types.CloudServerFlavorCSO32A64

	// Database — engine
	DatabaseEngineMySQL80    = types.DatabaseEngineMySQL80
	DatabaseEnginePostgres14 = types.DatabaseEnginePostgres14

	// Database — DBaaS flavor
	DBaaSFlavorDBO1A2   = types.DBaaSFlavorDBO1A2
	DBaaSFlavorDBO1A4   = types.DBaaSFlavorDBO1A4
	DBaaSFlavorDBO2A4   = types.DBaaSFlavorDBO2A4
	DBaaSFlavorDBO2A8   = types.DBaaSFlavorDBO2A8
	DBaaSFlavorDBO4A8   = types.DBaaSFlavorDBO4A8
	DBaaSFlavorDBO4A16  = types.DBaaSFlavorDBO4A16
	DBaaSFlavorDBO8A16  = types.DBaaSFlavorDBO8A16
	DBaaSFlavorDBO8A32  = types.DBaaSFlavorDBO8A32
	DBaaSFlavorDBO12A24 = types.DBaaSFlavorDBO12A24
	DBaaSFlavorDBO16A32 = types.DBaaSFlavorDBO16A32
	DBaaSFlavorDBO16A64 = types.DBaaSFlavorDBO16A64
	DBaaSFlavorDBO24A48 = types.DBaaSFlavorDBO24A48
	DBaaSFlavorDBO32A64 = types.DBaaSFlavorDBO32A64

	// Container — Kubernetes version
	KubernetesVersion1282 = types.KubernetesVersion1282
	KubernetesVersion1292 = types.KubernetesVersion1292
	KubernetesVersion1302 = types.KubernetesVersion1302
	KubernetesVersion1332 = types.KubernetesVersion1332

	// Container — node pool instance
	NodePoolInstanceK1A2   = types.NodePoolInstanceK1A2
	NodePoolInstanceK1A4R  = types.NodePoolInstanceK1A4R
	NodePoolInstanceK2A4   = types.NodePoolInstanceK2A4
	NodePoolInstanceK2A8R  = types.NodePoolInstanceK2A8R
	NodePoolInstanceK4A8   = types.NodePoolInstanceK4A8
	NodePoolInstanceK4A16R = types.NodePoolInstanceK4A16R
	NodePoolInstanceK8A16  = types.NodePoolInstanceK8A16
	NodePoolInstanceK8A32R = types.NodePoolInstanceK8A32R
	NodePoolInstanceK12A24 = types.NodePoolInstanceK12A24
	NodePoolInstanceK16A32 = types.NodePoolInstanceK16A32
	NodePoolInstanceK24A48 = types.NodePoolInstanceK24A48
	NodePoolInstanceK32A64 = types.NodePoolInstanceK32A64

	// Metrics / Alerts — action type
	ActionTypeNotificationPanel = types.ActionTypeNotificationPanel
	ActionTypeSendEmail         = types.ActionTypeSendEmail
	ActionTypeSendSMS           = types.ActionTypeSendSMS
	ActionTypeAutoscalingDBaaS  = types.ActionTypeAutoscalingDBaaS

	// VPN IKE — encryption algorithms (for IKESettings.Encryption)
	IKEEncryptionAES128            = types.IKEEncryptionAES128
	IKEEncryptionAES192            = types.IKEEncryptionAES192
	IKEEncryptionAES256            = types.IKEEncryptionAES256
	IKEEncryptionAES128CTR         = types.IKEEncryptionAES128CTR
	IKEEncryptionAES192CTR         = types.IKEEncryptionAES192CTR
	IKEEncryptionAES256CTR         = types.IKEEncryptionAES256CTR
	IKEEncryptionAES128CCM64       = types.IKEEncryptionAES128CCM64
	IKEEncryptionAES128CCM96       = types.IKEEncryptionAES128CCM96
	IKEEncryptionAES128CCM128      = types.IKEEncryptionAES128CCM128
	IKEEncryptionAES192CCM64       = types.IKEEncryptionAES192CCM64
	IKEEncryptionAES192CCM96       = types.IKEEncryptionAES192CCM96
	IKEEncryptionAES192CCM128      = types.IKEEncryptionAES192CCM128
	IKEEncryptionAES256CCM64       = types.IKEEncryptionAES256CCM64
	IKEEncryptionAES256CCM96       = types.IKEEncryptionAES256CCM96
	IKEEncryptionAES256CCM128      = types.IKEEncryptionAES256CCM128
	IKEEncryptionAES128GCM64       = types.IKEEncryptionAES128GCM64
	IKEEncryptionAES128GCM96       = types.IKEEncryptionAES128GCM96
	IKEEncryptionAES128GCM128      = types.IKEEncryptionAES128GCM128
	IKEEncryptionAES192GCM64       = types.IKEEncryptionAES192GCM64
	IKEEncryptionAES192GCM96       = types.IKEEncryptionAES192GCM96
	IKEEncryptionAES192GCM128      = types.IKEEncryptionAES192GCM128
	IKEEncryptionAES256GCM64       = types.IKEEncryptionAES256GCM64
	IKEEncryptionAES256GCM96       = types.IKEEncryptionAES256GCM96
	IKEEncryptionAES256GCM128      = types.IKEEncryptionAES256GCM128
	IKEEncryptionAES128GMAC        = types.IKEEncryptionAES128GMAC
	IKEEncryptionAES192GMAC        = types.IKEEncryptionAES192GMAC
	IKEEncryptionAES256GMAC        = types.IKEEncryptionAES256GMAC
	IKEEncryption3DES              = types.IKEEncryption3DES
	IKEEncryptionBlowfish128       = types.IKEEncryptionBlowfish128
	IKEEncryptionBlowfish192       = types.IKEEncryptionBlowfish192
	IKEEncryptionBlowfish256       = types.IKEEncryptionBlowfish256
	IKEEncryptionCamellia128       = types.IKEEncryptionCamellia128
	IKEEncryptionCamellia192       = types.IKEEncryptionCamellia192
	IKEEncryptionCamellia256       = types.IKEEncryptionCamellia256
	IKEEncryptionCamellia128CTR    = types.IKEEncryptionCamellia128CTR
	IKEEncryptionCamellia192CTR    = types.IKEEncryptionCamellia192CTR
	IKEEncryptionCamellia256CTR    = types.IKEEncryptionCamellia256CTR
	IKEEncryptionCamellia128CCM64  = types.IKEEncryptionCamellia128CCM64
	IKEEncryptionCamellia128CCM96  = types.IKEEncryptionCamellia128CCM96
	IKEEncryptionCamellia128CCM128 = types.IKEEncryptionCamellia128CCM128
	IKEEncryptionCamellia192CCM64  = types.IKEEncryptionCamellia192CCM64
	IKEEncryptionCamellia192CCM96  = types.IKEEncryptionCamellia192CCM96
	IKEEncryptionCamellia192CCM128 = types.IKEEncryptionCamellia192CCM128
	IKEEncryptionCamellia256CCM64  = types.IKEEncryptionCamellia256CCM64
	IKEEncryptionCamellia256CCM96  = types.IKEEncryptionCamellia256CCM96
	IKEEncryptionCamellia256CCM128 = types.IKEEncryptionCamellia256CCM128
	IKEEncryptionSerpent128        = types.IKEEncryptionSerpent128
	IKEEncryptionSerpent192        = types.IKEEncryptionSerpent192
	IKEEncryptionSerpent256        = types.IKEEncryptionSerpent256
	IKEEncryptionTwofish128        = types.IKEEncryptionTwofish128
	IKEEncryptionTwofish192        = types.IKEEncryptionTwofish192
	IKEEncryptionTwofish256        = types.IKEEncryptionTwofish256
	IKEEncryptionCAST128           = types.IKEEncryptionCAST128
	IKEEncryptionChaCha20Poly1305  = types.IKEEncryptionChaCha20Poly1305

	// VPN IKE — hash algorithms (for IKESettings.Hash)
	IKEHashMD5        = types.IKEHashMD5
	IKEHashMD5128     = types.IKEHashMD5128
	IKEHashSHA1       = types.IKEHashSHA1
	IKEHashSHA1160    = types.IKEHashSHA1160
	IKEHashSHA256     = types.IKEHashSHA256
	IKEHashSHA25696   = types.IKEHashSHA25696
	IKEHashSHA384     = types.IKEHashSHA384
	IKEHashSHA512     = types.IKEHashSHA512
	IKEHashAESXCBC    = types.IKEHashAESXCBC
	IKEHashAESCMAC    = types.IKEHashAESCMAC
	IKEHashAES128GMAC = types.IKEHashAES128GMAC
	IKEHashAES192GMAC = types.IKEHashAES192GMAC
	IKEHashAES256GMAC = types.IKEHashAES256GMAC

	// VPN IKE — DH groups (for IKESettings.DHGroup)
	IKEDHGroup1  = types.IKEDHGroup1
	IKEDHGroup2  = types.IKEDHGroup2
	IKEDHGroup5  = types.IKEDHGroup5
	IKEDHGroup14 = types.IKEDHGroup14
	IKEDHGroup15 = types.IKEDHGroup15
	IKEDHGroup16 = types.IKEDHGroup16
	IKEDHGroup17 = types.IKEDHGroup17
	IKEDHGroup18 = types.IKEDHGroup18
	IKEDHGroup19 = types.IKEDHGroup19
	IKEDHGroup20 = types.IKEDHGroup20
	IKEDHGroup21 = types.IKEDHGroup21
	IKEDHGroup22 = types.IKEDHGroup22
	IKEDHGroup23 = types.IKEDHGroup23
	IKEDHGroup24 = types.IKEDHGroup24
	IKEDHGroup25 = types.IKEDHGroup25
	IKEDHGroup26 = types.IKEDHGroup26
	IKEDHGroup27 = types.IKEDHGroup27
	IKEDHGroup28 = types.IKEDHGroup28
	IKEDHGroup29 = types.IKEDHGroup29
	IKEDHGroup30 = types.IKEDHGroup30
	IKEDHGroup31 = types.IKEDHGroup31
	IKEDHGroup32 = types.IKEDHGroup32

	// VPN IKE — DPD actions (for IKESettings.DPDAction)
	IKEDPDActionTrap    = types.IKEDPDActionTrap
	IKEDPDActionClear   = types.IKEDPDActionClear
	IKEDPDActionRestart = types.IKEDPDActionRestart

	// VPN ESP — encryption algorithms (for ESPSettings.Encryption)
	ESPEncryptionAES128            = types.ESPEncryptionAES128
	ESPEncryptionAES192            = types.ESPEncryptionAES192
	ESPEncryptionAES256            = types.ESPEncryptionAES256
	ESPEncryptionAES128CTR         = types.ESPEncryptionAES128CTR
	ESPEncryptionAES192CTR         = types.ESPEncryptionAES192CTR
	ESPEncryptionAES256CTR         = types.ESPEncryptionAES256CTR
	ESPEncryptionAES128CCM64       = types.ESPEncryptionAES128CCM64
	ESPEncryptionAES128CCM96       = types.ESPEncryptionAES128CCM96
	ESPEncryptionAES128CCM128      = types.ESPEncryptionAES128CCM128
	ESPEncryptionAES192CCM64       = types.ESPEncryptionAES192CCM64
	ESPEncryptionAES192CCM96       = types.ESPEncryptionAES192CCM96
	ESPEncryptionAES192CCM128      = types.ESPEncryptionAES192CCM128
	ESPEncryptionAES256CCM64       = types.ESPEncryptionAES256CCM64
	ESPEncryptionAES256CCM96       = types.ESPEncryptionAES256CCM96
	ESPEncryptionAES256CCM128      = types.ESPEncryptionAES256CCM128
	ESPEncryptionAES128GCM64       = types.ESPEncryptionAES128GCM64
	ESPEncryptionAES128GCM96       = types.ESPEncryptionAES128GCM96
	ESPEncryptionAES128GCM128      = types.ESPEncryptionAES128GCM128
	ESPEncryptionAES192GCM64       = types.ESPEncryptionAES192GCM64
	ESPEncryptionAES192GCM96       = types.ESPEncryptionAES192GCM96
	ESPEncryptionAES192GCM128      = types.ESPEncryptionAES192GCM128
	ESPEncryptionAES256GCM64       = types.ESPEncryptionAES256GCM64
	ESPEncryptionAES256GCM96       = types.ESPEncryptionAES256GCM96
	ESPEncryptionAES256GCM128      = types.ESPEncryptionAES256GCM128
	ESPEncryptionAES128GMAC        = types.ESPEncryptionAES128GMAC
	ESPEncryptionAES192GMAC        = types.ESPEncryptionAES192GMAC
	ESPEncryptionAES256GMAC        = types.ESPEncryptionAES256GMAC
	ESPEncryption3DES              = types.ESPEncryption3DES
	ESPEncryptionBlowfish128       = types.ESPEncryptionBlowfish128
	ESPEncryptionBlowfish192       = types.ESPEncryptionBlowfish192
	ESPEncryptionBlowfish256       = types.ESPEncryptionBlowfish256
	ESPEncryptionCamellia128       = types.ESPEncryptionCamellia128
	ESPEncryptionCamellia192       = types.ESPEncryptionCamellia192
	ESPEncryptionCamellia256       = types.ESPEncryptionCamellia256
	ESPEncryptionCamellia128CTR    = types.ESPEncryptionCamellia128CTR
	ESPEncryptionCamellia192CTR    = types.ESPEncryptionCamellia192CTR
	ESPEncryptionCamellia256CTR    = types.ESPEncryptionCamellia256CTR
	ESPEncryptionCamellia128CCM64  = types.ESPEncryptionCamellia128CCM64
	ESPEncryptionCamellia128CCM96  = types.ESPEncryptionCamellia128CCM96
	ESPEncryptionCamellia128CCM128 = types.ESPEncryptionCamellia128CCM128
	ESPEncryptionCamellia192CCM64  = types.ESPEncryptionCamellia192CCM64
	ESPEncryptionCamellia192CCM96  = types.ESPEncryptionCamellia192CCM96
	ESPEncryptionCamellia192CCM128 = types.ESPEncryptionCamellia192CCM128
	ESPEncryptionCamellia256CCM64  = types.ESPEncryptionCamellia256CCM64
	ESPEncryptionCamellia256CCM96  = types.ESPEncryptionCamellia256CCM96
	ESPEncryptionCamellia256CCM128 = types.ESPEncryptionCamellia256CCM128
	ESPEncryptionSerpent128        = types.ESPEncryptionSerpent128
	ESPEncryptionSerpent192        = types.ESPEncryptionSerpent192
	ESPEncryptionSerpent256        = types.ESPEncryptionSerpent256
	ESPEncryptionTwofish128        = types.ESPEncryptionTwofish128
	ESPEncryptionTwofish192        = types.ESPEncryptionTwofish192
	ESPEncryptionTwofish256        = types.ESPEncryptionTwofish256
	ESPEncryptionCAST128           = types.ESPEncryptionCAST128
	ESPEncryptionChaCha20Poly1305  = types.ESPEncryptionChaCha20Poly1305

	// VPN ESP — hash algorithms (for ESPSettings.Hash)
	ESPHashMD5        = types.ESPHashMD5
	ESPHashMD5128     = types.ESPHashMD5128
	ESPHashSHA1       = types.ESPHashSHA1
	ESPHashSHA1160    = types.ESPHashSHA1160
	ESPHashSHA256     = types.ESPHashSHA256
	ESPHashSHA25696   = types.ESPHashSHA25696
	ESPHashSHA384     = types.ESPHashSHA384
	ESPHashSHA512     = types.ESPHashSHA512
	ESPHashAESXCBC    = types.ESPHashAESXCBC
	ESPHashAESCMAC    = types.ESPHashAESCMAC
	ESPHashAES128GMAC = types.ESPHashAES128GMAC
	ESPHashAES192GMAC = types.ESPHashAES192GMAC
	ESPHashAES256GMAC = types.ESPHashAES256GMAC

	// VPN ESP — PFS groups (for ESPSettings.PFS)
	ESPPFSGroupEnable    = types.ESPPFSGroupEnable
	ESPPFSGroupDisable   = types.ESPPFSGroupDisable
	ESPPFSGroupDHGroup1  = types.ESPPFSGroupDHGroup1
	ESPPFSGroupDHGroup2  = types.ESPPFSGroupDHGroup2
	ESPPFSGroupDHGroup5  = types.ESPPFSGroupDHGroup5
	ESPPFSGroupDHGroup14 = types.ESPPFSGroupDHGroup14
	ESPPFSGroupDHGroup15 = types.ESPPFSGroupDHGroup15
	ESPPFSGroupDHGroup16 = types.ESPPFSGroupDHGroup16
	ESPPFSGroupDHGroup17 = types.ESPPFSGroupDHGroup17
	ESPPFSGroupDHGroup18 = types.ESPPFSGroupDHGroup18
	ESPPFSGroupDHGroup19 = types.ESPPFSGroupDHGroup19
	ESPPFSGroupDHGroup20 = types.ESPPFSGroupDHGroup20
	ESPPFSGroupDHGroup21 = types.ESPPFSGroupDHGroup21
	ESPPFSGroupDHGroup22 = types.ESPPFSGroupDHGroup22
	ESPPFSGroupDHGroup23 = types.ESPPFSGroupDHGroup23
	ESPPFSGroupDHGroup24 = types.ESPPFSGroupDHGroup24
	ESPPFSGroupDHGroup25 = types.ESPPFSGroupDHGroup25
	ESPPFSGroupDHGroup26 = types.ESPPFSGroupDHGroup26
	ESPPFSGroupDHGroup27 = types.ESPPFSGroupDHGroup27
	ESPPFSGroupDHGroup28 = types.ESPPFSGroupDHGroup28
	ESPPFSGroupDHGroup29 = types.ESPPFSGroupDHGroup29
	ESPPFSGroupDHGroup30 = types.ESPPFSGroupDHGroup30
	ESPPFSGroupDHGroup31 = types.ESPPFSGroupDHGroup31
	ESPPFSGroupDHGroup32 = types.ESPPFSGroupDHGroup32

	// VPN tunnel — type and client protocol
	VPNTypeSiteToSite      = types.VPNTypeSiteToSite
	VPNClientProtocolIKEv2 = types.VPNClientProtocolIKEv2

	// Schedule — HTTP verb
	HTTPVerbGET    = types.HTTPVerbGET
	HTTPVerbPOST   = types.HTTPVerbPOST
	HTTPVerbPUT    = types.HTTPVerbPUT
	HTTPVerbDELETE = types.HTTPVerbDELETE
	HTTPVerbPATCH  = types.HTTPVerbPATCH
)

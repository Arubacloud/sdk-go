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

	// VPN — encryption algorithms (for IKESettings.Encryption and ESPSettings.Encryption)
	VPNEncryptionAES128            = types.VPNEncryptionAES128
	VPNEncryptionAES192            = types.VPNEncryptionAES192
	VPNEncryptionAES256            = types.VPNEncryptionAES256
	VPNEncryptionAES128CTR         = types.VPNEncryptionAES128CTR
	VPNEncryptionAES192CTR         = types.VPNEncryptionAES192CTR
	VPNEncryptionAES256CTR         = types.VPNEncryptionAES256CTR
	VPNEncryptionAES128CCM64       = types.VPNEncryptionAES128CCM64
	VPNEncryptionAES128CCM96       = types.VPNEncryptionAES128CCM96
	VPNEncryptionAES128CCM128      = types.VPNEncryptionAES128CCM128
	VPNEncryptionAES192CCM64       = types.VPNEncryptionAES192CCM64
	VPNEncryptionAES192CCM96       = types.VPNEncryptionAES192CCM96
	VPNEncryptionAES192CCM128      = types.VPNEncryptionAES192CCM128
	VPNEncryptionAES256CCM64       = types.VPNEncryptionAES256CCM64
	VPNEncryptionAES256CCM96       = types.VPNEncryptionAES256CCM96
	VPNEncryptionAES256CCM128      = types.VPNEncryptionAES256CCM128
	VPNEncryptionAES128GCM64       = types.VPNEncryptionAES128GCM64
	VPNEncryptionAES128GCM96       = types.VPNEncryptionAES128GCM96
	VPNEncryptionAES128GCM128      = types.VPNEncryptionAES128GCM128
	VPNEncryptionAES192GCM64       = types.VPNEncryptionAES192GCM64
	VPNEncryptionAES192GCM96       = types.VPNEncryptionAES192GCM96
	VPNEncryptionAES192GCM128      = types.VPNEncryptionAES192GCM128
	VPNEncryptionAES256GCM64       = types.VPNEncryptionAES256GCM64
	VPNEncryptionAES256GCM96       = types.VPNEncryptionAES256GCM96
	VPNEncryptionAES256GCM128      = types.VPNEncryptionAES256GCM128
	VPNEncryptionAES128GMAC        = types.VPNEncryptionAES128GMAC
	VPNEncryptionAES192GMAC        = types.VPNEncryptionAES192GMAC
	VPNEncryptionAES256GMAC        = types.VPNEncryptionAES256GMAC
	VPNEncryption3DES              = types.VPNEncryption3DES
	VPNEncryptionBlowfish128       = types.VPNEncryptionBlowfish128
	VPNEncryptionBlowfish192       = types.VPNEncryptionBlowfish192
	VPNEncryptionBlowfish256       = types.VPNEncryptionBlowfish256
	VPNEncryptionCamellia128       = types.VPNEncryptionCamellia128
	VPNEncryptionCamellia192       = types.VPNEncryptionCamellia192
	VPNEncryptionCamellia256       = types.VPNEncryptionCamellia256
	VPNEncryptionCamellia128CTR    = types.VPNEncryptionCamellia128CTR
	VPNEncryptionCamellia192CTR    = types.VPNEncryptionCamellia192CTR
	VPNEncryptionCamellia256CTR    = types.VPNEncryptionCamellia256CTR
	VPNEncryptionCamellia128CCM64  = types.VPNEncryptionCamellia128CCM64
	VPNEncryptionCamellia128CCM96  = types.VPNEncryptionCamellia128CCM96
	VPNEncryptionCamellia128CCM128 = types.VPNEncryptionCamellia128CCM128
	VPNEncryptionCamellia192CCM64  = types.VPNEncryptionCamellia192CCM64
	VPNEncryptionCamellia192CCM96  = types.VPNEncryptionCamellia192CCM96
	VPNEncryptionCamellia192CCM128 = types.VPNEncryptionCamellia192CCM128
	VPNEncryptionCamellia256CCM64  = types.VPNEncryptionCamellia256CCM64
	VPNEncryptionCamellia256CCM96  = types.VPNEncryptionCamellia256CCM96
	VPNEncryptionCamellia256CCM128 = types.VPNEncryptionCamellia256CCM128
	VPNEncryptionSerpent128        = types.VPNEncryptionSerpent128
	VPNEncryptionSerpent192        = types.VPNEncryptionSerpent192
	VPNEncryptionSerpent256        = types.VPNEncryptionSerpent256
	VPNEncryptionTwofish128        = types.VPNEncryptionTwofish128
	VPNEncryptionTwofish192        = types.VPNEncryptionTwofish192
	VPNEncryptionTwofish256        = types.VPNEncryptionTwofish256
	VPNEncryptionCAST128           = types.VPNEncryptionCAST128
	VPNEncryptionChaCha20Poly1305  = types.VPNEncryptionChaCha20Poly1305

	// VPN — hash algorithms (for IKESettings.Hash and ESPSettings.Hash)
	VPNHashMD5        = types.VPNHashMD5
	VPNHashMD5128     = types.VPNHashMD5128
	VPNHashSHA1       = types.VPNHashSHA1
	VPNHashSHA1160    = types.VPNHashSHA1160
	VPNHashSHA256     = types.VPNHashSHA256
	VPNHashSHA25696   = types.VPNHashSHA25696
	VPNHashSHA384     = types.VPNHashSHA384
	VPNHashSHA512     = types.VPNHashSHA512
	VPNHashAESXCBC    = types.VPNHashAESXCBC
	VPNHashAESCMAC    = types.VPNHashAESCMAC
	VPNHashAES128GMAC = types.VPNHashAES128GMAC
	VPNHashAES192GMAC = types.VPNHashAES192GMAC
	VPNHashAES256GMAC = types.VPNHashAES256GMAC

	// VPN — DH groups (for IKESettings.DHGroup)
	VPNDHGroup1  = types.VPNDHGroup1
	VPNDHGroup2  = types.VPNDHGroup2
	VPNDHGroup5  = types.VPNDHGroup5
	VPNDHGroup14 = types.VPNDHGroup14
	VPNDHGroup15 = types.VPNDHGroup15
	VPNDHGroup16 = types.VPNDHGroup16
	VPNDHGroup17 = types.VPNDHGroup17
	VPNDHGroup18 = types.VPNDHGroup18
	VPNDHGroup19 = types.VPNDHGroup19
	VPNDHGroup20 = types.VPNDHGroup20
	VPNDHGroup21 = types.VPNDHGroup21
	VPNDHGroup22 = types.VPNDHGroup22
	VPNDHGroup23 = types.VPNDHGroup23
	VPNDHGroup24 = types.VPNDHGroup24
	VPNDHGroup25 = types.VPNDHGroup25
	VPNDHGroup26 = types.VPNDHGroup26
	VPNDHGroup27 = types.VPNDHGroup27
	VPNDHGroup28 = types.VPNDHGroup28
	VPNDHGroup29 = types.VPNDHGroup29
	VPNDHGroup30 = types.VPNDHGroup30
	VPNDHGroup31 = types.VPNDHGroup31
	VPNDHGroup32 = types.VPNDHGroup32

	// VPN — DPD actions (for IKESettings.DPDAction)
	VPNDPDActionTrap    = types.VPNDPDActionTrap
	VPNDPDActionClear   = types.VPNDPDActionClear
	VPNDPDActionRestart = types.VPNDPDActionRestart

	// VPN — PFS groups (for ESPSettings.PFS)
	VPNPFSEnable    = types.VPNPFSEnable
	VPNPFSDisable   = types.VPNPFSDisable
	VPNPFSDHGroup1  = types.VPNPFSDHGroup1
	VPNPFSDHGroup2  = types.VPNPFSDHGroup2
	VPNPFSDHGroup5  = types.VPNPFSDHGroup5
	VPNPFSDHGroup14 = types.VPNPFSDHGroup14
	VPNPFSDHGroup15 = types.VPNPFSDHGroup15
	VPNPFSDHGroup16 = types.VPNPFSDHGroup16
	VPNPFSDHGroup17 = types.VPNPFSDHGroup17
	VPNPFSDHGroup18 = types.VPNPFSDHGroup18
	VPNPFSDHGroup19 = types.VPNPFSDHGroup19
	VPNPFSDHGroup20 = types.VPNPFSDHGroup20
	VPNPFSDHGroup21 = types.VPNPFSDHGroup21
	VPNPFSDHGroup22 = types.VPNPFSDHGroup22
	VPNPFSDHGroup23 = types.VPNPFSDHGroup23
	VPNPFSDHGroup24 = types.VPNPFSDHGroup24
	VPNPFSDHGroup25 = types.VPNPFSDHGroup25
	VPNPFSDHGroup26 = types.VPNPFSDHGroup26
	VPNPFSDHGroup27 = types.VPNPFSDHGroup27
	VPNPFSDHGroup28 = types.VPNPFSDHGroup28
	VPNPFSDHGroup29 = types.VPNPFSDHGroup29
	VPNPFSDHGroup30 = types.VPNPFSDHGroup30
	VPNPFSDHGroup31 = types.VPNPFSDHGroup31
	VPNPFSDHGroup32 = types.VPNPFSDHGroup32
)

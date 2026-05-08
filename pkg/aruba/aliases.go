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

	// VPN — PFS groups (for ESPSettings.PFS; typed ESPPFSGroup* consts coming in next commit)
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

package security

// API path constants for security resources
const (
	// KMS paths
	KMSsPath = "/projects/%s/providers/Aruba.Security/kms"
	KMSPath  = "/projects/%s/providers/Aruba.Security/kms/%s"

	// KMIP paths (nested under KMS)
	KmipsPath        = "/projects/%s/providers/Aruba.Security/kms/%s/kmip"
	KmipPath         = "/projects/%s/providers/Aruba.Security/kms/%s/kmip/%s"
	KmipDownloadPath = "/projects/%s/providers/Aruba.Security/kms/%s/kmip/%s/download"

	// Key paths (nested under KMS)
	KeysPath = "/projects/%s/providers/Aruba.Security/kms/%s/keys"
	KeyPath  = "/projects/%s/providers/Aruba.Security/kms/%s/keys/%s"
)

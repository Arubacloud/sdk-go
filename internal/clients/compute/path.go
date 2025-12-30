package compute

// API path constants for compute resources
const (
	// CloudServer paths
	CloudServersPath        = "/projects/%s/providers/Aruba.Compute/cloudServers"
	CloudServerPath         = "/projects/%s/providers/Aruba.Compute/cloudServers/%s"
	CloudServerPowerOnPath  = "/projects/%s/providers/Aruba.Compute/cloudServers/%s/poweron"
	CloudServerPowerOffPath = "/projects/%s/providers/Aruba.Compute/cloudServers/%s/poweroff"
	CloudServerPasswordPath = "/projects/%s/providers/Aruba.Compute/cloudServers/%s/password"

	// KeyPair paths
	KeyPairsPath = "/projects/%s/providers/Aruba.Compute/keypairs"
	KeyPairPath  = "/projects/%s/providers/Aruba.Compute/keypairs/%s"
)

package aruba

// TestCompileTimeInterfaceGuards verifies at compile time that every
// internal client implementation satisfies its declared public interface.
// This is a pure compile-time check: if an impl drops a method or changes
// a signature, the test package will fail to build before any test runs.
//
// Placed in a test file so that (a) no constructors are called in
// production binaries, and (b) no allocations persist in production memory.
//
// Guards are in package aruba (not aruba_test) because builder.go already
// imports every internal client package; a _test.go file in the same package
// inherits those imports without introducing a new cycle.
//
// The security domain is intentionally omitted: KMSClient, KeyClient, and
// KmipClient are type aliases to concrete pointer types, not interfaces.

import (
	"testing"

	"github.com/Arubacloud/sdk-go/internal/clients/audit"
	"github.com/Arubacloud/sdk-go/internal/clients/container"
	"github.com/Arubacloud/sdk-go/internal/clients/metric"
	"github.com/Arubacloud/sdk-go/internal/clients/schedule"
)

func TestCompileTimeInterfaceGuards(_ *testing.T) {
	var (
		// Audit
		_ EventsClient = audit.NewEventsClientImpl(nil)

		// Compute
		_ CloudServersClient = newCloudServersClientAdapter(nil)
		_ KeyPairsClient     = newKeyPairsClientAdapter(nil)

		// Container
		_ KaaSClient              = container.NewKaaSClientImpl(nil)
		_ ContainerRegistryClient = newContainerRegistriesClientAdapter(nil)

		// Database
		_ DBaaSClient     = newDBaaSClientAdapter(nil)
		_ DatabasesClient = newDatabasesClientAdapter(nil)
		_ BackupsClient   = newDBaaSBackupsClientAdapter(nil)
		_ UsersClient     = newUsersClientAdapter(nil)
		_ GrantsClient    = newGrantsClientAdapter(nil)

		// Metric
		_ AlertsClient  = metric.NewAlertsClientImpl(nil)
		_ MetricsClient = metric.NewMetricsClientImpl(nil)

		// Network
		_ ElasticIPsClient         = newElasticIPsClientAdapter(nil)
		_ LoadBalancersClient      = newLoadBalancersClientAdapter(nil)
		_ VPCsClient               = newVPCsClientAdapter(nil)
		_ SecurityGroupsClient     = newSecurityGroupsClientAdapter(nil)
		_ SecurityGroupRulesClient = newSecurityGroupRulesClientAdapter(nil)
		_ SubnetsClient            = newSubnetsClientAdapter(nil)
		_ VPCPeeringsClient        = newVPCPeeringsClientAdapter(nil)
		_ VPCPeeringRoutesClient   = newVPCPeeringRoutesClientAdapter(nil)
		_ VPNRoutesClient          = newVPNRoutesClientAdapter(nil)
		_ VPNTunnelsClient         = newVPNTunnelsClientAdapter(nil)

		// Project
		_ ProjectClient = newProjectClientAdapter(nil)

		// Schedule
		_ JobsClient = schedule.NewJobsClientImpl(nil)

		// Storage
		_ VolumesClient        = newVolumesClientAdapter(nil)
		_ SnapshotsClient      = newSnapshotsClientAdapter(nil)
		_ StorageBackupsClient = newStorageBackupsClientAdapter(nil)
		_ StorageRestoreClient = newStorageRestoresClientAdapter(nil)
	)
}

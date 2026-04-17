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
	"github.com/Arubacloud/sdk-go/internal/clients/compute"
	"github.com/Arubacloud/sdk-go/internal/clients/container"
	"github.com/Arubacloud/sdk-go/internal/clients/database"
	"github.com/Arubacloud/sdk-go/internal/clients/metric"
	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/clients/project"
	"github.com/Arubacloud/sdk-go/internal/clients/schedule"
	"github.com/Arubacloud/sdk-go/internal/clients/storage"
)

func TestCompileTimeInterfaceGuards(_ *testing.T) {
	// Local variables — stack-allocated, zero production overhead.
	// Named intermediates are needed only to satisfy the nil-dep checks
	// added by TD-018; they exist only for the duration of this call.
	vpcsImpl := network.NewVPCsClientImpl(nil)
	sgImpl := network.NewSecurityGroupsClientImpl(nil, vpcsImpl)
	volImpl := storage.NewVolumesClientImpl(nil)
	bkpImpl := storage.NewBackupClientImpl(nil)

	var (
		// Audit
		_ EventsClient = audit.NewEventsClientImpl(nil)

		// Compute
		_ CloudServersClient = compute.NewCloudServersClientImpl(nil)
		_ KeyPairsClient     = compute.NewKeyPairsClientImpl(nil)

		// Container
		_ KaaSClient              = container.NewKaaSClientImpl(nil)
		_ ContainerRegistryClient = container.NewContainerRegistryClientImpl(nil)

		// Database
		_ DBaaSClient     = database.NewDBaaSClientImpl(nil)
		_ DatabasesClient = database.NewDatabasesClientImpl(nil)
		_ BackupsClient   = database.NewBackupsClientImpl(nil)
		_ UsersClient     = database.NewUsersClientImpl(nil)
		_ GrantsClient    = database.NewGrantsClientImpl(nil)

		// Metric
		_ AlertsClient  = metric.NewAlertsClientImpl(nil)
		_ MetricsClient = metric.NewMetricsClientImpl(nil)

		// Network
		_ ElasticIPsClient         = network.NewElasticIPsClientImpl(nil)
		_ LoadBalancersClient      = network.NewLoadBalancersClientImpl(nil)
		_ VPCsClient               = vpcsImpl
		_ SecurityGroupsClient     = sgImpl
		_ SecurityGroupRulesClient = network.NewSecurityGroupRulesClientImpl(nil, sgImpl)
		_ SubnetsClient            = network.NewSubnetsClientImpl(nil, vpcsImpl)
		_ VPCPeeringsClient        = network.NewVPCPeeringsClientImpl(nil)
		_ VPCPeeringRoutesClient   = network.NewVPCPeeringRoutesClientImpl(nil)
		_ VPNRoutesClient          = network.NewVPNRoutesClientImpl(nil)
		_ VPNTunnelsClient         = network.NewVPNTunnelsClientImpl(nil)

		// Project
		_ ProjectClient = project.NewProjectsClientImpl(nil)

		// Schedule
		_ JobsClient = schedule.NewJobsClientImpl(nil)

		// Storage
		_ VolumesClient        = volImpl
		_ SnapshotsClient      = storage.NewSnapshotsClientImpl(nil, volImpl)
		_ StorageBackupsClient = bkpImpl
		_ StorageRestoreClient = storage.NewRestoreClientImpl(nil, bkpImpl)
	)
}

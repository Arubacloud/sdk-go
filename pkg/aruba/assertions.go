package aruba

// Compile-time assertions that internal client implementations satisfy their
// public interfaces. A build failure here means an impl no longer implements
// its declared interface — surfaced at compile time rather than at runtime.
//
// Guards are placed here (pkg/aruba) rather than in the internal packages
// because pkg/aruba/builder.go already imports every internal client package;
// the reverse import would create a cycle.
//
// The security domain is intentionally omitted: KMSClient, KeyClient, and
// KmipClient are declared as type aliases to concrete pointer types, not
// interfaces, so satisfaction guards are degenerate.

import (
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

// Sentinel instances wired with a nil restclient. Used only to satisfy the
// multi-arg constructors that require non-nil dependency pointers (TD-018).
var (
	assertVPCsImpl = network.NewVPCsClientImpl(nil)
	assertSGImpl   = network.NewSecurityGroupsClientImpl(nil, assertVPCsImpl)
	assertVolImpl  = storage.NewVolumesClientImpl(nil)
	assertBkpImpl  = storage.NewBackupClientImpl(nil)
)

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
	_ VPCsClient               = assertVPCsImpl
	_ SecurityGroupsClient     = assertSGImpl
	_ SecurityGroupRulesClient = network.NewSecurityGroupRulesClientImpl(nil, assertSGImpl)
	_ SubnetsClient            = network.NewSubnetsClientImpl(nil, assertVPCsImpl)
	_ VPCPeeringsClient        = network.NewVPCPeeringsClientImpl(nil)
	_ VPCPeeringRoutesClient   = network.NewVPCPeeringRoutesClientImpl(nil)
	_ VPNRoutesClient          = network.NewVPNRoutesClientImpl(nil)
	_ VPNTunnelsClient         = network.NewVPNTunnelsClientImpl(nil)

	// Project
	_ ProjectClient = project.NewProjectsClientImpl(nil)

	// Schedule
	_ JobsClient = schedule.NewJobsClientImpl(nil)

	// Storage
	_ VolumesClient        = assertVolImpl
	_ SnapshotsClient      = storage.NewSnapshotsClientImpl(nil, assertVolImpl)
	_ StorageBackupsClient = assertBkpImpl
	_ StorageRestoreClient = storage.NewRestoreClientImpl(nil, assertBkpImpl)
)

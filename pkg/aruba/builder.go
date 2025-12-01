//nolint:unparam // TODO: better error handling
package aruba

import (
	"log"

	"github.com/Arubacloud/sdk-go/internal/clients/audit"
	"github.com/Arubacloud/sdk-go/internal/clients/compute"
	"github.com/Arubacloud/sdk-go/internal/clients/container"
	"github.com/Arubacloud/sdk-go/internal/clients/database"
	"github.com/Arubacloud/sdk-go/internal/clients/metric"
	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/clients/project"
	"github.com/Arubacloud/sdk-go/internal/clients/schedule"
	"github.com/Arubacloud/sdk-go/internal/clients/security"
	"github.com/Arubacloud/sdk-go/internal/clients/storage"
	std "github.com/Arubacloud/sdk-go/internal/impl/auth/tokenmanager/standard"
	"github.com/Arubacloud/sdk-go/internal/impl/auth/tokenrepository/file"
	"github.com/Arubacloud/sdk-go/internal/impl/auth/tokenrepository/redis"
	"github.com/Arubacloud/sdk-go/internal/impl/interceptor/standard"
	"github.com/Arubacloud/sdk-go/internal/impl/logger/native"
	"github.com/Arubacloud/sdk-go/internal/impl/logger/noop"
	"github.com/Arubacloud/sdk-go/internal/ports/auth"
	"github.com/Arubacloud/sdk-go/internal/ports/interceptor"
	"github.com/Arubacloud/sdk-go/internal/ports/logger"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	redis_client "github.com/redis/go-redis/v9"
)

// Client
func buildClient(config *restclient.Config) (Client, error) {
	restClient, err := buildRESTClient(config)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	tokenManager, err := buildTokenManager(config)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	restClient.Logger().Debugf("Using Token Manager: %T", tokenManager)

	auditClient, err := buildAuditClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	computeClient, err := buildComputeClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	containerClient, err := buildContainerClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	databaseClient, err := buildDetebaseClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	metricClient, err := buildMetricClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	networkClient, err := buildNetworkClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	projectClient, err := buildProjectClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	scheduleClient, err := buildScheduleClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	securityClient, err := buildSecurityClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	storageClient, err := buildStorageClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &clientImpl{
		auditClient:     auditClient,
		computeClient:   computeClient,
		containerClient: containerClient,
		databaseClient:  databaseClient,
		metricsClient:   metricClient,
		networkClient:   networkClient,
		projectClient:   projectClient,
		scheduleClient:  scheduleClient,
		securityClient:  securityClient,
		storageClient:   storageClient,
	}, nil
}

//
// Dependencies

func buildRESTClient(config *restclient.Config) (*restclient.Client, error) {
	logger, err := buildLogger(config)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	middleware, err := buildMiddleware(config)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return restclient.NewClient(config, logger, middleware)
}

func buildLogger(config *restclient.Config) (logger.Logger, error) {
	if config.Logger != nil {
		return config.Logger, nil
	}

	if config.Debug {
		return native.NewDefaultLogger(), nil
	}

	return &noop.NoOpLogger{}, nil
}

func buildMiddleware(config *restclient.Config) (interceptor.Interceptor, error) {
	return standard.NewInterceptor(), nil
}

//
// Token Manager

func buildTokenManager(config *restclient.Config) (*std.TokenManager, error) {

	var tr auth.TokenRepository
	if config.Redis != nil {
		opt, err := redis_client.ParseURL(config.Redis.RedisURI)

		if err != nil {
			log.Fatal("Cannot parse Redis URI", err)
		}

		rdb := redis_client.NewClient(opt)

		tr = redis.NeWRedisTokenRepository(config.ClientID, rdb)

	} else if config.File != nil {
		tr = file.NeWFileTokenRepository(config.File.BaseDir, config.ClientID)
	} else {
		tr = nil
	}

	tm := std.NewTokenManager(tr, nil)

	return tm, nil

}

//
// Audit domain clients

func buildAuditClient(restClient *restclient.Client) (AuditClient, error) {
	eventsClient, err := buildEventsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &auditClientImpl{eventsClient: eventsClient}, nil
}

func buildEventsClient(restClient *restclient.Client) (EventsClient, error) {
	return audit.NewEventsClientImpl(restClient), nil
}

//
// Compute domain clients

func buildComputeClient(restClient *restclient.Client) (ComputeClient, error) {
	cloudServerClient, err := buildCloudServersClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	keyPairClient, err := buildKeyPairsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &computeClientImpl{
		cloudServerClient: cloudServerClient,
		keyPairClient:     keyPairClient,
	}, nil
}

func buildCloudServersClient(restClient *restclient.Client) (CloudServersClient, error) {
	return compute.NewCloudServersClientImpl(restClient), nil
}

func buildKeyPairsClient(restClient *restclient.Client) (KeyPairsClient, error) {
	return compute.NewKeyPairsClientImpl(restClient), nil
}

//
// Container domain clients

func buildContainerClient(restClient *restclient.Client) (ContainerClient, error) {
	kaasClient, err := buildKaaSClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &containerClientImpl{kaasClient: kaasClient}, nil
}

func buildKaaSClient(restClient *restclient.Client) (KaaSClient, error) {
	return container.NewKaaSClientImpl(restClient), nil
}

//
// Database domain clients

func buildDetebaseClient(restClient *restclient.Client) (DatabaseClient, error) {
	dbaasClient, err := buildDBaaSClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	databasesClient, err := buildDatabasesClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	backupsClient, err := buildBackupsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	usersClient, err := buildUsersClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	grantsClient, err := buildGrantsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &databaseClientImpl{
		dbaasClient:     dbaasClient,
		databasesClient: databasesClient,
		backupsClient:   backupsClient,
		usersClient:     usersClient,
		grantsClient:    grantsClient,
	}, nil
}

func buildDBaaSClient(restClient *restclient.Client) (DBaaSClient, error) {
	return database.NewDBaaSClientImpl(restClient), nil
}

func buildDatabasesClient(restClient *restclient.Client) (DatabasesClient, error) {
	return database.NewDatabasesClientImpl(restClient), nil
}

func buildBackupsClient(restClient *restclient.Client) (BackupsClient, error) {
	return database.NewBackupsClientImpl(restClient), nil
}

func buildUsersClient(restClient *restclient.Client) (UsersClient, error) {
	return database.NewUsersClientImpl(restClient), nil
}

func buildGrantsClient(restClient *restclient.Client) (GrantsClient, error) {
	return database.NewGrantsClientImpl(restClient), nil
}

//
// Metric domain clients

func buildMetricClient(restClient *restclient.Client) (MetricClient, error) {
	alertsClient, err := buildAlertsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	metricsClient, err := buildMetricsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &metricClientImpl{
		alertsClient:  alertsClient,
		metricsClient: metricsClient,
	}, nil
}

func buildAlertsClient(restClient *restclient.Client) (AlertsClient, error) {
	return metric.NewAlertsClientImpl(restClient), nil
}

func buildMetricsClient(restClient *restclient.Client) (MetricsClient, error) {
	return metric.NewMetricsClientImpl(restClient), nil
}

//
// Network domain clients

func buildNetworkClient(restClient *restclient.Client) (NetworkClient, error) {
	elasticIPsClient, err := buildElasticIPsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	loadBalancersClient, err := buildLoadBalancersClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	securityGroupRulesClient, err := buildSecurityGroupRulesClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	securityGroupsClient, err := buildSecurityGroupsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	subnetsClient, err := buildSubnetsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	vpcPeeringRoutesClient, err := buildVPCPeeringRoutesClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	vpcPeeringsClient, err := buildVPCPeeringsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	vpcsClient, err := buildVPCsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	vpnRoutesClient, err := buildVPNRoutesClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	vpnTunnelsClient, err := buildVPNTunnelsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &networkClientImpl{
		elasticIPsClient:         elasticIPsClient,
		loadBalancersClient:      loadBalancersClient,
		securityGroupRulesClient: securityGroupRulesClient,
		securityGroupsClient:     securityGroupsClient,
		subnetsClient:            subnetsClient,
		vpcPeeringRoutesClient:   vpcPeeringRoutesClient,
		vpcPeeringsClient:        vpcPeeringsClient,
		vpcsClient:               vpcsClient,
		vpnRoutesClient:          vpnRoutesClient,
		vpnTunnelsClient:         vpnTunnelsClient,
	}, nil
}

func buildElasticIPsClient(restClient *restclient.Client) (ElasticIPsClient, error) {
	return network.NewElasticIPsClientImpl(restClient), nil
}

func buildLoadBalancersClient(restClient *restclient.Client) (LoadBalancersClient, error) {
	return network.NewLoadBalancersClientImpl(restClient), nil
}

func buildSecurityGroupRulesClient(restClient *restclient.Client) (SecurityGroupRulesClient, error) {
	return network.NewSecurityGroupRulesClientImpl(
		restClient,
		network.NewSecurityGroupsClientImpl(
			restClient,
			network.NewVPCsClientImpl(restClient),
		),
	), nil
}

func buildSecurityGroupsClient(restClient *restclient.Client) (SecurityGroupsClient, error) {
	return network.NewSecurityGroupsClientImpl(
		restClient,
		network.NewVPCsClientImpl(restClient),
	), nil
}

func buildSubnetsClient(restClient *restclient.Client) (SubnetsClient, error) {
	return network.NewSubnetsClientImpl(
		restClient,
		network.NewVPCsClientImpl(restClient),
	), nil
}

func buildVPCPeeringRoutesClient(restClient *restclient.Client) (VPCPeeringRoutesClient, error) {
	return network.NewVPCPeeringRoutesClientImpl(restClient), nil
}

func buildVPCPeeringsClient(restClient *restclient.Client) (VPCPeeringsClient, error) {
	return network.NewVPCPeeringsClientImpl(restClient), nil
}

func buildVPCsClient(restClient *restclient.Client) (VPCsClient, error) {
	return network.NewVPCsClientImpl(restClient), nil
}

func buildVPNRoutesClient(restClient *restclient.Client) (VPNRoutesClient, error) {
	return network.NewVPNRoutesClientImpl(restClient), nil
}

func buildVPNTunnelsClient(restClient *restclient.Client) (VPNTunnelsClient, error) {
	return network.NewVPNTunnelsClientImpl(restClient), nil
}

//
// Project domain clients

func buildProjectClient(restClient *restclient.Client) (ProjectClient, error) {
	return project.NewProjectsClientImpl(restClient), nil
}

//
// Schedule domain clients

func buildScheduleClient(restClient *restclient.Client) (ScheduleClient, error) {
	jobsClient, err := buildJobsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &scheduleClientImpl{
		jobsClient: jobsClient,
	}, nil
}

func buildJobsClient(restClient *restclient.Client) (JobsClient, error) {
	return schedule.NewJobsClientImpl(restClient), nil
}

//
// Security domain clients

func buildSecurityClient(restClient *restclient.Client) (SecurityClient, error) {
	kmsKeysClient, err := buildKMSKeysClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &securityClientImpl{
		kmsKeysClient: kmsKeysClient,
	}, nil
}

func buildKMSKeysClient(restClient *restclient.Client) (KMSKeysClient, error) {
	return security.NewKMSKeysClientImpl(restClient), nil
}

//
// Storage domain clients

func buildStorageClient(restClient *restclient.Client) (StorageClient, error) {
	snapshotsClient, err := buildSnapshotsClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	volumesClient, err := buildVolumesClient(restClient)
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &storageClientImpl{
		snapshotsClient: snapshotsClient,
		volumesClient:   volumesClient,
	}, nil
}

func buildSnapshotsClient(restClient *restclient.Client) (SnapshotsClient, error) {
	return storage.NewSnapshotsClientImpl(
		restClient,
		storage.NewVolumesClientImpl(restClient),
	), nil
}

func buildVolumesClient(restClient *restclient.Client) (VolumesClient, error) {
	return storage.NewVolumesClientImpl(restClient), nil
}

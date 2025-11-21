//nolint:unparam // TODO: better error handling
package aruba

import (
	"github.com/Arubacloud/sdk-go/internal/clients/audit"
	"github.com/Arubacloud/sdk-go/internal/clients/compute"
	"github.com/Arubacloud/sdk-go/internal/clients/container"
	"github.com/Arubacloud/sdk-go/internal/clients/database"
	"github.com/Arubacloud/sdk-go/internal/clients/metric"
	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/pkg/spec/project"
	"github.com/Arubacloud/sdk-go/pkg/spec/schedule"
	"github.com/Arubacloud/sdk-go/pkg/spec/security"
	"github.com/Arubacloud/sdk-go/pkg/spec/storage"
)

// Client
func buildClient(config *restclient.Config) (Client, error) {
	restClient, err := buildRESTClient(config)
	if err != nil {
		return nil, err // TODO: better error handling
	}

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

	return &clientImpl{
		auditClient:     auditClient,
		computeClient:   computeClient,
		containerClient: containerClient,
		databaseClient:  databaseClient,
		metricsClient:   metricClient,
		networkClient:   networkClient,
		// TODO: Replace all below for refactored servers
		projectClient:  project.NewService(restClient),
		scheduleClient: schedule.NewService(restClient),
		securityClient: security.NewService(restClient),
		storageClient:  storage.NewService(restClient),
	}, nil
}

//
// Dependencies

func buildRESTClient(config *restclient.Config) (*restclient.Client, error) {
	return restclient.NewClient(config)
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

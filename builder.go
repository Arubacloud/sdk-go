//nolint:unparam // TODO: better error handling
package aruba

import (
	"github.com/Arubacloud/sdk-go/internal/clients/audit"
	"github.com/Arubacloud/sdk-go/internal/clients/compute"
	"github.com/Arubacloud/sdk-go/internal/clients/container"
	"github.com/Arubacloud/sdk-go/internal/clients/database"
	"github.com/Arubacloud/sdk-go/internal/clients/metric"
	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/pkg/spec/network"
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

	return &clientImpl{
		auditClient:     auditClient,
		computeClient:   computeClient,
		containerClient: containerClient,
		databaseClient:  databaseClient,
		metricsClient:   metricClient,
		// TODO: Replace all below for refactored servers
		networkClient:  network.NewService(restClient),
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

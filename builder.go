package aruba

import (
	"github.com/Arubacloud/sdk-go/internal/clients/audit"
	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/pkg/spec/compute"
	"github.com/Arubacloud/sdk-go/pkg/spec/container"
	"github.com/Arubacloud/sdk-go/pkg/spec/database"
	"github.com/Arubacloud/sdk-go/pkg/spec/metric"
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

	return &clientImpl{
		auditClient: auditClient,
		// TODO: Replace all below for refactored servers
		computeClient:   compute.NewService(restClient),
		containerClient: container.NewService(restClient),
		databaseClient:  database.NewService(restClient),
		metricsClient:   metric.NewService(restClient),
		networkClient:   network.NewService(restClient),
		projectClient:   project.NewService(restClient),
		scheduleClient:  schedule.NewService(restClient),
		securityClient:  security.NewService(restClient),
		storageClient:   storage.NewService(restClient),
	}, nil
}

//
// Audit domain clients

func buildAuditClient(restClient *restclient.Client) (AuditClient, error) {
	eventsClient, err := buildEventsClient(restClient) //nolint:unparam
	if err != nil {
		return nil, err // TODO: better error handling
	}

	return &auditClientImpl{eventsClient: eventsClient}, nil
}

func buildEventsClient(restClient *restclient.Client) (EventsClient, error) {
	//nolint:unparam // TODO: better error handling
	return audit.NewEventsClientImpl(restClient), nil
}

//
// Dependencies

func buildRESTClient(config *restclient.Config) (*restclient.Client, error) {
	return restclient.NewClient(config)
}

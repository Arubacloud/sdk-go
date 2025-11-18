// Package sdkgo provides the main entry point for the Aruba Cloud SDK
package sdkgo

import (
	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/pkg/spec/audit"
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

// Client wraps the client.Client and provides direct access to all service interfaces
type Client struct {
	*restclient.Client

	// Service interfaces for all API categories
	Compute   compute.ComputeAPI
	Network   network.NetworkAPI
	Storage   storage.StorageAPI
	Database  database.DatabaseAPI
	Container container.ContainerAPI
	Security  security.SecurityAPI
	Metric    metric.MetricAPI
	Audit     audit.AuditAPI
	Schedule  schedule.ScheduleAPI
	Project   project.ProjectAPI
}

// NewClient creates a new SDK client with all services initialized
func NewClient(config *restclient.Config) (*Client, error) {
	baseClient, err := restclient.NewClient(config)
	if err != nil {
		return nil, err
	}

	sdkClient := &Client{
		Client:    baseClient,
		Compute:   compute.NewService(baseClient),
		Network:   network.NewService(baseClient),
		Storage:   storage.NewService(baseClient),
		Database:  database.NewService(baseClient),
		Container: container.NewService(baseClient),
		Security:  security.NewService(baseClient),
		Metric:    metric.NewService(baseClient),
		Audit:     audit.NewService(baseClient),
		Schedule:  schedule.NewService(baseClient),
		Project:   project.NewService(baseClient),
	}

	return sdkClient, nil
}

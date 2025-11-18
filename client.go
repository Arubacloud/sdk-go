// Package sdkgo provides the main entry point for the Aruba Cloud SDK
package aruba

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

type Client interface {
	FromAudit() audit.AuditAPI
	FromCompute() compute.ComputeAPI
	FromContainer() container.ContainerAPI
	FromDatabase() database.DatabaseAPI
	FromMetric() metric.MetricAPI
	FromNetwork() network.NetworkAPI
	FromProject() project.ProjectAPI
	FromSchedule() schedule.ScheduleAPI
	FromSecurity() security.SecurityAPI
	FromStorage() storage.StorageAPI
}

// ClientImpl wraps the client.ClientImpl and provides direct access to all service interfaces
type ClientImpl struct {
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

var _ Client = (*ClientImpl)(nil)

// NewClient creates a new SDK client with all services initialized
func NewClient(config *restclient.Config) (*ClientImpl, error) {
	baseClient, err := restclient.NewClient(config)
	if err != nil {
		return nil, err
	}

	sdkClient := &ClientImpl{
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

func (c *ClientImpl) FromAudit() audit.AuditAPI {
	return c.Audit
}

func (c *ClientImpl) FromCompute() compute.ComputeAPI {
	return c.Compute
}
func (c *ClientImpl) FromContainer() container.ContainerAPI {
	return c.Container
}
func (c *ClientImpl) FromDatabase() database.DatabaseAPI {
	return c.Database
}
func (c *ClientImpl) FromMetric() metric.MetricAPI {
	return c.Metric
}
func (c *ClientImpl) FromNetwork() network.NetworkAPI {
	return c.Network
}
func (c *ClientImpl) FromProject() project.ProjectAPI {
	return c.Project
}
func (c *ClientImpl) FromSchedule() schedule.ScheduleAPI {
	return c.Schedule
}
func (c *ClientImpl) FromSecurity() security.SecurityAPI {
	return c.Security
}
func (c *ClientImpl) FromStorage() storage.StorageAPI {
	return c.Storage
}

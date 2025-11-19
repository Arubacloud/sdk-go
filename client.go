// Package sdkgo provides the main entry point for the Aruba Cloud SDK
package aruba

import (
	"github.com/Arubacloud/sdk-go/pkg/spec/database"
	"github.com/Arubacloud/sdk-go/pkg/spec/metric"
	"github.com/Arubacloud/sdk-go/pkg/spec/network"
	"github.com/Arubacloud/sdk-go/pkg/spec/project"
	"github.com/Arubacloud/sdk-go/pkg/spec/schedule"
	"github.com/Arubacloud/sdk-go/pkg/spec/security"
	"github.com/Arubacloud/sdk-go/pkg/spec/storage"
)

type Client interface {
	FromAudit() AuditClient
	FromCompute() ComputeClient
	FromContainer() ContainerClient
	FromDatabase() database.DatabaseAPI
	FromMetric() metric.MetricAPI
	FromNetwork() network.NetworkAPI
	FromProject() project.ProjectAPI
	FromSchedule() schedule.ScheduleAPI
	FromSecurity() security.SecurityAPI
	FromStorage() storage.StorageAPI
}

type clientImpl struct {
	auditClient     AuditClient
	computeClient   ComputeClient
	containerClient ContainerClient
	databaseClient  database.DatabaseAPI
	metricsClient   metric.MetricAPI
	networkClient   network.NetworkAPI
	projectClient   project.ProjectAPI
	scheduleClient  schedule.ScheduleAPI
	securityClient  security.SecurityAPI
	storageClient   storage.StorageAPI
}

var _ Client = (*clientImpl)(nil)

func (c *clientImpl) FromAudit() AuditClient {
	return c.auditClient
}
func (c *clientImpl) FromCompute() ComputeClient {
	return c.computeClient
}
func (c *clientImpl) FromContainer() ContainerClient {
	return c.containerClient
}
func (c *clientImpl) FromDatabase() database.DatabaseAPI {
	return c.databaseClient
}
func (c *clientImpl) FromMetric() metric.MetricAPI {
	return c.metricsClient
}
func (c *clientImpl) FromNetwork() network.NetworkAPI {
	return c.networkClient
}
func (c *clientImpl) FromProject() project.ProjectAPI {
	return c.projectClient
}
func (c *clientImpl) FromSchedule() schedule.ScheduleAPI {
	return c.scheduleClient
}
func (c *clientImpl) FromSecurity() security.SecurityAPI {
	return c.securityClient
}
func (c *clientImpl) FromStorage() storage.StorageAPI {
	return c.storageClient
}

// Package sdkgo provides the main entry point for the Aruba Cloud SDK
package aruba

import (
	impl "github.com/Arubacloud/sdk-go/internal/client"
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

func NewClient(config *restclient.Config) (Client, error) {
	return impl.NewClient(config)
}

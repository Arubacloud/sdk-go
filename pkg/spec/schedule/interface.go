package schedule

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ScheduleAPI defines the unified interface for all Schedule operations
type ScheduleAPI interface {
	// ScheduleJob operations
	ListScheduleJobs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.JobList], error)
	GetScheduleJob(ctx context.Context, project string, scheduleJobId string, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error)
	CreateScheduleJob(ctx context.Context, project string, body schema.JobRequest, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error)
	UpdateScheduleJob(ctx context.Context, project string, scheduleJobId string, body schema.JobRequest, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error)
	DeleteScheduleJob(ctx context.Context, projectId string, scheduleJobId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

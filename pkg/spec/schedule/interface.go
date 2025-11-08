package schedule

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ScheduleJobAPI defines the interface for managing schedule jobs.
type ScheduleJobAPI interface {
	ListScheduleJobs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.JobList], error)
	GetScheduleJob(ctx context.Context, project string, scheduleJobId string, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error)
	CreateOrUpdateScheduleJob(ctx context.Context, project string, body schema.JobRequest, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error)
	DeleteScheduleJob(ctx context.Context, projectId string, scheduleJobId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

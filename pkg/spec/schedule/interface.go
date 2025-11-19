package schedule

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

// ScheduleAPI defines the unified interface for all Schedule operations
type ScheduleAPI interface {
	// ScheduleJob operations
	ListScheduleJobs(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.JobList], error)
	GetScheduleJob(ctx context.Context, project string, scheduleJobId string, params *types.RequestParameters) (*types.Response[types.JobResponse], error)
	CreateScheduleJob(ctx context.Context, project string, body types.JobRequest, params *types.RequestParameters) (*types.Response[types.JobResponse], error)
	UpdateScheduleJob(ctx context.Context, project string, scheduleJobId string, body types.JobRequest, params *types.RequestParameters) (*types.Response[types.JobResponse], error)
	DeleteScheduleJob(ctx context.Context, projectId string, scheduleJobId string, params *types.RequestParameters) (*types.Response[any], error)
}

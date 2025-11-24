package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

type ScheduleClient interface {
	Jobs() JobsClient
}

type scheduleClientImpl struct {
	jobsClient JobsClient
}

var _ ScheduleClient = (*scheduleClientImpl)(nil)

func (c *scheduleClientImpl) Jobs() JobsClient {
	return c.jobsClient
}

type JobsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.JobList], error)
	Get(ctx context.Context, projectID string, scheduleJobID string, params *types.RequestParameters) (*types.Response[types.JobResponse], error)
	Create(ctx context.Context, projectID string, body types.JobRequest, params *types.RequestParameters) (*types.Response[types.JobResponse], error)
	Update(ctx context.Context, projectID string, scheduleJobID string, body types.JobRequest, params *types.RequestParameters) (*types.Response[types.JobResponse], error)
	Delete(ctx context.Context, projectID string, scheduleJobID string, params *types.RequestParameters) (*types.Response[any], error)
}

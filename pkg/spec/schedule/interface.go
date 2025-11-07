package schedule

import (
	"context"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ScheduleJobAPI defines the interface for managing schedule jobs.
type ScheduleJobAPI interface {
	ListScheduleJobs(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetScheduleJob(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateScheduleJob(ctx context.Context, project string, body schema.ScheduleJobRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteScheduleJob(ctx context.Context, projectId string, vpcId string, params *schema.RequestParameters) (*http.Response, error)
}

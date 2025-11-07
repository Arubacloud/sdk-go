package schedule

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// JobService implements the JobAPI interface
type JobService struct {
	client *client.Client
}

// NewJobService creates a new JobService
func NewJobService(client *client.Client) *JobService {
	return &JobService{
		client: client,
	}
}

// ListJobs retrieves all scheduled jobs for a project
func (s *JobService) ListJobs(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(JobsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetJob retrieves a specific scheduled job by ID
func (s *JobService) GetJob(ctx context.Context, project string, jobId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if jobId == "" {
		return nil, fmt.Errorf("job ID cannot be empty")
	}

	path := fmt.Sprintf(JobPath, project, jobId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateJob creates or updates a scheduled job
func (s *JobService) CreateOrUpdateJob(ctx context.Context, project string, body schema.JobRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(JobsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteJob deletes a scheduled job by ID
func (s *JobService) DeleteJob(ctx context.Context, projectId string, jobId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if jobId == "" {
		return nil, fmt.Errorf("job ID cannot be empty")
	}

	path := fmt.Sprintf(JobPath, projectId, jobId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

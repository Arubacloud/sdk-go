package schedule

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// JobService implements the ScheduleJobAPI interface
type JobService struct {
	client *client.Client
}

// NewJobService creates a new JobService
func NewJobService(client *client.Client) *JobService {
	return &JobService{
		client: client,
	}
}

// ListScheduleJobs retrieves all schedule jobs for a project
func (s *JobService) ListScheduleJobs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.JobList], error) {
	s.client.Logger().Debugf("Listing schedule jobs for project: %s", project)

	if err := validateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ScheduleJobsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.JobList]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.JobList
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// GetScheduleJob retrieves a specific schedule job by ID
func (s *JobService) GetScheduleJob(ctx context.Context, project string, scheduleJobId string, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error) {
	s.client.Logger().Debugf("Getting schedule job: %s in project: %s", scheduleJobId, project)

	if err := validateProjectAndResource(project, scheduleJobId, "job ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(JobPath, project, scheduleJobId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.JobResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.JobResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// CreateScheduleJob creates a new schedule job
func (s *JobService) CreateScheduleJob(ctx context.Context, project string, body schema.JobRequest, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error) {
	s.client.Logger().Debugf("Creating schedule job in project: %s", project)

	if err := validateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ScheduleJobsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.JobResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.JobResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// UpdateScheduleJob updates an existing schedule job
func (s *JobService) UpdateScheduleJob(ctx context.Context, project string, scheduleJobId string, body schema.JobRequest, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error) {
	s.client.Logger().Debugf("Updating schedule job: %s in project: %s", scheduleJobId, project)

	if err := validateProjectAndResource(project, scheduleJobId, "job ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(JobPath, project, scheduleJobId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPut, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.JobResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.JobResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// DeleteScheduleJob deletes a schedule job by ID
func (s *JobService) DeleteScheduleJob(ctx context.Context, projectId string, scheduleJobId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting schedule job: %s in project: %s", scheduleJobId, projectId)

	if err := validateProjectAndResource(projectId, scheduleJobId, "job ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(JobPath, projectId, scheduleJobId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[any]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() && len(bodyBytes) > 0 {
		var data any
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

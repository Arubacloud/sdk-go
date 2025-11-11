package schedule

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ListScheduleJobs retrieves all schedule jobs for a project
func (s *Service) ListScheduleJobs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.JobList], error) {
	s.client.Logger().Debugf("Listing schedule jobs for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(JobsPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &ScheduleJobListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ScheduleJobListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.JobList](httpResp)
}

// GetScheduleJob retrieves a specific schedule job by ID
func (s *Service) GetScheduleJob(ctx context.Context, project string, scheduleJobId string, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error) {
	s.client.Logger().Debugf("Getting schedule job: %s in project: %s", scheduleJobId, project)

	if err := schema.ValidateProjectAndResource(project, scheduleJobId, "job ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(JobPath, project, scheduleJobId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &ScheduleJobGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ScheduleJobGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.JobResponse](httpResp)
}

// CreateScheduleJob creates a new schedule job
func (s *Service) CreateScheduleJob(ctx context.Context, project string, body schema.JobRequest, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error) {
	s.client.Logger().Debugf("Creating schedule job in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(JobsPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &ScheduleJobCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ScheduleJobCreateAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

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
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp schema.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// UpdateScheduleJob updates an existing schedule job
func (s *Service) UpdateScheduleJob(ctx context.Context, project string, scheduleJobId string, body schema.JobRequest, params *schema.RequestParameters) (*schema.Response[schema.JobResponse], error) {
	s.client.Logger().Debugf("Updating schedule job: %s in project: %s", scheduleJobId, project)

	if err := schema.ValidateProjectAndResource(project, scheduleJobId, "job ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(JobPath, project, scheduleJobId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &ScheduleJobUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ScheduleJobUpdateAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

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
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp schema.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// DeleteScheduleJob deletes a schedule job by ID
func (s *Service) DeleteScheduleJob(ctx context.Context, projectId string, scheduleJobId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting schedule job: %s in project: %s", scheduleJobId, projectId)

	if err := schema.ValidateProjectAndResource(projectId, scheduleJobId, "job ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(JobPath, projectId, scheduleJobId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &ScheduleJobDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ScheduleJobDeleteAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[any](httpResp)
}

package project

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/types"
)

// ListProjects retrieves all projects
func (s *Service) ListProjects(ctx context.Context, params *types.RequestParameters) (*types.Response[types.ProjectList], error) {
	s.client.Logger().Debugf("Listing projects")

	path := ProjectsPath

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ProjectListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ProjectListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.ProjectList](httpResp)
}

// GetProject retrieves a specific project by ID
func (s *Service) GetProject(ctx context.Context, projectId string, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error) {
	s.client.Logger().Debugf("Getting project: %s", projectId)

	if err := types.ValidateProject(projectId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ProjectPath, projectId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ProjectGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ProjectGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.ProjectResponse](httpResp)
}

// CreateProject creates a new project
func (s *Service) CreateProject(ctx context.Context, body types.ProjectRequest, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error) {
	s.client.Logger().Debugf("Creating project")

	path := ProjectsPath

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ProjectCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ProjectCreateAPIVersion
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

	response := &types.Response[types.ProjectResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.ProjectResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp types.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// UpdateProject updates an existing project
func (s *Service) UpdateProject(ctx context.Context, projectId string, body types.ProjectRequest, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error) {
	s.client.Logger().Debugf("Updating project: %s", projectId)

	if err := types.ValidateProject(projectId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ProjectPath, projectId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ProjectUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ProjectUpdateAPIVersion
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

	response := &types.Response[types.ProjectResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.ProjectResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp types.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// DeleteProject deletes a project by ID
func (s *Service) DeleteProject(ctx context.Context, projectId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting project: %s", projectId)

	if err := types.ValidateProject(projectId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ProjectPath, projectId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ProjectDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ProjectDeleteAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &types.Response[any]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	return response, nil
}

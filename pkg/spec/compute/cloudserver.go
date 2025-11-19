package compute

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/types"
)

// ListCloudServers retrieves all cloud servers for a project
func (s *Service) ListCloudServers(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.CloudServerList], error) {
	s.client.Logger().Debugf("Listing cloud servers for project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(CloudServersPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ComputeCloudServerList,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ComputeCloudServerList
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.CloudServerList](httpResp)
}

// GetCloudServer retrieves a specific cloud server by ID
func (s *Service) GetCloudServer(ctx context.Context, project string, cloudServerId string, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error) {
	s.client.Logger().Debugf("Getting cloud server: %s in project: %s", cloudServerId, project)

	if err := types.ValidateProjectAndResource(project, cloudServerId, "cloud server ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(CloudServerPath, project, cloudServerId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ComputeCloudServerGet,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ComputeCloudServerGet
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.CloudServerResponse](httpResp)
}

// CreateCloudServer creates a new cloud server
func (s *Service) CreateCloudServer(ctx context.Context, project string, body types.CloudServerRequest, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error) {
	s.client.Logger().Debugf("Creating cloud server in project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(CloudServersPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ComputeCloudServerCreate,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ComputeCloudServerCreate
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	// Marshal the request body to JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// Read the response body
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &types.Response[types.CloudServerResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.CloudServerResponse
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

// UpdateCloudServer updates an existing cloud server
func (s *Service) UpdateCloudServer(ctx context.Context, project string, cloudServerId string, body types.CloudServerRequest, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error) {
	s.client.Logger().Debugf("Updating cloud server: %s in project: %s", cloudServerId, project)

	if err := types.ValidateProjectAndResource(project, cloudServerId, "cloud server ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(CloudServerPath, project, cloudServerId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ComputeCloudServerUpdate,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ComputeCloudServerUpdate
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	// Marshal the request body to JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPut, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// Read the response body
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &types.Response[types.CloudServerResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.CloudServerResponse
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

// DeleteCloudServer deletes a cloud server by ID
func (s *Service) DeleteCloudServer(ctx context.Context, projectId string, cloudServerId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting cloud server: %s in project: %s", cloudServerId, projectId)

	if err := types.ValidateProjectAndResource(projectId, cloudServerId, "cloud server ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(CloudServerPath, projectId, cloudServerId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ComputeCloudServerDelete,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ComputeCloudServerDelete
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[any](httpResp)
}

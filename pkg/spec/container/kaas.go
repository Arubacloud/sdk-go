package container

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/types"
)

// ListKaaS retrieves all KaaS clusters for a project
func (s *Service) ListKaaS(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.KaaSList], error) {
	s.client.Logger().Debugf("Listing KaaS clusters for project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KaaSPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ContainerKaaSListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ContainerKaaSListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.KaaSList](httpResp)
}

// GetKaaS retrieves a specific KaaS cluster by ID
func (s *Service) GetKaaS(ctx context.Context, project string, kaasId string, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error) {
	s.client.Logger().Debugf("Getting KaaS cluster: %s in project: %s", kaasId, project)

	if err := types.ValidateProjectAndResource(project, kaasId, "KaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KaaSItemPath, project, kaasId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ContainerKaaSGetVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ContainerKaaSGetVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.KaaSResponse](httpResp)
}

// CreateKaaS creates a new KaaS cluster
func (s *Service) CreateKaaS(ctx context.Context, project string, body types.KaaSRequest, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error) {
	s.client.Logger().Debugf("Creating KaaS cluster in project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KaaSPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ContainerKaaSCreateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ContainerKaaSCreateVersion
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
	response := &types.Response[types.KaaSResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.KaaSResponse
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

// UpdateKaaS updates an existing KaaS cluster
func (s *Service) UpdateKaaS(ctx context.Context, project string, kaasId string, body types.KaaSRequest, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error) {
	s.client.Logger().Debugf("Updating KaaS cluster: %s in project: %s", kaasId, project)

	if err := types.ValidateProjectAndResource(project, kaasId, "KaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KaaSItemPath, project, kaasId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ContainerKaaSUpdateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ContainerKaaSUpdateVersion
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
	response := &types.Response[types.KaaSResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.KaaSResponse
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

// DeleteKaaS deletes a KaaS cluster by ID
func (s *Service) DeleteKaaS(ctx context.Context, projectId string, kaasId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting KaaS cluster: %s in project: %s", kaasId, projectId)

	if err := types.ValidateProjectAndResource(projectId, kaasId, "KaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KaaSItemPath, projectId, kaasId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ContainerKaaSDeleteVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ContainerKaaSDeleteVersion
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

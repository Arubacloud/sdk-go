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

// ListKeyPairs retrieves all key pairs for a project
func (s *Service) ListKeyPairs(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.KeyPairListResponse], error) {
	s.client.Logger().Debugf("Listing key pairs for project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KeyPairsPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ComputeKeyPairList,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ComputeKeyPairList
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.KeyPairListResponse](httpResp)
}

// GetKeyPair retrieves a specific key pair by ID
func (s *Service) GetKeyPair(ctx context.Context, project string, keyPairId string, params *types.RequestParameters) (*types.Response[types.KeyPairResponse], error) {
	s.client.Logger().Debugf("Getting key pair: %s in project: %s", keyPairId, project)

	if err := types.ValidateProjectAndResource(project, keyPairId, "key pair ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KeyPairPath, project, keyPairId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ComputeKeyPairGet,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ComputeKeyPairGet
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.KeyPairResponse](httpResp)
}

// CreateKeyPair creates a new key pair
func (s *Service) CreateKeyPair(ctx context.Context, project string, body types.KeyPairRequest, params *types.RequestParameters) (*types.Response[types.KeyPairResponse], error) {
	s.client.Logger().Debugf("Creating key pair in project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KeyPairsPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ComputeKeyPairCreate,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ComputeKeyPairCreate
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
	response := &types.Response[types.KeyPairResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.KeyPairResponse
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

// DeleteKeyPair deletes a key pair by ID
func (s *Service) DeleteKeyPair(ctx context.Context, projectId string, keyPairId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting key pair: %s in project: %s", keyPairId, projectId)

	if err := types.ValidateProjectAndResource(projectId, keyPairId, "key pair ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KeyPairPath, projectId, keyPairId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &ComputeKeyPairDelete,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ComputeKeyPairDelete
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

package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/types"
)

// ListGrants retrieves all grants for a database
func (s *Service) ListGrants(ctx context.Context, project string, dbaasId string, databaseId string, params *types.RequestParameters) (*types.Response[types.GrantList], error) {
	s.client.Logger().Debugf("Listing grants for database: %s in DBaaS: %s in project: %s", databaseId, dbaasId, project)

	if err := types.ValidateDBaaSResource(project, dbaasId, databaseId, "database ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(GrantsPath, project, dbaasId, databaseId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseGrantListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseGrantListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.GrantList](httpResp)
}

// GetGrant retrieves a specific grant by ID
func (s *Service) GetGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, params *types.RequestParameters) (*types.Response[types.GrantResponse], error) {
	s.client.Logger().Debugf("Getting grant: %s from database: %s in DBaaS: %s in project: %s", grantId, databaseId, dbaasId, project)

	if err := types.ValidateDatabaseGrant(project, dbaasId, databaseId, grantId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(GrantItemPath, project, dbaasId, databaseId, grantId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseGrantGetVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseGrantGetVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.GrantResponse](httpResp)
}

// CreateGrant creates a new grant for a database
func (s *Service) CreateGrant(ctx context.Context, project string, dbaasId string, databaseId string, body types.GrantRequest, params *types.RequestParameters) (*types.Response[types.GrantResponse], error) {
	s.client.Logger().Debugf("Creating grant in database: %s in DBaaS: %s in project: %s", databaseId, dbaasId, project)

	if err := types.ValidateDBaaSResource(project, dbaasId, databaseId, "database ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(GrantsPath, project, dbaasId, databaseId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseGrantCreateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseGrantCreateVersion
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
	response := &types.Response[types.GrantResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.GrantResponse
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

// UpdateGrant updates an existing grant
func (s *Service) UpdateGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, body types.GrantRequest, params *types.RequestParameters) (*types.Response[types.GrantResponse], error) {
	s.client.Logger().Debugf("Updating grant: %s in database: %s in DBaaS: %s in project: %s", grantId, databaseId, dbaasId, project)

	if err := types.ValidateDatabaseGrant(project, dbaasId, databaseId, grantId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(GrantItemPath, project, dbaasId, databaseId, grantId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseGrantUpdateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseGrantUpdateVersion
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
	response := &types.Response[types.GrantResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.GrantResponse
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

// DeleteGrant deletes a grant by ID
func (s *Service) DeleteGrant(ctx context.Context, projectId string, dbaasId string, databaseId string, grantId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting grant: %s from database: %s in DBaaS: %s in project: %s", grantId, databaseId, dbaasId, projectId)

	if err := types.ValidateDatabaseGrant(projectId, dbaasId, databaseId, grantId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(GrantItemPath, projectId, dbaasId, databaseId, grantId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseGrantDeleteVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseGrantDeleteVersion
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

package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ListGrants retrieves all grants for a database
func (s *Service) ListGrants(ctx context.Context, project string, dbaasId string, databaseId string, params *schema.RequestParameters) (*schema.Response[schema.GrantList], error) {
	s.client.Logger().Debugf("Listing grants for database: %s in DBaaS: %s in project: %s", databaseId, dbaasId, project)

	if err := schema.ValidateDBaaSResource(project, dbaasId, databaseId, "database ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(GrantsPath, project, dbaasId, databaseId)

	if params == nil {
		params = &schema.RequestParameters{
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

	return schema.ParseResponseBody[schema.GrantList](httpResp)
}

// GetGrant retrieves a specific grant by ID
func (s *Service) GetGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, params *schema.RequestParameters) (*schema.Response[schema.GrantResponse], error) {
	s.client.Logger().Debugf("Getting grant: %s from database: %s in DBaaS: %s in project: %s", grantId, databaseId, dbaasId, project)

	if err := schema.ValidateDatabaseGrant(project, dbaasId, databaseId, grantId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(GrantItemPath, project, dbaasId, databaseId, grantId)

	if params == nil {
		params = &schema.RequestParameters{
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

	return schema.ParseResponseBody[schema.GrantResponse](httpResp)
}

// CreateGrant creates a new grant for a database
func (s *Service) CreateGrant(ctx context.Context, project string, dbaasId string, databaseId string, body schema.GrantRequest, params *schema.RequestParameters) (*schema.Response[schema.GrantResponse], error) {
	s.client.Logger().Debugf("Creating grant in database: %s in DBaaS: %s in project: %s", databaseId, dbaasId, project)

	if err := schema.ValidateDBaaSResource(project, dbaasId, databaseId, "database ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(GrantsPath, project, dbaasId, databaseId)

	if params == nil {
		params = &schema.RequestParameters{
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
	response := &schema.Response[schema.GrantResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.GrantResponse
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

// UpdateGrant updates an existing grant
func (s *Service) UpdateGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, body schema.GrantRequest, params *schema.RequestParameters) (*schema.Response[schema.GrantResponse], error) {
	s.client.Logger().Debugf("Updating grant: %s in database: %s in DBaaS: %s in project: %s", grantId, databaseId, dbaasId, project)

	if err := schema.ValidateDatabaseGrant(project, dbaasId, databaseId, grantId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(GrantItemPath, project, dbaasId, databaseId, grantId)

	if params == nil {
		params = &schema.RequestParameters{
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
	response := &schema.Response[schema.GrantResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.GrantResponse
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

// DeleteGrant deletes a grant by ID
func (s *Service) DeleteGrant(ctx context.Context, projectId string, dbaasId string, databaseId string, grantId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting grant: %s from database: %s in DBaaS: %s in project: %s", grantId, databaseId, dbaasId, projectId)

	if err := schema.ValidateDatabaseGrant(projectId, dbaasId, databaseId, grantId); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(GrantItemPath, projectId, dbaasId, databaseId, grantId)

	if params == nil {
		params = &schema.RequestParameters{
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

	return schema.ParseResponseBody[any](httpResp)
}

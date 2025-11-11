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

// ListDatabases retrieves all databases for a DBaaS instance
func (s *Service) ListDatabases(ctx context.Context, project string, dbaasId string, params *schema.RequestParameters) (*schema.Response[schema.DatabaseList], error) {
	s.client.Logger().Debugf("Listing databases for DBaaS: %s in project: %s", dbaasId, project)

	if err := schema.ValidateProjectAndResource(project, dbaasId, "DBaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DatabaseInstancesPath, project, dbaasId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &DatabaseInstanceListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseInstanceListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.DatabaseList](httpResp)
}

// GetDatabase retrieves a specific database by ID
func (s *Service) GetDatabase(ctx context.Context, project string, dbaasId string, databaseId string, params *schema.RequestParameters) (*schema.Response[schema.DatabaseResponse], error) {
	s.client.Logger().Debugf("Getting database: %s from DBaaS: %s in project: %s", databaseId, dbaasId, project)

	if err := schema.ValidateDBaaSResource(project, dbaasId, databaseId, "database ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DatabaseInstancePath, project, dbaasId, databaseId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &DatabaseInstanceGetVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseInstanceGetVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.DatabaseResponse](httpResp)
}

// CreateDatabase creates a new database
func (s *Service) CreateDatabase(ctx context.Context, project string, dbaasId string, body schema.DatabaseRequest, params *schema.RequestParameters) (*schema.Response[schema.DatabaseResponse], error) {
	s.client.Logger().Debugf("Creating database in DBaaS: %s in project: %s", dbaasId, project)

	if err := schema.ValidateProjectAndResource(project, dbaasId, "DBaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DatabaseInstancesPath, project, dbaasId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &DatabaseInstanceCreateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseInstanceCreateVersion
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
	response := &schema.Response[schema.DatabaseResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.DatabaseResponse
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

// UpdateDatabase updates an existing database
func (s *Service) UpdateDatabase(ctx context.Context, project string, dbaasId string, databaseId string, body schema.DatabaseRequest, params *schema.RequestParameters) (*schema.Response[schema.DatabaseResponse], error) {
	s.client.Logger().Debugf("Updating database: %s in DBaaS: %s in project: %s", databaseId, dbaasId, project)

	if err := schema.ValidateDBaaSResource(project, dbaasId, databaseId, "database ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DatabaseInstancePath, project, dbaasId, databaseId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &DatabaseInstanceUpdateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseInstanceUpdateVersion
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
	response := &schema.Response[schema.DatabaseResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.DatabaseResponse
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

// DeleteDatabase deletes a database by ID
func (s *Service) DeleteDatabase(ctx context.Context, projectId string, dbaasId string, databaseId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting database: %s from DBaaS: %s in project: %s", databaseId, dbaasId, projectId)

	if err := schema.ValidateDBaaSResource(projectId, dbaasId, databaseId, "database ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DatabaseInstancePath, projectId, dbaasId, databaseId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &DatabaseInstanceDeleteVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseInstanceDeleteVersion
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

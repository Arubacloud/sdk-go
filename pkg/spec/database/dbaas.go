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

// ListDBaaS retrieves all DBaaS instances for a project
func (s *Service) ListDBaaS(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.DBaaSList], error) {
	s.client.Logger().Debugf("Listing DBaaS instances for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DBaaSPath, project)
	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &DatabaseDBaaSListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseDBaaSListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.DBaaSList](httpResp)
}

// GetDBaaS retrieves a specific DBaaS instance by ID
func (s *Service) GetDBaaS(ctx context.Context, project string, dbaasId string, params *schema.RequestParameters) (*schema.Response[schema.DBaaSResponse], error) {
	s.client.Logger().Debugf("Getting DBaaS instance: %s in project: %s", dbaasId, project)

	if err := schema.ValidateProjectAndResource(project, dbaasId, "DBaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DBaaSItemPath, project, dbaasId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &DatabaseDBaaSGetVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseDBaaSGetVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.DBaaSResponse](httpResp)
}

// CreateDBaaS creates a new DBaaS instance
func (s *Service) CreateDBaaS(ctx context.Context, project string, body schema.DBaaSRequest, params *schema.RequestParameters) (*schema.Response[schema.DBaaSResponse], error) {
	s.client.Logger().Debugf("Creating DBaaS instance in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DBaaSPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &DatabaseDBaaSCreateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseDBaaSCreateVersion
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
	response := &schema.Response[schema.DBaaSResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.DBaaSResponse
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

// UpdateDBaaS updates an existing DBaaS instance
func (s *Service) UpdateDBaaS(ctx context.Context, project string, databaseId string, body schema.DBaaSRequest, params *schema.RequestParameters) (*schema.Response[schema.DBaaSResponse], error) {
	s.client.Logger().Debugf("Updating DBaaS instance: %s in project: %s", databaseId, project)

	if err := schema.ValidateProjectAndResource(project, databaseId, "DBaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DBaaSItemPath, project, databaseId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &DatabaseDBaaSUpdateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseDBaaSUpdateVersion
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
	response := &schema.Response[schema.DBaaSResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.DBaaSResponse
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

// DeleteDBaaS deletes a DBaaS instance by ID
func (s *Service) DeleteDBaaS(ctx context.Context, projectId string, dbaasId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting DBaaS instance: %s in project: %s", dbaasId, projectId)

	if err := schema.ValidateProjectAndResource(projectId, dbaasId, "DBaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DBaaSItemPath, projectId, dbaasId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &DatabaseDBaaSDeleteVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseDBaaSDeleteVersion
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

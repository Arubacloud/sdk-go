package database

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

// DatabaseService implements the DatabaseAPI interface
type DatabaseService struct {
	client *client.Client
}

// NewDatabaseService creates a new DatabaseService
func NewDatabaseService(client *client.Client) *DatabaseService {
	return &DatabaseService{
		client: client,
	}
}

// ListDatabases retrieves all databases for a DBaaS instance
func (s *DatabaseService) ListDatabases(ctx context.Context, project string, dbaasId string, params *schema.RequestParameters) (*schema.Response[schema.DatabaseList], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}

	path := fmt.Sprintf(DatabaseInstancesPath, project, dbaasId)

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

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &schema.Response[schema.DatabaseList]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.DatabaseList
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// GetDatabase retrieves a specific database by ID
func (s *DatabaseService) GetDatabase(ctx context.Context, project string, dbaasId string, databaseId string, params *schema.RequestParameters) (*schema.Response[schema.DatabaseResponse], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseId == "" {
		return nil, fmt.Errorf("database ID cannot be empty")
	}

	path := fmt.Sprintf(DatabaseInstancePath, project, dbaasId, databaseId)

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

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &schema.Response[schema.DatabaseResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.DatabaseResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// CreateDatabase creates a new database
func (s *DatabaseService) CreateDatabase(ctx context.Context, project string, dbaasId string, body schema.DatabaseRequest, params *schema.RequestParameters) (*schema.Response[schema.DatabaseResponse], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}

	path := fmt.Sprintf(DatabaseInstancesPath, project, dbaasId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

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
	}

	return response, nil
}

// UpdateDatabase updates an existing database
func (s *DatabaseService) UpdateDatabase(ctx context.Context, project string, dbaasId string, databaseId string, body schema.DatabaseRequest, params *schema.RequestParameters) (*schema.Response[schema.DatabaseResponse], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseId == "" {
		return nil, fmt.Errorf("database ID cannot be empty")
	}

	path := fmt.Sprintf(DatabaseInstancePath, project, dbaasId, databaseId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

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
	}

	return response, nil
}

// DeleteDatabase deletes a database by ID
func (s *DatabaseService) DeleteDatabase(ctx context.Context, projectId string, dbaasId string, databaseId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseId == "" {
		return nil, fmt.Errorf("database ID cannot be empty")
	}

	path := fmt.Sprintf(DatabaseInstancePath, projectId, dbaasId, databaseId)

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

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &schema.Response[any]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// For DELETE operations, we typically don't parse the body unless there's content
	if response.IsSuccess() && len(bodyBytes) > 0 {
		var data any
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

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

// GrantService implements the GrantAPI interface
type GrantService struct {
	client *client.Client
}

// NewGrantService creates a new GrantService
func NewGrantService(client *client.Client) *GrantService {
	return &GrantService{
		client: client,
	}
}

// ListGrants retrieves all grants for a database
func (s *GrantService) ListGrants(ctx context.Context, project string, dbaasId string, databaseId string, params *schema.RequestParameters) (*schema.Response[schema.GrantList], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseId == "" {
		return nil, fmt.Errorf("database ID cannot be empty")
	}

	path := fmt.Sprintf(GrantsPath, project, dbaasId, databaseId)

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
	response := &schema.Response[schema.GrantList]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.GrantList
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// GetGrant retrieves a specific grant by ID
func (s *GrantService) GetGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, params *schema.RequestParameters) (*schema.Response[schema.GrantResponse], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseId == "" {
		return nil, fmt.Errorf("database ID cannot be empty")
	}
	if grantId == "" {
		return nil, fmt.Errorf("grant ID cannot be empty")
	}

	path := fmt.Sprintf(GrantItemPath, project, dbaasId, databaseId, grantId)

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
	response := &schema.Response[schema.GrantResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.GrantResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// CreateGrant creates a new grant
func (s *GrantService) CreateGrant(ctx context.Context, project string, dbaasId string, databaseId string, body schema.GrantRequest, params *schema.RequestParameters) (*schema.Response[schema.GrantResponse], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseId == "" {
		return nil, fmt.Errorf("database ID cannot be empty")
	}

	path := fmt.Sprintf(GrantsPath, project, dbaasId, databaseId)

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
	}

	return response, nil
}

// UpdateGrant updates an existing grant
func (s *GrantService) UpdateGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, body schema.GrantRequest, params *schema.RequestParameters) (*schema.Response[schema.GrantResponse], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseId == "" {
		return nil, fmt.Errorf("database ID cannot be empty")
	}
	if grantId == "" {
		return nil, fmt.Errorf("grant ID cannot be empty")
	}

	path := fmt.Sprintf(GrantItemPath, project, dbaasId, databaseId, grantId)

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
	}

	return response, nil
}

// DeleteGrant deletes a grant by ID
func (s *GrantService) DeleteGrant(ctx context.Context, projectId string, dbaasId string, databaseId string, grantId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseId == "" {
		return nil, fmt.Errorf("database ID cannot be empty")
	}
	if grantId == "" {
		return nil, fmt.Errorf("grant ID cannot be empty")
	}

	path := fmt.Sprintf(GrantItemPath, projectId, dbaasId, databaseId, grantId)

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

package compute

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

// CloudServerService implements the CloudServerAPI interface
type CloudServerService struct {
	client *client.Client
}

// NewCloudServerService creates a new CloudServerService
func NewCloudServerService(client *client.Client) *CloudServerService {
	return &CloudServerService{
		client: client,
	}
}

// ListCloudServers retrieves all cloud servers for a project
func (s *CloudServerService) ListCloudServers(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.CloudServerList], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(CloudServersPath, project)

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
	response := &schema.Response[schema.CloudServerList]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.CloudServerList
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// GetCloudServer retrieves a specific cloud server by ID
func (s *CloudServerService) GetCloudServer(ctx context.Context, project string, cloudServerId string, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if cloudServerId == "" {
		return nil, fmt.Errorf("cloud server ID cannot be empty")
	}

	path := fmt.Sprintf(CloudServerPath, project, cloudServerId)

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
	response := &schema.Response[schema.CloudServerResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.CloudServerResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// CreateCloudServer creates a new cloud server
func (s *CloudServerService) CreateCloudServer(ctx context.Context, project string, body schema.CloudServerRequest, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(CloudServersPath, project)

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
	response := &schema.Response[schema.CloudServerResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.CloudServerResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// UpdateCloudServer updates an existing cloud server
func (s *CloudServerService) UpdateCloudServer(ctx context.Context, project string, cloudServerId string, body schema.CloudServerRequest, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if cloudServerId == "" {
		return nil, fmt.Errorf("cloud server ID cannot be empty")
	}

	path := fmt.Sprintf(CloudServerPath, project, cloudServerId)

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
	response := &schema.Response[schema.CloudServerResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.CloudServerResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// DeleteCloudServer deletes a cloud server by ID
func (s *CloudServerService) DeleteCloudServer(ctx context.Context, projectId string, cloudServerId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if cloudServerId == "" {
		return nil, fmt.Errorf("cloud server ID cannot be empty")
	}

	path := fmt.Sprintf(CloudServerPath, projectId, cloudServerId)

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

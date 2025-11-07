package compute

import (
	"context"
	"fmt"
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
func (s *CloudServerService) ListCloudServers(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
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

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetCloudServer retrieves a specific cloud server by ID
func (s *CloudServerService) GetCloudServer(ctx context.Context, project string, cloudServerId string, params *schema.RequestParameters) (*http.Response, error) {
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

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateCloudServer creates or updates a cloud server
func (s *CloudServerService) CreateOrUpdateCloudServer(ctx context.Context, project string, body schema.CloudServerRequest, params *schema.RequestParameters) (*http.Response, error) {
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

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteCloudServer deletes a cloud server by ID
func (s *CloudServerService) DeleteCloudServer(ctx context.Context, projectId string, cloudServerId string, params *schema.RequestParameters) (*http.Response, error) {
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

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ElasticIPService implements the ElasticIPAPI interface
type ElasticIPService struct {
	client *client.Client
}

// NewElasticIPService creates a new ElasticIPService
func NewElasticIPService(client *client.Client) *ElasticIPService {
	return &ElasticIPService{
		client: client,
	}
}

// ListElasticIPs retrieves all elastic IPs for a project
func (s *ElasticIPService) ListElasticIPs(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(ElasticIPsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetElasticIP retrieves a specific elastic IP by ID
func (s *ElasticIPService) GetElasticIP(ctx context.Context, project string, elasticIPId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if elasticIPId == "" {
		return nil, fmt.Errorf("elastic IP ID cannot be empty")
	}

	path := fmt.Sprintf(ElasticIPPath, project, elasticIPId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateElasticIP creates or updates an elastic IP
func (s *ElasticIPService) CreateOrUpdateElasticIP(ctx context.Context, project string, body schema.ElasticIPRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(ElasticIPsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteElasticIP deletes an elastic IP by ID
func (s *ElasticIPService) DeleteElasticIP(ctx context.Context, projectId string, elasticIPId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if elasticIPId == "" {
		return nil, fmt.Errorf("elastic IP ID cannot be empty")
	}

	path := fmt.Sprintf(ElasticIPPath, projectId, elasticIPId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

package network

import (
	"bytes"
	"context"
	"encoding/json"
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
func (s *ElasticIPService) ListElasticIPs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.ElasticList], error) {
	s.client.Logger().Debugf("Listing elastic IPs for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ElasticIPsPath, project)

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

	return schema.ParseResponseBody[schema.ElasticList](httpResp)
}

// GetElasticIP retrieves a specific elastic IP by ID
func (s *ElasticIPService) GetElasticIP(ctx context.Context, project string, elasticIPId string, params *schema.RequestParameters) (*schema.Response[schema.ElasticIpResponse], error) {
	s.client.Logger().Debugf("Getting elastic IP: %s in project: %s", elasticIPId, project)

	if err := schema.ValidateProjectAndResource(project, elasticIPId, "elastic IP ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ElasticIPPath, project, elasticIPId)

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

	return schema.ParseResponseBody[schema.ElasticIpResponse](httpResp)
}

// CreateElasticIP creates a new elastic IP
func (s *ElasticIPService) CreateElasticIP(ctx context.Context, project string, body schema.ElasticIpRequest, params *schema.RequestParameters) (*schema.Response[schema.ElasticIpResponse], error) {
	s.client.Logger().Debugf("Creating elastic IP in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ElasticIPsPath, project)

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

	return schema.ParseResponseBody[schema.ElasticIpResponse](httpResp)
}

// UpdateElasticIP updates an existing elastic IP
func (s *ElasticIPService) UpdateElasticIP(ctx context.Context, project string, elasticIPId string, body schema.ElasticIpRequest, params *schema.RequestParameters) (*schema.Response[schema.ElasticIpResponse], error) {
	s.client.Logger().Debugf("Updating elastic IP: %s in project: %s", elasticIPId, project)

	if err := schema.ValidateProjectAndResource(project, elasticIPId, "elastic IP ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ElasticIPPath, project, elasticIPId)

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

	return schema.ParseResponseBody[schema.ElasticIpResponse](httpResp)
}

// DeleteElasticIP deletes an elastic IP by ID
func (s *ElasticIPService) DeleteElasticIP(ctx context.Context, projectId string, elasticIPId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting elastic IP: %s in project: %s", elasticIPId, projectId)

	if err := schema.ValidateProjectAndResource(projectId, elasticIPId, "elastic IP ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ElasticIPPath, projectId, elasticIPId)

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

	return schema.ParseResponseBody[any](httpResp)
}

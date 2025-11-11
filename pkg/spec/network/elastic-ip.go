package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ListElasticIPs retrieves all elastic IPs for a project
func (s *Service) ListElasticIPs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.ElasticList], error) {
	s.client.Logger().Debugf("Listing elastic IPs for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ElasticIPsPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &ElasticIpListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ElasticIpListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.ElasticList](httpResp)
}

// GetElasticIP retrieves a specific elastic IP by ID
func (s *Service) GetElasticIP(ctx context.Context, project string, elasticIPId string, params *schema.RequestParameters) (*schema.Response[schema.ElasticIpResponse], error) {
	s.client.Logger().Debugf("Getting elastic IP: %s in project: %s", elasticIPId, project)

	if err := schema.ValidateProjectAndResource(project, elasticIPId, "elastic IP ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ElasticIPPath, project, elasticIPId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &ElasticIpGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ElasticIpGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.ElasticIpResponse](httpResp)
}

// CreateElasticIP creates a new elastic IP
func (s *Service) CreateElasticIP(ctx context.Context, project string, body schema.ElasticIpRequest, params *schema.RequestParameters) (*schema.Response[schema.ElasticIpResponse], error) {
	s.client.Logger().Debugf("Creating elastic IP in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ElasticIPsPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &ElasticIpCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ElasticIpCreateAPIVersion
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

	return schema.ParseResponseBody[schema.ElasticIpResponse](httpResp)
}

// UpdateElasticIP updates an existing elastic IP
func (s *Service) UpdateElasticIP(ctx context.Context, project string, elasticIPId string, body schema.ElasticIpRequest, params *schema.RequestParameters) (*schema.Response[schema.ElasticIpResponse], error) {
	s.client.Logger().Debugf("Updating elastic IP: %s in project: %s", elasticIPId, project)

	if err := schema.ValidateProjectAndResource(project, elasticIPId, "elastic IP ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ElasticIPPath, project, elasticIPId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &ElasticIpUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ElasticIpUpdateAPIVersion
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

	return schema.ParseResponseBody[schema.ElasticIpResponse](httpResp)
}

// DeleteElasticIP deletes an elastic IP by ID
func (s *Service) DeleteElasticIP(ctx context.Context, projectId string, elasticIPId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting elastic IP: %s in project: %s", elasticIPId, projectId)

	if err := schema.ValidateProjectAndResource(projectId, elasticIPId, "elastic IP ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(ElasticIPPath, projectId, elasticIPId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &ElasticIpDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &ElasticIpDeleteAPIVersion
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

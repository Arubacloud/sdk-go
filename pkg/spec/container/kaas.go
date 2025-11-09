package container

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

// KaaSService implements the KaaSAPI interface
type KaaSService struct {
	client *client.Client
}

// NewKaaSService creates a new KaaSService
func NewKaaSService(client *client.Client) *KaaSService {
	return &KaaSService{
		client: client,
	}
}

// ListKaaS retrieves all KaaS clusters for a project
func (s *KaaSService) ListKaaS(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KaaSList], error) {
	s.client.Logger().Debugf("Listing KaaS clusters for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KaaSPath, project)

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

	return schema.ParseResponseBody[schema.KaaSList](httpResp)
}

// GetKaaS retrieves a specific KaaS cluster by ID
func (s *KaaSService) GetKaaS(ctx context.Context, project string, kaasId string, params *schema.RequestParameters) (*schema.Response[schema.KaaSResponse], error) {
	s.client.Logger().Debugf("Getting KaaS cluster: %s in project: %s", kaasId, project)

	if err := schema.ValidateProjectAndResource(project, kaasId, "KaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KaaSItemPath, project, kaasId)

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

	return schema.ParseResponseBody[schema.KaaSResponse](httpResp)
}

// CreateKaaS creates a new KaaS cluster
func (s *KaaSService) CreateKaaS(ctx context.Context, project string, body schema.KaaSRequest, params *schema.RequestParameters) (*schema.Response[schema.KaaSResponse], error) {
	s.client.Logger().Debugf("Creating KaaS cluster in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KaaSPath, project)

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
	response := &schema.Response[schema.KaaSResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.KaaSResponse
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

// UpdateKaaS updates an existing KaaS cluster
func (s *KaaSService) UpdateKaaS(ctx context.Context, project string, kaasId string, body schema.KaaSRequest, params *schema.RequestParameters) (*schema.Response[schema.KaaSResponse], error) {
	s.client.Logger().Debugf("Updating KaaS cluster: %s in project: %s", kaasId, project)

	if err := schema.ValidateProjectAndResource(project, kaasId, "KaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KaaSItemPath, project, kaasId)

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
	response := &schema.Response[schema.KaaSResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.KaaSResponse
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

// DeleteKaaS deletes a KaaS cluster by ID
func (s *KaaSService) DeleteKaaS(ctx context.Context, projectId string, kaasId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting KaaS cluster: %s in project: %s", kaasId, projectId)

	if err := schema.ValidateProjectAndResource(projectId, kaasId, "KaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KaaSItemPath, projectId, kaasId)

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

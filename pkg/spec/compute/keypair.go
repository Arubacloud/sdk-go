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

// KeyPairService implements the KeyPairAPI interface
type KeyPairService struct {
	client *client.Client
}

// NewKeyPairService creates a new KeyPairService
func NewKeyPairService(client *client.Client) *KeyPairService {
	return &KeyPairService{
		client: client,
	}
}

// ListKeyPairs retrieves all key pairs for a project
func (s *KeyPairService) ListKeyPairs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KeyPairListResponse], error) {
	s.client.Logger().Debugf("Listing key pairs for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KeyPairsPath, project)

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

	return schema.ParseResponseBody[schema.KeyPairListResponse](httpResp)
}

// GetKeyPair retrieves a specific key pair by ID
func (s *KeyPairService) GetKeyPair(ctx context.Context, project string, keyPairId string, params *schema.RequestParameters) (*schema.Response[schema.KeyPairResponse], error) {
	s.client.Logger().Debugf("Getting key pair: %s in project: %s", keyPairId, project)

	if err := schema.ValidateProjectAndResource(project, keyPairId, "key pair ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KeyPairPath, project, keyPairId)

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

	return schema.ParseResponseBody[schema.KeyPairResponse](httpResp)
}

// CreateKeyPair creates a new key pair
func (s *KeyPairService) CreateKeyPair(ctx context.Context, project string, body schema.KeyPairRequest, params *schema.RequestParameters) (*schema.Response[schema.KeyPairResponse], error) {
	s.client.Logger().Debugf("Creating key pair in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KeyPairsPath, project)

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
	response := &schema.Response[schema.KeyPairResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.KeyPairResponse
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

// DeleteKeyPair deletes a key pair by ID
func (s *KeyPairService) DeleteKeyPair(ctx context.Context, projectId string, keyPairId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting key pair: %s in project: %s", keyPairId, projectId)

	if err := schema.ValidateProjectAndResource(projectId, keyPairId, "key pair ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KeyPairPath, projectId, keyPairId)

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

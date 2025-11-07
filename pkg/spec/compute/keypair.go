package compute

import (
	"context"
	"fmt"
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
func (s *KeyPairService) ListKeyPairs(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(KeyPairsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetKeyPair retrieves a specific key pair by ID
func (s *KeyPairService) GetKeyPair(ctx context.Context, project string, keyPairId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if keyPairId == "" {
		return nil, fmt.Errorf("key pair ID cannot be empty")
	}

	path := fmt.Sprintf(KeyPairPath, project, keyPairId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateKeyPair creates or updates a key pair
func (s *KeyPairService) CreateOrUpdateKeyPair(ctx context.Context, project string, body schema.KeyPairRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(KeyPairsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteKeyPair deletes a key pair by ID
func (s *KeyPairService) DeleteKeyPair(ctx context.Context, projectId string, keyPairId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if keyPairId == "" {
		return nil, fmt.Errorf("key pair ID cannot be empty")
	}

	path := fmt.Sprintf(KeyPairPath, projectId, keyPairId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

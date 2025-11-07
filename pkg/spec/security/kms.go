package security

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// KmsKeyService implements the KmsKeyAPI interface
type KmsKeyService struct {
	client *client.Client
}

// NewKmsKeyService creates a new KmsKeyService
func NewKmsKeyService(client *client.Client) *KmsKeyService {
	return &KmsKeyService{
		client: client,
	}
}

// ListKmsKeys retrieves all KMS keys for a project
func (s *KmsKeyService) ListKmsKeys(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(KmsKeysPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetKmsKey retrieves a specific KMS key by ID
func (s *KmsKeyService) GetKmsKey(ctx context.Context, project string, kmsKeyId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if kmsKeyId == "" {
		return nil, fmt.Errorf("KMS key ID cannot be empty")
	}

	path := fmt.Sprintf(KmsKeyPath, project, kmsKeyId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateKmsKey creates or updates a KMS key
func (s *KmsKeyService) CreateOrUpdateKmsKey(ctx context.Context, project string, body schema.KmsKeyRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(KmsKeysPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteKmsKey deletes a KMS key by ID
func (s *KmsKeyService) DeleteKmsKey(ctx context.Context, projectId string, kmsKeyId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if kmsKeyId == "" {
		return nil, fmt.Errorf("KMS key ID cannot be empty")
	}

	path := fmt.Sprintf(KmsKeyPath, projectId, kmsKeyId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// BlockStorageService implements the BlockStorageAPI interface
type BlockStorageService struct {
	client *client.Client
}

// NewBlockStorageService creates a new BlockStorageService
func NewBlockStorageService(client *client.Client) *BlockStorageService {
	return &BlockStorageService{
		client: client,
	}
}

// ListBlockStorages retrieves all block storages for a project
func (s *BlockStorageService) ListBlockStorages(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(BlockStoragesPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetBlockStorage retrieves a specific block storage by ID
func (s *BlockStorageService) GetBlockStorage(ctx context.Context, project string, blockStorageId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if blockStorageId == "" {
		return nil, fmt.Errorf("block storage ID cannot be empty")
	}

	path := fmt.Sprintf(BlockStoragePath, project, blockStorageId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateBlockStorage creates or updates a block storage
func (s *BlockStorageService) CreateOrUpdateBlockStorage(ctx context.Context, project string, body schema.BlockStorageRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(BlockStoragesPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteBlockStorage deletes a block storage by ID
func (s *BlockStorageService) DeleteBlockStorage(ctx context.Context, projectId string, blockStorageId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if blockStorageId == "" {
		return nil, fmt.Errorf("block storage ID cannot be empty")
	}

	path := fmt.Sprintf(BlockStoragePath, projectId, blockStorageId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

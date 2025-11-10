package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ListBlockStorageVolumes retrieves all block storage volumes for a project
func (s *Service) ListBlockStorageVolumes(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.BlockStorageList], error) {
	s.client.Logger().Debugf("Listing block storage volumes for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BlockStoragesPath, project)

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

	return schema.ParseResponseBody[schema.BlockStorageList](httpResp)
}

// GetBlockStorageVolume retrieves a specific block storage volume by ID
func (s *Service) GetBlockStorageVolume(ctx context.Context, project string, volumeId string, params *schema.RequestParameters) (*schema.Response[schema.BlockStorageResponse], error) {
	s.client.Logger().Debugf("Getting block storage volume: %s in project: %s", volumeId, project)

	if err := schema.ValidateProjectAndResource(project, volumeId, "block storage ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BlockStoragePath, project, volumeId)

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

	return schema.ParseResponseBody[schema.BlockStorageResponse](httpResp)
}

// CreateBlockStorageVolume creates a new block storage volume
func (s *Service) CreateBlockStorageVolume(ctx context.Context, project string, body schema.BlockStorageRequest, params *schema.RequestParameters) (*schema.Response[schema.BlockStorageResponse], error) {
	s.client.Logger().Debugf("Creating block storage volume in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BlockStoragesPath, project)

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

	return schema.ParseResponseBody[schema.BlockStorageResponse](httpResp)
}

// DeleteBlockStorageVolume deletes a block storage volume by ID
func (s *Service) DeleteBlockStorageVolume(ctx context.Context, project string, volumeId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting block storage volume: %s in project: %s", volumeId, project)

	if err := schema.ValidateProjectAndResource(project, volumeId, "block storage ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BlockStoragePath, project, volumeId)

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

package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// ListBlockStorageVolumes retrieves all block storage volumes for a project
func (s *BlockStorageService) ListBlockStorageVolumes(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.BlockStorageList], error) {
	s.client.Logger().Debugf("Listing block storage volumes for project: %s", project)

	if err := validateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BlockStorageVolumesPath, project)

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

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.BlockStorageList]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.BlockStorageList
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// GetBlockStorageVolume retrieves a specific block storage volume by ID
func (s *BlockStorageService) GetBlockStorageVolume(ctx context.Context, project string, volumeId string, params *schema.RequestParameters) (*schema.Response[schema.BlockStorageResponse], error) {
	s.client.Logger().Debugf("Getting block storage volume: %s in project: %s", volumeId, project)

	if err := validateProjectAndResource(project, volumeId, "block storage ID"); err != nil {
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

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.BlockStorageResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.BlockStorageResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// CreateBlockStorageVolume creates a new block storage volume
func (s *BlockStorageService) CreateBlockStorageVolume(ctx context.Context, project string, body schema.BlockStorageRequest, params *schema.RequestParameters) (*schema.Response[schema.BlockStorageResponse], error) {
	s.client.Logger().Debugf("Creating block storage volume in project: %s", project)

	if err := validateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BlockStorageVolumesPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.BlockStorageResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.BlockStorageResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// DeleteBlockStorageVolume deletes a block storage volume by ID
func (s *BlockStorageService) DeleteBlockStorageVolume(ctx context.Context, project string, volumeId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting block storage volume: %s in project: %s", volumeId, project)

	if err := validateProjectAndResource(project, volumeId, "block storage ID"); err != nil {
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

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[any]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() && len(bodyBytes) > 0 {
		var data any
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

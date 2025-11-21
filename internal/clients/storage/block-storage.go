package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/types"
)

type volumesClientImpl struct {
	client *restclient.Client
}

// NewVolumesClientImpl creates a new unified Storage service
func NewVolumesClientImpl(client *restclient.Client) *volumesClientImpl {
	return &volumesClientImpl{
		client: client,
	}
}

// List retrieves all block storage volumes for a project
func (c *volumesClientImpl) List(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.BlockStorageList], error) {
	c.client.Logger().Debugf("Listing block storage volumes for project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BlockStoragesPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &BlockStorageListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &BlockStorageListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.BlockStorageList](httpResp)
}

// Get retrieves a specific block storage volume by ID
func (c *volumesClientImpl) Get(ctx context.Context, project string, volumeId string, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error) {
	c.client.Logger().Debugf("Getting block storage volume: %s in project: %s", volumeId, project)

	if err := types.ValidateProjectAndResource(project, volumeId, "block storage ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BlockStoragePath, project, volumeId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &BlockStorageGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &BlockStorageGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.BlockStorageResponse](httpResp)
}

// Create creates a new block storage volume
func (c *volumesClientImpl) Create(ctx context.Context, project string, body types.BlockStorageRequest, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error) {
	c.client.Logger().Debugf("Creating block storage volume in project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BlockStoragesPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &BlockStorageCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &BlockStorageCreateAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	// Marshal the request body to JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := c.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.BlockStorageResponse](httpResp)
}

// Delete deletes a block storage volume by ID
func (c *volumesClientImpl) Delete(ctx context.Context, project string, volumeId string, params *types.RequestParameters) (*types.Response[any], error) {
	c.client.Logger().Debugf("Deleting block storage volume: %s in project: %s", volumeId, project)

	if err := types.ValidateProjectAndResource(project, volumeId, "block storage ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BlockStoragePath, project, volumeId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &BlockStorageDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &BlockStorageDeleteAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[any](httpResp)
}

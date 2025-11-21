package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/types"
)

// snapshotsClientImpl implements the StorageAPI interface for all Storage operations
type snapshotsClientImpl struct {
	client        *restclient.Client
	volumesClient *volumesClientImpl
}

// NewSnapshotsClientImpl creates a new unified Storage service
func NewSnapshotsClientImpl(client *restclient.Client, volumesClient *volumesClientImpl) *snapshotsClientImpl {
	return &snapshotsClientImpl{
		client:        client,
		volumesClient: volumesClient,
	}
}

// List retrieves all snapshots for a project
func (c *snapshotsClientImpl) List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.SnapshotList], error) {
	c.client.Logger().Debugf("Listing snapshots for project: %s", projectID)

	if err := types.ValidateProject(projectID); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SnapshotsPath, projectID)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SnapshotListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SnapshotListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SnapshotList](httpResp)
}

// Get retrieves a specific snapshot by ID
func (c *snapshotsClientImpl) Get(ctx context.Context, projectID string, snapshotId string, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error) {
	c.client.Logger().Debugf("Getting snapshot: %s in project: %s", snapshotId, projectID)

	if err := types.ValidateProjectAndResource(projectID, snapshotId, "snapshot ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SnapshotPath, projectID, snapshotId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SnapshotGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SnapshotGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.SnapshotResponse](httpResp)
}

// Create creates a new snapshot
// The SDK automatically waits for the source BlockStorage volume to become Active or NotUsed before creating the snapshot
func (c *snapshotsClientImpl) Create(ctx context.Context, projectID string, body types.SnapshotRequest, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error) {
	c.client.Logger().Debugf("Creating snapshot in project: %s", projectID)

	if err := types.ValidateProject(projectID); err != nil {
		return nil, err
	}

	// Extract volume ID from the Volume URI if present
	if body.Properties.Volume.URI != "" {
		// Parse URI to get volume ID: /projects/{project}/providers/Aruba.Storage/blockstorages/{volumeId}
		volumeID, err := extractVolumeIDFromURI(body.Properties.Volume.URI)
		if err == nil && volumeID != "" {
			// Wait for BlockStorage to become Active or NotUsed before creating snapshot
			err := waitForBlockStorageActive(ctx, c.volumesClient, projectID, volumeID)
			if err != nil {
				return nil, fmt.Errorf("failed waiting for BlockStorage to become ready: %w", err)
			}
		}
	}

	path := fmt.Sprintf(SnapshotsPath, projectID)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SnapshotCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SnapshotCreateAPIVersion
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

	return types.ParseResponseBody[types.SnapshotResponse](httpResp)
}

// Delete deletes a snapshot by ID
func (c *snapshotsClientImpl) Delete(ctx context.Context, projectID string, snapshotId string, params *types.RequestParameters) (*types.Response[any], error) {
	c.client.Logger().Debugf("Deleting snapshot: %s in project: %s", snapshotId, projectID)

	if err := types.ValidateProjectAndResource(projectID, snapshotId, "snapshot ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SnapshotPath, projectID, snapshotId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &SnapshotDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SnapshotDeleteAPIVersion
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

// extractVolumeIDFromURI extracts the volume ID from a volume URI
// URI format: /projects/{project}/providers/Aruba.Storage/blockstorages/{volumeId}
func extractVolumeIDFromURI(uri string) (string, error) {
	parts := strings.Split(uri, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid URI format: %s", uri)
	}
	// The volume ID is the last part of the URI
	volumeID := parts[len(parts)-1]
	if volumeID == "" {
		return "", fmt.Errorf("could not extract volume ID from URI: %s", uri)
	}
	return volumeID, nil
}

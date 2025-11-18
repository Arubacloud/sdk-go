package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

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

// ListSnapshots retrieves all snapshots for a project
func (s *Service) ListSnapshots(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.SnapshotList], error) {
	s.client.Logger().Debugf("Listing snapshots for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SnapshotsPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &SnapshotListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SnapshotListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.SnapshotList](httpResp)
}

// GetSnapshot retrieves a specific snapshot by ID
func (s *Service) GetSnapshot(ctx context.Context, project string, snapshotId string, params *schema.RequestParameters) (*schema.Response[schema.SnapshotResponse], error) {
	s.client.Logger().Debugf("Getting snapshot: %s in project: %s", snapshotId, project)

	if err := schema.ValidateProjectAndResource(project, snapshotId, "snapshot ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SnapshotPath, project, snapshotId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &SnapshotGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SnapshotGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.SnapshotResponse](httpResp)
}

// CreateSnapshot creates a new snapshot
// The SDK automatically waits for the source BlockStorage volume to become Active or NotUsed before creating the snapshot
func (s *Service) CreateSnapshot(ctx context.Context, project string, body schema.SnapshotRequest, params *schema.RequestParameters) (*schema.Response[schema.SnapshotResponse], error) {
	s.client.Logger().Debugf("Creating snapshot in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	// Extract volume ID from the Volume URI if present
	if body.Properties.Volume.URI != "" {
		// Parse URI to get volume ID: /projects/{project}/providers/Aruba.Storage/blockstorages/{volumeId}
		volumeID, err := extractVolumeIDFromURI(body.Properties.Volume.URI)
		if err == nil && volumeID != "" {
			// Wait for BlockStorage to become Active or NotUsed before creating snapshot
			err := s.waitForBlockStorageActive(ctx, project, volumeID)
			if err != nil {
				return nil, fmt.Errorf("failed waiting for BlockStorage to become ready: %w", err)
			}
		}
	}

	path := fmt.Sprintf(SnapshotsPath, project)

	if params == nil {
		params = &schema.RequestParameters{
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

	httpResp, err := s.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.SnapshotResponse](httpResp)
}

// DeleteSnapshot deletes a snapshot by ID
func (s *Service) DeleteSnapshot(ctx context.Context, project string, snapshotId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting snapshot: %s in project: %s", snapshotId, project)

	if err := schema.ValidateProjectAndResource(project, snapshotId, "snapshot ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SnapshotPath, project, snapshotId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &SnapshotDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &SnapshotDeleteAPIVersion
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

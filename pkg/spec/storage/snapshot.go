package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// SnapshotService implements the SnapshotAPI interface
type SnapshotService struct {
	client *client.Client
}

// NewSnapshotService creates a new SnapshotService
func NewSnapshotService(client *client.Client) *SnapshotService {
	return &SnapshotService{
		client: client,
	}
}

// ListSnapshots retrieves all snapshots for a project
func (s *SnapshotService) ListSnapshots(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(SnapshotsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetSnapshot retrieves a specific snapshot by ID
func (s *SnapshotService) GetSnapshot(ctx context.Context, project string, snapshotId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if snapshotId == "" {
		return nil, fmt.Errorf("snapshot ID cannot be empty")
	}

	path := fmt.Sprintf(SnapshotPath, project, snapshotId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateSnapshot creates or updates a snapshot
func (s *SnapshotService) CreateOrUpdateSnapshot(ctx context.Context, project string, body schema.SnapshotRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(SnapshotsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteSnapshot deletes a snapshot by ID
func (s *SnapshotService) DeleteSnapshot(ctx context.Context, projectId string, snapshotId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if snapshotId == "" {
		return nil, fmt.Errorf("snapshot ID cannot be empty")
	}

	path := fmt.Sprintf(SnapshotPath, projectId, snapshotId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

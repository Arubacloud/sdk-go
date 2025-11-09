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
func (s *SnapshotService) ListSnapshots(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.SnapshotList], error) {
	s.client.Logger().Debugf("Listing snapshots for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SnapshotsPath, project)

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

	return schema.ParseResponseBody[schema.SnapshotList](httpResp)
}

// GetSnapshot retrieves a specific snapshot by ID
func (s *SnapshotService) GetSnapshot(ctx context.Context, project string, snapshotId string, params *schema.RequestParameters) (*schema.Response[schema.SnapshotResponse], error) {
	s.client.Logger().Debugf("Getting snapshot: %s in project: %s", snapshotId, project)

	if err := schema.ValidateProjectAndResource(project, snapshotId, "snapshot ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SnapshotPath, project, snapshotId)

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

	return schema.ParseResponseBody[schema.SnapshotResponse](httpResp)
}

// CreateSnapshot creates a new snapshot
func (s *SnapshotService) CreateSnapshot(ctx context.Context, project string, body schema.SnapshotRequest, params *schema.RequestParameters) (*schema.Response[schema.SnapshotResponse], error) {
	s.client.Logger().Debugf("Creating snapshot in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SnapshotsPath, project)

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

	return schema.ParseResponseBody[schema.SnapshotResponse](httpResp)
}

// DeleteSnapshot deletes a snapshot by ID
func (s *SnapshotService) DeleteSnapshot(ctx context.Context, project string, snapshotId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting snapshot: %s in project: %s", snapshotId, project)

	if err := schema.ValidateProjectAndResource(project, snapshotId, "snapshot ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SnapshotPath, project, snapshotId)

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

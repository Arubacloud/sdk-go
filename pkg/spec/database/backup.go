package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// BackupService implements the BackupAPI interface
type BackupService struct {
	client *client.Client
}

// NewBackupService creates a new BackupService
func NewBackupService(client *client.Client) *BackupService {
	return &BackupService{
		client: client,
	}
}

// ListBackups retrieves all backups for a project
func (s *BackupService) ListBackups(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.BackupList], error) {
	s.client.Logger().Debugf("Listing backups for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupsPath, project)

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

	return schema.ParseResponseBody[schema.BackupList](httpResp)
}

// GetBackup retrieves a specific backup by ID
func (s *BackupService) GetBackup(ctx context.Context, project string, backupId string, params *schema.RequestParameters) (*schema.Response[schema.BackupResponse], error) {
	s.client.Logger().Debugf("Getting backup: %s in project: %s", backupId, project)

	if err := schema.ValidateProjectAndResource(project, backupId, "backup ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupPath, project, backupId)

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

	return schema.ParseResponseBody[schema.BackupResponse](httpResp)
}

// CreateBackup creates a new backup
func (s *BackupService) CreateBackup(ctx context.Context, project string, body schema.BackupRequest, params *schema.RequestParameters) (*schema.Response[schema.BackupResponse], error) {
	s.client.Logger().Debugf("Creating backup in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupsPath, project)

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

	return schema.ParseResponseBody[schema.BackupResponse](httpResp)
}

// DeleteBackup deletes a backup by ID
func (s *BackupService) DeleteBackup(ctx context.Context, projectId string, backupId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting backup: %s in project: %s", backupId, projectId)

	if err := schema.ValidateProjectAndResource(projectId, backupId, "backup ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupPath, projectId, backupId)

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

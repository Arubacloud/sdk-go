package database

import (
	"context"
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
func (s *BackupService) ListBackups(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(BackupsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetBackup retrieves a specific backup by ID
func (s *BackupService) GetBackup(ctx context.Context, project string, backupId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if backupId == "" {
		return nil, fmt.Errorf("backup ID cannot be empty")
	}

	path := fmt.Sprintf(BackupPath, project, backupId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateBackup creates a new backup
func (s *BackupService) CreateBackup(ctx context.Context, project string, body schema.BackupRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(BackupsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPost, path, body, queryParams, headers)
}

// DeleteBackup deletes a backup by ID
func (s *BackupService) DeleteBackup(ctx context.Context, projectId string, backupId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if backupId == "" {
		return nil, fmt.Errorf("backup ID cannot be empty")
	}

	path := fmt.Sprintf(BackupPath, projectId, backupId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

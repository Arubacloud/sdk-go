package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/types"
)

// ListBackups retrieves all backups for a project
func (s *Service) ListBackups(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.BackupList], error) {
	s.client.Logger().Debugf("Listing backups for project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupsPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseBackupListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseBackupListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.BackupList](httpResp)
}

// GetBackup retrieves a specific backup by ID
func (s *Service) GetBackup(ctx context.Context, project string, backupId string, params *types.RequestParameters) (*types.Response[types.BackupResponse], error) {
	s.client.Logger().Debugf("Getting backup: %s in project: %s", backupId, project)

	if err := types.ValidateProjectAndResource(project, backupId, "backup ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupPath, project, backupId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseBackupGetVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseBackupGetVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.BackupResponse](httpResp)
}

// CreateBackup creates a new backup
func (s *Service) CreateBackup(ctx context.Context, project string, body types.BackupRequest, params *types.RequestParameters) (*types.Response[types.BackupResponse], error) {
	s.client.Logger().Debugf("Creating backup in project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupsPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseBackupCreateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseBackupCreateVersion
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

	return types.ParseResponseBody[types.BackupResponse](httpResp)
}

// DeleteBackup deletes a backup by ID
func (s *Service) DeleteBackup(ctx context.Context, projectId string, backupId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting backup: %s in project: %s", backupId, projectId)

	if err := types.ValidateProjectAndResource(projectId, backupId, "backup ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupPath, projectId, backupId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseBackupDeleteVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseBackupDeleteVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[any](httpResp)
}

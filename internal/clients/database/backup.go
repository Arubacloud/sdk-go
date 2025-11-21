package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/types"
)

type backupsClientImpl struct {
	client *restclient.Client
}

// NewService creates a new unified Database service
func NewBackupsClientImpl(client *restclient.Client) *backupsClientImpl {
	return &backupsClientImpl{
		client: client,
	}
}

// List retrieves all backups for a project
func (c *backupsClientImpl) List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.BackupList], error) {
	c.client.Logger().Debugf("Listing backups for project: %s", projectID)

	if err := types.ValidateProject(projectID); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupsPath, projectID)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseBackupListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseBackupListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.BackupList](httpResp)
}

// Get retrieves a specific backup by ID
func (c *backupsClientImpl) Get(ctx context.Context, projectID string, backupID string, params *types.RequestParameters) (*types.Response[types.BackupResponse], error) {
	c.client.Logger().Debugf("Getting backup: %s in project: %s", backupID, projectID)

	if err := types.ValidateProjectAndResource(projectID, backupID, "backup ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupPath, projectID, backupID)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseBackupGetVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseBackupGetVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.BackupResponse](httpResp)
}

// Create creates a new backup
func (c *backupsClientImpl) Create(ctx context.Context, projectID string, body types.BackupRequest, params *types.RequestParameters) (*types.Response[types.BackupResponse], error) {
	c.client.Logger().Debugf("Creating backup in project: %s", projectID)

	if err := types.ValidateProject(projectID); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupsPath, projectID)

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

	httpResp, err := c.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.BackupResponse](httpResp)
}

// Delete deletes a backup by ID
func (c *backupsClientImpl) Delete(ctx context.Context, projectID string, backupID string, params *types.RequestParameters) (*types.Response[any], error) {
	c.client.Logger().Debugf("Deleting backup: %s in project: %s", backupID, projectID)

	if err := types.ValidateProjectAndResource(projectID, backupID, "backup ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(BackupPath, projectID, backupID)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseBackupDeleteVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseBackupDeleteVersion
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

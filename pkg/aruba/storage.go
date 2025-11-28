package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

type StorageClient interface {
	Snapshots() SnapshotsClient
	Volumes() VolumesClient
	Backups() StorageBackupsClient
}

type storageClientImpl struct {
	snapshotsClient SnapshotsClient
	volumesClient   VolumesClient
	backupsClient   StorageBackupsClient
	restoresClient  StorageRestoreClient
}

var _ StorageClient = (*storageClientImpl)(nil)

func (c *storageClientImpl) Snapshots() SnapshotsClient {
	return c.snapshotsClient
}

func (c *storageClientImpl) Volumes() VolumesClient {
	return c.volumesClient
}

func (c *storageClientImpl) Backups() StorageBackupsClient {
	return c.backupsClient
}

func (c *storageClientImpl) Restores() StorageRestoreClient {
	return c.restoresClient
}

type SnapshotsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.SnapshotList], error)
	Get(ctx context.Context, projectID string, snapshotID string, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error)
	Update(ctx context.Context, projectID string, snapshotID string, body types.SnapshotRequest, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error)
	Create(ctx context.Context, projectID string, body types.SnapshotRequest, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error)
	Delete(ctx context.Context, projectID string, snapshotID string, params *types.RequestParameters) (*types.Response[any], error)
}

type VolumesClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.BlockStorageList], error)
	Get(ctx context.Context, projectID string, volumeID string, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error)
	Update(ctx context.Context, projectID string, volumeID string, body types.BlockStorageRequest, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error)
	Create(ctx context.Context, projectID string, body types.BlockStorageRequest, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error)
	Delete(ctx context.Context, projectID string, volumeID string, params *types.RequestParameters) (*types.Response[any], error)
}

type StorageBackupsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.StorageBackupList], error)
	Get(ctx context.Context, projectID string, backupID string, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	Update(ctx context.Context, projectID string, backupID string, body types.StorageBackupRequest, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	Create(ctx context.Context, projectID string, body types.StorageBackupRequest, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	Delete(ctx context.Context, projectID string, backupID string, params *types.RequestParameters) (*types.Response[any], error)
}

type StorageRestoreClient interface {
	List(ctx context.Context, projectID string, backupID string, params *types.RequestParameters) (*types.Response[types.RestoreList], error)
	Get(ctx context.Context, projectID string, backupID string, restoreID string, params *types.RequestParameters) (*types.Response[types.RestoreResponse], error)
	Update(ctx context.Context, projectID string, backupID string, restoreID string, body types.RestoreRequest, params *types.RequestParameters) (*types.Response[types.RestoreResponse], error)
	Create(ctx context.Context, projectID string, backupID string, body types.RestoreRequest, params *types.RequestParameters) (*types.Response[types.RestoreResponse], error)
	Delete(ctx context.Context, projectID string, backupID string, restoreID string, params *types.RequestParameters) (*types.Response[any], error)
}

package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

type StorageClient interface {
	Snapshots() SnapshotsClient
	Volumes() VolumesClient
}

type storageClientImpl struct {
	snapshotsClient SnapshotsClient
	volumesClient   VolumesClient
}

var _ StorageClient = (*storageClientImpl)(nil)

func (c *storageClientImpl) Snapshots() SnapshotsClient {
	return c.snapshotsClient
}

func (c *storageClientImpl) Volumes() VolumesClient {
	return c.volumesClient
}

type SnapshotsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.SnapshotList], error)
	Get(ctx context.Context, projectID string, snapshotID string, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error)
	Create(ctx context.Context, projectID string, body types.SnapshotRequest, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error)
	Delete(ctx context.Context, projectID string, snapshotID string, params *types.RequestParameters) (*types.Response[any], error)
}

type VolumesClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.BlockStorageList], error)
	Get(ctx context.Context, projectID string, volumeID string, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error)
	Create(ctx context.Context, projectID string, body types.BlockStorageRequest, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error)
	Delete(ctx context.Context, projectID string, volumeID string, params *types.RequestParameters) (*types.Response[any], error)
}

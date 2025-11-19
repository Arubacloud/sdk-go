package storage

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

// StorageAPI defines the unified interface for all Storage operations
type StorageAPI interface {
	// BlockStorage operations
	ListBlockStorageVolumes(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.BlockStorageList], error)
	GetBlockStorageVolume(ctx context.Context, project string, volumeId string, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error)
	CreateBlockStorageVolume(ctx context.Context, project string, body types.BlockStorageRequest, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error)
	DeleteBlockStorageVolume(ctx context.Context, project string, volumeId string, params *types.RequestParameters) (*types.Response[any], error)

	// Snapshot operations
	ListSnapshots(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.SnapshotList], error)
	GetSnapshot(ctx context.Context, project string, snapshotId string, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error)
	CreateSnapshot(ctx context.Context, project string, body types.SnapshotRequest, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error)
	DeleteSnapshot(ctx context.Context, project string, snapshotId string, params *types.RequestParameters) (*types.Response[any], error)
}

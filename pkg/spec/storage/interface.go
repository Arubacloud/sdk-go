package storage

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// StorageAPI defines the unified interface for all Storage operations
type StorageAPI interface {
	// BlockStorage operations
	ListBlockStorageVolumes(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.BlockStorageList], error)
	GetBlockStorageVolume(ctx context.Context, project string, volumeId string, params *schema.RequestParameters) (*schema.Response[schema.BlockStorageResponse], error)
	CreateBlockStorageVolume(ctx context.Context, project string, body schema.BlockStorageRequest, params *schema.RequestParameters) (*schema.Response[schema.BlockStorageResponse], error)
	DeleteBlockStorageVolume(ctx context.Context, project string, volumeId string, params *schema.RequestParameters) (*schema.Response[any], error)

	// Snapshot operations
	ListSnapshots(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.SnapshotList], error)
	GetSnapshot(ctx context.Context, project string, snapshotId string, params *schema.RequestParameters) (*schema.Response[schema.SnapshotResponse], error)
	CreateSnapshot(ctx context.Context, project string, body schema.SnapshotRequest, params *schema.RequestParameters) (*schema.Response[schema.SnapshotResponse], error)
	DeleteSnapshot(ctx context.Context, project string, snapshotId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

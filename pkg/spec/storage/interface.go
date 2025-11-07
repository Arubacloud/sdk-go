package storage

import (
	"context"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// BlockStorageAPI defines the interface for managing block storage volumes.
type BlockStorageAPI interface {
	ListBlockStorageVolumes(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetBlockStorageVolume(ctx context.Context, project string, volumeId string, params *schema.RequestParameters) (*http.Response, error)
	CreateBlockStorageVolume(ctx context.Context, project string, body schema.BlockStorageRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteBlockStorageVolume(ctx context.Context, project string, volumeId string, params *schema.RequestParameters) (*http.Response, error)
}

// SnapshotAPI defines the interface for managing storage snapshots.
type SnapshotAPI interface {
	ListSnapshots(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetSnapshot(ctx context.Context, project string, snapshotId string, params *schema.RequestParameters) (*http.Response, error)
	CreateSnapshot(ctx context.Context, project string, body schema.SnapshotRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteSnapshot(ctx context.Context, project string, snapshotId string, params *schema.RequestParameters) (*http.Response, error)
}

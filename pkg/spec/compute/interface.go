package compute

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

// ComputeAPI defines the unified interface for all Compute operations (CloudServers, KeyPairs)
type ComputeAPI interface {
	// CloudServer operations
	ListCloudServers(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.CloudServerList], error)
	GetCloudServer(ctx context.Context, project string, cloudServerId string, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	CreateCloudServer(ctx context.Context, project string, body types.CloudServerRequest, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	UpdateCloudServer(ctx context.Context, project string, cloudServerId string, body types.CloudServerRequest, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	DeleteCloudServer(ctx context.Context, projectId string, cloudServerId string, params *types.RequestParameters) (*types.Response[any], error)

	// KeyPair operations
	ListKeyPairs(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.KeyPairListResponse], error)
	GetKeyPair(ctx context.Context, project string, keyPairId string, params *types.RequestParameters) (*types.Response[types.KeyPairResponse], error)
	CreateKeyPair(ctx context.Context, project string, body types.KeyPairRequest, params *types.RequestParameters) (*types.Response[types.KeyPairResponse], error)
	DeleteKeyPair(ctx context.Context, projectId string, keyPairId string, params *types.RequestParameters) (*types.Response[any], error)
}

package compute

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ComputeAPI defines the unified interface for all Compute operations (CloudServers, KeyPairs)
type ComputeAPI interface {
	// CloudServer operations
	ListCloudServers(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.CloudServerList], error)
	GetCloudServer(ctx context.Context, project string, cloudServerId string, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
	CreateCloudServer(ctx context.Context, project string, body schema.CloudServerRequest, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
	UpdateCloudServer(ctx context.Context, project string, cloudServerId string, body schema.CloudServerRequest, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
	DeleteCloudServer(ctx context.Context, projectId string, cloudServerId string, params *schema.RequestParameters) (*schema.Response[any], error)

	// KeyPair operations
	ListKeyPairs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KeyPairListResponse], error)
	GetKeyPair(ctx context.Context, project string, keyPairId string, params *schema.RequestParameters) (*schema.Response[schema.KeyPairResponse], error)
	CreateKeyPair(ctx context.Context, project string, body schema.KeyPairRequest, params *schema.RequestParameters) (*schema.Response[schema.KeyPairResponse], error)
	DeleteKeyPair(ctx context.Context, projectId string, keyPairId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

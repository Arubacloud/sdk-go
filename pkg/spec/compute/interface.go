package compute

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// CloudServerAPI defines the interface for CloudServer operations
type CloudServerAPI interface {
	ListCloudServers(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.CloudServerList], error)
	GetCloudServer(ctx context.Context, project string, cloudServerId string, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
	CreateOrUpdateCloudServer(ctx context.Context, project string, body schema.CloudServerRequest, params *schema.RequestParameters) (*schema.Response[schema.CloudServerResponse], error)
	DeleteCloudServer(ctx context.Context, projectId string, cloudServerId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// KeyPairAPI defines the interface for KeyPair operations
type KeyPairAPI interface {
	ListKeyPairs(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KeyPairListResponse], error)
	GetKeyPair(ctx context.Context, project string, keyPairId string, params *schema.RequestParameters) (*schema.Response[schema.KeyPairResponse], error)
	CreateOrUpdateKeyPair(ctx context.Context, project string, body schema.KeyPairRequest, params *schema.RequestParameters) (*schema.Response[schema.KeyPairResponse], error)
	DeleteKeyPair(ctx context.Context, projectId string, keyPairId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

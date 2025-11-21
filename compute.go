package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

type ComputeClient interface {
	CloudServers() CloudServersClient
	KeyPairs() KeyPairsClient
}

type computeClientImpl struct {
	cloudServerClient CloudServersClient
	keyPairClient     KeyPairsClient
}

var _ ComputeClient = (*computeClientImpl)(nil)

func (c *computeClientImpl) CloudServers() CloudServersClient {
	return c.cloudServerClient
}

func (c *computeClientImpl) KeyPairs() KeyPairsClient {
	return c.keyPairClient
}

type CloudServersClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.CloudServerList], error)
	Get(ctx context.Context, projectID string, cloudServerID string, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	Create(ctx context.Context, projectID string, body types.CloudServerRequest, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	Update(ctx context.Context, projectID string, cloudServerID string, body types.CloudServerRequest, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	Delete(ctx context.Context, projectID string, cloudServerID string, params *types.RequestParameters) (*types.Response[any], error)
}

type KeyPairsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.KeyPairListResponse], error)
	Get(ctx context.Context, projectID string, keyPairID string, params *types.RequestParameters) (*types.Response[types.KeyPairResponse], error)
	Create(ctx context.Context, projectID string, body types.KeyPairRequest, params *types.RequestParameters) (*types.Response[types.KeyPairResponse], error)
	Delete(ctx context.Context, projectID string, keyPairID string, params *types.RequestParameters) (*types.Response[any], error)
}

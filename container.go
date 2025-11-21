package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

type ContainerClient interface {
	KaaS() KaaSClient
}

type containerClientImpl struct {
	kaasClient KaaSClient
}

var _ ContainerClient = (*containerClientImpl)(nil)

func (c *containerClientImpl) KaaS() KaaSClient {
	return c.kaasClient
}

type KaaSClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.KaaSList], error)
	Get(ctx context.Context, projectID string, kaasID string, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	Create(ctx context.Context, projectID string, body types.KaaSRequest, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	Update(ctx context.Context, projectID string, kaasID string, body types.KaaSRequest, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	Delete(ctx context.Context, projectID string, kaasID string, params *types.RequestParameters) (*types.Response[any], error)
}

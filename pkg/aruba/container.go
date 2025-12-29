package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

type ContainerClient interface {
	KaaS() KaaSClient
	ContainerRegistry() ContainerRegistryClient
}

type containerClientImpl struct {
	kaasClient              KaaSClient
	containerRegistryClient ContainerRegistryClient
}

// ContainerRegistry implements ContainerClient.
func (c *containerClientImpl) ContainerRegistry() ContainerRegistryClient {
	return c.containerRegistryClient
}

var _ ContainerClient = (*containerClientImpl)(nil)

func (c *containerClientImpl) KaaS() KaaSClient {
	return c.kaasClient
}

type KaaSClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.KaaSList], error)
	Get(ctx context.Context, projectID string, kaasID string, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	Create(ctx context.Context, projectID string, body types.KaaSRequest, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	Update(ctx context.Context, projectID string, kaasID string, body types.KaaSUpdateRequest, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	Delete(ctx context.Context, projectID string, kaasID string, params *types.RequestParameters) (*types.Response[any], error)
}

type ContainerRegistryClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.ContainerRegistryList], error)
	Get(ctx context.Context, projectID string, registryID string, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error)
	Create(ctx context.Context, projectID string, body types.ContainerRegistryRequest, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error)
	Update(ctx context.Context, projectID string, registryID string, body types.ContainerRegistryRequest, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error)
	Delete(ctx context.Context, projectID string, registryID string, params *types.RequestParameters) (*types.Response[any], error)
}

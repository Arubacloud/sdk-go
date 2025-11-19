package container

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

// ContainerAPI defines the unified interface for all Container operations
type ContainerAPI interface {
	// KaaS operations
	ListKaaS(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.KaaSList], error)
	GetKaaS(ctx context.Context, project string, kaasId string, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	CreateKaaS(ctx context.Context, project string, body types.KaaSRequest, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	UpdateKaaS(ctx context.Context, project string, kaasId string, body types.KaaSRequest, params *types.RequestParameters) (*types.Response[types.KaaSResponse], error)
	DeleteKaaS(ctx context.Context, projectId string, kaasId string, params *types.RequestParameters) (*types.Response[any], error)
}

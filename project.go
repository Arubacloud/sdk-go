package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

type ProjectClient interface {
	List(ctx context.Context, params *types.RequestParameters) (*types.Response[types.ProjectList], error)
	Get(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error)
	Create(ctx context.Context, body types.ProjectRequest, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error)
	Update(ctx context.Context, projectID string, body types.ProjectRequest, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error)
	Delete(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[any], error)
}

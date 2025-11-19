package project

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

// ProjectAPI defines the interface for project operations
type ProjectAPI interface {
	ListProjects(ctx context.Context, params *types.RequestParameters) (*types.Response[types.ProjectList], error)
	GetProject(ctx context.Context, projectId string, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error)
	CreateProject(ctx context.Context, body types.ProjectRequest, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error)
	UpdateProject(ctx context.Context, projectId string, body types.ProjectRequest, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error)
	DeleteProject(ctx context.Context, projectId string, params *types.RequestParameters) (*types.Response[any], error)
}

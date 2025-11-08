package project

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ProjectAPI defines the interface for project operations
type ProjectAPI interface {
	ListProjects(ctx context.Context, params *schema.RequestParameters) (*schema.Response[schema.ProjectList], error)
	GetProject(ctx context.Context, projectId string, params *schema.RequestParameters) (*schema.Response[schema.ProjectResponse], error)
	CreateProject(ctx context.Context, body schema.ProjectRequest, params *schema.RequestParameters) (*schema.Response[schema.ProjectResponse], error)
	UpdateProject(ctx context.Context, projectId string, body schema.ProjectRequest, params *schema.RequestParameters) (*schema.Response[schema.ProjectResponse], error)
	DeleteProject(ctx context.Context, projectId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

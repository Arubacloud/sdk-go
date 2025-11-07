package project

import (
	"context"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ProjectAPI defines the interface for project operations
type ProjectAPI interface {
	ListProjects(ctx context.Context, params *schema.RequestParameters) (*http.Response, error)
	GetProject(ctx context.Context, projectId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateProject(ctx context.Context, body schema.ProjectRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteProject(ctx context.Context, projectId string, params *schema.RequestParameters) (*http.Response, error)
}

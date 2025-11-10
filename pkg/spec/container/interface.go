package container

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ContainerAPI defines the unified interface for all Container operations
type ContainerAPI interface {
	// KaaS operations
	ListKaaS(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KaaSList], error)
	GetKaaS(ctx context.Context, project string, kaasId string, params *schema.RequestParameters) (*schema.Response[schema.KaaSResponse], error)
	CreateKaaS(ctx context.Context, project string, body schema.KaaSRequest, params *schema.RequestParameters) (*schema.Response[schema.KaaSResponse], error)
	UpdateKaaS(ctx context.Context, project string, kaasId string, body schema.KaaSRequest, params *schema.RequestParameters) (*schema.Response[schema.KaaSResponse], error)
	DeleteKaaS(ctx context.Context, projectId string, kaasId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

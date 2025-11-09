package container

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// KaaSAPI defines the interface for KaaS operations
type KaaSAPI interface {
	// ListKaaS retrieves all KaaS clusters for a project
	ListKaaS(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KaaSList], error)

	// GetKaaS retrieves a specific KaaS cluster by ID
	GetKaaS(ctx context.Context, project string, kaasId string, params *schema.RequestParameters) (*schema.Response[schema.KaaSResponse], error)

	// CreateKaaS creates a new KaaS cluster
	CreateKaaS(ctx context.Context, project string, body schema.KaaSRequest, params *schema.RequestParameters) (*schema.Response[schema.KaaSResponse], error)

	// UpdateKaaS updates an existing KaaS cluster
	UpdateKaaS(ctx context.Context, project string, kaasId string, body schema.KaaSRequest, params *schema.RequestParameters) (*schema.Response[schema.KaaSResponse], error)

	// DeleteKaaS deletes a KaaS cluster by ID
	DeleteKaaS(ctx context.Context, projectId string, kaasId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

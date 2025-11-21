package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

type SecurityClient interface {
	KMSKeys() KMSKeysClient
}

type securityClientImpl struct {
	kmsKeysClient KMSKeysClient
}

var _ SecurityClient = (*securityClientImpl)(nil)

func (c *securityClientImpl) KMSKeys() KMSKeysClient {
	return c.kmsKeysClient
}

type KMSKeysClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.KmsList], error)
	Get(ctx context.Context, projectID string, kmsKeyID string, params *types.RequestParameters) (*types.Response[types.KmsResponse], error)
	Create(ctx context.Context, projectID string, body types.KmsRequest, params *types.RequestParameters) (*types.Response[types.KmsResponse], error)
	Update(ctx context.Context, projectID string, kmsKeyID string, body types.KmsRequest, params *types.RequestParameters) (*types.Response[types.KmsResponse], error)
	Delete(ctx context.Context, projectID string, kmsKeyID string, params *types.RequestParameters) (*types.Response[any], error)
}

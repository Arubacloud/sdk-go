package security

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

// SecurityAPI defines the unified interface for all Security operations
type SecurityAPI interface {
	// KMS operations
	ListKMSKeys(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.KmsList], error)
	GetKMSKey(ctx context.Context, project string, kmsKeyId string, params *types.RequestParameters) (*types.Response[types.KmsResponse], error)
	CreateKMSKey(ctx context.Context, project string, body types.KmsRequest, params *types.RequestParameters) (*types.Response[types.KmsResponse], error)
	UpdateKMSKey(ctx context.Context, project string, kmsKeyId string, body types.KmsRequest, params *types.RequestParameters) (*types.Response[types.KmsResponse], error)
	DeleteKMSKey(ctx context.Context, projectId string, kmsKeyId string, params *types.RequestParameters) (*types.Response[any], error)
}

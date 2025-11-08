package security

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// KMSAPI defines the interface for managing KMS keys.
type KMSAPI interface {
	ListKMSKeys(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KmsList], error)
	GetKMSKey(ctx context.Context, project string, kmsKeyId string, params *schema.RequestParameters) (*schema.Response[schema.KmsResponse], error)
	CreateKMSKey(ctx context.Context, project string, body schema.KmsRequest, params *schema.RequestParameters) (*schema.Response[schema.KmsResponse], error)
	UpdateKMSKey(ctx context.Context, project string, kmsKeyId string, body schema.KmsRequest, params *schema.RequestParameters) (*schema.Response[schema.KmsResponse], error)
	DeleteKMSKey(ctx context.Context, projectId string, kmsKeyId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

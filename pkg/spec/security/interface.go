package security

import (
	"context"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// KMSAPI defines the interface for managing KMS keys.
type KMSAPI interface {
	ListKMSKeys(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetKMSKey(ctx context.Context, project string, kmsKeyId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateKMSKey(ctx context.Context, project string, body schema.KmsRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteKMSKey(ctx context.Context, projectId string, kmsKeyId string, params *schema.RequestParameters) (*http.Response, error)
}

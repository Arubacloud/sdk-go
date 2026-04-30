package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/internal/clients/security"
)

type SecurityClient interface {
	KMS() KMSClient
	Keys() KeysClient
	Kmips() KmipsClient
}

type securityClientImpl struct {
	kmsClient   KMSClient
	keysClient  KeysClient
	kmipsClient KmipsClient
}

var _ SecurityClient = (*securityClientImpl)(nil)

func (c *securityClientImpl) KMS() KMSClient     { return c.kmsClient }
func (c *securityClientImpl) Keys() KeysClient   { return c.keysClient }
func (c *securityClientImpl) Kmips() KmipsClient { return c.kmipsClient }

type KMSClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*KMS], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*KMS, error)
	Create(ctx context.Context, k *KMS, opts ...CallOption) (*KMS, error)
	Update(ctx context.Context, k *KMS, opts ...CallOption) (*KMS, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

// Sub-client aliases — pluralized to match repo convention.
type (
	KeysClient  = *security.KeyClientImpl
	KmipsClient = *security.KmipClientImpl
)

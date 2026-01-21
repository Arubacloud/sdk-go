package aruba

import (
	"github.com/Arubacloud/sdk-go/internal/clients/security"
)

type SecurityClient interface {
	KMS() KMSClient
}

type securityClientImpl struct {
	kmsClient KMSClient
}

var _ SecurityClient = (*securityClientImpl)(nil)

func (c *securityClientImpl) KMS() KMSClient {
	return c.kmsClient
}

// Type aliases to internal implementations
type (
	KMSClient  = *security.KMSClientWrapper
	KeyClient  = *security.KeyClientImpl
	KmipClient = *security.KmipClientImpl
)

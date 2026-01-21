package security

import (
	"github.com/Arubacloud/sdk-go/internal/restclient"
)

// KMSClientWrapper wraps the KMS client and provides access to nested resources
type KMSClientWrapper struct {
	*kmsClientImpl
	keyClient  *KeyClientImpl
	kmipClient *KmipClientImpl
}

// NewKMSClientWrapper creates a new KMS client wrapper with nested resources
func NewKMSClientWrapper(client *restclient.Client) *KMSClientWrapper {
	return &KMSClientWrapper{
		kmsClientImpl: NewKMSClientImpl(client),
		keyClient:     NewKeyClientImpl(client),
		kmipClient:    NewKmipClientImpl(client),
	}
}

// Keys returns the Key client for managing KMS keys
func (w *KMSClientWrapper) Keys() *KeyClientImpl {
	return w.keyClient
}

// Kmips returns the KMIP client for managing KMIP services
func (w *KMSClientWrapper) Kmips() *KmipClientImpl {
	return w.kmipClient
}

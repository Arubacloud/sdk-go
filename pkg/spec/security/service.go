package security

import (
	"github.com/Arubacloud/sdk-go/pkg/restclient"
)

// Service implements the SecurityAPI interface for all Security operations
type Service struct {
	client *restclient.Client
}

// NewService creates a new unified Security service
func NewService(client *restclient.Client) *Service {
	return &Service{
		client: client,
	}
}

package security

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the SecurityAPI interface for all Security operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Security service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

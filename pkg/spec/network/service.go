package network

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the NetworkAPI interface for all Network operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Network service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

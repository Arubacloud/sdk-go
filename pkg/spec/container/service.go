package container

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the ContainerAPI interface for all Container operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Container service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

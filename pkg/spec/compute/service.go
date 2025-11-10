package compute

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the ComputeAPI interface for all Compute operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Compute service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

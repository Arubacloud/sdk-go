package compute

import (
	"github.com/Arubacloud/sdk-go/pkg/restclient"
)

// Service implements the ComputeAPI interface for all Compute operations
type Service struct {
	client *restclient.Client
}

// NewService creates a new unified Compute service
func NewService(client *restclient.Client) *Service {
	return &Service{
		client: client,
	}
}

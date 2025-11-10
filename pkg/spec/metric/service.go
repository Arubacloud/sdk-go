package metric

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the MetricAPI interface for all Metric operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Metric service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

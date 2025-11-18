package metric

import (
	"github.com/Arubacloud/sdk-go/pkg/restclient"
)

// Service implements the MetricAPI interface for all Metric operations
type Service struct {
	client *restclient.Client
}

// NewService creates a new unified Metric service
func NewService(client *restclient.Client) *Service {
	return &Service{
		client: client,
	}
}

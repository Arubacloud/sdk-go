package schedule

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the ScheduleAPI interface for all Schedule operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Schedule service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

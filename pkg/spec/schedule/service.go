package schedule

import (
	"github.com/Arubacloud/sdk-go/pkg/restclient"
)

// Service implements the ScheduleAPI interface for all Schedule operations
type Service struct {
	client *restclient.Client
}

// NewService creates a new unified Schedule service
func NewService(client *restclient.Client) *Service {
	return &Service{
		client: client,
	}
}

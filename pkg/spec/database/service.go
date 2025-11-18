package database

import (
	"github.com/Arubacloud/sdk-go/pkg/restclient"
)

// Service implements the DatabaseAPI interface for all Database operations
type Service struct {
	client *restclient.Client
}

// NewService creates a new unified Database service
func NewService(client *restclient.Client) *Service {
	return &Service{
		client: client,
	}
}

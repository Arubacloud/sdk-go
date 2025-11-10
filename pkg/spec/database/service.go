package database

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the DatabaseAPI interface for all Database operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Database service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

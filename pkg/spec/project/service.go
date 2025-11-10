package project

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the ProjectAPI interface for all Project operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Project service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

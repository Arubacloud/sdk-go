package storage

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the StorageAPI interface for all Storage operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Storage service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

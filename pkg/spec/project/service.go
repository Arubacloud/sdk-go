package project

import (
	"github.com/Arubacloud/sdk-go/pkg/restclient"
)

// Service implements the ProjectAPI interface for all Project operations
type Service struct {
	client *restclient.Client
}

// NewService creates a new unified Project service
func NewService(client *restclient.Client) *Service {
	return &Service{
		client: client,
	}
}

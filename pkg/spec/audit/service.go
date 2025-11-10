package audit

import (
	"github.com/Arubacloud/sdk-go/pkg/client"
)

// Service implements the AuditAPI interface for all Audit operations
type Service struct {
	client *client.Client
}

// NewService creates a new unified Audit service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

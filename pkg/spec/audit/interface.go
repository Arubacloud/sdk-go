package audit

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// AuditAPI defines the unified interface for all Audit operations
type AuditAPI interface {
	// Event operations
	ListEvents(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.AuditEventListResponse], error)
}

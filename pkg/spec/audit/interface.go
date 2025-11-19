package audit

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

// AuditAPI defines the unified interface for all Audit operations
type AuditAPI interface {
	// Event operations
	ListEvents(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.AuditEventListResponse], error)
}

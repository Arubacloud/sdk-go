package audit

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

type EventAPI interface {
	ListEvents(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.AuditEventListResponse], error)
}

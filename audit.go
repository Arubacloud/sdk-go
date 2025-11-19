package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

type EventsClient interface {
	List(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.AuditEventListResponse], error)
}

type AuditClient interface {
	Events() EventsClient
}

type auditClientImpl struct {
	eventsClient EventsClient
}

var _ AuditClient = (*auditClientImpl)(nil)

func (c auditClientImpl) Events() EventsClient {
	return c.eventsClient
}

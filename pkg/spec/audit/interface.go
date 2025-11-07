package audit

import (
	"context"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

type EventAPI interface {
	ListEvents(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
}

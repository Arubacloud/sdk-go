package metric

import (
	"context"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// MetricAPI defines the methods for interacting with metrics.
type MetricAPI interface {
	ListMetrics(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
}

// AlertAPI defines the methods for interacting with alerts.
type AlertAPI interface {
	ListAlerts(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
}

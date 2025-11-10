package metric

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// MetricAPI defines the unified interface for all Metric operations
type MetricAPI interface {
	// Metric operations
	ListMetrics(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.MetricListResponse], error)

	// Alert operations
	ListAlerts(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.AlertsListResponse], error)
}

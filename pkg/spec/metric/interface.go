package metric

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

// MetricAPI defines the unified interface for all Metric operations
type MetricAPI interface {
	// Metric operations
	ListMetrics(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.MetricListResponse], error)

	// Alert operations
	ListAlerts(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.AlertsListResponse], error)
}

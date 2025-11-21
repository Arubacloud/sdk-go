package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

type MetricClient interface {
	Alerts() AlertsClient
	Metrics() MetricsClient
}

type metricClientImpl struct {
	alertsClient  AlertsClient
	metricsClient MetricsClient
}

var _ MetricClient = (*metricClientImpl)(nil)

func (c metricClientImpl) Alerts() AlertsClient {
	return c.alertsClient
}

func (c metricClientImpl) Metrics() MetricsClient {
	return c.metricsClient
}

type AlertsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.AlertsListResponse], error)
}

type MetricsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.MetricListResponse], error)
}

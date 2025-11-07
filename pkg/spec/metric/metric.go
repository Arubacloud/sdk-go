package metric

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// MetricService implements the MetricAPI interface
type MetricService struct {
	client *client.Client
}

// NewMetricService creates a new MetricService
func NewMetricService(client *client.Client) *MetricService {
	return &MetricService{
		client: client,
	}
}

// ListMetrics retrieves all metrics for a project
func (s *MetricService) ListMetrics(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(MetricsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

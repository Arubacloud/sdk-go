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
func (s *MetricService) ListMetrics(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.MetricListResponse], error) {
	s.client.Logger().Debugf("Listing metrics for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(MetricsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.MetricListResponse](httpResp)
}

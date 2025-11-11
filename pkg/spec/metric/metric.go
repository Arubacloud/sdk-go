package metric

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ListMetrics retrieves all metrics for a project
func (s *Service) ListMetrics(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.MetricListResponse], error) {
	s.client.Logger().Debugf("Listing metrics for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(MetricsPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &MetricListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &MetricListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.MetricListResponse](httpResp)
}

package metric

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ListAlerts retrieves all alerts for a project
func (s *Service) ListAlerts(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.AlertsListResponse], error) {
	s.client.Logger().Debugf("Listing alerts for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(AlertsPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &AlertListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &AlertListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.AlertsListResponse](httpResp)
}

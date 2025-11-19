package metric

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/types"
)

// ListAlerts retrieves all alerts for a project
func (s *Service) ListAlerts(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.AlertsListResponse], error) {
	s.client.Logger().Debugf("Listing alerts for project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(AlertsPath, project)

	if params == nil {
		params = &types.RequestParameters{
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

	return types.ParseResponseBody[types.AlertsListResponse](httpResp)
}

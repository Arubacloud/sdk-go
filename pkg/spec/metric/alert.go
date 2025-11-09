package metric

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// AlertService implements the AlertAPI interface
type AlertService struct {
	client *client.Client
}

// NewAlertService creates a new AlertService
func NewAlertService(client *client.Client) *AlertService {
	return &AlertService{
		client: client,
	}
}

// ListAlerts retrieves all alerts for a project
func (s *AlertService) ListAlerts(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.AlertsListResponse], error) {
	s.client.Logger().Debugf("Listing alerts for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(AlertsPath, project)

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

	return schema.ParseResponseBody[schema.AlertsListResponse](httpResp)
}

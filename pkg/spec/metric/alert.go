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
func (s *AlertService) ListAlerts(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(AlertsPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

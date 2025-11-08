package metric

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	if err := validateProject(project); err != nil {
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

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &schema.Response[schema.AlertsListResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.AlertsListResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

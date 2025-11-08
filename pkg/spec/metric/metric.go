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
	response := &schema.Response[schema.MetricListResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.MetricListResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

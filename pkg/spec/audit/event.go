package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// EventService implements the EventAPI interface
type EventService struct {
	client *client.Client
}

// NewEventService creates a new EventService
func NewEventService(client *client.Client) *EventService {
	return &EventService{
		client: client,
	}
}

// ListEvents retrieves all audit events for a project
func (s *EventService) ListEvents(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.AuditEventListResponse], error) {
	s.client.Logger().Debugf("Listing audit events for project: %s", project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(EventsPath, project)

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
	response := &schema.Response[schema.AuditEventListResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.AuditEventListResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

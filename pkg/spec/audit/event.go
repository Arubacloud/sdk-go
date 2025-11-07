package audit

import (
	"context"
	"fmt"
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
func (s *EventService) ListEvents(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
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

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

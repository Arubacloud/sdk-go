package audit

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ListEvents retrieves all audit events for a project
func (s *Service) ListEvents(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.AuditEventListResponse], error) {
	s.client.Logger().Debugf("Listing audit events for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
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

	return schema.ParseResponseBody[schema.AuditEventListResponse](httpResp)
}

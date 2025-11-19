package audit

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/types"
)

// ListEvents retrieves all audit events for a project
func (s *Service) ListEvents(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.AuditEventListResponse], error) {
	s.client.Logger().Debugf("Listing audit events for project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(EventsPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &AuditLogListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &AuditLogListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.AuditEventListResponse](httpResp)
}

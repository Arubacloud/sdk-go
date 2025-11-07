package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// GrantService implements the GrantAPI interface
type GrantService struct {
    client *client.Client
}

// NewGrantService creates a new GrantService
func NewGrantService(client *client.Client) *GrantService {
    return &GrantService{
        client: client,
    }
}

// ListGrants retrieves all grants for a database
func (s *GrantService) ListGrants(ctx context.Context, project string, dbaasId string, databaseId string, params *schema.RequestParameters) (*http.Response, error) {
    if project == "" {
        return nil, fmt.Errorf("project cannot be empty")
    }
    if dbaasId == "" {
        return nil, fmt.Errorf("DBaaS ID cannot be empty")
    }
    if databaseId == "" {
        return nil, fmt.Errorf("database ID cannot be empty")
    }

    path := fmt.Sprintf(GrantsPath, project, dbaasId, databaseId)

    var queryParams map[string]string
    var headers map[string]string

    if params != nil {
        queryParams = params.ToQueryParams()
        headers = params.ToHeaders()
    }

    return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetGrant retrieves a specific grant by ID
func (s *GrantService) GetGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, params *schema.RequestParameters) (*http.Response, error) {
    if project == "" {
        return nil, fmt.Errorf("project cannot be empty")
    }
    if dbaasId == "" {
        return nil, fmt.Errorf("DBaaS ID cannot be empty")
    }
    if databaseId == "" {
        return nil, fmt.Errorf("database ID cannot be empty")
    }
    if grantId == "" {
        return nil, fmt.Errorf("grant ID cannot be empty")
    }

    path := fmt.Sprintf(GrantPath, project, dbaasId, databaseId, grantId)

    var queryParams map[string]string
    var headers map[string]string

    if params != nil {
        queryParams = params.ToQueryParams()
        headers = params.ToHeaders()
    }

    return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateGrant creates or updates a grant
func (s *GrantService) CreateOrUpdateGrant(ctx context.Context, project string, dbaasId string, databaseId string, body schema.GrantRequest, params *schema.RequestParameters) (*http.Response, error) {
    if project == "" {
        return nil, fmt.Errorf("project cannot be empty")
    }
    if dbaasId == "" {
        return nil, fmt.Errorf("DBaaS ID cannot be empty")
    }
    if databaseId == "" {
        return nil, fmt.Errorf("database ID cannot be empty")
    }

    path := fmt.Sprintf(GrantsPath, project, dbaasId, databaseId)

    var queryParams map[string]string
    var headers map[string]string

    if params != nil {
        queryParams = params.ToQueryParams()
        headers = params.ToHeaders()
    }

    return s.client.DoRequest(ctx, http.MethodPut, path, body, queryParams, headers)
}

// DeleteGrant deletes a grant by ID
func (s *GrantService) DeleteGrant(ctx context.Context, projectId string, dbaasId string, databaseId string, grantId string, params *schema.RequestParameters) (*http.Response, error) {
    if projectId == "" {
        return nil, fmt.Errorf("project ID cannot be empty")
    }
    if dbaasId == "" {
        return nil, fmt.Errorf("DBaaS ID cannot be empty")
    }
    if databaseId == "" {
        return nil, fmt.Errorf("database ID cannot be empty")
    }
    if grantId == "" {
        return nil, fmt.Errorf("grant ID cannot be empty")
    }

    path := fmt.Sprintf(GrantPath, projectId, dbaasId, databaseId, grantId)

    var queryParams map[string]string
    var headers map[string]string

    if params != nil {
        queryParams = params.ToQueryParams()
        headers = params.ToHeaders()
    }

    return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}
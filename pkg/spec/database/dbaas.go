package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// DBaaSService implements the DBaaSAPI interface
type DBaaSService struct {
	client *client.Client
}

// NewDBaaSService creates a new DBaaSService
func NewDBaaSService(client *client.Client) *DBaaSService {
	return &DBaaSService{
		client: client,
	}
}

// ListDBaaS retrieves all DBaaS instances for a project
func (s *DBaaSService) ListDBaaS(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(DBaaSPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetDBaaS retrieves a specific DBaaS instance by ID
func (s *DBaaSService) GetDBaaS(ctx context.Context, project string, dbaasId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}

	path := fmt.Sprintf(DBaaSItemPath, project, dbaasId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateDBaaS creates or updates a DBaaS instance
func (s *DBaaSService) CreateOrUpdateDBaaS(ctx context.Context, project string, body schema.DBaaSRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(DBaaSPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteDBaaS deletes a DBaaS instance by ID
func (s *DBaaSService) DeleteDBaaS(ctx context.Context, projectId string, dbaasId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}

	path := fmt.Sprintf(DBaaSItemPath, projectId, dbaasId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

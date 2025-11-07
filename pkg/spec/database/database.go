package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// DatabaseService implements the DatabaseAPI interface
type DatabaseService struct {
	client *client.Client
}

// NewDatabaseService creates a new DatabaseService
func NewDatabaseService(client *client.Client) *DatabaseService {
	return &DatabaseService{
		client: client,
	}
}

// ListDatabases retrieves all databases for a DBaaS instance
func (s *DatabaseService) ListDatabases(ctx context.Context, project string, dbaasId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}

	path := fmt.Sprintf(DatabaseInstancesPath, project, dbaasId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetDatabase retrieves a specific database by ID
func (s *DatabaseService) GetDatabase(ctx context.Context, project string, dbaasId string, databaseId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseId == "" {
		return nil, fmt.Errorf("database ID cannot be empty")
	}

	path := fmt.Sprintf(DatabaseInstancePath, project, dbaasId, databaseId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateDatabase creates or updates a database
func (s *DatabaseService) CreateOrUpdateDatabase(ctx context.Context, project string, dbaasId string, body schema.DatabaseRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}

	path := fmt.Sprintf(DatabaseInstancesPath, project, dbaasId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, body, queryParams, headers)
}

// DeleteDatabase deletes a database by ID
func (s *DatabaseService) DeleteDatabase(ctx context.Context, projectId string, dbaasId string, databaseId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseId == "" {
		return nil, fmt.Errorf("database ID cannot be empty")
	}

	path := fmt.Sprintf(DatabaseInstancePath, projectId, dbaasId, databaseId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// UserService implements the UserAPI interface
type UserService struct {
	client *client.Client
}

// NewUserService creates a new UserService
func NewUserService(client *client.Client) *UserService {
	return &UserService{
		client: client,
	}
}

// ListUsers retrieves all users for a project
func (s *UserService) ListUsers(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(UsersPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetUser retrieves a specific user by ID
func (s *UserService) GetUser(ctx context.Context, project string, userId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if userId == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	path := fmt.Sprintf(UserItemPath, project, userId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateUser creates or updates a user
func (s *UserService) CreateOrUpdateUser(ctx context.Context, project string, body schema.UserRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	path := fmt.Sprintf(UsersPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, body, queryParams, headers)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, projectId string, userId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if userId == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	path := fmt.Sprintf(UserItemPath, projectId, userId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

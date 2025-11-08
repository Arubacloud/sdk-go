package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// ListUsers retrieves all users for a DBaaS instance
func (s *UserService) ListUsers(ctx context.Context, project string, dbaasId string, params *schema.RequestParameters) (*schema.Response[schema.UserList], error) {
	s.client.Logger().Debugf("Listing users for DBaaS: %s in project: %s", dbaasId, project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}

	path := fmt.Sprintf(UsersPath, project, dbaasId)

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
	response := &schema.Response[schema.UserList]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.UserList
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// GetUser retrieves a specific user by ID
func (s *UserService) GetUser(ctx context.Context, project string, dbaasId string, userId string, params *schema.RequestParameters) (*schema.Response[schema.UserResponse], error) {
	s.client.Logger().Debugf("Getting user: %s from DBaaS: %s in project: %s", userId, dbaasId, project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if userId == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	path := fmt.Sprintf(UserItemPath, project, dbaasId, userId)

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
	response := &schema.Response[schema.UserResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.UserResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// CreateUser creates a new user in a DBaaS instance
func (s *UserService) CreateUser(ctx context.Context, project string, dbaasId string, body schema.UserRequest, params *schema.RequestParameters) (*schema.Response[schema.UserResponse], error) {
	s.client.Logger().Debugf("Creating user in DBaaS: %s in project: %s", dbaasId, project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}

	path := fmt.Sprintf(UsersPath, project, dbaasId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	// Marshal the request body to JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// Read the response body
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &schema.Response[schema.UserResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.UserResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, project string, dbaasId string, userId string, body schema.UserRequest, params *schema.RequestParameters) (*schema.Response[schema.UserResponse], error) {
	s.client.Logger().Debugf("Updating user: %s in DBaaS: %s in project: %s", userId, dbaasId, project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if userId == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	path := fmt.Sprintf(UserItemPath, project, dbaasId, userId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	// Marshal the request body to JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPut, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// Read the response body
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &schema.Response[schema.UserResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data schema.UserResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, projectId string, dbaasId string, userId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting user: %s from DBaaS: %s in project: %s", userId, dbaasId, projectId)

	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if dbaasId == "" {
		return nil, fmt.Errorf("DBaaS ID cannot be empty")
	}
	if userId == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	path := fmt.Sprintf(UserItemPath, projectId, dbaasId, userId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
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
	response := &schema.Response[any]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// For DELETE operations, we typically don't parse the body unless there's content
	if response.IsSuccess() && len(bodyBytes) > 0 {
		var data any
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

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

	if err := schema.ValidateProjectAndResource(project, dbaasId, "DBaaS ID"); err != nil {
		return nil, err
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

	return schema.ParseResponseBody[schema.UserList](httpResp)
}

// GetUser retrieves a specific user by ID
func (s *UserService) GetUser(ctx context.Context, project string, dbaasId string, userId string, params *schema.RequestParameters) (*schema.Response[schema.UserResponse], error) {
	s.client.Logger().Debugf("Getting user: %s from DBaaS: %s in project: %s", userId, dbaasId, project)

	if err := schema.ValidateDBaaSResource(project, dbaasId, userId, "user ID"); err != nil {
		return nil, err
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

	return schema.ParseResponseBody[schema.UserResponse](httpResp)
}

// CreateUser creates a new user in a DBaaS instance
func (s *UserService) CreateUser(ctx context.Context, project string, dbaasId string, body schema.UserRequest, params *schema.RequestParameters) (*schema.Response[schema.UserResponse], error) {
	s.client.Logger().Debugf("Creating user in DBaaS: %s in project: %s", dbaasId, project)

	if err := schema.ValidateProjectAndResource(project, dbaasId, "DBaaS ID"); err != nil {
		return nil, err
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
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp schema.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, project string, dbaasId string, userId string, body schema.UserRequest, params *schema.RequestParameters) (*schema.Response[schema.UserResponse], error) {
	s.client.Logger().Debugf("Updating user: %s in DBaaS: %s in project: %s", userId, dbaasId, project)

	if err := schema.ValidateDBaaSResource(project, dbaasId, userId, "user ID"); err != nil {
		return nil, err
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
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp schema.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, projectId string, dbaasId string, userId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting user: %s from DBaaS: %s in project: %s", userId, dbaasId, projectId)

	if err := schema.ValidateDBaaSResource(projectId, dbaasId, userId, "user ID"); err != nil {
		return nil, err
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

	return schema.ParseResponseBody[any](httpResp)
}

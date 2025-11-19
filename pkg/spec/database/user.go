package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/types"
)

// ListUsers retrieves all users for a DBaaS instance
func (s *Service) ListUsers(ctx context.Context, project string, dbaasId string, params *types.RequestParameters) (*types.Response[types.UserList], error) {
	s.client.Logger().Debugf("Listing users for DBaaS: %s in project: %s", dbaasId, project)

	if err := types.ValidateProjectAndResource(project, dbaasId, "DBaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(UsersPath, project, dbaasId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseUserListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseUserListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.UserList](httpResp)
}

// GetUser retrieves a specific user by ID
func (s *Service) GetUser(ctx context.Context, project string, dbaasId string, userId string, params *types.RequestParameters) (*types.Response[types.UserResponse], error) {
	s.client.Logger().Debugf("Getting user: %s from DBaaS: %s in project: %s", userId, dbaasId, project)

	if err := types.ValidateDBaaSResource(project, dbaasId, userId, "user ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(UserItemPath, project, dbaasId, userId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseUserGetVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseUserGetVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.UserResponse](httpResp)
}

// CreateUser creates a new user in a DBaaS instance
func (s *Service) CreateUser(ctx context.Context, project string, dbaasId string, body types.UserRequest, params *types.RequestParameters) (*types.Response[types.UserResponse], error) {
	s.client.Logger().Debugf("Creating user in DBaaS: %s in project: %s", dbaasId, project)

	if err := types.ValidateProjectAndResource(project, dbaasId, "DBaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(UsersPath, project, dbaasId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseUserCreateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseUserCreateVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

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
	response := &types.Response[types.UserResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.UserResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp types.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(ctx context.Context, project string, dbaasId string, userId string, body types.UserRequest, params *types.RequestParameters) (*types.Response[types.UserResponse], error) {
	s.client.Logger().Debugf("Updating user: %s in DBaaS: %s in project: %s", userId, dbaasId, project)

	if err := types.ValidateDBaaSResource(project, dbaasId, userId, "user ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(UserItemPath, project, dbaasId, userId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseUserUpdateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseUserUpdateVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

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
	response := &types.Response[types.UserResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.UserResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp types.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// DeleteUser deletes a user by ID
func (s *Service) DeleteUser(ctx context.Context, projectId string, dbaasId string, userId string, params *types.RequestParameters) (*types.Response[any], error) {
	s.client.Logger().Debugf("Deleting user: %s from DBaaS: %s in project: %s", userId, dbaasId, projectId)

	if err := types.ValidateDBaaSResource(projectId, dbaasId, userId, "user ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(UserItemPath, projectId, dbaasId, userId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseUserDeleteVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseUserDeleteVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[any](httpResp)
}

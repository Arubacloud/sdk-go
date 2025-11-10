package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ListSubnets retrieves all subnets for a VPC
func (s *Service) ListSubnets(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.SubnetList], error) {
	s.client.Logger().Debugf("Listing subnets for VPC: %s in project: %s", vpcId, project)

	if err := schema.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SubnetsPath, project, vpcId)

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

	return schema.ParseResponseBody[schema.SubnetList](httpResp)
}

// GetSubnet retrieves a specific subnet by ID
func (s *Service) GetSubnet(ctx context.Context, project string, vpcId string, subnetId string, params *schema.RequestParameters) (*schema.Response[schema.SubnetResponse], error) {
	s.client.Logger().Debugf("Getting subnet: %s from VPC: %s in project: %s", subnetId, vpcId, project)

	if err := schema.ValidateVPCResource(project, vpcId, subnetId, "subnet ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SubnetPath, project, vpcId, subnetId)

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

	return schema.ParseResponseBody[schema.SubnetResponse](httpResp)
}

// CreateSubnet creates a new subnet in a VPC
// The SDK automatically waits for the VPC to become Active before creating the subnet
func (s *Service) CreateSubnet(ctx context.Context, project string, vpcId string, body schema.SubnetRequest, params *schema.RequestParameters) (*schema.Response[schema.SubnetResponse], error) {
	s.client.Logger().Debugf("Creating subnet in VPC: %s in project: %s", vpcId, project)

	if err := schema.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	// Wait for VPC to become Active before creating subnet
	err := s.waitForVPCActive(ctx, project, vpcId)
	if err != nil {
		return nil, fmt.Errorf("failed waiting for VPC to become active: %w", err)
	}

	path := fmt.Sprintf(SubnetsPath, project, vpcId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.SubnetResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.SubnetResponse
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

// UpdateSubnet updates an existing subnet
func (s *Service) UpdateSubnet(ctx context.Context, project string, vpcId string, subnetId string, body schema.SubnetRequest, params *schema.RequestParameters) (*schema.Response[schema.SubnetResponse], error) {
	s.client.Logger().Debugf("Updating subnet: %s in VPC: %s in project: %s", subnetId, vpcId, project)

	if err := schema.ValidateVPCResource(project, vpcId, subnetId, "subnet ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SubnetPath, project, vpcId, subnetId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodPut, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.SubnetResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.SubnetResponse
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

// DeleteSubnet deletes a subnet by ID
func (s *Service) DeleteSubnet(ctx context.Context, projectId string, vpcId string, subnetId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting subnet: %s from VPC: %s in project: %s", subnetId, vpcId, projectId)

	if err := schema.ValidateVPCResource(projectId, vpcId, subnetId, "subnet ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(SubnetPath, projectId, vpcId, subnetId)

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

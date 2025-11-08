package network

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

// VpcPeeringRouteService implements the VpcPeeringRouteAPI interface
type VpcPeeringRouteService struct {
	client *client.Client
}

// NewVpcPeeringRouteService creates a new VpcPeeringRouteService
func NewVpcPeeringRouteService(client *client.Client) *VpcPeeringRouteService {
	return &VpcPeeringRouteService{
		client: client,
	}
}

// ListVpcPeeringRoutes retrieves all VPC peering routes for a VPC peering connection
func (s *VpcPeeringRouteService) ListVpcPeeringRoutes(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringRouteList], error) {
	s.client.Logger().Debugf("Listing VPC peering routes for VPC peering: %s in VPC: %s in project: %s", vpcPeeringId, vpcId, project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringRoutesPath, project, vpcId, vpcPeeringId)

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

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.VpcPeeringRouteList]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.VpcPeeringRouteList
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// GetVpcPeeringRoute retrieves a specific VPC peering route by ID
func (s *VpcPeeringRouteService) GetVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringRouteResponse], error) {
	s.client.Logger().Debugf("Getting VPC peering route: %s from VPC peering: %s in VPC: %s in project: %s", vpcPeeringRouteId, vpcPeeringId, vpcId, project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}
	if vpcPeeringRouteId == "" {
		return nil, fmt.Errorf("VPC peering route ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringRoutePath, project, vpcId, vpcPeeringId, vpcPeeringRouteId)

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

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.VpcPeeringRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() {
		var data schema.VpcPeeringRouteResponse
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// CreateVpcPeeringRoute creates a new VPC peering route
func (s *VpcPeeringRouteService) CreateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, body schema.VpcPeeringRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringRouteResponse], error) {
	s.client.Logger().Debugf("Creating VPC peering route in VPC peering: %s in VPC: %s in project: %s", vpcPeeringId, vpcId, project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringRoutesPath, project, vpcId, vpcPeeringId)

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

	response := &schema.Response[schema.VpcPeeringRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VpcPeeringRouteResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// UpdateVpcPeeringRoute updates an existing VPC peering route
func (s *VpcPeeringRouteService) UpdateVpcPeeringRoute(ctx context.Context, project string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, body schema.VpcPeeringRouteRequest, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringRouteResponse], error) {
	s.client.Logger().Debugf("Updating VPC peering route: %s in VPC peering: %s in VPC: %s in project: %s", vpcPeeringRouteId, vpcPeeringId, vpcId, project)

	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}
	if vpcPeeringRouteId == "" {
		return nil, fmt.Errorf("VPC peering route ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringRoutePath, project, vpcId, vpcPeeringId, vpcPeeringRouteId)

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

	response := &schema.Response[schema.VpcPeeringRouteResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VpcPeeringRouteResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

// DeleteVpcPeeringRoute deletes a VPC peering route by ID
func (s *VpcPeeringRouteService) DeleteVpcPeeringRoute(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, vpcPeeringRouteId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting VPC peering route: %s from VPC peering: %s in VPC: %s in project: %s", vpcPeeringRouteId, vpcPeeringId, vpcId, projectId)

	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if vpcId == "" {
		return nil, fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringId == "" {
		return nil, fmt.Errorf("VPC peering ID cannot be empty")
	}
	if vpcPeeringRouteId == "" {
		return nil, fmt.Errorf("VPC peering route ID cannot be empty")
	}

	path := fmt.Sprintf(VpcPeeringRoutePath, projectId, vpcId, vpcPeeringId, vpcPeeringRouteId)

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

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[any]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	if response.IsSuccess() && len(bodyBytes) > 0 {
		var data any
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	}

	return response, nil
}

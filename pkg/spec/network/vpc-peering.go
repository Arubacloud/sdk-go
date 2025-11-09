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

// VpcPeeringService implements the VpcPeeringAPI interface
type VpcPeeringService struct {
	client *client.Client
}

// NewVpcPeeringService creates a new VpcPeeringService
func NewVpcPeeringService(client *client.Client) *VpcPeeringService {
	return &VpcPeeringService{
		client: client,
	}
}

// ListVpcPeerings retrieves all VPC peerings for a VPC
func (s *VpcPeeringService) ListVpcPeerings(ctx context.Context, project string, vpcId string, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringList], error) {
	s.client.Logger().Debugf("Listing VPC peerings for VPC: %s in project: %s", vpcId, project)

	if err := schema.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpcPeeringsPath, project, vpcId)

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

	return schema.ParseResponseBody[schema.VpcPeeringList](httpResp)
}

// GetVpcPeering retrieves a specific VPC peering by ID
func (s *VpcPeeringService) GetVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringResponse], error) {
	s.client.Logger().Debugf("Getting VPC peering: %s from VPC: %s in project: %s", vpcPeeringId, vpcId, project)

	if err := schema.ValidateVPCResource(project, vpcId, vpcPeeringId, "VPC peering ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpcPeeringPath, project, vpcId, vpcPeeringId)

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

	return schema.ParseResponseBody[schema.VpcPeeringResponse](httpResp)
}

// CreateVpcPeering creates a new VPC peering
func (s *VpcPeeringService) CreateVpcPeering(ctx context.Context, project string, vpcId string, body schema.VpcPeeringRequest, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringResponse], error) {
	s.client.Logger().Debugf("Creating VPC peering in VPC: %s in project: %s", vpcId, project)

	if err := schema.ValidateProjectAndResource(project, vpcId, "VPC ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpcPeeringsPath, project, vpcId)

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

	response := &schema.Response[schema.VpcPeeringResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VpcPeeringResponse
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

// UpdateVpcPeering updates an existing VPC peering
func (s *VpcPeeringService) UpdateVpcPeering(ctx context.Context, project string, vpcId string, vpcPeeringId string, body schema.VpcPeeringRequest, params *schema.RequestParameters) (*schema.Response[schema.VpcPeeringResponse], error) {
	s.client.Logger().Debugf("Updating VPC peering: %s in VPC: %s in project: %s", vpcPeeringId, vpcId, project)

	if err := schema.ValidateVPCResource(project, vpcId, vpcPeeringId, "VPC peering ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpcPeeringPath, project, vpcId, vpcPeeringId)

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

	response := &schema.Response[schema.VpcPeeringResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.VpcPeeringResponse
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

// DeleteVpcPeering deletes a VPC peering by ID
func (s *VpcPeeringService) DeleteVpcPeering(ctx context.Context, projectId string, vpcId string, vpcPeeringId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting VPC peering: %s from VPC: %s in project: %s", vpcPeeringId, vpcId, projectId)

	if err := schema.ValidateVPCResource(projectId, vpcId, vpcPeeringId, "VPC peering ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(VpcPeeringPath, projectId, vpcId, vpcPeeringId)

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

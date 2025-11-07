package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// VpnRouteService implements the VpnRouteAPI interface
type VpnRouteService struct {
	client *client.Client
}

// NewVpnRouteService creates a new VpnRouteService
func NewVpnRouteService(client *client.Client) *VpnRouteService {
	return &VpnRouteService{
		client: client,
	}
}

// ListVpnRoutes retrieves all VPN routes for a VPN tunnel
func (s *VpnRouteService) ListVpnRoutes(ctx context.Context, project string, vpnTunnelId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpnTunnelId == "" {
		return nil, fmt.Errorf("VPN tunnel ID cannot be empty")
	}

	path := fmt.Sprintf(VpnRoutesPath, project, vpnTunnelId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// GetVpnRoute retrieves a specific VPN route by ID
func (s *VpnRouteService) GetVpnRoute(ctx context.Context, project string, vpnTunnelId string, vpnRouteId string, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpnTunnelId == "" {
		return nil, fmt.Errorf("VPN tunnel ID cannot be empty")
	}
	if vpnRouteId == "" {
		return nil, fmt.Errorf("VPN route ID cannot be empty")
	}

	path := fmt.Sprintf(VpnRoutePath, project, vpnTunnelId, vpnRouteId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
}

// CreateOrUpdateVpnRoute creates or updates a VPN route
func (s *VpnRouteService) CreateOrUpdateVpnRoute(ctx context.Context, project string, vpnTunnelId string, body schema.VpnRouteRequest, params *schema.RequestParameters) (*http.Response, error) {
	if project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}
	if vpnTunnelId == "" {
		return nil, fmt.Errorf("VPN tunnel ID cannot be empty")
	}

	path := fmt.Sprintf(VpnRoutesPath, project, vpnTunnelId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodPut, path, nil, queryParams, headers)
}

// DeleteVpnRoute deletes a VPN route by ID
func (s *VpnRouteService) DeleteVpnRoute(ctx context.Context, projectId string, vpnTunnelId string, vpnRouteId string, params *schema.RequestParameters) (*http.Response, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if vpnTunnelId == "" {
		return nil, fmt.Errorf("VPN tunnel ID cannot be empty")
	}
	if vpnRouteId == "" {
		return nil, fmt.Errorf("VPN route ID cannot be empty")
	}

	path := fmt.Sprintf(VpnRoutePath, projectId, vpnTunnelId, vpnRouteId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	return s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
}

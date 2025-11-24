package network

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/impl/logger/noop"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

func TestListVpcPeeringRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpcs/vpc-123/vpcPeerings/peering-1/routes" {
			w.WriteHeader(http.StatusOK)
			resp := types.VPCPeeringRouteList{
				ListResponse: types.ListResponse{Total: 1},
				Values: []types.VPCPeeringRouteResponse{
					{
						Metadata: types.RegionalResourceMetadataRequest{
							ResourceMetadataRequest: types.ResourceMetadataRequest{
								Name: "route-1",
							},
						},
						Properties: types.VPCPeeringRoutePropertiesResponse{
							LocalNetworkAddress:  "10.0.0.0/16",
							RemoteNetworkAddress: "10.1.0.0/16",
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &restclient.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &noop.NoOpLogger{},
	}
	c, err := restclient.NewClient(cfg, cfg.Logger)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewVPCPeeringRoutesClientImpl(c)

	resp, err := svc.List(context.Background(), "test-project", "vpc-123", "peering-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
		t.Errorf("expected 1 vpc peering route")
	}
}

func TestGetVpcPeeringRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpcs/vpc-123/vpcPeerings/peering-1/routes/route-1" {
			w.WriteHeader(http.StatusOK)
			resp := types.VPCPeeringRouteResponse{
				Metadata: types.RegionalResourceMetadataRequest{
					ResourceMetadataRequest: types.ResourceMetadataRequest{
						Name: "route-1",
					},
				},
				Properties: types.VPCPeeringRoutePropertiesResponse{
					LocalNetworkAddress:  "10.0.0.0/16",
					RemoteNetworkAddress: "10.1.0.0/16",
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &restclient.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &noop.NoOpLogger{},
	}
	c, err := restclient.NewClient(cfg, cfg.Logger)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewVPCPeeringRoutesClientImpl(c)

	resp, err := svc.Get(context.Background(), "test-project", "vpc-123", "peering-1", "route-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Data == nil || resp.Data.Metadata.Name != "route-1" {
		t.Errorf("expected route name 'route-1'")
	}
}

func TestDeleteVpcPeeringRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpcs/vpc-123/vpcPeerings/peering-1/routes/route-1" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &restclient.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &noop.NoOpLogger{},
	}
	c, err := restclient.NewClient(cfg, cfg.Logger)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewVPCPeeringRoutesClientImpl(c)

	_, err = svc.Delete(context.Background(), "test-project", "vpc-123", "peering-1", "route-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

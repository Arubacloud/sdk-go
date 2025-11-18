package network

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func TestListVpcPeerings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpcs/vpc-123/vpcPeerings" {
			w.WriteHeader(http.StatusOK)
			resp := schema.VPCPeeringList{
				ListResponse: schema.ListResponse{Total: 1},
				Values: []schema.VPCPeeringResponse{
					{Metadata: schema.ResourceMetadataResponse{Name: schema.StringPtr("test-peering")}},
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &client.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &client.NoOpLogger{},
	}
	c, err := client.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewService(c)

	resp, err := svc.ListVpcPeerings(context.Background(), "test-project", "vpc-123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
		t.Errorf("expected 1 peering")
	}
}

func TestGetVpcPeering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpcs/vpc-123/vpcPeerings/peering-1" {
			w.WriteHeader(http.StatusOK)
			resp := schema.VPCPeeringResponse{
				Metadata: schema.ResourceMetadataResponse{Name: schema.StringPtr("test-peering")},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &client.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &client.NoOpLogger{},
	}
	c, err := client.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewService(c)

	resp, err := svc.GetVpcPeering(context.Background(), "test-project", "vpc-123", "peering-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Data == nil || resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "test-peering" {
		t.Errorf("expected peering name 'test-peering'")
	}
}

func TestDeleteVpcPeering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpcs/vpc-123/vpcPeerings/peering-1" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &client.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &client.NoOpLogger{},
	}
	c, err := client.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewService(c)

	_, err = svc.DeleteVpcPeering(context.Background(), "test-project", "vpc-123", "peering-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

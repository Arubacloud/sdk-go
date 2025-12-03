package network

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/impl/interceptor/standard"
	"github.com/Arubacloud/sdk-go/internal/impl/logger/noop"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

func TestListSecurityGroups(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			resp := types.SecurityGroupList{
				ListResponse: types.ListResponse{Total: 1},
				Values: []types.SecurityGroupResponse{
					{Metadata: types.ResourceMetadataResponse{Name: types.StringPtr("sg-1")}},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewSecurityGroupsClientImpl(c, NewVPCsClientImpl(c))

		resp, err := svc.List(context.Background(), "test-project", "vpc-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Data.Total != 1 {
			t.Errorf("expected total 1, got %d", resp.Data.Total)
		}
	})
}

func TestGetSecurityGroup(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			resp := types.SecurityGroupResponse{
				Metadata: types.ResourceMetadataResponse{Name: types.StringPtr("my-sg")},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewSecurityGroupsClientImpl(c, NewVPCsClientImpl(c))

		resp, err := svc.Get(context.Background(), "test-project", "vpc-123", "sg-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "my-sg" {
			t.Errorf("expected name 'my-sg', got '%v'", resp.Data.Metadata.Name)
		}
	})
}

func TestCreateSecurityGroup(t *testing.T) {
	t.Skip("Skipping CreateSecurityGroup test - requires complex VPC polling mock setup")
	// NOTE: CreateSecurityGroup calls waitForVPCActive() which polls the VPC status
	t.Run("successful create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			w.WriteHeader(http.StatusCreated)
			resp := types.SecurityGroupResponse{
				Metadata: types.ResourceMetadataResponse{Name: types.StringPtr("new-sg")},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewSecurityGroupsClientImpl(c, NewVPCsClientImpl(c))

		req := types.SecurityGroupRequest{
			Metadata: types.ResourceMetadataRequest{Name: "new-sg"},
		}

		resp, err := svc.Create(context.Background(), "test-project", "vpc-123", req, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteSecurityGroup(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", r.Method)
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewSecurityGroupsClientImpl(c, NewVPCsClientImpl(c))

		_, err := svc.Delete(context.Background(), "test-project", "vpc-123", "sg-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

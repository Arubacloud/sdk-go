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

func TestListSecurityGroupRules(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			resp := types.SecurityRuleList{
				ListResponse: types.ListResponse{Total: 1},
				Values: []types.SecurityRuleResponse{
					{Metadata: types.ResourceMetadataResponse{Name: types.StringPtr("rule-1")}},
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

		svc := NewSecurityGroupRulesClientImpl(c, NewSecurityGroupsClientImpl(c, NewVPCsClientImpl(c)))

		resp, err := svc.List(context.Background(), "test-project", "vpc-123", "sg-456", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Data.Total != 1 {
			t.Errorf("expected total 1, got %d", resp.Data.Total)
		}
	})
}

func TestGetSecurityGroupRule(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			resp := types.SecurityRuleResponse{
				Metadata: types.ResourceMetadataResponse{Name: types.StringPtr("my-rule")},
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

		svc := NewSecurityGroupRulesClientImpl(c, NewSecurityGroupsClientImpl(c, NewVPCsClientImpl(c)))

		resp, err := svc.Get(context.Background(), "test-project", "vpc-123", "sg-456", "rule-789", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "my-rule" {
			t.Errorf("expected name 'my-rule', got '%v'", resp.Data.Metadata.Name)
		}
	})
}

// Note: TestCreateSecurityGroupRule is skipped because it requires polling SecurityGroup status
// which would need complex mock server logic to handle waitForSecurityGroupActive()

func TestDeleteSecurityGroupRule(t *testing.T) {
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

		svc := NewSecurityGroupRulesClientImpl(c, NewSecurityGroupsClientImpl(c, NewVPCsClientImpl(c)))

		_, err := svc.Delete(context.Background(), "test-project", "vpc-123", "sg-456", "rule-789", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

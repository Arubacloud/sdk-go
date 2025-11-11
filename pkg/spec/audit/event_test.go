package audit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func TestListEvents(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Audit/events" {
				w.WriteHeader(http.StatusOK)
				resp := schema.AuditEventListResponse{
					ListResponse: schema.ListResponse{Total: 1},
					Values: []schema.AuditEvent{
						{
							SeverityLevel: "Information",
							LogFormat: schema.LogFormatVersion{
								Version: "1.0",
							},
							Timestamp: time.Now(),
							Operation: schema.Operation{
								Id:    "Microsoft.Compute/virtualMachines/start/action",
								Value: schema.StringPtr("Start Virtual Machine"),
							},
							Event: schema.EventInfo{
								Id:    "event-123",
								Value: schema.StringPtr("Virtual Machine Started"),
								Type:  "operational",
							},
							Category: schema.EventCategory{
								Value:       "Administrative",
								Description: schema.StringPtr("Administrative operations"),
							},
							Origin:  "user",
							Channel: "Operation",
							Status: schema.Status{
								Value:       "Succeeded",
								Description: schema.StringPtr("Operation completed successfully"),
								Code:        schema.Int32Ptr(200),
							},
							Identity: schema.Identity{
								Caller: schema.Caller{
									Subject:  "user@example.com",
									Username: schema.StringPtr("testuser"),
									Company:  schema.StringPtr("TestCompany"),
									TenantId: schema.StringPtr("tenant-123"),
								},
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

		resp, err := svc.ListEvents(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 audit event")
		}
		if resp.Data.Values[0].SeverityLevel != "Information" {
			t.Errorf("expected severity level 'Information', got %s", resp.Data.Values[0].SeverityLevel)
		}
		if resp.Data.Values[0].Operation.Id != "Microsoft.Compute/virtualMachines/start/action" {
			t.Errorf("expected operation id 'Microsoft.Compute/virtualMachines/start/action', got %s", resp.Data.Values[0].Operation.Id)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Audit/events" {
				w.WriteHeader(http.StatusOK)
				resp := schema.AuditEventListResponse{
					ListResponse: schema.ListResponse{Total: 0},
					Values:       []schema.AuditEvent{},
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

		resp, err := svc.ListEvents(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 0 {
			t.Errorf("expected 0 audit events")
		}
	})

	t.Run("with pagination", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Audit/events" {
				// Check pagination query params
				limit := r.URL.Query().Get("limit")
				offset := r.URL.Query().Get("offset")

				if limit != "10" || offset != "5" {
					t.Errorf("expected limit=10 and offset=5, got limit=%s and offset=%s", limit, offset)
				}

				w.WriteHeader(http.StatusOK)
				resp := schema.AuditEventListResponse{
					ListResponse: schema.ListResponse{
						Total: 100,
					},
					Values: []schema.AuditEvent{
						{
							SeverityLevel: "Warning",
							LogFormat: schema.LogFormatVersion{
								Version: "1.0",
							},
							Timestamp: time.Now(),
							Operation: schema.Operation{
								Id: "test-operation",
							},
							Event: schema.EventInfo{
								Id:   "event-456",
								Type: "operational",
							},
							Category: schema.EventCategory{
								Value: "Security",
							},
							Origin:  "system",
							Channel: "Security",
							Status: schema.Status{
								Value: "Failed",
								Code:  schema.Int32Ptr(403),
							},
							Identity: schema.Identity{
								Caller: schema.Caller{
									Subject: "system",
								},
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

		params := &schema.RequestParameters{
			Limit:  schema.Int32Ptr(10),
			Offset: schema.Int32Ptr(5),
		}

		resp, err := svc.ListEvents(context.Background(), "test-project", params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Total != 100 {
			t.Errorf("expected total 100, got %d", resp.Data.Total)
		}
		if len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 audit event, got %d", len(resp.Data.Values))
		}
	})
}

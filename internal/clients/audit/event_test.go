package audit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Arubacloud/sdk-go/internal/impl/interceptor/standard"
	"github.com/Arubacloud/sdk-go/internal/impl/logger/noop"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

var ()

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
				resp := types.AuditEventListResponse{
					ListResponse: types.ListResponse{Total: 1},
					Values: []types.AuditEvent{
						{
							SeverityLevel: "Information",
							LogFormat: types.LogFormatVersion{
								Version: "1.0",
							},
							Timestamp: time.Now(),
							Operation: types.Operation{
								ID:    "Microsoft.Compute/virtualMachines/start/action",
								Value: types.StringPtr("Start Virtual Machine"),
							},
							Event: types.EventInfo{
								ID:    "event-123",
								Value: types.StringPtr("Virtual Machine Started"),
								Type:  "operational",
							},
							Category: types.EventCategory{
								Value:       "Administrative",
								Description: types.StringPtr("Administrative operations"),
							},
							Origin:  "user",
							Channel: "Operation",
							Status: types.Status{
								Value:       "Succeeded",
								Description: types.StringPtr("Operation completed successfully"),
								Code:        types.Int32Ptr(200),
							},
							Identity: types.Identity{
								Caller: types.Caller{
									Subject:  "user@example.com",
									Username: types.StringPtr("testuser"),
									Company:  types.StringPtr("TestCompany"),
									TenantID: types.StringPtr("tenant-123"),
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

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewEventsClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 audit event")
		}
		if resp.Data.Values[0].SeverityLevel != "Information" {
			t.Errorf("expected severity level 'Information', got %s", resp.Data.Values[0].SeverityLevel)
		}
		if resp.Data.Values[0].Operation.ID != "Microsoft.Compute/virtualMachines/start/action" {
			t.Errorf("expected operation id 'Microsoft.Compute/virtualMachines/start/action', got %s", resp.Data.Values[0].Operation.ID)
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
				resp := types.AuditEventListResponse{
					ListResponse: types.ListResponse{Total: 0},
					Values:       []types.AuditEvent{},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewEventsClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
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
				resp := types.AuditEventListResponse{
					ListResponse: types.ListResponse{
						Total: 100,
					},
					Values: []types.AuditEvent{
						{
							SeverityLevel: "Warning",
							LogFormat: types.LogFormatVersion{
								Version: "1.0",
							},
							Timestamp: time.Now(),
							Operation: types.Operation{
								ID: "test-operation",
							},
							Event: types.EventInfo{
								ID:   "event-456",
								Type: "operational",
							},
							Category: types.EventCategory{
								Value: "Security",
							},
							Origin:  "system",
							Channel: "Security",
							Status: types.Status{
								Value: "Failed",
								Code:  types.Int32Ptr(403),
							},
							Identity: types.Identity{
								Caller: types.Caller{
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

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewEventsClientImpl(c)

		params := &types.RequestParameters{
			Limit:  types.Int32Ptr(10),
			Offset: types.Int32Ptr(5),
		}

		resp, err := svc.List(context.Background(), "test-project", params)
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

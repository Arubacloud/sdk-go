package metric

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

func TestListMetrics(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Insight/metrics" {
				w.WriteHeader(http.StatusOK)
				resp := schema.MetricListResponse{
					ListResponse: schema.ListResponse{Total: 2},
					Values: []schema.MetricResponse{
						{
							ReferenceId:   "resource-123",
							Name:          "cpu_usage",
							ReferenceName: "test-server",
							Metadata: []schema.MetricMetadata{
								{
									Field: "unit",
									Value: "percent",
								},
							},
							Data: []schema.MetricData{
								{
									Time:    "2024-01-01T00:00:00Z",
									Measure: "45.5",
								},
								{
									Time:    "2024-01-01T00:01:00Z",
									Measure: "48.2",
								},
							},
						},
						{
							ReferenceId:   "resource-123",
							Name:          "memory_usage",
							ReferenceName: "test-server",
							Metadata: []schema.MetricMetadata{
								{
									Field: "unit",
									Value: "MB",
								},
							},
							Data: []schema.MetricData{
								{
									Time:    "2024-01-01T00:00:00Z",
									Measure: "2048",
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

		resp, err := svc.ListMetrics(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 metrics")
		}
		if resp.Data.Values[0].Name != "cpu_usage" {
			t.Errorf("expected metric name 'cpu_usage', got %s", resp.Data.Values[0].Name)
		}
		if resp.Data.Values[0].ReferenceId != "resource-123" {
			t.Errorf("expected reference ID 'resource-123', got %s", resp.Data.Values[0].ReferenceId)
		}
		if len(resp.Data.Values[0].Data) != 2 {
			t.Errorf("expected 2 data points for cpu_usage")
		}
		if len(resp.Data.Values[1].Metadata) != 1 {
			t.Errorf("expected 1 metadata field for memory_usage")
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

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Insight/metrics" {
				w.WriteHeader(http.StatusOK)
				resp := schema.MetricListResponse{
					ListResponse: schema.ListResponse{Total: 0},
					Values:       []schema.MetricResponse{},
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

		resp, err := svc.ListMetrics(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 0 {
			t.Errorf("expected 0 metrics")
		}
	})
}

func TestListAlerts(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Insight/alerts" {
				w.WriteHeader(http.StatusOK)
				resp := schema.AlertsListResponse{
					ListResponse: schema.ListResponse{Total: 2},
					Values: []schema.AlertResponse{
						{
							Id:                 "alert-123",
							EventId:            "event-456",
							EventName:          "High CPU Usage",
							Username:           "user@example.com",
							ServiceCategory:    "Compute",
							ServiceTypology:    "VirtualMachine",
							ResourceId:         "vm-789",
							ServiceName:        "test-vm",
							ResourceTypology:   "CloudServer",
							Metric:             "cpu_usage",
							LastReception:      time.Now(),
							Rule:               "greater_than",
							Theshold:           80,
							Um:                 "%",
							Duration:           "5m",
							ThesholdExceedence: "85%",
							Component:          "CPU",
							Email:              true,
							Panel:              true,
							Sms:                false,
							Hidden:             false,
							ExecutedAlertActions: []schema.ExecutedAlertAction{
								{
									ActionType:   schema.ActionTypeSendEmail,
									Success:      true,
									ErrorMessage: "",
								},
							},
							Actions: []schema.AlertAction{
								{
									Key:        "acknowledge",
									Disabled:   false,
									Executable: true,
								},
							},
						},
						{
							Id:                 "alert-456",
							EventId:            "event-789",
							EventName:          "Low Disk Space",
							Username:           "admin@example.com",
							ServiceCategory:    "Storage",
							ServiceTypology:    "BlockStorage",
							ResourceId:         "disk-321",
							ServiceName:        "test-disk",
							ResourceTypology:   "Volume",
							Metric:             "disk_usage",
							LastReception:      time.Now(),
							Rule:               "greater_than",
							Theshold:           90,
							Um:                 "%",
							Duration:           "10m",
							ThesholdExceedence: "95%",
							Email:              true,
							Panel:              true,
							Sms:                true,
							Hidden:             false,
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

		resp, err := svc.ListAlerts(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 alerts")
		}
		if resp.Data.Values[0].Id != "alert-123" {
			t.Errorf("expected alert ID 'alert-123', got %s", resp.Data.Values[0].Id)
		}
		if resp.Data.Values[0].EventName != "High CPU Usage" {
			t.Errorf("expected event name 'High CPU Usage', got %s", resp.Data.Values[0].EventName)
		}
		if resp.Data.Values[0].Metric != "cpu_usage" {
			t.Errorf("expected metric 'cpu_usage', got %s", resp.Data.Values[0].Metric)
		}
		if resp.Data.Values[0].Theshold != 80 {
			t.Errorf("expected threshold 80, got %d", resp.Data.Values[0].Theshold)
		}
		if !resp.Data.Values[0].Email {
			t.Errorf("expected email notification to be enabled")
		}
		if len(resp.Data.Values[0].ExecutedAlertActions) != 1 {
			t.Errorf("expected 1 executed action")
		}
		if resp.Data.Values[1].Sms != true {
			t.Errorf("expected SMS notification to be enabled for second alert")
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

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Insight/alerts" {
				w.WriteHeader(http.StatusOK)
				resp := schema.AlertsListResponse{
					ListResponse: schema.ListResponse{Total: 0},
					Values:       []schema.AlertResponse{},
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

		resp, err := svc.ListAlerts(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 0 {
			t.Errorf("expected 0 alerts")
		}
	})

	t.Run("with filtering", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Insight/alerts" {
				// Verify filter parameters were passed
				filter := r.URL.Query().Get("filter")
				if filter != "resourceId eq 'vm-789'" {
					t.Errorf("expected filter 'resourceId eq 'vm-789'', got %s", filter)
				}

				w.WriteHeader(http.StatusOK)
				resp := schema.AlertsListResponse{
					ListResponse: schema.ListResponse{Total: 1},
					Values: []schema.AlertResponse{
						{
							Id:            "alert-123",
							EventName:     "Filtered Alert",
							ResourceId:    "vm-789",
							Metric:        "cpu_usage",
							LastReception: time.Now(),
							Email:         true,
							Panel:         true,
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
			Filter: schema.StringPtr("resourceId eq 'vm-789'"),
		}

		resp, err := svc.ListAlerts(context.Background(), "test-project", params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 filtered alert")
		}
		if resp.Data.Values[0].ResourceId != "vm-789" {
			t.Errorf("expected resource ID 'vm-789' in filtered result")
		}
	})
}

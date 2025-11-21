package schedule

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

func TestListScheduleJobs(t *testing.T) {
	t.Run("successful_list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs" {
				w.WriteHeader(http.StatusOK)
				resp := types.JobList{
					ListResponse: types.ListResponse{Total: 2},
					Values: []types.JobResponse{
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("daily-backup"),
								ID:   types.StringPtr("job-123"),
							},
							Properties: types.JobPropertiesResponse{
								Enabled:       true,
								JobType:       types.TypeJobRecurring,
								Cron:          types.StringPtr("0 2 * * *"),
								ExecuteUntil:  types.StringPtr("2025-12-31T23:59:59Z"),
								NextExecution: types.StringPtr("2025-11-12T02:00:00Z"),
								Steps: []types.JobStepResponse{
									{
										Name:        types.StringPtr("backup-step"),
										ResourceURI: types.StringPtr("/projects/test-project/providers/Aruba.Storage/block-storages/vol-123"),
										ActionURI:   types.StringPtr("/snapshot"),
										HttpVerb:    types.StringPtr("POST"),
									},
								},
							},
							Status: types.ResourceStatus{
								State: types.StringPtr("active"),
							},
						},
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("oneshot-task"),
								ID:   types.StringPtr("job-456"),
							},
							Properties: types.JobPropertiesResponse{
								Enabled:    true,
								JobType:    types.TypeJobOneShot,
								ScheduleAt: types.StringPtr("2025-11-15T10:00:00Z"),
								Steps: []types.JobStepResponse{
									{
										Name:        types.StringPtr("shutdown-step"),
										ResourceURI: types.StringPtr("/projects/test-project/providers/Aruba.Compute/cloudservers/vm-123"),
										ActionURI:   types.StringPtr("/stop"),
										HttpVerb:    types.StringPtr("POST"),
									},
								},
							},
							Status: types.ResourceStatus{
								State: types.StringPtr("scheduled"),
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
			Logger:         &restclient.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewJobsClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 jobs")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "daily-backup" {
			t.Errorf("expected name 'daily-backup'")
		}
		if resp.Data.Values[0].Properties.JobType != types.TypeJobRecurring {
			t.Errorf("expected recurring job type")
		}
		if resp.Data.Values[1].Properties.JobType != types.TypeJobOneShot {
			t.Errorf("expected one-shot job type")
		}
	})

	t.Run("empty_list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs" {
				w.WriteHeader(http.StatusOK)
				resp := types.JobList{
					ListResponse: types.ListResponse{Total: 0},
					Values:       []types.JobResponse{},
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
			Logger:         &restclient.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewJobsClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 0 {
			t.Errorf("expected empty list")
		}
	})
}

func TestGetScheduleJob(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs/job-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.JobResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("daily-backup"),
						ID:   types.StringPtr("job-123"),
					},
					Properties: types.JobPropertiesResponse{
						Enabled:       true,
						JobType:       types.TypeJobRecurring,
						Cron:          types.StringPtr("0 2 * * *"),
						ExecuteUntil:  types.StringPtr("2025-12-31T23:59:59Z"),
						NextExecution: types.StringPtr("2025-11-12T02:00:00Z"),
						Recurrency:    (*types.RecurrenceType)(types.StringPtr(string(types.RecurrenceTypeDaily))),
						Steps: []types.JobStepResponse{
							{
								Name:         types.StringPtr("backup-step"),
								ResourceURI:  types.StringPtr("/projects/test-project/providers/Aruba.Storage/block-storages/vol-123"),
								ActionURI:    types.StringPtr("/snapshot"),
								ActionName:   types.StringPtr("CreateSnapshot"),
								Typology:     types.StringPtr("Aruba.Storage/block-storages"),
								TypologyName: types.StringPtr("Block Storage"),
								HttpVerb:     types.StringPtr("POST"),
								Body:         types.StringPtr(`{"name":"daily-snapshot"}`),
							},
						},
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("active"),
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
			Logger:         &restclient.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewJobsClientImpl(c)

		resp, err := svc.Get(context.Background(), "test-project", "job-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "daily-backup" {
			t.Errorf("expected name 'daily-backup'")
		}
		if resp.Data.Properties.JobType != types.TypeJobRecurring {
			t.Errorf("expected recurring job type")
		}
		if resp.Data.Properties.Cron == nil || *resp.Data.Properties.Cron != "0 2 * * *" {
			t.Errorf("expected cron '0 2 * * *'")
		}
		if len(resp.Data.Properties.Steps) != 1 {
			t.Errorf("expected 1 step")
		}
	})
}

func TestCreateScheduleJob(t *testing.T) {
	t.Run("successful_create_recurring", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs" {
				w.WriteHeader(http.StatusCreated)
				resp := types.JobResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("weekly-cleanup"),
						ID:   types.StringPtr("job-789"),
					},
					Properties: types.JobPropertiesResponse{
						Enabled:       true,
						JobType:       types.TypeJobRecurring,
						Cron:          types.StringPtr("0 3 * * 0"),
						ExecuteUntil:  types.StringPtr("2026-01-01T00:00:00Z"),
						NextExecution: types.StringPtr("2025-11-17T03:00:00Z"),
						Recurrency:    (*types.RecurrenceType)(types.StringPtr(string(types.RecurrenceTypeWeekly))),
						Steps: []types.JobStepResponse{
							{
								Name:        types.StringPtr("cleanup-old-snapshots"),
								ResourceURI: types.StringPtr("/projects/test-project/providers/Aruba.Storage/snapshots"),
								ActionURI:   types.StringPtr("/cleanup"),
								HttpVerb:    types.StringPtr("POST"),
							},
						},
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("active"),
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
			Logger:         &restclient.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewJobsClientImpl(c)

		body := types.JobRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "weekly-cleanup",
				},
				Location: types.LocationRequest{Value: "it-eur"},
			},
			Properties: types.JobPropertiesRequest{
				Enabled:      true,
				JobType:      types.TypeJobRecurring,
				Cron:         types.StringPtr("0 3 * * 0"),
				ExecuteUntil: types.StringPtr("2026-01-01T00:00:00Z"),
				Steps: []types.JobStep{
					{
						Name:        types.StringPtr("cleanup-old-snapshots"),
						ResourceURI: "/projects/test-project/providers/Aruba.Storage/snapshots",
						ActionURI:   "/cleanup",
						HttpVerb:    "POST",
					},
				},
			},
		}

		resp, err := svc.Create(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "weekly-cleanup" {
			t.Errorf("expected name 'weekly-cleanup'")
		}
		if resp.Data.Properties.JobType != types.TypeJobRecurring {
			t.Errorf("expected recurring job type")
		}
	})

	t.Run("successful_create_oneshot", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs" {
				w.WriteHeader(http.StatusCreated)
				resp := types.JobResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("maintenance-window"),
						ID:   types.StringPtr("job-999"),
					},
					Properties: types.JobPropertiesResponse{
						Enabled:    true,
						JobType:    types.TypeJobOneShot,
						ScheduleAt: types.StringPtr("2025-11-20T22:00:00Z"),
						Steps: []types.JobStepResponse{
							{
								Name:        types.StringPtr("stop-servers"),
								ResourceURI: types.StringPtr("/projects/test-project/providers/Aruba.Compute/cloudservers/vm-123"),
								ActionURI:   types.StringPtr("/stop"),
								HttpVerb:    types.StringPtr("POST"),
							},
						},
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("scheduled"),
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
			Logger:         &restclient.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewJobsClientImpl(c)

		body := types.JobRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "maintenance-window",
				},
				Location: types.LocationRequest{Value: "it-eur"},
			},
			Properties: types.JobPropertiesRequest{
				Enabled:    true,
				JobType:    types.TypeJobOneShot,
				ScheduleAt: types.StringPtr("2025-11-20T22:00:00Z"),
				Steps: []types.JobStep{
					{
						Name:        types.StringPtr("stop-servers"),
						ResourceURI: "/projects/test-project/providers/Aruba.Compute/cloudservers/vm-123",
						ActionURI:   "/stop",
						HttpVerb:    "POST",
					},
				},
			},
		}

		resp, err := svc.Create(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "maintenance-window" {
			t.Errorf("expected name 'maintenance-window'")
		}
		if resp.Data.Properties.JobType != types.TypeJobOneShot {
			t.Errorf("expected one-shot job type")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "scheduled" {
			t.Errorf("expected state 'scheduled'")
		}
	})
}

func TestUpdateScheduleJob(t *testing.T) {
	t.Run("successful_update", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "PUT" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs/job-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.JobResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("updated-backup"),
						ID:   types.StringPtr("job-123"),
					},
					Properties: types.JobPropertiesResponse{
						Enabled:        false,
						JobType:        types.TypeJobRecurring,
						Cron:           types.StringPtr("0 4 * * *"),
						ExecuteUntil:   types.StringPtr("2025-12-31T23:59:59Z"),
						NextExecution:  types.StringPtr("2025-11-12T04:00:00Z"),
						DeactiveReason: (*types.DeactiveReasonDto)(types.StringPtr(string(types.DeactiveReasonManual))),
						Steps: []types.JobStepResponse{
							{
								Name:        types.StringPtr("updated-backup-step"),
								ResourceURI: types.StringPtr("/projects/test-project/providers/Aruba.Storage/block-storages/vol-456"),
								ActionURI:   types.StringPtr("/snapshot"),
								HttpVerb:    types.StringPtr("POST"),
							},
						},
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("inactive"),
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
			Logger:         &restclient.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewJobsClientImpl(c)

		body := types.JobRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "updated-backup",
				},
				Location: types.LocationRequest{Value: "it-eur"},
			},
			Properties: types.JobPropertiesRequest{
				Enabled:      false,
				JobType:      types.TypeJobRecurring,
				Cron:         types.StringPtr("0 4 * * *"),
				ExecuteUntil: types.StringPtr("2025-12-31T23:59:59Z"),
				Steps: []types.JobStep{
					{
						Name:        types.StringPtr("updated-backup-step"),
						ResourceURI: "/projects/test-project/providers/Aruba.Storage/block-storages/vol-456",
						ActionURI:   "/snapshot",
						HttpVerb:    "POST",
					},
				},
			},
		}

		resp, err := svc.Update(context.Background(), "test-project", "job-123", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "updated-backup" {
			t.Errorf("expected name 'updated-backup'")
		}
		if resp.Data.Properties.Enabled {
			t.Errorf("expected job to be disabled")
		}
		if resp.Data.Properties.Cron == nil || *resp.Data.Properties.Cron != "0 4 * * *" {
			t.Errorf("expected updated cron '0 4 * * *'")
		}
	})
}

func TestDeleteScheduleJob(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Schedule/jobs/job-123" {
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
			Logger:         &restclient.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewJobsClientImpl(c)

		_, err = svc.Delete(context.Background(), "test-project", "job-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

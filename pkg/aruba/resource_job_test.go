package aruba

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/Arubacloud/sdk-go/internal/testutil"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// --------------------------------------------------------------------------
// Compile-time interface satisfaction
// --------------------------------------------------------------------------

var (
	_ Ref     = (*Job)(nil)
	_ Wrapper = (*Job)(nil)
)

// --------------------------------------------------------------------------
// Fluent setters
// --------------------------------------------------------------------------

func TestJob_FluentSetters(t *testing.T) {
	proj := &Project{}
	proj.fromResponse(projectTestResponse("p-1", "my-proj", "/projects/p-1"))

	j := NewJob().
		IntoProject(proj).
		WithName("my-job").
		AddTag("env:prod").
		AddTag("schedule").
		AddTag("env:prod"). // dedupe
		WithLocation("ITBG-Bergamo").
		WithEnabled(true)

	if j.Name() != "my-job" {
		t.Errorf("Name() = %q", j.Name())
	}
	if tags := j.Tags(); len(tags) != 2 || tags[0] != "env:prod" || tags[1] != "schedule" {
		t.Errorf("Tags() = %v", tags)
	}
	if j.Region() != "ITBG-Bergamo" {
		t.Errorf("Region() = %q", j.Region())
	}
	if !j.Enabled() {
		t.Error("Enabled() should be true")
	}
	if j.ProjectID() != "p-1" {
		t.Errorf("ProjectID() = %q", j.ProjectID())
	}
	if j.Err() != nil {
		t.Errorf("Err() = %v", j.Err())
	}
}

// --------------------------------------------------------------------------
// IntoProject
// --------------------------------------------------------------------------

func TestJob_IntoProject_TypedRef(t *testing.T) {
	proj := &Project{}
	proj.fromResponse(projectTestResponse("p-42", "proj", "/projects/p-42"))
	j := NewJob().IntoProject(proj)
	if j.ProjectID() != "p-42" {
		t.Errorf("ProjectID() = %q", j.ProjectID())
	}
	if j.Err() != nil {
		t.Errorf("Err() = %v", j.Err())
	}
}

func TestJob_IntoProject_URIRef(t *testing.T) {
	j := NewJob().IntoProject(URI("/projects/p-uri"))
	if j.ProjectID() != "p-uri" {
		t.Errorf("ProjectID() = %q", j.ProjectID())
	}
}

func TestJob_IntoProject_BadRef(t *testing.T) {
	j := NewJob().IntoProject(URI("not-a-project-uri"))
	if j.Err() == nil {
		t.Error("expected Err() != nil for non-project URI")
	}
}

// --------------------------------------------------------------------------
// WithEnabled
// --------------------------------------------------------------------------

func TestJob_WithEnabled_True(t *testing.T) {
	j := NewJob().WithEnabled(true)
	req := j.RawRequest()
	if !req.Properties.Enabled {
		t.Error("Enabled should be true")
	}
}

func TestJob_WithEnabled_False(t *testing.T) {
	j := NewJob().WithEnabled(false)
	if j.enabled == nil || *j.enabled != false {
		t.Error("enabled *bool should be set to false")
	}
}

// --------------------------------------------------------------------------
// Schedule setters — happy paths
// --------------------------------------------------------------------------

func TestJob_OneShotAt(t *testing.T) {
	ts := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	j := NewJob().OneShotAt(ts)
	if j.Err() != nil {
		t.Fatalf("Err() = %v", j.Err())
	}
	if j.JobType() != types.JobTypeOneShot {
		t.Errorf("JobType() = %q", j.JobType())
	}
	if j.scheduleAt == nil || *j.scheduleAt != "2026-05-01T12:00:00Z" {
		t.Errorf("scheduleAt = %v", j.scheduleAt)
	}
}

func TestJob_WithCron(t *testing.T) {
	j := NewJob().WithCron("0 8 * * 1-5")
	if j.Err() != nil {
		t.Fatalf("Err() = %v", j.Err())
	}
	if j.JobType() != types.JobTypeRecurring {
		t.Errorf("JobType() = %q", j.JobType())
	}
	if j.Cron() != "0 8 * * 1-5" {
		t.Errorf("Cron() = %q", j.Cron())
	}
}

func TestJob_RecurringUntil(t *testing.T) {
	ts := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	j := NewJob().RecurringUntil(ts)
	if j.Err() != nil {
		t.Fatalf("Err() = %v", j.Err())
	}
	if j.JobType() != types.JobTypeRecurring {
		t.Errorf("JobType() = %q", j.JobType())
	}
	if j.executeUntil == nil || *j.executeUntil != "2026-12-31T00:00:00Z" {
		t.Errorf("executeUntil = %v", j.executeUntil)
	}
}

func TestJob_WithCron_And_RecurringUntil(t *testing.T) {
	ts := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	j := NewJob().WithCron("0 8 * * 1-5").RecurringUntil(ts)
	if j.Err() != nil {
		t.Fatalf("Cron+Until should not error, got: %v", j.Err())
	}
	if j.JobType() != types.JobTypeRecurring {
		t.Errorf("JobType() = %q", j.JobType())
	}
	if j.Cron() != "0 8 * * 1-5" {
		t.Errorf("Cron() = %q", j.Cron())
	}
	if j.executeUntil == nil {
		t.Error("executeUntil should be set")
	}
}

// --------------------------------------------------------------------------
// Schedule setters — mode-conflict errors
// --------------------------------------------------------------------------

func TestJob_OneShotAt_Then_WithCron_Errors(t *testing.T) {
	ts := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	j := NewJob().OneShotAt(ts).WithCron("0 8 * * 1-5")
	if j.Err() == nil {
		t.Error("expected error mixing OneShot and Recurring")
	}
}

func TestJob_OneShotAt_Then_RecurringUntil_Errors(t *testing.T) {
	ts := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	j := NewJob().OneShotAt(ts).RecurringUntil(ts)
	if j.Err() == nil {
		t.Error("expected error mixing OneShot and Recurring")
	}
}

func TestJob_WithCron_Then_OneShotAt_Errors(t *testing.T) {
	ts := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	j := NewJob().WithCron("0 8 * * 1-5").OneShotAt(ts)
	if j.Err() == nil {
		t.Error("expected error mixing Recurring and OneShot")
	}
}

// --------------------------------------------------------------------------
// JobStep sub-builder
// --------------------------------------------------------------------------

func TestJobStep_Build_Basic(t *testing.T) {
	s := NewJobStep().
		Named("restart").
		OfResource(URI("/projects/p/providers/Aruba.Compute/cloudServers/srv-1")).
		WithAction("/projects/p/providers/Aruba.Compute/cloudServers/srv-1/providers/Aruba.Compute/actions/reboot").
		WithVerb("POST").
		WithBody(`{"force":true}`)

	out := s.build()
	if out.Name == nil || *out.Name != "restart" {
		t.Errorf("Name = %v", out.Name)
	}
	if out.ResourceURI != "/projects/p/providers/Aruba.Compute/cloudServers/srv-1" {
		t.Errorf("ResourceURI = %q", out.ResourceURI)
	}
	if out.HttpVerb != "POST" {
		t.Errorf("HttpVerb = %q", out.HttpVerb)
	}
	if out.Body == nil || *out.Body != `{"force":true}` {
		t.Errorf("Body = %v", out.Body)
	}
}

func TestJobStep_OfResource_EmptyURI_Errors(t *testing.T) {
	s := NewJobStep().OfResource(URI(""))
	if s.Err() == nil {
		t.Error("expected Err() != nil for empty resource URI")
	}
}

func TestJob_AddStep_DrainErrors(t *testing.T) {
	step := NewJobStep().OfResource(URI("")) // adds error to step
	j := NewJob().AddStep(step)
	if j.Err() == nil {
		t.Error("expected step errors to be drained into job")
	}
}

func TestJob_AddStep_Nil(t *testing.T) {
	j := NewJob().AddStep(nil)
	if j.Err() != nil {
		t.Errorf("AddStep(nil) should not error: %v", j.Err())
	}
}

// --------------------------------------------------------------------------
// toRequest round-trip
// --------------------------------------------------------------------------

func TestJob_ToRequest_OneShot(t *testing.T) {
	ts := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	j := NewJob().
		IntoProject(URI("/projects/p")).
		WithName("one-shot-job").
		WithLocation("ITBG-Bergamo").
		WithEnabled(true).
		OneShotAt(ts).
		AddStep(NewJobStep().
			Named("step-1").
			OfResource(URI("/projects/p/providers/Aruba.Compute/cloudServers/s-1")).
			WithAction("/projects/p/providers/Aruba.Compute/cloudServers/s-1/actions/start").
			WithVerb("POST"))

	req := j.RawRequest()
	if req.Metadata.Name != "one-shot-job" {
		t.Errorf("Metadata.Name = %q", req.Metadata.Name)
	}
	if req.Properties.JobType != types.JobTypeOneShot {
		t.Errorf("JobType = %q", req.Properties.JobType)
	}
	if req.Properties.ScheduleAt == nil || *req.Properties.ScheduleAt != "2026-05-01T12:00:00Z" {
		t.Errorf("ScheduleAt = %v", req.Properties.ScheduleAt)
	}
	if !req.Properties.Enabled {
		t.Error("Enabled should be true")
	}
	if len(req.Properties.Steps) != 1 {
		t.Fatalf("Steps len = %d", len(req.Properties.Steps))
	}
	if req.Properties.Steps[0].HttpVerb != "POST" {
		t.Errorf("Step.HttpVerb = %q", req.Properties.Steps[0].HttpVerb)
	}
}

func TestJob_ToRequest_Recurring(t *testing.T) {
	until := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	j := NewJob().
		IntoProject(URI("/projects/p")).
		WithName("cron-job").
		WithCron("0 8 * * 1-5").
		RecurringUntil(until)

	req := j.RawRequest()
	if req.Properties.JobType != types.JobTypeRecurring {
		t.Errorf("JobType = %q", req.Properties.JobType)
	}
	if req.Properties.Cron == nil || *req.Properties.Cron != "0 8 * * 1-5" {
		t.Errorf("Cron = %v", req.Properties.Cron)
	}
	if req.Properties.ExecuteUntil == nil || *req.Properties.ExecuteUntil != "2026-12-31T00:00:00Z" {
		t.Errorf("ExecuteUntil = %v", req.Properties.ExecuteUntil)
	}
}

// --------------------------------------------------------------------------
// fromResponse hydration
// --------------------------------------------------------------------------

func jobTestResponse(name string) *types.JobResponse {
	id := "job-1"
	uri := "/projects/p/providers/Aruba.Schedule/jobs/job-1"
	state := "Active"
	schedAt := "2026-05-01T12:00:00Z"
	cronExpr := "0 8 * * 1-5"
	execUntil := "2026-12-31T00:00:00Z"
	jt := types.JobTypeOneShot
	resURI := "/projects/p/providers/Aruba.Compute/cloudServers/s-1"
	actionURI := "/projects/p/providers/Aruba.Compute/cloudServers/s-1/actions/start"
	verb := "POST"
	stepName := "step-1"
	return &types.JobResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:               &id,
			URI:              &uri,
			Name:             func() *string { s := name; return &s }(),
			Tags:             []string{"tag1"},
			LocationResponse: &types.LocationResponse{Value: "ITBG-Bergamo"},
			ProjectResponseMetadata: &types.ProjectResponseMetadata{
				ID: "p",
			},
		},
		Properties: types.JobPropertiesResponse{
			Enabled:      true,
			JobType:      jt,
			ScheduleAt:   &schedAt,
			ExecuteUntil: &execUntil,
			Cron:         &cronExpr,
			Steps: []types.JobStepResponse{
				{
					Name:        &stepName,
					ResourceURI: &resURI,
					ActionURI:   &actionURI,
					HttpVerb:    &verb,
				},
			},
		},
		Status: types.ResourceStatus{State: &state},
	}
}

func TestJob_FromResponseHydration(t *testing.T) {
	j := &Job{}
	j.fromResponse(jobTestResponse("my-job"))

	if j.Name() != "my-job" {
		t.Errorf("Name() = %q", j.Name())
	}
	if j.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", j.ProjectID())
	}
	if j.ID() != "job-1" {
		t.Errorf("ID() = %q", j.ID())
	}
	if j.Region() != "ITBG-Bergamo" {
		t.Errorf("Region() = %q", j.Region())
	}
	if !j.Enabled() {
		t.Error("Enabled() should be true")
	}
	if j.JobType() != types.JobTypeOneShot {
		t.Errorf("JobType() = %q", j.JobType())
	}
	if j.scheduleAt == nil || *j.scheduleAt != "2026-05-01T12:00:00Z" {
		t.Errorf("scheduleAt = %v", j.scheduleAt)
	}
	if j.cron == nil || *j.cron != "0 8 * * 1-5" {
		t.Errorf("cron = %v", j.cron)
	}
	if j.executeUntil == nil || *j.executeUntil != "2026-12-31T00:00:00Z" {
		t.Errorf("executeUntil = %v", j.executeUntil)
	}
	if len(j.steps) != 1 {
		t.Fatalf("steps len = %d", len(j.steps))
	}
	step := j.steps[0]
	if step.name == nil || *step.name != "step-1" {
		t.Errorf("steps[0].name = %v", step.name)
	}
	if step.resourceURI == nil || *step.resourceURI != "/projects/p/providers/Aruba.Compute/cloudServers/s-1" {
		t.Errorf("steps[0].resourceURI = %v", step.resourceURI)
	}
	if step.httpVerb == nil || *step.httpVerb != "POST" {
		t.Errorf("steps[0].httpVerb = %v", step.httpVerb)
	}
}

func TestJob_FromResponse_BackfillsProjectID_FromURI(t *testing.T) {
	id := "job-x"
	uri := "/projects/proj-abc/providers/Aruba.Schedule/jobs/job-x"
	resp := &types.JobResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:  &id,
			URI: &uri,
		},
		Properties: types.JobPropertiesResponse{},
	}
	j := &Job{}
	j.fromResponse(resp)
	if j.ProjectID() != "proj-abc" {
		t.Errorf("ProjectID() backfilled from URI = %q", j.ProjectID())
	}
}

func TestJob_FromResponse_Nil(t *testing.T) {
	j := &Job{}
	j.fromResponse(nil) // must not panic
}

// --------------------------------------------------------------------------
// jobIDsFromRef
// --------------------------------------------------------------------------

func TestJobIDsFromRef_URIRef(t *testing.T) {
	ref := URI("/projects/proj-1/providers/Aruba.Schedule/jobs/job-42")
	pid, jid, err := jobIDsFromRef(ref)
	if err != nil {
		t.Fatalf("jobIDsFromRef error: %v", err)
	}
	if pid != "proj-1" {
		t.Errorf("projectID = %q", pid)
	}
	if jid != "job-42" {
		t.Errorf("jobID = %q", jid)
	}
}

func TestJobIDsFromRef_TypedRef(t *testing.T) {
	j := &Job{}
	j.fromResponse(jobTestResponse("j"))
	pid, jid, err := jobIDsFromRef(j)
	if err != nil {
		t.Fatalf("jobIDsFromRef error: %v", err)
	}
	if pid != "p" {
		t.Errorf("projectID = %q", pid)
	}
	if jid != "job-1" {
		t.Errorf("jobID = %q", jid)
	}
}

func TestJobIDsFromRef_BadURI_NoJob(t *testing.T) {
	_, _, err := jobIDsFromRef(URI("/projects/p/providers/Aruba.Schedule"))
	if err == nil {
		t.Error("expected error when job segment missing")
	}
}

func TestJobIDsFromRef_BadURI_NoProject(t *testing.T) {
	_, _, err := jobIDsFromRef(URI("/providers/Aruba.Schedule/jobs/j"))
	if err == nil {
		t.Error("expected error when project segment missing")
	}
}

// --------------------------------------------------------------------------
// HTTP-mock adapter helper
// --------------------------------------------------------------------------

func buildJobsTestAdapter(t *testing.T, handler http.HandlerFunc) *jobsClientAdapter {
	t.Helper()
	server := testutil.NewMockServer(t, handler)
	return newJobsClientAdapter(testutil.NewClient(t, server.URL))
}

const jobSuccessBody = `{` +
	`"metadata":{"id":"job-1","name":"my-job","uri":"/projects/p/providers/Aruba.Schedule/jobs/job-1","project":{"id":"p"}},` +
	`"properties":{"enabled":true,"scheduleJobType":"OneShot","scheduleAt":"2026-05-01T12:00:00Z"},` +
	`"status":{"state":"Active"}}`

// --------------------------------------------------------------------------
// Create adapter tests
// --------------------------------------------------------------------------

func TestJobsClientAdapter_Create_Success(t *testing.T) {
	var gotBody types.JobRequest
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		if !containsSubstring(r.URL.Path, "jobs") {
			t.Errorf("path %q should contain 'jobs'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, jobSuccessBody)
	})

	ts := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	j := NewJob().
		IntoProject(URI("/projects/p")).
		WithName("my-job").
		WithLocation("ITBG-Bergamo").
		WithEnabled(true).
		OneShotAt(ts)

	result, err := adapter.Create(context.Background(), j)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if result.ID() != "job-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.Name() != "my-job" {
		t.Errorf("Name() = %q", result.Name())
	}
	if result.StatusCode() != http.StatusCreated {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	if gotBody.Metadata.Name != "my-job" {
		t.Errorf("request Metadata.Name = %q", gotBody.Metadata.Name)
	}
	if gotBody.Properties.JobType != types.JobTypeOneShot {
		t.Errorf("request JobType = %q", gotBody.Properties.JobType)
	}
}

func TestJobsClientAdapter_Create_NoProject(t *testing.T) {
	callCount := 0
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusCreated)
	})
	_, err := adapter.Create(context.Background(), NewJob().WithName("x"))
	if err == nil {
		t.Fatal("expected error when Job has no project")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without project")
	}
}

func TestJobsClientAdapter_Create_MetadataValidationError(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Missing "id" — triggers MetadataValidationError from low-level Validate()
		fmt.Fprint(w, `{"metadata":{"name":"j","uri":"/projects/p/providers/Aruba.Schedule/jobs/x"},"properties":{},"status":{}}`)
	})

	j := NewJob().IntoProject(URI("/projects/p")).WithName("j")
	result, err := adapter.Create(context.Background(), j)
	if err == nil {
		t.Fatal("expected MetadataValidationError, got nil")
	}
	var mvErr *types.MetadataValidationError
	if !errors.As(err, &mvErr) {
		t.Fatalf("expected *types.MetadataValidationError, got %T: %v", err, err)
	}
	if result == nil {
		t.Error("result wrapper should not be nil even on error")
	}
}

func TestJobsClientAdapter_Create_NonTwoXX(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"bad request"}`)
	})
	_, err := adapter.Create(context.Background(), NewJob().IntoProject(URI("/projects/p")).WithName("j"))
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
	if httpErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d", httpErr.StatusCode)
	}
}

// --------------------------------------------------------------------------
// Update adapter tests
// --------------------------------------------------------------------------

func TestJobsClientAdapter_Update_Success(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, jobSuccessBody)
	})

	// Load from response, then only update non-schedule fields to avoid mode conflict.
	j := &Job{}
	j.fromResponse(jobTestResponse("my-job"))
	j.WithName("my-job-updated").WithEnabled(false)

	result, err := adapter.Update(context.Background(), j)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if result.ID() != "job-1" {
		t.Errorf("ID() = %q", result.ID())
	}
}

func TestJobsClientAdapter_Update_NoID(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	j := NewJob().IntoProject(URI("/projects/p")).WithName("x")
	_, err := adapter.Update(context.Background(), j)
	if err == nil {
		t.Fatal("expected error when Job has no ID")
	}
}

func TestJobsClientAdapter_Update_NoProject(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	j := &Job{}
	id := "job-1"
	j.setMeta(&types.ResourceMetadataResponse{ID: &id})
	_, err := adapter.Update(context.Background(), j)
	if err == nil {
		t.Fatal("expected error when Job has no project")
	}
}

func TestJobsClientAdapter_Update_NonTwoXX(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"not found"}`)
	})
	j := &Job{}
	j.fromResponse(jobTestResponse("j"))
	_, err := adapter.Update(context.Background(), j)
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
}

// --------------------------------------------------------------------------
// Get adapter tests
// --------------------------------------------------------------------------

func TestJobsClientAdapter_Get_URIRef(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, jobSuccessBody)
	})

	ref := URI("/projects/p/providers/Aruba.Schedule/jobs/job-1")
	result, err := adapter.Get(context.Background(), ref)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "job-1" {
		t.Errorf("ID() = %q", result.ID())
	}
}

func TestJobsClientAdapter_Get_TypedRef(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, jobSuccessBody)
	})

	j := &Job{}
	j.fromResponse(jobTestResponse("j"))
	result, err := adapter.Get(context.Background(), j)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "job-1" {
		t.Errorf("ID() = %q", result.ID())
	}
}

// --------------------------------------------------------------------------
// Delete adapter tests
// --------------------------------------------------------------------------

func TestJobsClientAdapter_Delete_Success(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	j := &Job{}
	j.fromResponse(jobTestResponse("j"))
	if err := adapter.Delete(context.Background(), j); err != nil {
		t.Errorf("Delete error: %v", err)
	}
}

func TestJobsClientAdapter_Delete_NonTwoXX(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"not found"}`)
	})
	j := &Job{}
	j.fromResponse(jobTestResponse("j"))
	err := adapter.Delete(context.Background(), j)
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
}

// --------------------------------------------------------------------------
// List adapter tests
// --------------------------------------------------------------------------

const jobListBody = `{"total":2,"values":[` +
	`{"metadata":{"id":"job-1","name":"job-one","uri":"/projects/p/providers/Aruba.Schedule/jobs/job-1","project":{"id":"p"}},"properties":{},"status":{}},` +
	`{"metadata":{"id":"job-2","name":"job-two","uri":"/projects/p/providers/Aruba.Schedule/jobs/job-2","project":{"id":"p"}},"properties":{},"status":{}}` +
	`]}`

func TestJobsClientAdapter_List_TwoItems(t *testing.T) {
	adapter := buildJobsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, jobListBody)
	})

	list, err := adapter.List(context.Background(), URI("/projects/p"))
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if list.Total() != 2 {
		t.Errorf("Total() = %d", list.Total())
	}
	items := list.Items()
	if len(items) != 2 {
		t.Fatalf("Items() len = %d", len(items))
	}
	if items[0].Name() != "job-one" {
		t.Errorf("items[0].Name() = %q", items[0].Name())
	}
	if items[1].Name() != "job-two" {
		t.Errorf("items[1].Name() = %q", items[1].Name())
	}
}

// --------------------------------------------------------------------------
// Reflective guard
// --------------------------------------------------------------------------

func TestJobsClient_HasUpdateMethod(t *testing.T) {
	iface := reflect.TypeOf((*JobsClient)(nil)).Elem()
	if _, ok := iface.MethodByName("Update"); !ok {
		t.Error("JobsClient interface is missing the Update method")
	}
}

package aruba

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/testutil"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// --------------------------------------------------------------------------
// Compile-time Ref satisfaction
// --------------------------------------------------------------------------

var _ Ref = (*StorageBackup)(nil)

// --------------------------------------------------------------------------
// Fluent setters
// --------------------------------------------------------------------------

func TestStorageBackup_FluentSetters(t *testing.T) {
	proj := &Project{}
	proj.fromResponse(projectTestResponse("p-1", "my-proj", "/projects/p-1"))

	bkp := NewStorageBackup().
		IntoProject(proj).
		WithName("my-backup").
		AddTag("backup").
		AddTag("storage").
		AddTag("backup"). // dedupe
		InRegion("ITBG-Bergamo").
		WithType(types.StorageBackupTypeFull).
		WithRetentionDays(30).
		WithBillingPeriod("Monthly")

	if bkp.Name() != "my-backup" {
		t.Errorf("Name() = %q", bkp.Name())
	}
	if tags := bkp.Tags(); len(tags) != 2 || tags[0] != "backup" || tags[1] != "storage" {
		t.Errorf("Tags() = %v", tags)
	}
	if bkp.Region() != "ITBG-Bergamo" {
		t.Errorf("Region() = %q", bkp.Region())
	}
	if bkp.Type() != types.StorageBackupTypeFull {
		t.Errorf("Type() = %q", bkp.Type())
	}
	if bkp.RetentionDays() != 30 {
		t.Errorf("RetentionDays() = %d", bkp.RetentionDays())
	}
	if bkp.BillingPeriod() != "Monthly" {
		t.Errorf("BillingPeriod() = %q", bkp.BillingPeriod())
	}
	if bkp.ProjectID() != "p-1" {
		t.Errorf("ProjectID() = %q", bkp.ProjectID())
	}
	if bkp.Err() != nil {
		t.Errorf("Err() = %v", bkp.Err())
	}

	bkp.RemoveTag("backup")
	if tags := bkp.Tags(); len(tags) != 1 || tags[0] != "storage" {
		t.Errorf("after RemoveTag Tags() = %v", tags)
	}

	bkp.ReplaceTags("x", "y")
	if tags := bkp.Tags(); len(tags) != 2 || tags[0] != "x" || tags[1] != "y" {
		t.Errorf("after ReplaceTags Tags() = %v", tags)
	}
}

// --------------------------------------------------------------------------
// IntoProject — typed / URI / bad Ref
// --------------------------------------------------------------------------

func TestStorageBackup_IntoProject_TypedRef(t *testing.T) {
	proj := &Project{}
	proj.fromResponse(projectTestResponse("p-42", "n", "/projects/p-42"))
	bkp := NewStorageBackup().IntoProject(proj)
	if bkp.ProjectID() != "p-42" {
		t.Errorf("ProjectID() = %q", bkp.ProjectID())
	}
	if bkp.Err() != nil {
		t.Errorf("Err() = %v", bkp.Err())
	}
}

func TestStorageBackup_IntoProject_URIRef(t *testing.T) {
	bkp := NewStorageBackup().IntoProject(URI("/projects/p-uri"))
	if bkp.ProjectID() != "p-uri" {
		t.Errorf("ProjectID() = %q", bkp.ProjectID())
	}
	if bkp.Err() != nil {
		t.Errorf("Err() = %v", bkp.Err())
	}
}

func TestStorageBackup_IntoProject_BadRef(t *testing.T) {
	bkp := NewStorageBackup().IntoProject(URI("/garbage"))
	if bkp.Err() == nil {
		t.Error("expected Err() != nil for unresolvable Ref")
	}
}

// --------------------------------------------------------------------------
// WithOrigin
// --------------------------------------------------------------------------

func TestStorageBackup_WithOrigin_URIRef(t *testing.T) {
	volURI := "/projects/p/providers/Aruba.Storage/blockstorages/bs-1"
	bkp := NewStorageBackup().WithOrigin(URI(volURI))
	if bkp.OriginURI() != volURI {
		t.Errorf("OriginURI() = %q", bkp.OriginURI())
	}
	if bkp.Err() != nil {
		t.Errorf("Err() = %v", bkp.Err())
	}
}

func TestStorageBackup_WithOrigin_TypedRef(t *testing.T) {
	bs := &BlockStorage{}
	bs.fromResponse(blockStorageTestResponse("bs-42", "n", "/projects/p/providers/Aruba.Storage/blockstorages/bs-42", "p"))

	bkp := NewStorageBackup().WithOrigin(bs)
	if bkp.OriginURI() != bs.URI() {
		t.Errorf("OriginURI() = %q, want %q", bkp.OriginURI(), bs.URI())
	}
	if bkp.Err() != nil {
		t.Errorf("Err() = %v", bkp.Err())
	}
}

func TestStorageBackup_WithOrigin_EmptyURI(t *testing.T) {
	bkp := NewStorageBackup().WithOrigin(URI(""))
	if bkp.Err() == nil {
		t.Error("expected Err() != nil for empty origin URI")
	}
	if bkp.OriginURI() != "" {
		t.Errorf("OriginURI() should remain empty, got %q", bkp.OriginURI())
	}
}

// --------------------------------------------------------------------------
// WithType typed enum
// --------------------------------------------------------------------------

func TestStorageBackup_WithType_Enum(t *testing.T) {
	bkp := NewStorageBackup().WithType(types.StorageBackupTypeFull)
	if bkp.Type() != types.StorageBackupTypeFull {
		t.Errorf("Type() = %q", bkp.Type())
	}
	req := bkp.RawRequest()
	if req.Properties.StorageBackupType != types.StorageBackupTypeFull {
		t.Errorf("RawRequest().Properties.StorageBackupType = %q", req.Properties.StorageBackupType)
	}
}

// --------------------------------------------------------------------------
// toRequest round-trip
// --------------------------------------------------------------------------

func TestStorageBackup_ToRequestRoundTrip(t *testing.T) {
	volURI := "/projects/p/providers/Aruba.Storage/blockstorages/bs-1"
	bkp := NewStorageBackup().
		WithName("bkp-rt").
		AddTag("t1").AddTag("t2").
		WithLocation("ITBG-Bergamo").
		WithType(types.StorageBackupTypeFull).
		WithOrigin(URI(volURI)).
		WithRetentionDays(14).
		WithBillingPeriod("Hour")

	req := bkp.RawRequest()

	if req.Metadata.Name != "bkp-rt" {
		t.Errorf("Metadata.Name = %q", req.Metadata.Name)
	}
	if len(req.Metadata.Tags) != 2 {
		t.Errorf("Metadata.Tags = %v", req.Metadata.Tags)
	}
	if req.Metadata.Location.Value != "ITBG-Bergamo" {
		t.Errorf("Location.Value = %q", req.Metadata.Location.Value)
	}
	if req.Properties.StorageBackupType != types.StorageBackupTypeFull {
		t.Errorf("StorageBackupType = %q", req.Properties.StorageBackupType)
	}
	if req.Properties.Origin.URI != volURI {
		t.Errorf("Origin.URI = %q", req.Properties.Origin.URI)
	}
	if req.Properties.RetentionDays == nil || *req.Properties.RetentionDays != 14 {
		t.Errorf("RetentionDays = %v", req.Properties.RetentionDays)
	}
	if req.Properties.BillingPeriod == nil || *req.Properties.BillingPeriod != "Hour" {
		t.Errorf("BillingPeriod = %v", req.Properties.BillingPeriod)
	}
}

func TestStorageBackup_ToRequest_UnsetOptionals_AreNilOrEmpty(t *testing.T) {
	bkp := NewStorageBackup().WithName("bare")
	req := bkp.RawRequest()

	if req.Properties.RetentionDays != nil {
		t.Errorf("RetentionDays should be nil, got %v", req.Properties.RetentionDays)
	}
	if req.Properties.BillingPeriod != nil {
		t.Errorf("BillingPeriod should be nil, got %v", req.Properties.BillingPeriod)
	}
	if req.Properties.Origin.URI != "" {
		t.Errorf("Origin.URI should be empty, got %q", req.Properties.Origin.URI)
	}
	if req.Properties.StorageBackupType != "" {
		t.Errorf("StorageBackupType should be empty, got %q", req.Properties.StorageBackupType)
	}
}

// --------------------------------------------------------------------------
// fromResponse hydration
// --------------------------------------------------------------------------

func storageBackupTestResponse(id, name, uri, projectID string) *types.StorageBackupResponse {
	loc := &types.LocationResponse{Value: "ITBG-Bergamo"}
	state := "Active"
	billingPeriod := "Monthly"
	retentionDays := 30
	originURI := "/projects/p/providers/Aruba.Storage/blockstorages/bs-1"
	return &types.StorageBackupResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:               &id,
			URI:              &uri,
			Name:             &name,
			Tags:             []string{"tag1"},
			LocationResponse: loc,
			ProjectResponseMetadata: &types.ProjectResponseMetadata{
				ID: projectID,
			},
		},
		Properties: types.StorageBackupPropertiesResult{
			Type:          types.StorageBackupTypeFull,
			Origin:        types.ReferenceResource{URI: originURI},
			RetentionDays: &retentionDays,
			BillingPeriod: &billingPeriod,
		},
		Status: types.ResourceStatus{
			State: &state,
			DisableStatusInfo: &types.DisableStatusInfo{
				IsDisabled: false,
			},
		},
	}
}

func TestStorageBackup_FromResponseHydration(t *testing.T) {
	bkp := &StorageBackup{}
	resp := storageBackupTestResponse("bkp-1", "my-backup", "/projects/p1/providers/Aruba.Storage/backups/bkp-1", "p1")
	bkp.fromResponse(resp)

	if bkp.ID() != "bkp-1" {
		t.Errorf("ID() = %q", bkp.ID())
	}
	if bkp.URI() != "/projects/p1/providers/Aruba.Storage/backups/bkp-1" {
		t.Errorf("URI() = %q", bkp.URI())
	}
	if bkp.BackupID() != "bkp-1" {
		t.Errorf("BackupID() = %q", bkp.BackupID())
	}
	if bkp.Name() != "my-backup" {
		t.Errorf("Name() = %q", bkp.Name())
	}
	if tags := bkp.Tags(); len(tags) != 1 || tags[0] != "tag1" {
		t.Errorf("Tags() = %v", tags)
	}
	if bkp.Region() != "ITBG-Bergamo" {
		t.Errorf("Region() = %q", bkp.Region())
	}
	if bkp.State() != "Active" {
		t.Errorf("State() = %q", bkp.State())
	}
	if bkp.IsDisabled() {
		t.Error("IsDisabled() should be false")
	}
	if bkp.Type() != types.StorageBackupTypeFull {
		t.Errorf("Type() = %q", bkp.Type())
	}
	if bkp.OriginURI() != "/projects/p/providers/Aruba.Storage/blockstorages/bs-1" {
		t.Errorf("OriginURI() = %q", bkp.OriginURI())
	}
	if bkp.RetentionDays() != 30 {
		t.Errorf("RetentionDays() = %d", bkp.RetentionDays())
	}
	if bkp.BillingPeriod() != "Monthly" {
		t.Errorf("BillingPeriod() = %q", bkp.BillingPeriod())
	}
	if bkp.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", bkp.ProjectID())
	}
	if bkp.Raw() != resp {
		t.Error("Raw() should return the hydrated response pointer")
	}
}

func TestStorageBackup_FromResponsePartial(t *testing.T) {
	bkp := &StorageBackup{}
	bkp.fromResponse(nil)
	if bkp.ID() != "" || bkp.URI() != "" || bkp.Name() != "" {
		t.Error("fromResponse(nil) should be a no-op")
	}

	bkp2 := &StorageBackup{}
	bkp2.fromResponse(&types.StorageBackupResponse{})
	if bkp2.ID() != "" || bkp2.URI() != "" || bkp2.State() != "" || bkp2.BillingPeriod() != "" || bkp2.OriginURI() != "" {
		t.Error("empty response should yield zero accessor values")
	}
	if bkp2.RetentionDays() != 0 {
		t.Errorf("RetentionDays() should be 0 for empty response, got %d", bkp2.RetentionDays())
	}
}

func TestStorageBackup_FromResponseURIBackfill(t *testing.T) {
	id := "bkp-99"
	uri := "/projects/p-uri/providers/Aruba.Storage/backups/bkp-99"
	state := ""
	resp := &types.StorageBackupResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:  &id,
			URI: &uri,
		},
		Status: types.ResourceStatus{State: &state},
	}
	bkp := &StorageBackup{}
	bkp.fromResponse(resp)

	// ProjectResponseMetadata is nil → should backfill from URI.
	if bkp.ProjectID() != "p-uri" {
		t.Errorf("ProjectID() via URI backfill = %q", bkp.ProjectID())
	}
}

// --------------------------------------------------------------------------
// Ref + ancestor ID satisfaction (runtime)
// --------------------------------------------------------------------------

func TestStorageBackup_RefSatisfaction(t *testing.T) {
	bkp := &StorageBackup{}
	bkp.fromResponse(storageBackupTestResponse("bkp-99", "n", "/projects/p99/providers/Aruba.Storage/backups/bkp-99", "p99"))

	// withBackupID typed path
	bid, ok := extractID(bkp, func(r Ref) (string, bool) {
		if w, ok := r.(withBackupID); ok {
			return w.BackupID(), true
		}
		return "", false
	}, "backups")
	if !ok || bid != "bkp-99" {
		t.Errorf("extractID via withBackupID = (%q, %v)", bid, ok)
	}

	// withProjectID typed path
	pid, ok := extractID(bkp, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid != "p99" {
		t.Errorf("extractID via withProjectID = (%q, %v)", pid, ok)
	}
}

// --------------------------------------------------------------------------
// backupIDsFromRef helper
// --------------------------------------------------------------------------

func TestBackupIDsFromRef_TypedRef(t *testing.T) {
	bkp := &StorageBackup{}
	bkp.fromResponse(storageBackupTestResponse("bid", "n", "/projects/p/providers/Aruba.Storage/backups/bid", "p"))
	pid, bid, err := backupIDsFromRef(bkp)
	if err != nil || pid != "p" || bid != "bid" {
		t.Errorf("backupIDsFromRef typed = (%q, %q, %v)", pid, bid, err)
	}
}

func TestBackupIDsFromRef_URIRef(t *testing.T) {
	ref := URI("/projects/p/providers/Aruba.Storage/backups/bkp-1")
	pid, bid, err := backupIDsFromRef(ref)
	if err != nil || pid != "p" || bid != "bkp-1" {
		t.Errorf("backupIDsFromRef URI = (%q, %q, %v)", pid, bid, err)
	}
}

func TestBackupIDsFromRef_BadURI_MissingBackup(t *testing.T) {
	_, _, err := backupIDsFromRef(URI("/projects/p/providers/Aruba.Storage/something/else"))
	if err == nil {
		t.Error("expected error for URI without /backups/<id>")
	}
}

func TestBackupIDsFromRef_BadURI_MissingProject(t *testing.T) {
	_, _, err := backupIDsFromRef(URI("/providers/Aruba.Storage/backups/bkp-1"))
	if err == nil {
		t.Error("expected error for URI without /projects/<id>")
	}
}

func TestBackupIDsFromRef_BadURI_MissingAll(t *testing.T) {
	_, _, err := backupIDsFromRef(URI("/something/else"))
	if err == nil {
		t.Error("expected error for completely unrelated URI")
	}
}

// --------------------------------------------------------------------------
// storageBackupsClientAdapter — fake low-level client for body tests
// --------------------------------------------------------------------------

type fakeStorageBackupLowLevel struct {
	createFunc func(ctx context.Context, projectID string, body types.StorageBackupRequest, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	getFunc    func(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	updateFunc func(ctx context.Context, projectID, backupID string, body types.StorageBackupRequest, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	deleteFunc func(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[any], error)
	listFunc   func(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.StorageBackupList], error)
}

func (f *fakeStorageBackupLowLevel) Create(ctx context.Context, projectID string, body types.StorageBackupRequest, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error) {
	return f.createFunc(ctx, projectID, body, params)
}
func (f *fakeStorageBackupLowLevel) Get(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error) {
	return f.getFunc(ctx, projectID, backupID, params)
}
func (f *fakeStorageBackupLowLevel) Update(ctx context.Context, projectID, backupID string, body types.StorageBackupRequest, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error) {
	return f.updateFunc(ctx, projectID, backupID, body, params)
}
func (f *fakeStorageBackupLowLevel) Delete(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[any], error) {
	return f.deleteFunc(ctx, projectID, backupID, params)
}
func (f *fakeStorageBackupLowLevel) List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.StorageBackupList], error) {
	return f.listFunc(ctx, projectID, params)
}

// --------------------------------------------------------------------------
// storageBackupsClientAdapter — HTTP mock tests
// --------------------------------------------------------------------------

func buildStorageBackupsTestAdapter(t *testing.T, handler http.HandlerFunc) *storageBackupsClientAdapter {
	t.Helper()
	server := testutil.NewMockServer(t, handler)
	return newStorageBackupsClientAdapter(testutil.NewClient(t, server.URL))
}

const storageBackupSuccessBody = `{` +
	`"metadata":{"id":"bkp-1","name":"my-backup","uri":"/projects/p/providers/Aruba.Storage/backups/bkp-1","project":{"id":"p"}},` +
	`"properties":{"type":"Full","sourceVolume":{"uri":"/projects/p/providers/Aruba.Storage/blockstorages/bs-1"},"retentionDays":30,"billingPeriod":"Monthly"},` +
	`"status":{"state":"Active"}}`

func TestStorageBackupsClientAdapter_Create_Success(t *testing.T) {
	var gotBody types.StorageBackupRequest
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		if !containsSubstring(r.URL.Path, "backups") {
			t.Errorf("path %q should contain 'backups'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, storageBackupSuccessBody)
	})

	bkp := NewStorageBackup().
		IntoProject(URI("/projects/p")).
		WithName("my-backup").
		InRegion("ITBG-Bergamo").
		WithType(types.StorageBackupTypeFull).
		WithRetentionDays(30).
		WithBillingPeriod("Monthly")

	result, err := adapter.Create(context.Background(), bkp)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if result.ID() != "bkp-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.Name() != "my-backup" {
		t.Errorf("Name() = %q", result.Name())
	}
	if result.StatusCode() != http.StatusCreated {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	if gotBody.Metadata.Name != "my-backup" {
		t.Errorf("request Name = %q", gotBody.Metadata.Name)
	}
	if gotBody.Properties.StorageBackupType != types.StorageBackupTypeFull {
		t.Errorf("request StorageBackupType = %q", gotBody.Properties.StorageBackupType)
	}
}

func TestStorageBackupsClientAdapter_Create_WithOrigin(t *testing.T) {
	volURI := "/projects/p/providers/Aruba.Storage/blockstorages/bs-1"
	var capturedBody types.StorageBackupRequest

	bkpResp := storageBackupTestResponse("bkp-1", "my-backup", "/projects/p/providers/Aruba.Storage/backups/bkp-1", "p")
	resp := &types.Response[types.StorageBackupResponse]{
		StatusCode: http.StatusCreated,
		Data:       bkpResp,
	}

	fake := &fakeStorageBackupLowLevel{
		createFunc: func(_ context.Context, _ string, body types.StorageBackupRequest, _ *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error) {
			capturedBody = body
			return resp, nil
		},
	}
	adapter := &storageBackupsClientAdapter{low: fake}

	bkp := NewStorageBackup().
		IntoProject(URI("/projects/p")).
		WithName("my-backup").
		WithOrigin(URI(volURI))

	result, err := adapter.Create(context.Background(), bkp)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if capturedBody.Properties.Origin.URI != volURI {
		t.Errorf("Origin.URI in request = %q", capturedBody.Properties.Origin.URI)
	}
	if result.ID() != "bkp-1" {
		t.Errorf("ID() = %q", result.ID())
	}
}

func TestStorageBackupsClientAdapter_Create_NoProject(t *testing.T) {
	callCount := 0
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusCreated)
	})

	_, err := adapter.Create(context.Background(), NewStorageBackup().WithName("x"))
	if err == nil {
		t.Fatal("expected error when StorageBackup has no project")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without project")
	}
}

func TestStorageBackupsClientAdapter_Create_MetadataValidationError(t *testing.T) {
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Missing "id" field — triggers MetadataValidationError
		fmt.Fprint(w, `{"metadata":{"name":"bkp","uri":"/projects/p/providers/Aruba.Storage/backups/x"},"properties":{},"status":{}}`)
	})

	bkp := NewStorageBackup().IntoProject(URI("/projects/p")).WithName("bkp")
	result, err := adapter.Create(context.Background(), bkp)
	if err == nil {
		t.Fatal("expected MetadataValidationError, got nil")
	}
	var mvErr *types.MetadataValidationError
	if !errors.As(err, &mvErr) {
		t.Fatalf("expected *types.MetadataValidationError, got %T: %v", err, err)
	}
	if result == nil {
		t.Fatal("result must be non-nil alongside MetadataValidationError")
	}
}

func TestStorageBackupsClientAdapter_Create_NonTwoXX(t *testing.T) {
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Validation Failed", "name is required", 422))
	})

	bkp := NewStorageBackup().IntoProject(URI("/projects/p"))
	result, err := adapter.Create(context.Background(), bkp)
	if err == nil {
		t.Fatal("expected error on 422")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
	if httpErr.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("HTTPError.StatusCode = %d", httpErr.StatusCode)
	}
	if result == nil {
		t.Fatal("result must be non-nil on non-2xx")
	}
}

func TestStorageBackupsClientAdapter_Get_URIRef(t *testing.T) {
	var capturedPath string
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, storageBackupSuccessBody)
	})

	ref := URI("/projects/p/providers/Aruba.Storage/backups/bkp-1")
	result, err := adapter.Get(context.Background(), ref)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "bkp-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", result.ProjectID())
	}
	if !containsSubstring(capturedPath, "backups") {
		t.Errorf("path %q should contain 'backups'", capturedPath)
	}
}

func TestStorageBackupsClientAdapter_Get_TypedRef(t *testing.T) {
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, storageBackupSuccessBody)
	})

	existing := &StorageBackup{}
	existing.fromResponse(storageBackupTestResponse("bkp-1", "n", "/projects/p/providers/Aruba.Storage/backups/bkp-1", "p"))

	result, err := adapter.Get(context.Background(), existing)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "bkp-1" {
		t.Errorf("ID() = %q", result.ID())
	}
}

func TestStorageBackupsClientAdapter_Update_Success(t *testing.T) {
	var capturedBody types.StorageBackupRequest
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&capturedBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"metadata":{"id":"bkp-1","name":"my-backup","uri":"/projects/p/providers/Aruba.Storage/backups/bkp-1","project":{"id":"p"}},"properties":{"type":"Full","sourceVolume":{},"retentionDays":30},"status":{}}`)
	})

	bkp := &StorageBackup{}
	bkp.fromResponse(storageBackupTestResponse("bkp-1", "my-backup", "/projects/p/providers/Aruba.Storage/backups/bkp-1", "p"))
	bkp.WithRetentionDays(30)

	result, err := adapter.Update(context.Background(), bkp)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if result.RetentionDays() != 30 {
		t.Errorf("RetentionDays() = %d", result.RetentionDays())
	}
	if capturedBody.Properties.RetentionDays == nil || *capturedBody.Properties.RetentionDays != 30 {
		t.Errorf("request RetentionDays = %v", capturedBody.Properties.RetentionDays)
	}
}

func TestStorageBackupsClientAdapter_Update_NoID(t *testing.T) {
	callCount := 0
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	bkp := NewStorageBackup().IntoProject(URI("/projects/p")).WithName("x")
	_, err := adapter.Update(context.Background(), bkp)
	if err == nil {
		t.Fatal("expected error when StorageBackup has no ID")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made when ID is missing")
	}
}

func TestStorageBackupsClientAdapter_Update_NoProject(t *testing.T) {
	callCount := 0
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	bkp := &StorageBackup{}
	bkp.fromResponse(&types.StorageBackupResponse{
		Metadata: types.ResourceMetadataResponse{
			ID: strPtr("bkp-1"),
		},
		Status: types.ResourceStatus{},
	})

	_, err := adapter.Update(context.Background(), bkp)
	if err == nil {
		t.Fatal("expected error when StorageBackup has no project")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without project")
	}
}

func TestStorageBackupsClientAdapter_Delete_Success(t *testing.T) {
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := adapter.Delete(context.Background(), URI("/projects/p/providers/Aruba.Storage/backups/bkp-1"))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestStorageBackupsClientAdapter_Delete_NonTwoXX(t *testing.T) {
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "backup not found", 404))
	})

	err := adapter.Delete(context.Background(), URI("/projects/p/providers/Aruba.Storage/backups/missing"))
	if err == nil {
		t.Fatal("expected error on 404")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T", err)
	}
	if httpErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d", httpErr.StatusCode)
	}
}

func TestStorageBackupsClientAdapter_List_TwoItems(t *testing.T) {
	adapter := buildStorageBackupsTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"total":2,"self":"","prev":"","next":"","first":"","last":"","values":[`+
			`{"metadata":{"id":"bkp-1","name":"n1","uri":"/projects/p/providers/Aruba.Storage/backups/bkp-1","project":{"id":"p"}},"properties":{"type":"Full","sourceVolume":{"uri":"/projects/p/providers/Aruba.Storage/blockstorages/bs-1"},"retentionDays":10},"status":{}},`+
			`{"metadata":{"id":"bkp-2","name":"n2","uri":"/projects/p/providers/Aruba.Storage/backups/bkp-2","project":{"id":"p"}},"properties":{"type":"Incremental","sourceVolume":{"uri":"/projects/p/providers/Aruba.Storage/blockstorages/bs-2"},"retentionDays":20},"status":{}}`+
			`]}`)
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
	if items[0].ID() != "bkp-1" || items[0].Name() != "n1" {
		t.Errorf("items[0] = {%q, %q}", items[0].ID(), items[0].Name())
	}
	if items[0].Type() != types.StorageBackupTypeFull {
		t.Errorf("items[0].Type() = %q", items[0].Type())
	}
	if items[0].RetentionDays() != 10 {
		t.Errorf("items[0].RetentionDays() = %d", items[0].RetentionDays())
	}
	if items[1].ID() != "bkp-2" || items[1].Type() != types.StorageBackupTypeIncremental {
		t.Errorf("items[1] ID=%q Type=%q", items[1].ID(), items[1].Type())
	}
	if items[0].ProjectID() != "p" {
		t.Errorf("items[0].ProjectID() = %q", items[0].ProjectID())
	}
}

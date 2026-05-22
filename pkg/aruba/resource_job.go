package aruba

import (
	"context"
	"fmt"
	"time"

	"github.com/Arubacloud/sdk-go/internal/clients/schedule"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// ---- Wrapper ----

// Job is the wrapper for an Aruba Cloud scheduled job (a direct child of a Project).
// Construct with aruba.NewJob() and bind via IntoProject(project).
//
// Family A: regional, Metadata/Properties envelope, location-aware.
// Supports full CRUD. Create and Update share the same request type (JobRequest) —
// there is no narrower JobUpdateRequest.
//
// Schedule modes:
//   - One-shot:  OneShotAt(t)                → JobType=OneShot, ScheduleAt=RFC3339
//   - Recurring: WithCron(expr) [+ RecurringUntil(t)] → JobType=Recurring
//
// Setter-time error if you mix the two modes.
//
// Path: /projects/{projectID}/providers/Aruba.Schedule/jobs[/{jobID}]
type Job struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	httpEnvelopeMixin

	// Schedule cache (request-side).
	enabled      *bool          // *bool so we can distinguish unset from false
	jobType      *types.JobType // implied from setter usage
	scheduleAt   *string        // RFC3339; set by OneShotAt
	cron         *string        // set by WithCron
	executeUntil *string        // RFC3339; set by RecurringUntil

	// Sub-builders.
	steps []*JobStep

	response *types.JobResponse
}

// NewJob returns a fresh *Job ready for fluent setters and a Create call.
// Binds projectScopedMixin's error sink so IntoProject failures surface via Err().
func NewJob() *Job {
	j := &Job{}
	j.projectScopedMixin = bindProjectScoped(&j.errMixin)
	return j
}

// Setters — chainable, general → specific

// IntoProject binds this Job to its parent project. Required before Create.
func (j *Job) IntoProject(p Ref) *Job { j.intoProject(p); return j }

// Named sets the resource name. Required by the API.
func (j *Job) Named(n string) *Job { j.named(n); return j }

// AddTag appends a tag for filtering and accounting.
func (j *Job) AddTag(t string) *Job { j.addTag(t); return j }

// RemoveTag removes a previously-added tag. No-op if absent.
func (j *Job) RemoveTag(t string) *Job { j.removeTag(t); return j }

// ReplaceTags replaces the entire tag set with the given values.
func (j *Job) ReplaceTags(ts ...string) *Job { j.replaceTags(ts...); return j }

// InRegion sets the region for this resource.
func (j *Job) InRegion(region Region) *Job { j.inRegion(region); return j }

// WithEnabled sets whether the job is active. Pass false to disable the job via Update.
func (j *Job) WithEnabled(enabled bool) *Job { j.enabled = &enabled; return j }

// OneShotAt schedules a one-time execution at t (UTC, RFC3339).
// Returns an error if a Recurring schedule has already been configured.
func (j *Job) OneShotAt(t time.Time) *Job {
	if !j.requireMode(types.JobTypeOneShot, "OneShotAt") {
		return j
	}
	s := t.UTC().Format(time.RFC3339)
	j.scheduleAt = &s
	return j
}

// WithCron sets the cron expression for a recurring job.
// Returns an error if a OneShot schedule has already been configured.
func (j *Job) WithCron(expr string) *Job {
	if !j.requireMode(types.JobTypeRecurring, "WithCron") {
		return j
	}
	j.cron = &expr
	return j
}

// RecurringUntil sets the end date for a recurring job (UTC, RFC3339).
// Returns an error if a OneShot schedule has already been configured.
func (j *Job) RecurringUntil(t time.Time) *Job {
	if !j.requireMode(types.JobTypeRecurring, "RecurringUntil") {
		return j
	}
	s := t.UTC().Format(time.RFC3339)
	j.executeUntil = &s
	return j
}

func (j *Job) requireMode(want types.JobType, label string) bool {
	if j.jobType != nil && *j.jobType != want {
		j.addErr(fmt.Errorf("%s: cannot mix %s and %s schedule modes", label, *j.jobType, want))
		return false
	}
	j.jobType = &want
	return true
}

// AddStep appends step to the job's step list.
// Errors accumulated on step are drained into j at attachment time.
func (j *Job) AddStep(step *JobStep) *Job {
	if step == nil {
		return j
	}
	for _, e := range step.errs {
		j.addErr(e)
	}
	j.steps = append(j.steps, step)
	return j
}

// Getters — general → specific

// URI satisfies Ref by returning the server-assigned canonical URI, or "" if Create hasn't run yet.
func (j *Job) URI() string { return j.RespURI() }

// JobID satisfies withJobID so child wrappers can extract this ID by typed assertion.
func (j *Job) JobID() string { return j.ID() }

// Raw shadows responseMetadataMixin.Raw() with the typed Job response.
func (j *Job) Raw() *types.JobResponse { return j.response }

// RawRequest returns what toRequest() would emit right now.
func (j *Job) RawRequest() types.JobRequest { return j.toRequest() }

// Enabled returns whether the job is active.
func (j *Job) Enabled() bool {
	if j.response != nil {
		return j.response.Properties.Enabled
	}
	if j.enabled != nil {
		return *j.enabled
	}
	return false
}

// JobType returns the schedule type (OneShot or Recurring) from the response or local state.
func (j *Job) JobType() types.JobType {
	if j.response != nil && j.response.Properties.JobType != "" {
		return j.response.Properties.JobType
	}
	if j.jobType != nil {
		return *j.jobType
	}
	return ""
}

// Cron returns the cron expression from the response or local state, or "" if unset.
func (j *Job) Cron() string {
	if j.response != nil && j.response.Properties.Cron != nil {
		return *j.response.Properties.Cron
	}
	if j.cron != nil {
		return *j.cron
	}
	return ""
}

// Wire converters

// toRequest assembles the Create/Update body from current setter state. Defaults are applied at the wire boundary.
func (j *Job) toRequest() types.JobRequest {
	props := types.JobPropertiesRequest{
		ScheduleAt:   j.scheduleAt,
		ExecuteUntil: j.executeUntil,
		Cron:         j.cron,
		Enabled:      j.enabled,
	}
	if j.jobType != nil {
		props.JobType = *j.jobType
	}
	if len(j.steps) > 0 {
		props.Steps = make([]types.JobStep, 0, len(j.steps))
		for _, s := range j.steps {
			props.Steps = append(props.Steps, s.build())
		}
	}
	return types.JobRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: j.toMetadata(),
			Location:                j.toLocation(),
		},
		Properties: props,
	}
}

// fromResponse hydrates the wrapper from a server reply. Nil-safe.
func (j *Job) fromResponse(resp *types.JobResponse) {
	if resp == nil {
		return
	}
	j.response = resp
	j.setMeta(&resp.Metadata)
	j.named(jobDeref(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		j.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		j.inRegion(resp.Metadata.LocationResponse.Value)
	}
	j.setStatus(&resp.Status)

	// Hydrate request-side cache.
	e := resp.Properties.Enabled
	j.enabled = &e
	if resp.Properties.JobType != "" {
		jt := resp.Properties.JobType
		j.jobType = &jt
	}
	if resp.Properties.ScheduleAt != nil {
		v := *resp.Properties.ScheduleAt
		j.scheduleAt = &v
	}
	if resp.Properties.Cron != nil {
		v := *resp.Properties.Cron
		j.cron = &v
	}
	if resp.Properties.ExecuteUntil != nil {
		v := *resp.Properties.ExecuteUntil
		j.executeUntil = &v
	}
	j.steps = jobRebuildSteps(resp.Properties.Steps)

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		j.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if j.projectID == "" && j.RespURI() != "" {
		if pid := parseURIIDs(j.RespURI())["projects"]; pid != "" {
			j.projectID = pid
		}
	}
}

// jobRebuildSteps converts response-side steps (all *string) back to sub-builders,
// dropping response-only enrichments (ActionName, Typology, TypologyName).
func jobRebuildSteps(steps []types.JobStepResponse) []*JobStep {
	if steps == nil {
		return nil
	}
	result := make([]*JobStep, 0, len(steps))
	for _, rs := range steps {
		s := &JobStep{}
		if rs.Name != nil {
			v := *rs.Name
			s.name = &v
		}
		if rs.ResourceURI != nil {
			v := *rs.ResourceURI
			s.resourceURI = &v
		}
		if rs.ActionURI != nil {
			v := *rs.ActionURI
			s.actionURI = &v
		}
		if rs.HttpVerb != nil {
			v := HTTPVerb(*rs.HttpVerb)
			s.httpVerb = &v
		}
		if rs.Body != nil {
			v := *rs.Body
			s.body = &v
		}
		result = append(result, s)
	}
	return result
}

func jobDeref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func jobIDsFromRef(ref Ref) (projectID, jobID string, err error) {
	jid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withJobID); ok {
			return w.JobID(), true
		}
		return "", false
	}, "jobs")
	if !ok || jid == "" {
		return "", "", fmt.Errorf("cannot determine Job ID from Ref %q", ref.URI())
	}
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid == "" {
		return "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return pid, jid, nil
}

// ---- Low-level client interface ----

// jobsLowLevelClient is the contract the wrapper depends on. Returning
// *types.Response[T] preserves HTTP envelope details (status code, headers,
// raw body) for the wrapper's diagnostics.
type jobsLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.JobList], error)
	Get(ctx context.Context, projectID, jobID string, params *types.RequestParameters) (*types.Response[types.JobResponse], error)
	Create(ctx context.Context, projectID string, body types.JobRequest, params *types.RequestParameters) (*types.Response[types.JobResponse], error)
	Update(ctx context.Context, projectID, jobID string, body types.JobRequest, params *types.RequestParameters) (*types.Response[types.JobResponse], error)
	Delete(ctx context.Context, projectID, jobID string, params *types.RequestParameters) (*types.Response[any], error)
}

// ---- Adapter ----

// jobsClientAdapter bridges the wrapper API (chainable, error-accumulating,
// wire-shape-hidden) to the low-level client (parameter-explicit, returning
// typed wire structs). Translates Job ↔ types.JobRequest/Response and
// surfaces HTTP errors as *aruba.HTTPError.
type jobsClientAdapter struct {
	low  jobsLowLevelClient
	rest *restclient.Client
}

var _ JobsClient = (*jobsClientAdapter)(nil)

func newJobsClientAdapter(rest *restclient.Client) *jobsClientAdapter {
	if rest == nil {
		return &jobsClientAdapter{}
	}
	return &jobsClientAdapter{low: schedule.NewJobsClientImpl(rest), rest: rest}
}

// Create posts a new Job to the API and hydrates the wrapper from the response.
func (a *jobsClientAdapter) Create(ctx context.Context, j *Job, opts ...CallOption) (*Job, error) {
	if err := j.Err(); err != nil {
		return j, err
	}
	if j.ProjectID() == "" {
		return j, fmt.Errorf("Create: Job has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, j.ProjectID(), j.toRequest(), rp)
	populateHTTPEnvelope(&j.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		j.fromResponse(resp.Data)
		j.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, j)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				j.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return j, err
	}
	if resp != nil && !resp.IsSuccess() {
		return j, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return j, nil
}

// Update sends a PUT for the current wrapper state. Requires ID and parent.
func (a *jobsClientAdapter) Update(ctx context.Context, j *Job, opts ...CallOption) (*Job, error) {
	if err := j.Err(); err != nil {
		return j, err
	}
	if j.JobID() == "" {
		return j, fmt.Errorf("Update: Job has no ID")
	}
	if j.ProjectID() == "" {
		return j, fmt.Errorf("Update: Job has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, j.ProjectID(), j.JobID(), j.toRequest(), rp)
	populateHTTPEnvelope(&j.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		j.fromResponse(resp.Data)
		j.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, j)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				j.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return j, err
	}
	if resp != nil && !resp.IsSuccess() {
		return j, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return j, nil
}

// Get fetches a Job by Ref and returns a freshly hydrated wrapper.
func (a *jobsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*Job, error) {
	projectID, jobID, err := jobIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, jobID, rp)
	out := &Job{}
	out.projectID = projectID
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
		out.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, out)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				out.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if out.projectID == "" {
		out.projectID = projectID
	}
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

// Delete removes the Job identified by Ref.
func (a *jobsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, jobID, err := jobIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, jobID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

// List returns a paginated list of Jobs in the given parent scope.
func (a *jobsClientAdapter) List(ctx context.Context, parent Ref, opts ...CallOption) (*List[*Job], error) {
	projectID, err := projectIDFromRef(parent)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*Job
	if resp != nil && resp.Data != nil {
		items = make([]*Job, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			j := &Job{}
			j.projectID = projectID
			j.fromResponse(&resp.Data.Values[i])
			j.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, j)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					j.fromResponse(fresh.Raw())
				}
				return nil
			})
			if j.projectID == "" {
				j.projectID = projectID
			}
			items = append(items, j)
		}
	}
	var refetch func(ctx context.Context, pageURL string) (*List[*Job], error)
	refetch = func(ctx context.Context, pageURL string) (*List[*Job], error) {
		fetch := listPageFetch[types.JobList](a.rest, opts)
		pageResp, fetchErr := fetch(ctx, pageURL)
		if fetchErr != nil {
			return nil, fetchErr
		}
		if pageResp != nil && !pageResp.IsSuccess() {
			return nil, &HTTPError{StatusCode: pageResp.StatusCode, Body: pageResp.RawBody, ErrResp: pageResp.Error}
		}
		var pageItems []*Job
		if pageResp != nil && pageResp.Data != nil {
			pageItems = make([]*Job, 0, len(pageResp.Data.Values))
			for i := range pageResp.Data.Values {
				j := &Job{}
				j.projectID = projectID
				j.fromResponse(&pageResp.Data.Values[i])
				j.setRefresh(func(ctx context.Context) error {
					fresh, err := a.Get(ctx, j)
					if err != nil {
						return err
					}
					if fresh != nil && fresh.Raw() != nil {
						j.fromResponse(fresh.Raw())
					}
					return nil
				})
				if j.projectID == "" {
					j.projectID = projectID
				}
				pageItems = append(pageItems, j)
			}
		}
		var total2 int64
		var self2, prev2, next2, first2, last2 string
		if pageResp != nil && pageResp.Data != nil {
			total2 = pageResp.Data.Total
			self2 = pageResp.Data.Self
			prev2 = pageResp.Data.Prev
			next2 = pageResp.Data.Next
			first2 = pageResp.Data.First
			last2 = pageResp.Data.Last
		}
		return newList(pageItems, total2, self2, prev2, next2, first2, last2, pageResp, opts, refetch), nil
	}
	var total int64
	var self, prev, next, first, last string
	if resp != nil && resp.Data != nil {
		total = resp.Data.Total
		self = resp.Data.Self
		prev = resp.Data.Prev
		next = resp.Data.Next
		first = resp.Data.First
		last = resp.Data.Last
	}
	return newList(items, total, self, prev, next, first, last, resp, opts, refetch), nil
}

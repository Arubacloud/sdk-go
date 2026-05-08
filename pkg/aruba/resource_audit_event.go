package aruba

import (
	"context"
	"fmt"
	"time"

	"github.com/Arubacloud/sdk-go/internal/clients/audit"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// AuditEvent is the wrapper for an Aruba Cloud audit event.
// Instances are read-only and can only be obtained via Client.FromAudit().Events().List.
// There is no factory, no setters, and no individual-fetch endpoint.
type AuditEvent struct {
	responseMetadataMixin // shadowed: ID(), URI(), CreatedAt(), Raw()
	httpEnvelopeMixin     // RawHTTP(), StatusCode(), Headers(), RawError()

	projectID string            // back-filled from the parent Ref at List time
	response  *types.AuditEvent // backs Raw() and all field accessors
}

// ID returns the event identifier (from Event.ID), or "" when unset.
// Shadows responseMetadataMixin.ID().
func (e *AuditEvent) ID() string {
	if e.response == nil {
		return ""
	}
	return e.response.Event.ID
}

// URI returns "" — audit events have no individual fetch endpoint.
func (e *AuditEvent) URI() string { return "" }

// CreatedAt returns the event timestamp. Shadows responseMetadataMixin.CreatedAt().
func (e *AuditEvent) CreatedAt() time.Time {
	if e.response == nil {
		return time.Time{}
	}
	return e.response.Timestamp
}

// Raw returns the underlying wire payload. Shadows responseMetadataMixin.Raw().
func (e *AuditEvent) Raw() *types.AuditEvent { return e.response }

// ProjectID returns the project that owns this event.
// The value is back-filled from the parent Ref used in the List call.
func (e *AuditEvent) ProjectID() string { return e.projectID }

// SeverityLevel returns the event severity level, or "" when unset.
func (e *AuditEvent) SeverityLevel() string {
	if e.response == nil {
		return ""
	}
	return e.response.SeverityLevel
}

// Origin returns the event origin, or "" when unset.
func (e *AuditEvent) Origin() string {
	if e.response == nil {
		return ""
	}
	return e.response.Origin
}

// Channel returns the event channel, or "" when unset.
func (e *AuditEvent) Channel() string {
	if e.response == nil {
		return ""
	}
	return e.response.Channel
}

// LogFormat returns the log format version.
func (e *AuditEvent) LogFormat() types.LogFormatVersion {
	if e.response == nil {
		return types.LogFormatVersion{}
	}
	return e.response.LogFormat
}

// Operation returns the operation associated with this event.
func (e *AuditEvent) Operation() types.Operation {
	if e.response == nil {
		return types.Operation{}
	}
	return e.response.Operation
}

// Event returns the event information (ID, type, value).
func (e *AuditEvent) Event() types.EventInfo {
	if e.response == nil {
		return types.EventInfo{}
	}
	return e.response.Event
}

// Category returns the event category.
func (e *AuditEvent) Category() types.EventCategory {
	if e.response == nil {
		return types.EventCategory{}
	}
	return e.response.Category
}

// Region returns the optional region for this event, or nil when absent.
func (e *AuditEvent) Region() *types.RegionInfo {
	if e.response == nil {
		return nil
	}
	return e.response.Region
}

// Status returns the event status.
func (e *AuditEvent) Status() types.Status {
	if e.response == nil {
		return types.Status{}
	}
	return e.response.Status
}

// SubStatus returns the optional sub-status, or nil when absent.
func (e *AuditEvent) SubStatus() *types.SubStatus {
	if e.response == nil {
		return nil
	}
	return e.response.SubStatus
}

// Identity returns the caller identity for this event.
func (e *AuditEvent) Identity() types.Identity {
	if e.response == nil {
		return types.Identity{}
	}
	return e.response.Identity
}

// Properties returns the arbitrary properties map, or nil when absent.
func (e *AuditEvent) Properties() map[string]interface{} {
	if e.response == nil {
		return nil
	}
	return e.response.Properties
}

// Actions returns the available actions for this event, or nil when absent.
func (e *AuditEvent) Actions() []types.Action {
	if e.response == nil {
		return nil
	}
	return e.response.Actions
}

// CategoryID returns the category ID, or "" when absent.
func (e *AuditEvent) CategoryID() string {
	if e.response == nil || e.response.CategoryID == nil {
		return ""
	}
	return *e.response.CategoryID
}

// TypologyID returns the typology ID, or "" when absent.
func (e *AuditEvent) TypologyID() string {
	if e.response == nil || e.response.TypologyID == nil {
		return ""
	}
	return *e.response.TypologyID
}

// Title returns the event title, or "" when absent.
func (e *AuditEvent) Title() string {
	if e.response == nil || e.response.Title == nil {
		return ""
	}
	return *e.response.Title
}

func (e *AuditEvent) fromResponse(resp *types.AuditEvent) {
	if resp == nil {
		return
	}
	e.response = resp
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type auditEventsLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.AuditEventListResponse], error)
}

type auditEventsClientAdapter struct{ low auditEventsLowLevelClient }

func newAuditEventsClientAdapter(rest *restclient.Client) *auditEventsClientAdapter {
	if rest == nil {
		return &auditEventsClientAdapter{}
	}
	return &auditEventsClientAdapter{low: audit.NewEventsClientImpl(rest)}
}

func (a *auditEventsClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*AuditEvent], error) {
	projectID, err := projectIDFromRef(project)
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
	var items []*AuditEvent
	if resp != nil && resp.Data != nil {
		items = make([]*AuditEvent, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			ev := &AuditEvent{}
			ev.fromResponse(&resp.Data.Values[i])
			ev.projectID = projectID
			populateHTTPEnvelope(&ev.httpEnvelopeMixin, resp)
			items = append(items, ev)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*AuditEvent], error) {
		return nil, fmt.Errorf("List pagination by URL not yet wired; re-call List with adjusted CallOptions")
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

package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/metric"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// Metric is the wrapper for an Aruba Cloud metric.
// Instances are read-only and can only be obtained via Client.FromMetric().Metrics().List.
// There is no factory, no setters, and no individual-fetch endpoint.
type Metric struct {
	responseMetadataMixin // shadowed: ID(), Raw(); no metadata envelope on the wire
	httpEnvelopeMixin     // RawHTTP(), StatusCode(), Headers(), RawError()

	projectID string // back-filled from the parent Ref at List time
	response  *types.MetricResponse
}

// ID returns the metric reference ID, or "" when unset. Shadows responseMetadataMixin.ID().
func (m *Metric) ID() string {
	if m.response == nil {
		return ""
	}
	return m.response.ReferenceID
}

// URI returns "" — metrics have no individual fetch endpoint.
func (m *Metric) URI() string { return "" }

// Raw returns the underlying wire payload. Shadows responseMetadataMixin.Raw().
func (m *Metric) Raw() *types.MetricResponse { return m.response }

// ProjectID returns the project that owns this metric.
// The value is back-filled from the parent Ref used in the List call.
func (m *Metric) ProjectID() string { return m.projectID }

// ReferenceID returns the metric reference ID, or "" when unset.
func (m *Metric) ReferenceID() string {
	if m.response == nil {
		return ""
	}
	return m.response.ReferenceID
}

// Name returns the metric name, or "" when unset.
func (m *Metric) Name() string {
	if m.response == nil {
		return ""
	}
	return m.response.Name
}

// ReferenceName returns the metric reference name, or "" when unset.
func (m *Metric) ReferenceName() string {
	if m.response == nil {
		return ""
	}
	return m.response.ReferenceName
}

// Metadata returns the metric metadata entries, or nil when absent.
// Each entry contains a Field and Value string.
func (m *Metric) Metadata() []types.MetricMetadata {
	if m.response == nil {
		return nil
	}
	return m.response.Metadata
}

// Data returns the metric datapoints, or nil when absent.
// Each datapoint contains a Time and Measure string.
func (m *Metric) Data() []types.MetricData {
	if m.response == nil {
		return nil
	}
	return m.response.Data
}

func (m *Metric) fromResponse(resp *types.MetricResponse) {
	if resp == nil {
		return
	}
	m.response = resp
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type metricsLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.MetricListResponse], error)
}

type metricsClientAdapter struct{ low metricsLowLevelClient }

func newMetricsClientAdapter(rest *restclient.Client) *metricsClientAdapter {
	if rest == nil {
		return &metricsClientAdapter{}
	}
	return &metricsClientAdapter{low: metric.NewMetricsClientImpl(rest)}
}

func (a *metricsClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*Metric], error) {
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
	var items []*Metric
	if resp != nil && resp.Data != nil {
		items = make([]*Metric, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			met := &Metric{}
			met.fromResponse(&resp.Data.Values[i])
			met.projectID = projectID
			populateHTTPEnvelope(&met.httpEnvelopeMixin, resp)
			items = append(items, met)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*Metric], error) {
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

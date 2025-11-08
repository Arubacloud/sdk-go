package schema

// MetricMetadata represents metadata for a metric
type MetricMetadata struct {
	Field string `json:"field,omitempty"`
	Value string `json:"value,omitempty"`
}

// MetricData represents data points for a metric
type MetricData struct {
	Time    string `json:"time,omitempty"`
	Measure string `json:"measure,omitempty"`
}

// Metric represents a metric response
type MetricResponse struct {
	ReferenceId   string           `json:"referenceId,omitempty"`
	Name          string           `json:"name,omitempty"`
	ReferenceName string           `json:"referenceName,omitempty"`
	Metadata      []MetricMetadata `json:"metadata,omitempty"`
	Data          []MetricData     `json:"data,omitempty"`
}

// MetricList represents a list of metrics
type MetricListResponse struct {
	ListResponse
	Values []MetricResponse `json:"values"`
}

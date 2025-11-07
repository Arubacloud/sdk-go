package schema

// AcceptHeader defines model for acceptHeader.
type AcceptHeader string

type RequestParameters struct {
	Filter     *string       `json:"filter,omitempty"`
	Sort       *string       `json:"sort,omitempty"`
	Projection *string       `json:"projection,omitempty"`
	Accept     *AcceptHeader `json:"-"`
	Offset     *int32        `json:"offset,omitempty"`
	Limit      *int32        `json:"limit,omitempty"`
	APIVersion *string       `json:"api-version,omitempty"`
}

package types

import (
	"net/http"
	"time"
)

// Resource Metadata Request
type ResourceMetadataRequest struct {
	Name string   `json:"name"`
	Tags []string `json:"tags,omitempty"`
}

// Regional Resource Metadata Request
type RegionalResourceMetadataRequest struct {
	ResourceMetadataRequest
	Location LocationRequest `json:"location"`
}

type LocationRequest struct {
	Value string `json:"value"`
}

// Resource Metadata Response
type LocationResponse struct {
	Code    string `json:"code,omitempty"`
	Country string `json:"country,omitempty"`
	Name    string `json:"region,omitempty"`
	City    string `json:"city,omitempty"`
	Value   string `json:"value,omitempty"`
}

type ProjectResponseMetadata struct {
	ID string `json:"id,omitempty"`
}

type ResourceRequest struct {
	Metadata *ResourceMetadataRequest `json:"metadata"`
}

type TypologyResponseMetadata struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type CategoryResponseMetadata struct {
	Name     string                   `json:"name,omitempty"`
	Provider string                   `json:"provider,omitempty"`
	Typology TypologyResponseMetadata `json:"typology,omitempty"`
}

type ResourceMetadataResponse struct {
	ID                      *string                   `json:"id,omitempty"`
	URI                     *string                   `json:"uri,omitempty"`
	Name                    *string                   `json:"name,omitempty"`
	LocationResponse        *LocationResponse         `json:"location,omitempty"`
	ProjectResponseMetadata *ProjectResponseMetadata  `json:"project,omitempty"`
	Tags                    []string                  `json:"tags,omitempty"`
	Category                *CategoryResponseMetadata `json:"category,omitempty"`
	CreationDate            *time.Time                `json:"creationDate,omitempty"`
	CreatedBy               *string                   `json:"createdBy,omitempty"`
	UpdateDate              *time.Time                `json:"updateDate,omitempty"`
	UpdatedBy               *string                   `json:"updatedBy,omitempty"`
	Version                 *string                   `json:"version,omitempty"`
	CreatedUser             *string                   `json:"createdUser,omitempty"`
	UpdatedUser             *string                   `json:"updatedUser,omitempty"`
}

// Status
type PreviousStatus struct {
	State        *string    `json:"state,omitempty"`
	CreationDate *time.Time `json:"creationDate,omitempty"`
}

type DisableStatusInfo struct {
	IsDisabled bool     `json:"isDisabled,omitempty"`
	Reasons    []string `json:"reasons,omitempty"`
}

type ResourceStatus struct {
	State             *string            `json:"state,omitempty"`
	CreationDate      *time.Time         `json:"creationDate,omitempty"`
	DisableStatusInfo *DisableStatusInfo `json:"disableStatusInfo,omitempty"`
	PreviousStatus    *PreviousStatus    `json:"previousStatus,omitempty"`
	FailureReason     *string            `json:"failureReason,omitempty"`
}

// LinkedResource represents a resource linked
type LinkedResource struct {
	// URI of the linked resource
	URI string `json:"uri"`

	// StrictCorrelation indicates strict correlation with the resource
	StrictCorrelation bool `json:"strictCorrelation"`
}

type BillingPeriodResource struct {
	BillingPeriod string `json:"billingPeriod"`
}

type ReferenceResource struct {
	URI string `json:"uri"`
}

type ListResponse struct {
	// Total number of items
	Total int64 `json:"total"`

	// Self link to current page
	Self string `json:"self"`

	// Prev link to previous page
	Prev string `json:"prev"`

	// Next link to next page
	Next string `json:"next"`

	// First link to first page
	First string `json:"first"`

	// Last link to last page
	Last string `json:"last"`
}

// Response wraps an HTTP response with parsed data
type Response[T any] struct {
	// Data contains the parsed response body (for 2xx responses)
	Data *T

	// Error contains the parsed error response (for 4xx/5xx responses)
	Error *ErrorResponse

	// HTTPResponse is the underlying HTTP response
	HTTPResponse *http.Response

	// StatusCode is the HTTP status code
	StatusCode int

	// Headers contains the response headers
	Headers http.Header

	// RawBody contains the raw response body (if requested)
	RawBody []byte
}

// IsSuccess returns true if the status code is 2xx
func (r *Response[T]) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError returns true if the status code is 4xx or 5xx
func (r *Response[T]) IsError() bool {
	return r.StatusCode >= 400
}

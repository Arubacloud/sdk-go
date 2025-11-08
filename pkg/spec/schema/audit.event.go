package schema

import "time"

// Operation represents an operation in the audit log
type Operation struct {
	Id    string  `json:"id"`
	Value *string `json:"value,omitempty"`
}

// EventInfo represents event information
type EventInfo struct {
	Id    string  `json:"id"`
	Value *string `json:"value,omitempty"`
	Type  string  `json:"type"`
}

// EventCategory represents the event category
type EventCategory struct {
	Value       string  `json:"value"`
	Description *string `json:"description,omitempty"`
}

// Region represents the region information
type Region struct {
	Name             *string `json:"name,omitempty"`
	AvailabilityZone *string `json:"availabilityZone,omitempty"`
}

// Status represents the status of the event
type Status struct {
	Value       string                 `json:"value"`
	Description *string                `json:"description,omitempty"`
	Code        *int32                 `json:"code,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

// SubStatus represents the sub-status of the event
type SubStatus struct {
	Value       *string                `json:"value,omitempty"`
	Description *string                `json:"description,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

// Caller represents the caller identity
type Caller struct {
	Subject  string  `json:"subject"`
	Username *string `json:"username,omitempty"`
	Company  *string `json:"company,omitempty"`
	TenantId *string `json:"tenantId,omitempty"`
}

// Identity represents the identity information
type Identity struct {
	Caller     Caller                 `json:"caller"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Action represents an available action
type Action struct {
	Key        *string `json:"key,omitempty"`
	Disabled   *bool   `json:"disabled,omitempty"`
	Executable *bool   `json:"executable,omitempty"`
}

// LogFormatVersion represents the log format version
type LogFormatVersion struct {
	Version string `json:"version"`
}

// AuditEvent represents the complete audit event response
type AuditEvent struct {
	SeverityLevel string                 `json:"severityLevel"`
	LogFormat     LogFormatVersion       `json:"logFormat"`
	Timestamp     time.Time              `json:"@timestamp"`
	Operation     Operation              `json:"operation"`
	Event         EventInfo              `json:"event"`
	Category      EventCategory          `json:"category"`
	Region        *Region                `json:"region,omitempty"`
	Origin        string                 `json:"origin"`
	Channel       string                 `json:"channel"`
	Status        Status                 `json:"status"`
	SubStatus     *SubStatus             `json:"subStatus,omitempty"`
	Identity      Identity               `json:"identity"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
	Actions       []Action               `json:"actions,omitempty"`
	CategoryId    *string                `json:"categoryId,omitempty"`
	TypologyId    *string                `json:"typologyId,omitempty"`
	Title         *string                `json:"title,omitempty"`
}

// AuditEventListResponse represents a paginated list of audit events
type AuditEventListResponse struct {
	ListResponse
	Values []AuditEvent `json:"values"`
}

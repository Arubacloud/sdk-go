package schema

import "encoding/json"

// ErrorResponse represents an error response following RFC 7807 Problem Details
type ErrorResponse struct {
	// Type A URI reference that identifies the problem type (nullable)
	Type *string `json:"type,omitempty"`

	// Title A short, human-readable summary of the problem type (nullable)
	Title *string `json:"title,omitempty"`

	// Status The HTTP status code (nullable)
	Status *int32 `json:"status,omitempty"`

	// Detail A human-readable explanation specific to this occurrence of the problem (nullable)
	Detail *string `json:"detail,omitempty"`

	// Instance A URI reference that identifies the specific occurrence of the problem (nullable)
	Instance *string `json:"instance,omitempty"`

	// Extensions Additional properties for extensibility
	// Allows for dynamic properties with any name and value
	Extensions map[string]interface{} `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle dynamic properties
func (e *ErrorResponse) UnmarshalJSON(data []byte) error {
	type Alias ErrorResponse
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	// Unmarshal into a map to capture all fields
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Unmarshal known fields
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Extract unknown fields into Extensions
	e.Extensions = make(map[string]interface{})
	knownFields := map[string]bool{
		"type":     true,
		"title":    true,
		"status":   true,
		"detail":   true,
		"instance": true,
	}

	for key, value := range raw {
		if !knownFields[key] {
			e.Extensions[key] = value
		}
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling to include dynamic properties
func (e *ErrorResponse) MarshalJSON() ([]byte, error) {
	type Alias ErrorResponse
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	// Marshal known fields
	data, err := json.Marshal(aux)
	if err != nil {
		return nil, err
	}

	// If no extensions, return as-is
	if len(e.Extensions) == 0 {
		return data, nil
	}

	// Merge extensions
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	for key, value := range e.Extensions {
		result[key] = value
	}

	return json.Marshal(result)
}

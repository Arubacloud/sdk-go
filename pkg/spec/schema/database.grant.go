package schema

import "time"

type GrantUser struct {
	Username string `json:"username"`
}

type GrantRole struct {
	Name string `json:"name"`
}

type GrantDatabaseResponse struct {
	Name string `json:"name"`
}
type GrantRequest struct {
	User GrantUser `json:"user"`
	Role GrantRole `json:"role"`
}

type GrantResponse struct {
	User         GrantUser             `json:"user"`
	Role         GrantRole             `json:"role"`
	Database     GrantDatabaseResponse `json:"database"`
	CreationDate *time.Time            `json:"creationDate,omitempty"`
	CreatedBy    *string               `json:"createdBy,omitempty"`
}

type GrantList struct {
	ListResponse
	Values []GrantResponse `json:"values"`
}

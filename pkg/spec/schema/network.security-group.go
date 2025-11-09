package schema

type SecurityGroupPropertiesRequest struct {
	// Indicates if the security group must be a default subnet. Only one default security group for vpc is admissible.
	Default *bool `json:"default,omitempty"`
}
type SecurityGroupPropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	Default bool `json:"default,omitempty"`
}

type SecurityGroupRequest struct {
	// Metadata of the Security Group
	Metadata ResourceMetadataRequest `json:"metadata"`

	// Spec contains the Security Group specification
	Properties SecurityGroupPropertiesRequest `json:"properties"`
}

type SecurityGroupResponse struct {
	// Metadata of the Security Group
	Metadata RegionalResourceMetadataRequest `json:"metadata"`
	// Spec contains the Security Group specification
	Properties SecurityGroupPropertiesResponse `json:"properties"`

	Status ResourceStatus `json:"status,omitempty"`
}

type SecurityGroupList struct {
	ListResponse
	Values []SecurityGroupResponse `json:"values"`
}

package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// Project is the wrapper for an Aruba Cloud project.
// Construct with aruba.Project(); pass it to ProjectClient.Create / .Update,
// or receive it from .Get / .List.
type Project struct {
	errMixin
	metadataMixin
	projectScopedMixin    // self-referential: ProjectID() returns this project's own ID after Get/Create
	responseMetadataMixin // promotes ID(), RespURI(), CreatedAt(), …
	httpEnvelopeMixin

	description *string
	defaultProj bool // ProjectPropertiesRequest.Default is plain bool — no tri-state needed
}

// WithName sets the project name.
func (p *Project) WithName(n string) *Project { p.withName(n); return p }

// AddTag adds a tag (deduped).
func (p *Project) AddTag(t string) *Project { p.addTag(t); return p }

// RemoveTag removes a tag.
func (p *Project) RemoveTag(t string) *Project { p.removeTag(t); return p }

// ReplaceTags overwrites the tag list.
func (p *Project) ReplaceTags(ts ...string) *Project { p.replaceTags(ts...); return p }

// WithDescription sets the project description.
func (p *Project) WithDescription(d string) *Project { p.description = &d; return p }

// WithDefault marks the project as the account default.
func (p *Project) WithDefault(b bool) *Project { p.defaultProj = b; return p }

// URI satisfies Ref; returns the server-assigned URI, or "" before the first reply.
// ID() is promoted from responseMetadataMixin and needs no override.
func (p *Project) URI() string { return p.RespURI() }

// Description returns the description set via WithDescription, or "" if unset.
func (p *Project) Description() string {
	if p.description == nil {
		return ""
	}
	return *p.description
}

// IsDefault returns true if this project is marked as the account default.
func (p *Project) IsDefault() bool { return p.defaultProj }

// toRequest builds the wire-level request from the wrapper's slots.
func (p *Project) toRequest() types.ProjectRequest {
	return types.ProjectRequest{
		Metadata: p.toMetadata(),
		Properties: types.ProjectPropertiesRequest{
			Description: p.description,
			Default:     p.defaultProj,
		},
	}
}

// fromResponse hydrates the wrapper from a server reply. Nil-safe.
func (p *Project) fromResponse(resp *types.ProjectResponse) {
	if resp == nil {
		return
	}
	p.setMeta(&resp.Metadata)
	p.withName(projectDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		p.replaceTags(resp.Metadata.Tags...)
	}
	p.description = resp.Properties.Description
	p.defaultProj = resp.Properties.Default
	// Seed our own projectID so that ProjectID() works immediately after Create/Get,
	// enabling child-resource setters that assert withProjectID.
	if resp.Metadata.ID != nil {
		p.projectID = *resp.Metadata.ID
	}
}

func projectDerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

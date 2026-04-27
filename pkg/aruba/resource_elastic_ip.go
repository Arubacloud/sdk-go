package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// ElasticIP is the wrapper for an Aruba Cloud Elastic IP (a direct child of a Project).
// Construct with aruba.NewElasticIP() and bind it to a parent project via IntoProject(project).
type ElasticIP struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	billingPeriod *string                  // Properties.BillingPlan.BillingPeriod
	address       *string                  // Properties.Address (read-only from response)
	response      *types.ElasticIPResponse // backs Raw()
}

func (e *ElasticIP) IntoProject(p Ref) *ElasticIP          { e.intoProject(p); return e }
func (e *ElasticIP) WithName(n string) *ElasticIP          { e.withName(n); return e }
func (e *ElasticIP) AddTag(t string) *ElasticIP            { e.addTag(t); return e }
func (e *ElasticIP) RemoveTag(t string) *ElasticIP         { e.removeTag(t); return e }
func (e *ElasticIP) ReplaceTags(ts ...string) *ElasticIP   { e.replaceTags(ts...); return e }
func (e *ElasticIP) WithLocation(loc string) *ElasticIP    { e.withLocation(loc); return e }
func (e *ElasticIP) InRegion(region string) *ElasticIP     { e.withLocation(region); return e }
func (e *ElasticIP) WithBillingPeriod(p string) *ElasticIP { e.billingPeriod = &p; return e }

// URI satisfies Ref.
func (e *ElasticIP) URI() string { return e.RespURI() }

// ElasticIPID satisfies withElasticIPID so adapters can extract this ID typed.
func (e *ElasticIP) ElasticIPID() string { return e.ID() }

// Raw shadows responseMetadataMixin.Raw() with the full ElasticIP response.
func (e *ElasticIP) Raw() *types.ElasticIPResponse { return e.response }

// RawRequest returns what toRequest() would emit right now.
func (e *ElasticIP) RawRequest() types.ElasticIPRequest { return e.toRequest() }

func (e *ElasticIP) BillingPeriod() string {
	if e.billingPeriod == nil {
		return ""
	}
	return *e.billingPeriod
}

func (e *ElasticIP) Address() string {
	if e.address == nil {
		return ""
	}
	return *e.address
}

func (e *ElasticIP) toRequest() types.ElasticIPRequest {
	var bp string
	if e.billingPeriod != nil {
		bp = *e.billingPeriod
	}
	return types.ElasticIPRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: e.toMetadata(),
			Location:                e.toLocation(),
		},
		Properties: types.ElasticIPPropertiesRequest{
			BillingPlan: types.BillingPeriodResource{BillingPeriod: bp},
		},
	}
}

func (e *ElasticIP) fromResponse(resp *types.ElasticIPResponse) {
	if resp == nil {
		return
	}
	e.response = resp
	e.setMeta(&resp.Metadata)
	e.withName(elasticIPDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		e.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		e.withLocation(resp.Metadata.LocationResponse.Value)
	}
	e.setStatus(&resp.Status)
	e.setLinked(resp.Properties.LinkedResources)

	if resp.Properties.BillingPlan.BillingPeriod != "" {
		bp := resp.Properties.BillingPlan.BillingPeriod
		e.billingPeriod = &bp
	}
	if resp.Properties.Address != nil && *resp.Properties.Address != "" {
		addr := *resp.Properties.Address
		e.address = &addr
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		e.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if e.projectID == "" && e.RespURI() != "" {
		if pid := parseURIIDs(e.RespURI())["projects"]; pid != "" {
			e.projectID = pid
		}
	}
}

func elasticIPDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

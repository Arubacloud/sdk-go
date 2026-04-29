package aruba

import (
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

// DBaaS is the wrapper for an Aruba Cloud Database-as-a-Service instance
// (a direct child of a Project). Construct with aruba.NewDBaaS() and bind
// it via IntoProject(project), WithEngine, WithFlavor, WithStorage,
// WithNetworking(vpc, subnet, sg, eip), etc.
//
// Schema asymmetries:
//   - Engine: request emits Engine.ID; response returns a full
//     DBaaSEngineResponse{Type,Name,Version,...}. Engine() reads .Type
//     (the human-meaningful identifier on the response side).
//   - Flavor: request emits Flavor.Name; response returns the full
//     DBaaSFlavorResponse struct.
//   - Networking: request emits 4 raw URI strings (VPCURI/SubnetURI/
//     SecurityGroupURI/ElasticIPURI); response returns 4 *ReferenceResource
//     objects. Read-back getters (VPC/Subnet/SecurityGroup/ElasticIP)
//     prefer the response side, falling back to the locally-set URI.
//   - Zone: Go field is "Zone" but the wire JSON tag is "dataCenter".
type DBaaS struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	// Request-side scalars.
	zone          *string
	engine        *string // wire: Engine.ID
	flavor        *string // wire: Flavor.Name
	storageGB     *int32  // wire: Storage.SizeGB
	autoscaling   *bool   // wire: Autoscaling.Enabled
	billingPeriod *string // wire: BillingPlan.BillingPeriod

	// Networking refs.
	vpcRef           *string
	subnetRef        *string
	securityGroupRef *string
	elasticIPRef     *string

	// Hydrated response.
	response *types.DBaaSResponse
}

// Setters (chainable).

func (d *DBaaS) IntoProject(p Ref) *DBaaS        { d.intoProject(p); return d }
func (d *DBaaS) WithName(n string) *DBaaS        { d.withName(n); return d }
func (d *DBaaS) AddTag(t string) *DBaaS          { d.addTag(t); return d }
func (d *DBaaS) RemoveTag(t string) *DBaaS       { d.removeTag(t); return d }
func (d *DBaaS) ReplaceTags(ts ...string) *DBaaS { d.replaceTags(ts...); return d }
func (d *DBaaS) WithLocation(loc string) *DBaaS  { d.withLocation(loc); return d }
func (d *DBaaS) InRegion(region string) *DBaaS   { d.withLocation(region); return d }

func (d *DBaaS) InZone(zone string) *DBaaS              { d.zone = &zone; return d }
func (d *DBaaS) WithEngine(engine string) *DBaaS        { d.engine = &engine; return d }
func (d *DBaaS) WithFlavor(flavor string) *DBaaS        { d.flavor = &flavor; return d }
func (d *DBaaS) WithStorage(gb int) *DBaaS              { v := int32(gb); d.storageGB = &v; return d }
func (d *DBaaS) WithAutoscaling(enabled bool) *DBaaS    { d.autoscaling = &enabled; return d }
func (d *DBaaS) WithBillingPeriod(period string) *DBaaS { d.billingPeriod = &period; return d }

// WithNetworking sets VPC, Subnet, SecurityGroup, and ElasticIP in a single call.
// Each Ref is validated independently; an empty URI records an error but does not
// prevent the other three from being set.
func (d *DBaaS) WithNetworking(vpc, subnet, sg, eip Ref) *DBaaS {
	d.setSingleRef("WithNetworking[vpc]", vpc, &d.vpcRef)
	d.setSingleRef("WithNetworking[subnet]", subnet, &d.subnetRef)
	d.setSingleRef("WithNetworking[sg]", sg, &d.securityGroupRef)
	d.setSingleRef("WithNetworking[eip]", eip, &d.elasticIPRef)
	return d
}

func (d *DBaaS) setSingleRef(label string, r Ref, dst **string) {
	uri := r.URI()
	if uri == "" {
		d.addErr(fmt.Errorf("%s: empty URI", label))
		return
	}
	*dst = &uri
}

// Ref + ID accessors.

func (d *DBaaS) URI() string     { return d.RespURI() }
func (d *DBaaS) DBaaSID() string { return d.ID() }

// Accessors.

func (d *DBaaS) Raw() *types.DBaaSResponse      { return d.response }
func (d *DBaaS) RawRequest() types.DBaaSRequest { return d.toRequest() }

func (d *DBaaS) Zone() string { return dbaasDerefString(d.zone) }

// Engine returns the engine identifier. On a hydrated response the value comes
// from Engine.Type; before hydration it returns what was passed to WithEngine.
func (d *DBaaS) Engine() string {
	if d.response != nil && d.response.Properties.Engine != nil && d.response.Properties.Engine.Type != nil {
		return *d.response.Properties.Engine.Type
	}
	return dbaasDerefString(d.engine)
}

// EngineRaw returns the full engine struct from the last response, or nil.
func (d *DBaaS) EngineRaw() *types.DBaaSEngineResponse {
	if d.response == nil {
		return nil
	}
	return d.response.Properties.Engine
}

// Flavor returns the flavor name. On a hydrated response the value comes from
// Flavor.Name; before hydration it returns what was passed to WithFlavor.
func (d *DBaaS) Flavor() string {
	if d.response != nil && d.response.Properties.Flavor != nil && d.response.Properties.Flavor.Name != nil {
		return *d.response.Properties.Flavor.Name
	}
	return dbaasDerefString(d.flavor)
}

// FlavorRaw returns the full flavor struct from the last response, or nil.
func (d *DBaaS) FlavorRaw() *types.DBaaSFlavorResponse {
	if d.response == nil {
		return nil
	}
	return d.response.Properties.Flavor
}

// Storage returns the storage size in GB. On a hydrated response the value comes
// from Storage.SizeGB; before hydration it returns what was passed to WithStorage.
func (d *DBaaS) Storage() int32 {
	if d.response != nil && d.response.Properties.Storage != nil && d.response.Properties.Storage.SizeGB != nil {
		return *d.response.Properties.Storage.SizeGB
	}
	if d.storageGB != nil {
		return *d.storageGB
	}
	return 0
}

// BillingPeriod returns the billing period. On a hydrated response the value comes
// from BillingPlan.BillingPeriod; before hydration it returns what was passed to
// WithBillingPeriod.
func (d *DBaaS) BillingPeriod() string {
	if d.response != nil && d.response.Properties.BillingPlan != nil && d.response.Properties.BillingPlan.BillingPeriod != nil {
		return *d.response.Properties.BillingPlan.BillingPeriod
	}
	return dbaasDerefString(d.billingPeriod)
}

// Autoscaling returns whether autoscaling is enabled. The response side carries no
// Enabled field (only Status/AvailableSpace/StepSize/RuleID), so this always reads
// from the locally-set value.
func (d *DBaaS) Autoscaling() bool {
	if d.autoscaling != nil {
		return *d.autoscaling
	}
	return false
}

func (d *DBaaS) VPC() string {
	return dbaasNetworkingURI(d.response, func(n *types.DBaaSNetworkingResponse) *types.ReferenceResource { return n.VPC }, d.vpcRef)
}

func (d *DBaaS) Subnet() string {
	return dbaasNetworkingURI(d.response, func(n *types.DBaaSNetworkingResponse) *types.ReferenceResource { return n.Subnet }, d.subnetRef)
}

func (d *DBaaS) SecurityGroup() string {
	return dbaasNetworkingURI(d.response, func(n *types.DBaaSNetworkingResponse) *types.ReferenceResource { return n.SecurityGroup }, d.securityGroupRef)
}

func (d *DBaaS) ElasticIP() string {
	return dbaasNetworkingURI(d.response, func(n *types.DBaaSNetworkingResponse) *types.ReferenceResource { return n.ElasticIP }, d.elasticIPRef)
}

// Wire conversions.

func (d *DBaaS) toRequest() types.DBaaSRequest {
	props := types.DBaaSPropertiesRequest{}
	if d.zone != nil {
		v := *d.zone
		props.Zone = &v
	}
	if d.engine != nil {
		v := *d.engine
		props.Engine = &types.DBaaSEngine{ID: &v}
	}
	if d.flavor != nil {
		v := *d.flavor
		props.Flavor = &types.DBaaSFlavor{Name: &v}
	}
	if d.storageGB != nil {
		v := *d.storageGB
		props.Storage = &types.DBaaSStorage{SizeGB: &v}
	}
	if d.autoscaling != nil {
		v := *d.autoscaling
		props.Autoscaling = &types.DBaaSAutoscaling{Enabled: &v}
	}
	if d.billingPeriod != nil {
		v := *d.billingPeriod
		props.BillingPlan = &types.DBaaSBillingPlan{BillingPeriod: &v}
	}
	if d.vpcRef != nil || d.subnetRef != nil || d.securityGroupRef != nil || d.elasticIPRef != nil {
		net := &types.DBaaSNetworking{}
		if d.vpcRef != nil {
			net.VPCURI = d.vpcRef
		}
		if d.subnetRef != nil {
			net.SubnetURI = d.subnetRef
		}
		if d.securityGroupRef != nil {
			net.SecurityGroupURI = d.securityGroupRef
		}
		if d.elasticIPRef != nil {
			net.ElasticIPURI = d.elasticIPRef
		}
		props.Networking = net
	}
	return types.DBaaSRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: d.toMetadata(),
			Location:                d.toLocation(),
		},
		Properties: props,
	}
}

func (d *DBaaS) fromResponse(resp *types.DBaaSResponse) {
	if resp == nil {
		return
	}
	d.response = resp
	d.setMeta(&resp.Metadata)
	d.withName(dbaasDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		d.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		d.withLocation(resp.Metadata.LocationResponse.Value)
	}
	d.setLinked(resp.Properties.LinkedResources)
	d.setStatus(&resp.Status)

	// Hydrate request-side fields from response for round-trip Update support.
	if resp.Properties.Engine != nil && resp.Properties.Engine.Type != nil {
		v := *resp.Properties.Engine.Type
		d.engine = &v
	}
	if resp.Properties.Flavor != nil && resp.Properties.Flavor.Name != nil {
		v := *resp.Properties.Flavor.Name
		d.flavor = &v
	}
	if resp.Properties.Storage != nil && resp.Properties.Storage.SizeGB != nil {
		v := *resp.Properties.Storage.SizeGB
		d.storageGB = &v
	}
	if resp.Properties.BillingPlan != nil && resp.Properties.BillingPlan.BillingPeriod != nil {
		v := *resp.Properties.BillingPlan.BillingPeriod
		d.billingPeriod = &v
	}
	if resp.Properties.Networking != nil {
		if resp.Properties.Networking.VPC != nil && resp.Properties.Networking.VPC.URI != "" {
			v := resp.Properties.Networking.VPC.URI
			d.vpcRef = &v
		}
		if resp.Properties.Networking.Subnet != nil && resp.Properties.Networking.Subnet.URI != "" {
			v := resp.Properties.Networking.Subnet.URI
			d.subnetRef = &v
		}
		if resp.Properties.Networking.SecurityGroup != nil && resp.Properties.Networking.SecurityGroup.URI != "" {
			v := resp.Properties.Networking.SecurityGroup.URI
			d.securityGroupRef = &v
		}
		if resp.Properties.Networking.ElasticIP != nil && resp.Properties.Networking.ElasticIP.URI != "" {
			v := resp.Properties.Networking.ElasticIP.URI
			d.elasticIPRef = &v
		}
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		d.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if d.projectID == "" && d.RespURI() != "" {
		ids := parseURIIDs(d.RespURI())
		d.projectID = ids["projects"]
	}
}

func dbaasDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func dbaasNetworkingURI(resp *types.DBaaSResponse, pick func(*types.DBaaSNetworkingResponse) *types.ReferenceResource, fallback *string) string {
	if resp != nil && resp.Properties.Networking != nil {
		if r := pick(resp.Properties.Networking); r != nil && r.URI != "" {
			return r.URI
		}
	}
	return dbaasDerefString(fallback)
}

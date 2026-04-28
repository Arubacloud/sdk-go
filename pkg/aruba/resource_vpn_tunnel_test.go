package aruba

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/testutil"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// --------------------------------------------------------------------------
// Compile-time Ref satisfaction
// --------------------------------------------------------------------------

var _ Ref = (*VPNTunnel)(nil)

// --------------------------------------------------------------------------
// Fluent setters
// --------------------------------------------------------------------------

func TestVPNTunnel_FluentSetters(t *testing.T) {
	proj := &Project{}
	proj.fromResponse(projectTestResponse("p1", "my-project", "/projects/p1"))

	tun := NewVPNTunnel().
		IntoProject(proj).
		WithName("my-tunnel").
		AddTag("vpn").
		AddTag("ipsec").
		AddTag("vpn"). // dedupe
		InRegion("ITBG-Bergamo").
		WithVPNType("ipsec").
		WithVPNClientProtocol("ikev2").
		WithBillingPeriod("monthly").
		WithPeerClientPublicIP("1.2.3.4")

	if tun.Name() != "my-tunnel" {
		t.Errorf("Name() = %q", tun.Name())
	}
	if tags := tun.Tags(); len(tags) != 2 || tags[0] != "vpn" || tags[1] != "ipsec" {
		t.Errorf("Tags() = %v", tags)
	}
	if tun.Region() != "ITBG-Bergamo" {
		t.Errorf("Region() = %q", tun.Region())
	}
	if tun.VPNType() != "ipsec" {
		t.Errorf("VPNType() = %q", tun.VPNType())
	}
	if tun.VPNClientProtocol() != "ikev2" {
		t.Errorf("VPNClientProtocol() = %q", tun.VPNClientProtocol())
	}
	if tun.BillingPeriod() != "monthly" {
		t.Errorf("BillingPeriod() = %q", tun.BillingPeriod())
	}
	if tun.PeerClientPublicIP() != "1.2.3.4" {
		t.Errorf("PeerClientPublicIP() = %q", tun.PeerClientPublicIP())
	}
	if tun.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", tun.ProjectID())
	}
	if tun.Err() != nil {
		t.Errorf("Err() = %v", tun.Err())
	}

	tun.RemoveTag("vpn")
	if tags := tun.Tags(); len(tags) != 1 || tags[0] != "ipsec" {
		t.Errorf("after RemoveTag Tags() = %v", tags)
	}

	tun.ReplaceTags("x", "y")
	if tags := tun.Tags(); len(tags) != 2 || tags[0] != "x" || tags[1] != "y" {
		t.Errorf("after ReplaceTags Tags() = %v", tags)
	}
}

// --------------------------------------------------------------------------
// IntoProject — typed Ref
// --------------------------------------------------------------------------

func TestVPNTunnel_IntoProject_TypedRef(t *testing.T) {
	proj := &Project{}
	proj.fromResponse(projectTestResponse("p1", "my-project", "/projects/p1"))

	tun := NewVPNTunnel().IntoProject(proj)

	if tun.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", tun.ProjectID())
	}
	if tun.Err() != nil {
		t.Errorf("Err() = %v", tun.Err())
	}
}

func TestVPNTunnel_IntoProject_URIRef(t *testing.T) {
	tun := NewVPNTunnel().IntoProject(URI("/projects/p"))

	if tun.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", tun.ProjectID())
	}
	if tun.Err() != nil {
		t.Errorf("Err() = %v", tun.Err())
	}
}

func TestVPNTunnel_IntoProject_BadRef(t *testing.T) {
	tun := NewVPNTunnel().IntoProject(URI("/garbage"))
	if tun.Err() == nil {
		t.Error("expected Err() != nil for unresolvable Ref, got nil")
	}
}

// --------------------------------------------------------------------------
// VPNIPConfig sub-builder
// --------------------------------------------------------------------------

func TestVPNIPConfig_FluentSetters(t *testing.T) {
	cfg := NewVPNIPConfig().
		WithVPC(URI("/projects/p/providers/Aruba.Network/vpcs/v")).
		WithPublicIP(URI("/projects/p/providers/Aruba.Network/elasticIps/eip-1")).
		WithSubnet("my-subnet", "10.0.0.0/24")

	if cfg.Err() != nil {
		t.Errorf("Err() = %v", cfg.Err())
	}
	built := cfg.build()
	if built == nil {
		t.Fatal("build() returned nil")
	}
	if built.VPC == nil || built.VPC.URI != "/projects/p/providers/Aruba.Network/vpcs/v" {
		t.Errorf("VPC.URI = %q", func() string {
			if built.VPC != nil {
				return built.VPC.URI
			}
			return "<nil>"
		}())
	}
	if built.PublicIP == nil || built.PublicIP.URI != "/projects/p/providers/Aruba.Network/elasticIps/eip-1" {
		t.Errorf("PublicIP.URI = %q", func() string {
			if built.PublicIP != nil {
				return built.PublicIP.URI
			}
			return "<nil>"
		}())
	}
	if built.Subnet == nil || built.Subnet.Name != "my-subnet" || built.Subnet.CIDR != "10.0.0.0/24" {
		t.Errorf("Subnet = %+v", built.Subnet)
	}
}

func TestVPNIPConfig_WithVPC_EmptyURI(t *testing.T) {
	cfg := NewVPNIPConfig().WithVPC(URI(""))
	if cfg.Err() == nil {
		t.Fatal("expected Err() != nil for empty URI VPC")
	}
	if !strings.Contains(cfg.Err().Error(), "empty URI") {
		t.Errorf("error = %q, expected 'empty URI'", cfg.Err().Error())
	}
	built := cfg.build()
	if built != nil && built.VPC != nil {
		t.Error("VPC must remain nil when URI was empty")
	}
}

func TestVPNIPConfig_WithPublicIP_EmptyURI(t *testing.T) {
	cfg := NewVPNIPConfig().WithPublicIP(URI(""))
	if cfg.Err() == nil {
		t.Fatal("expected Err() != nil for empty URI PublicIP")
	}
	built := cfg.build()
	if built != nil && built.PublicIP != nil {
		t.Error("PublicIP must remain nil when URI was empty")
	}
}

func TestVPNIPConfig_Build_SubnetOmittedWhenUnset(t *testing.T) {
	cfg := NewVPNIPConfig().WithVPC(URI("/projects/p/providers/Aruba.Network/vpcs/v"))
	built := cfg.build()
	if built.Subnet != nil {
		t.Error("Subnet should be nil when WithSubnet was not called")
	}
}

func TestVPNIPConfig_Build_NilReceiver(t *testing.T) {
	var cfg *VPNIPConfig
	if cfg.build() != nil {
		t.Error("nil receiver build() must return nil")
	}
}

// --------------------------------------------------------------------------
// VPNIKE sub-builder
// --------------------------------------------------------------------------

func TestVPNIKE_FluentSetters(t *testing.T) {
	ike := NewVPNIKE().
		WithLifetime(28800).
		WithEncryption(types.VPNEncryptionAES256).
		WithHash(types.VPNHashSHA256).
		WithDHGroup(types.VPNDHGroup14).
		WithDPDAction(types.VPNDPDActionRestart).
		WithDPDInterval(30).
		WithDPDTimeout(120)

	built := ike.build()
	if built == nil {
		t.Fatal("build() returned nil")
	}
	if built.Lifetime != 28800 {
		t.Errorf("Lifetime = %d", built.Lifetime)
	}
	if built.Encryption == nil || *built.Encryption != types.VPNEncryptionAES256 {
		t.Errorf("Encryption = %v", built.Encryption)
	}
	if built.Hash == nil || *built.Hash != types.VPNHashSHA256 {
		t.Errorf("Hash = %v", built.Hash)
	}
	if built.DHGroup == nil || *built.DHGroup != types.VPNDHGroup14 {
		t.Errorf("DHGroup = %v", built.DHGroup)
	}
	if built.DPDAction == nil || *built.DPDAction != types.VPNDPDActionRestart {
		t.Errorf("DPDAction = %v", built.DPDAction)
	}
	if built.DPDInterval != 30 {
		t.Errorf("DPDInterval = %d", built.DPDInterval)
	}
	if built.DPDTimeout != 120 {
		t.Errorf("DPDTimeout = %d", built.DPDTimeout)
	}
}

func TestVPNIKE_Build_PartialFields(t *testing.T) {
	ike := NewVPNIKE().WithLifetime(3600)
	built := ike.build()
	if built.Encryption != nil {
		t.Error("Encryption should be nil when not set")
	}
	if built.Hash != nil {
		t.Error("Hash should be nil when not set")
	}
	if built.DHGroup != nil {
		t.Error("DHGroup should be nil when not set")
	}
}

func TestVPNIKE_Build_NilReceiver(t *testing.T) {
	var k *VPNIKE
	if k.build() != nil {
		t.Error("nil receiver build() must return nil")
	}
}

// --------------------------------------------------------------------------
// VPNESP sub-builder
// --------------------------------------------------------------------------

func TestVPNESP_FluentSetters(t *testing.T) {
	esp := NewVPNESP().
		WithLifetime(3600).
		WithEncryption(types.VPNEncryptionAES128).
		WithHash(types.VPNHashSHA1).
		WithPFS(types.VPNPFSDHGroup14)

	built := esp.build()
	if built == nil {
		t.Fatal("build() returned nil")
	}
	if built.Lifetime != 3600 {
		t.Errorf("Lifetime = %d", built.Lifetime)
	}
	if built.Encryption == nil || *built.Encryption != types.VPNEncryptionAES128 {
		t.Errorf("Encryption = %v", built.Encryption)
	}
	if built.Hash == nil || *built.Hash != types.VPNHashSHA1 {
		t.Errorf("Hash = %v", built.Hash)
	}
	if built.PFS == nil || *built.PFS != types.VPNPFSDHGroup14 {
		t.Errorf("PFS = %v", built.PFS)
	}
}

func TestVPNESP_Build_PartialFields(t *testing.T) {
	esp := NewVPNESP()
	built := esp.build()
	if built.Encryption != nil || built.Hash != nil || built.PFS != nil {
		t.Error("unset pointer fields should be nil")
	}
	if built.Lifetime != 0 {
		t.Errorf("Lifetime = %d", built.Lifetime)
	}
}

func TestVPNESP_Build_NilReceiver(t *testing.T) {
	var e *VPNESP
	if e.build() != nil {
		t.Error("nil receiver build() must return nil")
	}
}

// --------------------------------------------------------------------------
// VPNPSK sub-builder
// --------------------------------------------------------------------------

func TestVPNPSK_FluentSetters(t *testing.T) {
	psk := NewVPNPSK().
		WithCloudSite("cloud-site-A").
		WithOnPremSite("on-prem-site-B").
		WithKey("s3cr3t")

	built := psk.build()
	if built == nil {
		t.Fatal("build() returned nil")
	}
	if built.CloudSite == nil || *built.CloudSite != "cloud-site-A" {
		t.Errorf("CloudSite = %v", built.CloudSite)
	}
	if built.OnPremSite == nil || *built.OnPremSite != "on-prem-site-B" {
		t.Errorf("OnPremSite = %v", built.OnPremSite)
	}
	if built.Secret == nil || *built.Secret != "s3cr3t" {
		t.Errorf("Secret = %v", built.Secret)
	}
}

func TestVPNPSK_Build_PartialFields(t *testing.T) {
	psk := NewVPNPSK()
	built := psk.build()
	if built.CloudSite != nil || built.OnPremSite != nil || built.Secret != nil {
		t.Error("unset pointer fields should be nil")
	}
}

func TestVPNPSK_Build_NilReceiver(t *testing.T) {
	var p *VPNPSK
	if p.build() != nil {
		t.Error("nil receiver build() must return nil")
	}
}

// --------------------------------------------------------------------------
// toRequest round-trip
// --------------------------------------------------------------------------

func TestVPNTunnel_ToRequestRoundTrip(t *testing.T) {
	tun := NewVPNTunnel().
		WithName("my-tunnel").
		AddTag("t1").
		InRegion("ITBG-Bergamo").
		WithVPNType("ipsec").
		WithVPNClientProtocol("ikev2").
		WithBillingPeriod("monthly").
		WithPeerClientPublicIP("1.2.3.4").
		WithIPConfig(
			NewVPNIPConfig().
				WithVPC(URI("/projects/p/providers/Aruba.Network/vpcs/v")).
				WithSubnet("subnet-1", "10.0.0.0/24"),
		).
		WithIKESettings(
			NewVPNIKE().
				WithLifetime(28800).
				WithEncryption(types.VPNEncryptionAES256).
				WithHash(types.VPNHashSHA256).
				WithDHGroup(types.VPNDHGroup14),
		).
		WithESPSettings(
			NewVPNESP().
				WithLifetime(3600).
				WithPFS(types.VPNPFSEnable),
		).
		WithPSKSettings(
			NewVPNPSK().
				WithCloudSite("cloud").
				WithKey("secret"),
		)

	req := tun.RawRequest()

	if req.Metadata.Name != "my-tunnel" {
		t.Errorf("Metadata.Name = %q", req.Metadata.Name)
	}
	if len(req.Metadata.Tags) != 1 {
		t.Errorf("Metadata.Tags = %v", req.Metadata.Tags)
	}
	if req.Metadata.Location.Value != "ITBG-Bergamo" {
		t.Errorf("Metadata.Location.Value = %q", req.Metadata.Location.Value)
	}
	if req.Properties.VPNType == nil || *req.Properties.VPNType != "ipsec" {
		t.Errorf("Properties.VPNType = %v", req.Properties.VPNType)
	}
	if req.Properties.VPNClientProtocol == nil || *req.Properties.VPNClientProtocol != "ikev2" {
		t.Errorf("Properties.VPNClientProtocol = %v", req.Properties.VPNClientProtocol)
	}
	if req.Properties.BillingPlan == nil || req.Properties.BillingPlan.BillingPeriod != "monthly" {
		t.Errorf("BillingPlan = %v", req.Properties.BillingPlan)
	}
	if cs := req.Properties.VPNClientSettings; cs == nil {
		t.Fatal("VPNClientSettings must be set")
	} else {
		if cs.PeerClientPublicIP == nil || *cs.PeerClientPublicIP != "1.2.3.4" {
			t.Errorf("PeerClientPublicIP = %v", cs.PeerClientPublicIP)
		}
		if cs.IKE == nil {
			t.Fatal("IKE must be set")
		}
		if cs.IKE.Encryption == nil || *cs.IKE.Encryption != types.VPNEncryptionAES256 {
			t.Errorf("IKE.Encryption = %v", cs.IKE.Encryption)
		}
		if cs.ESP == nil {
			t.Fatal("ESP must be set")
		}
		if cs.PSK == nil {
			t.Fatal("PSK must be set")
		}
	}
	if ip := req.Properties.IPConfigurations; ip == nil {
		t.Fatal("IPConfigurations must be set")
	} else {
		if ip.VPC == nil || ip.VPC.URI != "/projects/p/providers/Aruba.Network/vpcs/v" {
			t.Errorf("IPConfig.VPC.URI = %q", func() string {
				if ip.VPC != nil {
					return ip.VPC.URI
				}
				return "<nil>"
			}())
		}
		if ip.Subnet == nil || ip.Subnet.Name != "subnet-1" {
			t.Errorf("IPConfig.Subnet = %v", ip.Subnet)
		}
	}
}

func TestVPNTunnel_ToRequest_NoBillingPeriod_OmitsBillingPlan(t *testing.T) {
	tun := NewVPNTunnel().WithName("bare")
	req := tun.RawRequest()
	if req.Properties.BillingPlan != nil {
		t.Errorf("BillingPlan should be nil when not set, got %+v", req.Properties.BillingPlan)
	}
}

func TestVPNTunnel_ToRequest_NoVPNClientSettings_OmitsObject(t *testing.T) {
	tun := NewVPNTunnel().WithName("bare")
	req := tun.RawRequest()
	if req.Properties.VPNClientSettings != nil {
		t.Errorf("VPNClientSettings should be nil when IKE/ESP/PSK/PeerIP all unset")
	}
}

func TestVPNTunnel_ToRequest_PeerClientPublicIPOnly_EmitsClientSettings(t *testing.T) {
	tun := NewVPNTunnel().WithPeerClientPublicIP("5.6.7.8")
	req := tun.RawRequest()
	if req.Properties.VPNClientSettings == nil {
		t.Fatal("VPNClientSettings must be non-nil when PeerClientPublicIP is set")
	}
	if req.Properties.VPNClientSettings.IKE != nil {
		t.Error("IKE should be nil")
	}
	if req.Properties.VPNClientSettings.ESP != nil {
		t.Error("ESP should be nil")
	}
	if req.Properties.VPNClientSettings.PSK != nil {
		t.Error("PSK should be nil")
	}
}

// --------------------------------------------------------------------------
// Sub-builder error absorption
// --------------------------------------------------------------------------

func TestVPNTunnel_AbsorbsSubBuilderErrors(t *testing.T) {
	tun := NewVPNTunnel().
		WithIPConfig(NewVPNIPConfig().WithVPC(URI("")))

	if tun.Err() == nil {
		t.Fatal("tunnel.Err() must be non-nil when sub-builder has errors")
	}
}

// --------------------------------------------------------------------------
// fromResponse hydration
// --------------------------------------------------------------------------

func vpnTunnelTestResponse(id, name, uri, projectID string) *types.VPNTunnelResponse {
	state := "Active"
	vpnType := "ipsec"
	proto := "ikev2"
	loc := &types.LocationResponse{Value: "ITBG-Bergamo"}
	peerIP := "1.2.3.4"
	bp := "monthly"
	return &types.VPNTunnelResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:               &id,
			URI:              &uri,
			Name:             &name,
			Tags:             []string{"vpn-tag"},
			LocationResponse: loc,
			ProjectResponseMetadata: &types.ProjectResponseMetadata{
				ID: projectID,
			},
		},
		Properties: types.VPNTunnelPropertiesResponse{
			VPNType:           &vpnType,
			VPNClientProtocol: &proto,
			BillingPlan:       &types.BillingPeriodResource{BillingPeriod: bp},
			VPNClientSettings: &types.VPNClientSettings{
				PeerClientPublicIP: &peerIP,
			},
		},
		Status: types.ResourceStatus{
			State: &state,
		},
	}
}

func TestVPNTunnel_FromResponseHydration(t *testing.T) {
	tun := &VPNTunnel{}
	resp := vpnTunnelTestResponse("t-1", "my-tunnel",
		"/projects/p1/providers/Aruba.Network/vpnTunnels/t-1", "p1")
	tun.fromResponse(resp)

	if tun.ID() != "t-1" {
		t.Errorf("ID() = %q", tun.ID())
	}
	if tun.URI() != "/projects/p1/providers/Aruba.Network/vpnTunnels/t-1" {
		t.Errorf("URI() = %q", tun.URI())
	}
	if tun.VPNTunnelID() != "t-1" {
		t.Errorf("VPNTunnelID() = %q", tun.VPNTunnelID())
	}
	if tun.Name() != "my-tunnel" {
		t.Errorf("Name() = %q", tun.Name())
	}
	if tags := tun.Tags(); len(tags) != 1 || tags[0] != "vpn-tag" {
		t.Errorf("Tags() = %v", tags)
	}
	if tun.Region() != "ITBG-Bergamo" {
		t.Errorf("Region() = %q", tun.Region())
	}
	if tun.State() != "Active" {
		t.Errorf("State() = %q", tun.State())
	}
	if tun.VPNType() != "ipsec" {
		t.Errorf("VPNType() = %q", tun.VPNType())
	}
	if tun.VPNClientProtocol() != "ikev2" {
		t.Errorf("VPNClientProtocol() = %q", tun.VPNClientProtocol())
	}
	if tun.BillingPeriod() != "monthly" {
		t.Errorf("BillingPeriod() = %q", tun.BillingPeriod())
	}
	if tun.PeerClientPublicIP() != "1.2.3.4" {
		t.Errorf("PeerClientPublicIP() = %q", tun.PeerClientPublicIP())
	}
	if tun.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", tun.ProjectID())
	}
	if tun.Raw() != resp {
		t.Error("Raw() should return the hydrated response pointer")
	}
}

func TestVPNTunnel_FromResponsePartial(t *testing.T) {
	tun := &VPNTunnel{}
	tun.fromResponse(nil)
	if tun.ID() != "" || tun.URI() != "" || tun.Name() != "" {
		t.Error("fromResponse(nil) should be a no-op")
	}
	if tun.Raw() != nil {
		t.Error("Raw() should be nil before hydration")
	}

	tun2 := &VPNTunnel{}
	tun2.fromResponse(&types.VPNTunnelResponse{})
	if tun2.ID() != "" || tun2.URI() != "" || tun2.State() != "" {
		t.Error("empty response should yield zero accessor values")
	}
	if tun2.VPNType() != "" || tun2.BillingPeriod() != "" {
		t.Error("empty Properties should yield empty strings")
	}
}

func TestVPNTunnel_FromResponseURIBackfill(t *testing.T) {
	uri := "/projects/p2/providers/Aruba.Network/vpnTunnels/t-2"
	id := "t-2"
	name := "uri-tunnel"
	resp := &types.VPNTunnelResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:   &id,
			URI:  &uri,
			Name: &name,
			// ProjectResponseMetadata intentionally nil
		},
	}
	tun := &VPNTunnel{}
	tun.fromResponse(resp)

	if tun.ProjectID() != "p2" {
		t.Errorf("ProjectID() via URI fallback = %q", tun.ProjectID())
	}
}

// --------------------------------------------------------------------------
// Ref + ancestor ID satisfaction (runtime)
// --------------------------------------------------------------------------

func TestVPNTunnel_RefSatisfaction(t *testing.T) {
	tun := &VPNTunnel{}
	tun.fromResponse(vpnTunnelTestResponse("t-99", "n",
		"/projects/p99/providers/Aruba.Network/vpnTunnels/t-99", "p99"))

	// withVPNTunnelID typed path
	tid, ok := extractID(tun, func(r Ref) (string, bool) {
		if w, ok := r.(withVPNTunnelID); ok {
			return w.VPNTunnelID(), true
		}
		return "", false
	}, "vpn-tunnels")
	if !ok || tid != "t-99" {
		t.Errorf("extractID via withVPNTunnelID = (%q, %v)", tid, ok)
	}

	// withProjectID typed path
	projID, ok := extractID(tun, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || projID != "p99" {
		t.Errorf("extractID via withProjectID = (%q, %v)", projID, ok)
	}
}

// --------------------------------------------------------------------------
// vpnTunnelIDsFromRef helper
// --------------------------------------------------------------------------

func TestVPNTunnelIDsFromRef_TypedRef(t *testing.T) {
	tun := &VPNTunnel{}
	tun.fromResponse(vpnTunnelTestResponse("t-1", "n",
		"/projects/p/providers/Aruba.Network/vpnTunnels/t-1", "p"))
	pid, tid, err := vpnTunnelIDsFromRef(tun)
	if err != nil || pid != "p" || tid != "t-1" {
		t.Errorf("vpnTunnelIDsFromRef typed = (%q, %q, %v)", pid, tid, err)
	}
}

func TestVPNTunnelIDsFromRef_URIRef_CamelCase(t *testing.T) {
	ref := URI("/projects/p/providers/Aruba.Network/vpnTunnels/t-1")
	pid, tid, err := vpnTunnelIDsFromRef(ref)
	if err != nil || pid != "p" || tid != "t-1" {
		t.Errorf("vpnTunnelIDsFromRef camelCase = (%q, %q, %v)", pid, tid, err)
	}
}

func TestVPNTunnelIDsFromRef_URIRef_KebabCase(t *testing.T) {
	ref := URI("/projects/p/providers/Aruba.Network/vpn-tunnels/t-1")
	pid, tid, err := vpnTunnelIDsFromRef(ref)
	if err != nil || pid != "p" || tid != "t-1" {
		t.Errorf("vpnTunnelIDsFromRef kebab-case = (%q, %q, %v)", pid, tid, err)
	}
}

func TestVPNTunnelIDsFromRef_BadURI_MissingTunnel(t *testing.T) {
	_, _, err := vpnTunnelIDsFromRef(URI("/projects/p/providers/Aruba.Network"))
	if err == nil {
		t.Error("expected error for URI without tunnel segment")
	}
}

func TestVPNTunnelIDsFromRef_BadURI_MissingProject(t *testing.T) {
	_, _, err := vpnTunnelIDsFromRef(URI("/providers/Aruba.Network/vpnTunnels/t-1"))
	if err == nil {
		t.Error("expected error for URI without project segment")
	}
}

func TestVPNTunnelIDsFromRef_BadURI_MissingAll(t *testing.T) {
	_, _, err := vpnTunnelIDsFromRef(URI("/something/else"))
	if err == nil {
		t.Error("expected error for totally invalid URI")
	}
}

// --------------------------------------------------------------------------
// vpnTunnelsClientAdapter — CRUD integration tests
// --------------------------------------------------------------------------

func buildVPNTunnelTestAdapter(t *testing.T, handler http.HandlerFunc) *vpnTunnelsClientAdapter {
	t.Helper()
	server := testutil.NewMockServer(t, handler)
	return newVPNTunnelsClientAdapter(testutil.NewClient(t, server.URL))
}

const vpnTunnelSuccessBody = `{` +
	`"metadata":{` +
	`"id":"t-1","name":"my-tunnel",` +
	`"uri":"/projects/p/providers/Aruba.Network/vpnTunnels/t-1",` +
	`"project":{"id":"p"}` +
	`},` +
	`"properties":{` +
	`"vpnType":"ipsec","vpnClientProtocol":"ikev2"` +
	`},` +
	`"status":{"state":"Active"}}`

func TestVPNTunnelsClientAdapter_Create_Success(t *testing.T) {
	var gotBody types.VPNTunnelRequest
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, vpnTunnelSuccessBody)
	})

	tun := NewVPNTunnel().
		IntoProject(URI("/projects/p")).
		WithName("my-tunnel").
		InRegion("ITBG-Bergamo").
		WithVPNType("ipsec").
		WithVPNClientProtocol("ikev2").
		WithIKESettings(NewVPNIKE().WithEncryption(types.VPNEncryptionAES256))

	result, err := adapter.Create(context.Background(), tun)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if result.ID() != "t-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.Name() != "my-tunnel" {
		t.Errorf("Name() = %q", result.Name())
	}
	if result.StatusCode() != http.StatusCreated {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	if gotBody.Metadata.Name != "my-tunnel" {
		t.Errorf("request Name = %q", gotBody.Metadata.Name)
	}
	if gotBody.Metadata.Location.Value != "ITBG-Bergamo" {
		t.Errorf("request Location = %q", gotBody.Metadata.Location.Value)
	}
	if gotBody.Properties.VPNType == nil || *gotBody.Properties.VPNType != "ipsec" {
		t.Errorf("request VPNType = %v", gotBody.Properties.VPNType)
	}
	if gotBody.Properties.VPNClientSettings == nil || gotBody.Properties.VPNClientSettings.IKE == nil {
		t.Error("request IKE must be present")
	}
}

func TestVPNTunnelsClientAdapter_Create_NoProject(t *testing.T) {
	callCount := 0
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusCreated)
	})

	_, err := adapter.Create(context.Background(), NewVPNTunnel().WithName("x"))
	if err == nil {
		t.Fatal("expected error when tunnel has no project")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without project")
	}
}

func TestVPNTunnelsClientAdapter_Create_MetadataValidationError(t *testing.T) {
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Missing "id" field — triggers MetadataValidationError
		fmt.Fprint(w, `{"metadata":{"name":"tunnel","uri":"/projects/p/providers/Aruba.Network/vpnTunnels/x"},"properties":{},"status":{}}`)
	})

	tun := NewVPNTunnel().IntoProject(URI("/projects/p")).WithName("tunnel")
	result, err := adapter.Create(context.Background(), tun)
	if err == nil {
		t.Fatal("expected MetadataValidationError, got nil")
	}
	var mvErr *types.MetadataValidationError
	if !errors.As(err, &mvErr) {
		t.Fatalf("expected *types.MetadataValidationError, got %T: %v", err, err)
	}
	if result == nil {
		t.Fatal("result must be non-nil alongside MetadataValidationError")
	}
}

func TestVPNTunnelsClientAdapter_Create_NonTwoXX(t *testing.T) {
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Validation Failed", "name is required", 422))
	})

	tun := NewVPNTunnel().IntoProject(URI("/projects/p"))
	result, err := adapter.Create(context.Background(), tun)
	if err == nil {
		t.Fatal("expected error on 422")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
	if httpErr.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("HTTPError.StatusCode = %d", httpErr.StatusCode)
	}
	if result == nil {
		t.Fatal("result must be non-nil on non-2xx")
	}
}

func TestVPNTunnelsClientAdapter_Get_URIRef(t *testing.T) {
	var capturedPath string
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, vpnTunnelSuccessBody)
	})

	ref := URI("/projects/p/providers/Aruba.Network/vpnTunnels/t-1")
	result, err := adapter.Get(context.Background(), ref)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "t-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", result.ProjectID())
	}
	if result.StatusCode() != http.StatusOK {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	if !strings.Contains(capturedPath, "vpnTunnels") {
		t.Errorf("path = %q, expected vpnTunnels segment", capturedPath)
	}
}

func TestVPNTunnelsClientAdapter_Get_TypedRef(t *testing.T) {
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, vpnTunnelSuccessBody)
	})

	existing := &VPNTunnel{}
	existing.fromResponse(vpnTunnelTestResponse("t-1", "n",
		"/projects/p/providers/Aruba.Network/vpnTunnels/t-1", "p"))

	result, err := adapter.Get(context.Background(), existing)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "t-1" {
		t.Errorf("ID() = %q", result.ID())
	}
}

func TestVPNTunnelsClientAdapter_Update_Success(t *testing.T) {
	var capturedBody types.VPNTunnelRequest
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&capturedBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"metadata":{"id":"t-1","name":"renamed","uri":"/projects/p/providers/Aruba.Network/vpnTunnels/t-1","project":{"id":"p"}},"properties":{},"status":{}}`)
	})

	tun := &VPNTunnel{}
	tun.fromResponse(vpnTunnelTestResponse("t-1", "orig",
		"/projects/p/providers/Aruba.Network/vpnTunnels/t-1", "p"))
	tun.WithName("renamed")

	result, err := adapter.Update(context.Background(), tun)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if result.Name() != "renamed" {
		t.Errorf("Name() = %q", result.Name())
	}
	if capturedBody.Metadata.Name != "renamed" {
		t.Errorf("request Name = %q", capturedBody.Metadata.Name)
	}
}

func TestVPNTunnelsClientAdapter_Update_NoID(t *testing.T) {
	callCount := 0
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	tun := NewVPNTunnel().IntoProject(URI("/projects/p")).WithName("x")
	_, err := adapter.Update(context.Background(), tun)
	if err == nil {
		t.Fatal("expected error when tunnel has no ID")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made when ID is missing")
	}
}

func TestVPNTunnelsClientAdapter_Update_NoProject(t *testing.T) {
	callCount := 0
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	tun := &VPNTunnel{}
	id := "t-1"
	tun.fromResponse(&types.VPNTunnelResponse{
		Metadata: types.ResourceMetadataResponse{
			ID: &id,
		},
	})

	_, err := adapter.Update(context.Background(), tun)
	if err == nil {
		t.Fatal("expected error when tunnel has no project")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without project")
	}
}

func TestVPNTunnelsClientAdapter_Delete_Success(t *testing.T) {
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := adapter.Delete(context.Background(), URI("/projects/p/providers/Aruba.Network/vpnTunnels/t-1"))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestVPNTunnelsClientAdapter_Delete_NonTwoXX(t *testing.T) {
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "vpn tunnel not found", 404))
	})

	err := adapter.Delete(context.Background(), URI("/projects/p/providers/Aruba.Network/vpnTunnels/missing"))
	if err == nil {
		t.Fatal("expected error on 404")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T", err)
	}
	if httpErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d", httpErr.StatusCode)
	}
}

func TestVPNTunnelsClientAdapter_List_TwoItems(t *testing.T) {
	adapter := buildVPNTunnelTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"total":2,"self":"","prev":"","next":"","first":"","last":"","values":[`+
			`{"metadata":{"id":"t-1","name":"tunnel-a","uri":"/projects/p/providers/Aruba.Network/vpnTunnels/t-1","project":{"id":"p"}},"properties":{},"status":{"state":"Active"}},`+
			`{"metadata":{"id":"t-2","name":"tunnel-b","uri":"/projects/p/providers/Aruba.Network/vpnTunnels/t-2","project":{"id":"p"}},"properties":{},"status":{"state":"Inactive"}}`+
			`]}`)
	})

	list, err := adapter.List(context.Background(), URI("/projects/p"))
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if list.Total() != 2 {
		t.Errorf("Total() = %d", list.Total())
	}
	items := list.Items()
	if len(items) != 2 {
		t.Fatalf("Items() len = %d", len(items))
	}
	if items[0].ID() != "t-1" || items[0].Name() != "tunnel-a" {
		t.Errorf("items[0] = {%q, %q}", items[0].ID(), items[0].Name())
	}
	if items[1].ID() != "t-2" || items[1].State() != "Inactive" {
		t.Errorf("items[1] ID=%q State=%q", items[1].ID(), items[1].State())
	}
	if items[0].ProjectID() != "p" {
		t.Errorf("items[0].ProjectID() = %q", items[0].ProjectID())
	}
}

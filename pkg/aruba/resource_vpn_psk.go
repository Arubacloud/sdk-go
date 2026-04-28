package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// VPNPSK is a fluent builder for the PSKSettings block of a VPNTunnel.
// Construct with NewVPNPSK() and attach via VPNTunnel.WithPSKSettings.
type VPNPSK struct {
	errMixin
	cloudSite  *string
	onPremSite *string
	secret     *string
}

func (p *VPNPSK) WithCloudSite(v string) *VPNPSK  { p.cloudSite = &v; return p }
func (p *VPNPSK) WithOnPremSite(v string) *VPNPSK { p.onPremSite = &v; return p }
func (p *VPNPSK) WithKey(v string) *VPNPSK        { p.secret = &v; return p }

func (p *VPNPSK) build() *types.PSKSettings {
	if p == nil {
		return nil
	}
	return &types.PSKSettings{
		CloudSite:  p.cloudSite,
		OnPremSite: p.onPremSite,
		Secret:     p.secret,
	}
}

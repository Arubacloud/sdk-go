package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// VPNIKE is a fluent builder for the IKESettings block of a VPNTunnel.
// Construct with NewVPNIKE() and attach via VPNTunnel.WithIKESettings.
type VPNIKE struct {
	errMixin
	lifetime    int32
	encryption  *string
	hash        *string
	dhGroup     *string
	dpdAction   *string
	dpdInterval int32
	dpdTimeout  int32
}

func (k *VPNIKE) WithLifetime(s int32) *VPNIKE    { k.lifetime = s; return k }
func (k *VPNIKE) WithEncryption(v string) *VPNIKE { k.encryption = &v; return k }
func (k *VPNIKE) WithHash(v string) *VPNIKE       { k.hash = &v; return k }
func (k *VPNIKE) WithDHGroup(v string) *VPNIKE    { k.dhGroup = &v; return k }
func (k *VPNIKE) WithDPDAction(v string) *VPNIKE  { k.dpdAction = &v; return k }
func (k *VPNIKE) WithDPDInterval(s int32) *VPNIKE { k.dpdInterval = s; return k }
func (k *VPNIKE) WithDPDTimeout(s int32) *VPNIKE  { k.dpdTimeout = s; return k }

func (k *VPNIKE) build() *types.IKESettings {
	if k == nil {
		return nil
	}
	return &types.IKESettings{
		Lifetime:    k.lifetime,
		Encryption:  k.encryption,
		Hash:        k.hash,
		DHGroup:     k.dhGroup,
		DPDAction:   k.dpdAction,
		DPDInterval: k.dpdInterval,
		DPDTimeout:  k.dpdTimeout,
	}
}

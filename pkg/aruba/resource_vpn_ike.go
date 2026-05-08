package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// VPNIKE is a fluent builder for the IKESettings block of a VPNTunnel.
// Construct with NewVPNIKE() and attach via VPNTunnel.WithIKESettings.
type VPNIKE struct {
	errMixin
	lifetime    int32
	encryption  *IKEEncryption
	hash        *IKEHash
	dhGroup     *IKEDHGroup
	dpdAction   *IKEDPDAction
	dpdInterval int32
	dpdTimeout  int32
}

func (k *VPNIKE) WithLifetimeSeconds(s int) *VPNIKE      { k.lifetime = int32(s); return k }
func (k *VPNIKE) WithEncryption(v IKEEncryption) *VPNIKE { k.encryption = &v; return k }
func (k *VPNIKE) WithHash(v IKEHash) *VPNIKE             { k.hash = &v; return k }
func (k *VPNIKE) WithDHGroup(v IKEDHGroup) *VPNIKE       { k.dhGroup = &v; return k }
func (k *VPNIKE) WithDPDAction(v IKEDPDAction) *VPNIKE   { k.dpdAction = &v; return k }
func (k *VPNIKE) WithDPDIntervalSeconds(s int) *VPNIKE   { k.dpdInterval = int32(s); return k }
func (k *VPNIKE) WithDPDTimeoutSeconds(s int) *VPNIKE    { k.dpdTimeout = int32(s); return k }

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

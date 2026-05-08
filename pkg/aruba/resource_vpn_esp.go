package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// VPNESP is a fluent builder for the ESPSettings block of a VPNTunnel.
// Construct with NewVPNESP() and attach via VPNTunnel.WithESPSettings.
type VPNESP struct {
	errMixin
	lifetime   int32
	encryption *ESPEncryption
	hash       *ESPHash
	pfs        *ESPPFSGroup
}

func (e *VPNESP) WithLifetimeSeconds(s int) *VPNESP      { e.lifetime = int32(s); return e }
func (e *VPNESP) WithEncryption(v ESPEncryption) *VPNESP { e.encryption = &v; return e }
func (e *VPNESP) WithHash(v ESPHash) *VPNESP             { e.hash = &v; return e }
func (e *VPNESP) WithPFS(v ESPPFSGroup) *VPNESP          { e.pfs = &v; return e }

func (e *VPNESP) build() *types.ESPSettings {
	if e == nil {
		return nil
	}
	return &types.ESPSettings{
		Lifetime:   e.lifetime,
		Encryption: e.encryption,
		Hash:       e.hash,
		PFS:        e.pfs,
	}
}

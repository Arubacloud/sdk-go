package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// VPNESP is a fluent builder for the ESPSettings block of a VPNTunnel.
// Construct with NewVPNESP() and attach via VPNTunnel.WithESPSettings.
type VPNESP struct {
	errMixin
	lifetime   int32
	encryption *string
	hash       *string
	pfs        *string
}

func (e *VPNESP) WithLifetime(s int32) *VPNESP    { e.lifetime = s; return e }
func (e *VPNESP) WithEncryption(v string) *VPNESP { e.encryption = &v; return e }
func (e *VPNESP) WithHash(v string) *VPNESP       { e.hash = &v; return e }
func (e *VPNESP) WithPFS(v string) *VPNESP        { e.pfs = &v; return e }

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

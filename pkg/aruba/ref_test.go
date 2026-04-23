package aruba

import (
	"testing"
)

func TestURIFactory(t *testing.T) {
	r := URI("/projects/p/network/vpcs/v")
	if got := r.URI(); got != "/projects/p/network/vpcs/v" {
		t.Errorf("URI() = %q, want %q", got, "/projects/p/network/vpcs/v")
	}
	if got := r.ID(); got != "" {
		t.Errorf("ID() = %q, want empty string", got)
	}
}

func TestParseURIIDs(t *testing.T) {
	cases := []struct {
		uri  string
		want map[string]string
	}{
		{
			uri:  "/projects/p",
			want: map[string]string{"projects": "p"},
		},
		{
			uri:  "/projects/p/network/vpcs/v",
			want: map[string]string{"projects": "p", "vpcs": "v"},
		},
		{
			uri:  "/projects/p/network/vpcs/v/security-groups/s",
			want: map[string]string{"projects": "p", "vpcs": "v", "security-groups": "s"},
		},
		{
			uri:  "/projects/p/network/vpcs/v/security-groups/s/rules/r",
			want: map[string]string{"projects": "p", "vpcs": "v", "security-groups": "s", "rules": "r"},
		},
		{
			uri:  "/projects/p/network/vpcs/v/peerings/pr",
			want: map[string]string{"projects": "p", "vpcs": "v", "peerings": "pr"},
		},
		{
			uri:  "/projects/p/network/vpcs/v/peerings/pr/routes/r",
			want: map[string]string{"projects": "p", "vpcs": "v", "peerings": "pr", "routes": "r"},
		},
		{
			uri:  "/projects/p/network/vpn-tunnels/t",
			want: map[string]string{"projects": "p", "vpn-tunnels": "t"},
		},
		{
			uri:  "/projects/p/database/dbaas/d/databases/db",
			want: map[string]string{"projects": "p", "dbaas": "d", "databases": "db"},
		},
		{
			uri:  "/projects/p/storage/backups/b",
			want: map[string]string{"projects": "p", "backups": "b"},
		},
		{
			uri:  "/projects/p/security/kms/k",
			want: map[string]string{"projects": "p", "kms": "k"},
		},
		{
			uri:  "",
			want: map[string]string{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.uri, func(t *testing.T) {
			got := parseURIIDs(tc.uri)
			for k, v := range tc.want {
				if got[k] != v {
					t.Errorf("parseURIIDs(%q)[%q] = %q, want %q", tc.uri, k, got[k], v)
				}
			}
			for k := range got {
				if _, ok := tc.want[k]; !ok {
					t.Errorf("parseURIIDs(%q) unexpected key %q", tc.uri, k)
				}
			}
		})
	}
}

// missingSegmentReturnsEmpty verifies that a URI lacking an expected segment returns an empty entry.
func TestParseURIIDsMissingSegment(t *testing.T) {
	got := parseURIIDs("/projects/p/network/vpcs/v")
	if _, ok := got["security-groups"]; ok {
		t.Error("expected no security-groups key in map")
	}
	if got["vpcs"] != "v" {
		t.Errorf("vpcs = %q, want %q", got["vpcs"], "v")
	}
}

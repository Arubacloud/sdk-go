package aruba

import "strings"

// Ref is a cross-resource reference. Every typed wrapper satisfies Ref;
// the aruba.URI(string) factory produces an opaque ref backed by a raw URI.
type Ref interface {
	// URI returns the resource's absolute URI path (e.g. "/projects/p/network/vpcs/v").
	URI() string
	// ID returns the resource's ID segment, or "" for opaque URI-only refs.
	ID() string
}

// uriRef is an opaque Ref backed by a raw URI string.
type uriRef struct{ uri string }

func (r uriRef) URI() string { return r.uri }
func (r uriRef) ID() string  { return "" }

// namespaceSegments are URI path segments that are category prefixes, not resource-type/id pairs.
var namespaceSegments = map[string]bool{
	"network":   true,
	"database":  true,
	"storage":   true,
	"container": true,
	"security":  true,
	"compute":   true,
	"schedule":  true,
	"metrics":   true,
	"audit":     true,
}

// parseURIIDs splits a URI path into resource-type → id pairs, skipping namespace prefixes.
//
// Example: "/projects/p/network/vpcs/v/security-groups/s"
// → {"projects":"p","vpcs":"v","security-groups":"s"}
func parseURIIDs(uri string) map[string]string {
	result := make(map[string]string)
	parts := strings.Split(strings.TrimPrefix(uri, "/"), "/")
	i := 0
	for i < len(parts) {
		seg := parts[i]
		if seg == "" {
			i++
			continue
		}
		if namespaceSegments[seg] {
			i++
			continue
		}
		if i+1 < len(parts) && parts[i+1] != "" {
			result[seg] = parts[i+1]
			i += 2
		} else {
			i++
		}
	}
	return result
}

// extractID tries to read an ID from parent first via typed interface assertion, then via URI parsing.
// Returns the id and true on success; ("", false) when not found.
func extractID(parent Ref, typedKey func(Ref) (string, bool), uriSegment string) (string, bool) {
	if id, ok := typedKey(parent); ok && id != "" {
		return id, true
	}
	m := parseURIIDs(parent.URI())
	id := m[uriSegment]
	return id, id != ""
}

// Internal interface assertions used by scoped mixins to read ancestor IDs from typed parents.
// These are not exported; external packages interact only via the public Ref interface.

type withProjectID interface{ ProjectID() string }
type withVPCID interface{ VPCID() string }
type withSecurityGroupID interface{ SecurityGroupID() string }
type withDBaaSID interface{ DBaaSID() string }
type withDatabaseID interface{ DatabaseID() string }
type withVPCPeeringID interface{ VPCPeeringID() string }
type withVPNTunnelID interface{ VPNTunnelID() string }
type withBackupID interface{ BackupID() string }
type withKMSID interface{ KMSID() string }

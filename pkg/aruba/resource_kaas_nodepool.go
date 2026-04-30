package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

// NodePool is a fluent builder for a single KaaS node pool.
// Construct with NewNodePool() and attach via KaaS.AddNodePool.
//
// Schema note: Instance and Zone are plain strings in the Create/Update request
// (wire fields "instance" and "dataCenter"); the response side uses object types
// (InstanceResponse, DataCenterResponse). AddNodePool flattens the response
// representation back to strings so Update round-trips correctly.
type NodePool struct {
	errMixin
	name        *string
	nodes       *int32
	instance    *string // wire JSON: "instance"
	zone        *string // wire JSON: "dataCenter"
	minCount    *int32
	maxCount    *int32
	autoscaling *bool
}

func (n *NodePool) Named(name string) *NodePool          { n.name = &name; return n }
func (n *NodePool) OfInstance(instance string) *NodePool { n.instance = &instance; return n }
func (n *NodePool) InZone(zone string) *NodePool         { n.zone = &zone; return n }

func (n *NodePool) WithCount(count int32) *NodePool { n.nodes = &count; return n }

// WithAutoscaling enables autoscaling with the given min and max node counts.
func (n *NodePool) WithAutoscaling(min, max int32) *NodePool {
	t := true
	n.autoscaling = &t
	n.minCount = &min
	n.maxCount = &max
	return n
}

func (n *NodePool) build() types.NodePoolProperties {
	out := types.NodePoolProperties{}
	if n.name != nil {
		out.Name = *n.name
	}
	if n.nodes != nil {
		out.Nodes = *n.nodes
	}
	if n.instance != nil {
		out.Instance = *n.instance
	}
	if n.zone != nil {
		out.Zone = *n.zone
	}
	if n.minCount != nil {
		out.MinCount = n.minCount
	}
	if n.maxCount != nil {
		out.MaxCount = n.maxCount
	}
	if n.autoscaling != nil {
		out.Autoscaling = *n.autoscaling
	}
	return out
}

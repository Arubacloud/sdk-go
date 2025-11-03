package schema

// RuleDirection represents the direction of a security rule
type RuleDirection string

const (
	RuleDirectionIngress RuleDirection = "Ingress"
	RuleDirectionEgress  RuleDirection = "Egress"
)

// EndpointTypeDto represents the type of target endpoint
type EndpointTypeDto string

const (
	EndpointTypeIP            EndpointTypeDto = "Ip"
	EndpointTypeSecurityGroup EndpointTypeDto = "SecurityGroup"
)

// RuleTarget represents the target of the rule (source or destination according to the direction)
type RuleTarget struct {
	// Kind Type of the target. Admissible values: Ip, SecurityGroup
	Kind EndpointTypeDto `json:"kind,omitempty"`

	// Value of the target.
	// If kind = "Ip", the value must be a valid network address in CIDR notation (included 0.0.0.0/0)
	// If kind = "SecurityGroup", the value must be a valid uri of any security group within the same vpc
	Value string `json:"value,omitempty"`
}

// SecurityRuleProperties contains the properties of a security rule
type SecurityRulePropertiesRequest struct {
	// Direction of the rule. Admissible values: Ingress, Egress
	Direction RuleDirection `json:"direction,omitempty"`

	// Protocol Name of the protocol. Admissible values: ANY, TCP, UDP, ICMP
	Protocol string `json:"protocol,omitempty"`

	// Port can be set with different values, according to the protocol.
	// - ANY and ICMP must not have a port
	// - TCP and UDP can have:
	//   - a single numeric port. For instance "80", "443" etc.
	//   - a port range. For instance "80-100"
	//   - the "*" value indicating any ports
	Port string `json:"port,omitempty"`

	// Target The target of the rule (source or destination according to the direction)
	Target *RuleTarget `json:"target,omitempty"`
}

type SecurityRulePropertiesResponse struct {
	LinkedResources []LinkedResource `json:"linkedResources,omitempty"`

	// Direction of the rule. Admissible values: Ingress, Egress
	Direction RuleDirection `json:"direction,omitempty"`

	// Protocol Name of the protocol. Admissible values: ANY, TCP, UDP, ICMP
	Protocol string `json:"protocol,omitempty"`

	// Port can be set with different values, according to the protocol.
	// - ANY and ICMP must not have a port
	// - TCP and UDP can have:
	//   - a single numeric port. For instance "80", "443" etc.
	//   - a port range. For instance "80-100"
	//   - the "*" value indicating any ports
	Port string `json:"port,omitempty"`

	// Target The target of the rule (source or destination according to the direction)
	Target *RuleTarget `json:"target,omitempty"`
}

type SecurityRuleRequest struct {
	Metadata ResourceMetadataRequest `json:"metadata"`
	// Properties of the security rule (nullable object)
	Properties SecurityRulePropertiesRequest `json:"properties"`
}

type SecurityRuleResponse struct {
	Metadata   ResourceMetadataResponse       `json:"metadata"`
	Status     ResourceStatus                 `json:"status"`
	Properties SecurityRulePropertiesResponse `json:"properties"`
}

type SecurityRuleList struct {
	ListResponse
	Values []SecurityRuleResponse `json:"values"`
}

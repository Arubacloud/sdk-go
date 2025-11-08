package network

import "fmt"

// validateProject checks if project ID is not empty
func validateProject(project string) error {
	if project == "" {
		return fmt.Errorf("project cannot be empty")
	}
	return nil
}

// validateProjectAndResource checks if both project and resource ID are not empty
func validateProjectAndResource(project, resourceID, resourceType string) error {
	if project == "" {
		return fmt.Errorf("project cannot be empty")
	}
	if resourceID == "" {
		return fmt.Errorf("%s cannot be empty", resourceType)
	}
	return nil
}

// validateVPCResource checks project, VPC ID and resource ID
func validateVPCResource(project, vpcID, resourceID, resourceType string) error {
	if project == "" {
		return fmt.Errorf("project cannot be empty")
	}
	if vpcID == "" {
		return fmt.Errorf("VPC ID cannot be empty")
	}
	if resourceID == "" {
		return fmt.Errorf("%s cannot be empty", resourceType)
	}
	return nil
}

// validateSecurityGroupRule checks all IDs for security group rule operations
func validateSecurityGroupRule(project, vpcID, securityGroupID, securityGroupRuleID string) error {
	if project == "" {
		return fmt.Errorf("project cannot be empty")
	}
	if vpcID == "" {
		return fmt.Errorf("VPC ID cannot be empty")
	}
	if securityGroupID == "" {
		return fmt.Errorf("security group ID cannot be empty")
	}
	if securityGroupRuleID == "" {
		return fmt.Errorf("security group rule ID cannot be empty")
	}
	return nil
}

// validateVPCPeeringRoute checks all IDs for VPC peering route operations
func validateVPCPeeringRoute(project, vpcID, vpcPeeringID, vpcPeeringRouteID string) error {
	if project == "" {
		return fmt.Errorf("project cannot be empty")
	}
	if vpcID == "" {
		return fmt.Errorf("VPC ID cannot be empty")
	}
	if vpcPeeringID == "" {
		return fmt.Errorf("VPC peering ID cannot be empty")
	}
	if vpcPeeringRouteID == "" {
		return fmt.Errorf("VPC peering route ID cannot be empty")
	}
	return nil
}

// validateVPNRoute checks all IDs for VPN route operations
func validateVPNRoute(project, vpnTunnelID, vpnRouteID string) error {
	if project == "" {
		return fmt.Errorf("project cannot be empty")
	}
	if vpnTunnelID == "" {
		return fmt.Errorf("VPN tunnel ID cannot be empty")
	}
	if vpnRouteID == "" {
		return fmt.Errorf("VPN route ID cannot be empty")
	}
	return nil
}

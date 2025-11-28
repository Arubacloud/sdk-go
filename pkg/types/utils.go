package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ParseResponseBody reads and parses the HTTP response body into the Response struct.
// For 2xx responses, it unmarshals into Data field.
// For 4xx/5xx responses, it unmarshals into Error field.
// Always stores the raw body in RawBody field.
func ParseResponseBody[T any](httpResp *http.Response) (*Response[T], error) {
	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &Response[T]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      bodyBytes,
	}

	// Parse the response body based on status code
	if response.IsSuccess() && len(bodyBytes) > 0 {
		var data T
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(bodyBytes) > 0 {
		var errorResp ErrorResponse
		if err := json.Unmarshal(bodyBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// Validation helper functions

// ValidateProject checks if project ID is not empty
func ValidateProject(projectID string) error {
	if projectID == "" {
		return fmt.Errorf("project cannot be empty")
	}
	return nil
}

// ValidateProjectAndResource checks if both project and resource ID are not empty
func ValidateProjectAndResource(project, resourceID, resourceType string) error {
	if project == "" {
		return fmt.Errorf("project cannot be empty")
	}
	if resourceID == "" {
		return fmt.Errorf("%s cannot be empty", resourceType)
	}
	return nil
}

// ValidateDBaaSResource checks project, DBaaS ID and resource ID
func ValidateDBaaSResource(project, dbaasID, resourceID, resourceType string) error {
	if project == "" {
		return fmt.Errorf("project cannot be empty")
	}
	if dbaasID == "" {
		return fmt.Errorf("DBaaS ID cannot be empty")
	}
	if resourceID == "" {
		return fmt.Errorf("%s cannot be empty", resourceType)
	}
	return nil
}

// ValidateDatabaseGrant checks all IDs for grant operations
func ValidateDatabaseGrant(project, dbaasID, databaseID, grantID string) error {
	if project == "" {
		return fmt.Errorf("project cannot be empty")
	}
	if dbaasID == "" {
		return fmt.Errorf("DBaaS ID cannot be empty")
	}
	if databaseID == "" {
		return fmt.Errorf("database ID cannot be empty")
	}
	if grantID == "" {
		return fmt.Errorf("grant ID cannot be empty")
	}
	return nil
}

// ValidateVPCResource checks project, VPC ID and resource ID
func ValidateVPCResource(project, vpcID, resourceID, resourceType string) error {
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

// ValidateSecurityGroupRule checks all IDs for security group rule operations
func ValidateSecurityGroupRule(project, vpcID, securityGroupID, securityGroupRuleID string) error {
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

// ValidateVPCPeeringRoute checks all IDs for VPC peering route operations
func ValidateVPCPeeringRoute(project, vpcID, vpcPeeringID, vpcPeeringRouteID string) error {
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

// ValidateVPNRoute checks all IDs for VPN route operations
func ValidateVPNRoute(project, vpnTunnelID, vpnRouteID string) error {
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

func ValidateStorageRestore(projectID, backupID string, restoreID *string) error {
	if projectID == "" {
		return fmt.Errorf("project cannot be empty")
	}
	if backupID == "" {
		return fmt.Errorf("backup ID cannot be empty")
	}

	if restoreID == nil {
		return nil
	}

	if *restoreID == "" {
		return fmt.Errorf("restore ID cannot be empty")
	}
	return nil
}

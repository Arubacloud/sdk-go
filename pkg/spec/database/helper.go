package database

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

// validateDBaaSResource checks project, DBaaS ID and resource ID
func validateDBaaSResource(project, dbaasID, resourceID, resourceType string) error {
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

// validateDatabaseGrant checks all IDs for grant operations
func validateDatabaseGrant(project, dbaasID, databaseID, grantID string) error {
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

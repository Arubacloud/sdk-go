package compute

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

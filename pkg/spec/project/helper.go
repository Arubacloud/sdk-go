package project

import "fmt"

// validateProjectID checks if project ID is not empty
func validateProjectID(projectID string) error {
	if projectID == "" {
		return fmt.Errorf("project ID cannot be empty")
	}
	return nil
}

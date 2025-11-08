package metric

import "fmt"

// validateProject checks if project ID is not empty
func validateProject(project string) error {
	if project == "" {
		return fmt.Errorf("project cannot be empty")
	}
	return nil
}

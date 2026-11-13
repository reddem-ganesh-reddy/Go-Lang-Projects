package utilities

import (
	"fmt"
	"strings"
)

// ValidateDepartmentName checks the department name.
func ValidateDepartmentName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("department name is required")
	}

	if len(name) < 2 {
		return fmt.Errorf("department name must contain at least 2 characters")
	}

	if len(name) > 100 {
		return fmt.Errorf("department name cannot exceed 100 characters")
	}

	return nil
}
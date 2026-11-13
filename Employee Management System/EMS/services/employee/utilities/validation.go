
package utilities

import (
	"fmt"
	"strings"

	"EMS/services/employee/dto"
)

func ValidateEmployee(request dto.CreateEmployeeRequest) error {
	name := strings.TrimSpace(request.Name)
	email := strings.TrimSpace(request.Email)

	if name == "" {
		return fmt.Errorf("employee name is required")
	}

	if len(name) < 2 || len(name) > 100 {
		return fmt.Errorf("employee name must contain 2 to 100 characters")
	}

	if email == "" {
		return fmt.Errorf("employee email is required")
	}

	if request.DepartmentID == 0 {
		return fmt.Errorf("department ID is required")
	}

	return nil
}

func ValidateEmployeeUpdate(request dto.UpdateEmployeeRequest) error {
	name := strings.TrimSpace(request.Name)
	email := strings.TrimSpace(request.Email)

	if name == "" {
		return fmt.Errorf("employee name is required")
	}

	if len(name) < 2 || len(name) > 100 {
		return fmt.Errorf("employee name must contain 2 to 100 characters")
	}

	if email == "" {
		return fmt.Errorf("employee email is required")
	}

	if request.DepartmentID == 0 {
		return fmt.Errorf("department ID is required")
	}

	return nil
}

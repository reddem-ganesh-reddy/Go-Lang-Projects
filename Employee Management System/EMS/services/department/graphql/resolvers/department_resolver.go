package resolvers

import (
	"strconv"

	"EMS/services/department/dto"
	"EMS/services/department/service"
)

// UserRoleContextKey identifies the authenticated user's role in GraphQL context.
const UserRoleContextKey = "user_role"

// DepartmentResolver handles GraphQL department operations.
type DepartmentResolver struct {
	service service.DepartmentService
}

// NewDepartmentResolver creates a department resolver.
func NewDepartmentResolver(
	departmentService service.DepartmentService,
) *DepartmentResolver {
	return &DepartmentResolver{
		service: departmentService,
	}
}

// GetAll retrieves all departments.
func (r *DepartmentResolver) GetAll() ([]dto.DepartmentResponse, error) {
	return r.service.GetAll()
}

// GetByID retrieves a department by ID.
func (r *DepartmentResolver) GetByID(id string) (dto.DepartmentResponse, error) {
	departmentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || departmentID == 0 {
		return dto.DepartmentResponse{}, strconv.ErrSyntax
	}

	return r.service.GetByID(uint(departmentID))
}

// Create creates a department.
func (r *DepartmentResolver) Create(
	name string,
) (dto.DepartmentResponse, error) {
	request := dto.CreateDepartmentRequest{
		Name: name,
	}

	return r.service.Create(request)
}

// Update updates a department.
func (r *DepartmentResolver) Update(
	id string,
	name string,
) (dto.DepartmentResponse, error) {
	departmentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || departmentID == 0 {
		return dto.DepartmentResponse{}, strconv.ErrSyntax
	}

	request := dto.UpdateDepartmentRequest{
		Name: name,
	}

	return r.service.Update(uint(departmentID), request)
}

// Delete deletes a department.
func (r *DepartmentResolver) Delete(id string) (bool, error) {
	departmentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || departmentID == 0 {
		return false, strconv.ErrSyntax
	}

	err = r.service.Delete(uint(departmentID))
	if err != nil {
		return false, err
	}

	return true, nil
}


package service

import (
	"context"
	"fmt"
	"math"
	"strings"

	"EMS/services/employee/dto"
	"EMS/services/employee/model"
	"EMS/services/employee/repository"
	"EMS/services/employee/utilities"

	"gorm.io/gorm"
)

type employeeServiceImpl struct {
	repository       repository.EmployeeRepository
	departmentClient *utilities.DepartmentClient
}

func NewEmployeeService(
	repo repository.EmployeeRepository,
	departmentClient *utilities.DepartmentClient,
) EmployeeService {
	return &employeeServiceImpl{
		repository:       repo,
		departmentClient: departmentClient,
	}
}

func (s *employeeServiceImpl) Create(
	request dto.CreateEmployeeRequest,
) (dto.EmployeeResponse, error) {

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)

	if err := utilities.ValidateEmployee(request); err != nil {
		return dto.EmployeeResponse{}, err
	}

	// Check whether the department exists.
	_, err := s.departmentClient.GetByID(context.Background(), request.DepartmentID)
	if err != nil {
		return dto.EmployeeResponse{}, fmt.Errorf("invalid department: %w", err)
	}

	// Check for duplicate email.
	_, err = s.repository.GetByEmail(request.Email)
	if err == nil {
		return dto.EmployeeResponse{}, fmt.Errorf("employee email already exists")
	}

	if err != gorm.ErrRecordNotFound {
		return dto.EmployeeResponse{}, fmt.Errorf("failed to check employee email: %w", err)
	}

	employee := model.Employee{
		Name:         request.Name,
		Email:        request.Email,
		DepartmentID: request.DepartmentID,
	}

	if err := s.repository.Create(&employee); err != nil {
		return dto.EmployeeResponse{}, fmt.Errorf("failed to create employee: %w", err)
	}

	return toEmployeeResponse(&employee), nil
}

func (s *employeeServiceImpl) GetAll(
	query dto.EmployeeListQuery,
) (dto.EmployeeListResponse, error) {

	// Apply default pagination values.
	if query.Page == 0 {
		query.Page = 1
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	// Validate pagination.
	if query.Page < 1 {
		return dto.EmployeeListResponse{}, fmt.Errorf("page must be at least 1")
	}

	if query.Limit < 1 || query.Limit > 100 {
		return dto.EmployeeListResponse{}, fmt.Errorf("limit must be between 1 and 100")
	}

	// Apply default sorting.
	if query.SortBy == "" {
		query.SortBy = "id"
	}

	if query.SortOrder == "" {
		query.SortOrder = "asc"
	}

	// Validate sorting column.
	allowedSortColumns := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"created_at": true,
	}

	if !allowedSortColumns[query.SortBy] {
		return dto.EmployeeListResponse{}, fmt.Errorf("invalid sort_by value")
	}

	// Validate sorting order.
	query.SortOrder = strings.ToLower(query.SortOrder)

	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		return dto.EmployeeListResponse{}, fmt.Errorf("sort_order must be asc or desc")
	}

	query.Name = strings.TrimSpace(query.Name)
	query.Email = strings.TrimSpace(query.Email)

	employees, totalRecords, err := s.repository.GetAll(query)
	if err != nil {
		return dto.EmployeeListResponse{}, fmt.Errorf("failed to retrieve employees: %w", err)
	}

	employeeResponses := make([]dto.EmployeeResponse, 0, len(employees))

	for i := range employees {
		employeeResponses = append(
			employeeResponses,
			toEmployeeResponse(&employees[i]),
		)
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(query.Limit)))

	return dto.EmployeeListResponse{
		Employees: employeeResponses,
		Pagination: dto.PaginationResponse{
			TotalRecords: totalRecords,
			CurrentPage:  query.Page,
			PageSize:     query.Limit,
			TotalPages:   totalPages,
		},
	}, nil
}

func (s *employeeServiceImpl) GetByID(id uint) (dto.EmployeeResponse, error) {
	if id == 0 {
		return dto.EmployeeResponse{}, fmt.Errorf("employee ID is required")
	}

	employee, err := s.repository.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.EmployeeResponse{}, fmt.Errorf("employee not found")
		}

		return dto.EmployeeResponse{}, fmt.Errorf("failed to retrieve employee: %w", err)
	}

	return toEmployeeResponse(employee), nil
}

func (s *employeeServiceImpl) Update(
	id uint,
	request dto.UpdateEmployeeRequest,
) (dto.EmployeeResponse, error) {

	if id == 0 {
		return dto.EmployeeResponse{}, fmt.Errorf("employee ID is required")
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)

	if err := utilities.ValidateEmployeeUpdate(request); err != nil {
		return dto.EmployeeResponse{}, err
	}

	employee, err := s.repository.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.EmployeeResponse{}, fmt.Errorf("employee not found")
		}

		return dto.EmployeeResponse{}, fmt.Errorf("failed to retrieve employee: %w", err)
	}

	// Check whether the requested department exists.
	_, err = s.departmentClient.GetByID(context.Background(), request.DepartmentID)
	if err != nil {
		return dto.EmployeeResponse{}, fmt.Errorf("invalid department: %w", err)
	}

	// Check whether another employee uses this email.
	existingEmployee, err := s.repository.GetByEmail(request.Email)

	if err == nil && existingEmployee.ID != id {
		return dto.EmployeeResponse{}, fmt.Errorf("employee email already exists")
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return dto.EmployeeResponse{}, fmt.Errorf("failed to check employee email: %w", err)
	}

	employee.Name = request.Name
	employee.Email = request.Email
	employee.DepartmentID = request.DepartmentID

	if err := s.repository.Update(employee); err != nil {
		return dto.EmployeeResponse{}, fmt.Errorf("failed to update employee: %w", err)
	}

	return toEmployeeResponse(employee), nil
}

func (s *employeeServiceImpl) Delete(id uint) error {
	if id == 0 {
		return fmt.Errorf("employee ID is required")
	}

	err := s.repository.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("employee not found")
		}

		return fmt.Errorf("failed to delete employee: %w", err)
	}

	return nil
}

func toEmployeeResponse(employee *model.Employee) dto.EmployeeResponse {
	return dto.EmployeeResponse{
		ID:           employee.ID,
		Name:         employee.Name,
		Email:        employee.Email,
		DepartmentID: employee.DepartmentID,
		CreatedAt:    employee.CreatedAt,
		UpdatedAt:    employee.UpdatedAt,
	}
}

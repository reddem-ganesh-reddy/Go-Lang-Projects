package service

import (
	"fmt"
	"strings"

	"EMS/services/department/dto"
	"EMS/services/department/model"
	"EMS/services/department/repository"
	"EMS/services/department/utilities"

	"gorm.io/gorm"
)

// departmentService implements DepartmentService.
type departmentService struct {
	repository repository.DepartmentRepository
}

// NewDepartmentService creates a department service.
func NewDepartmentService(
	repo repository.DepartmentRepository,
) DepartmentService {
	return &departmentService{
		repository: repo,
	}
}

// Create creates a new department.
func (s *departmentService) Create(
	request dto.CreateDepartmentRequest,
) (dto.DepartmentResponse, error) {

	request.Name = strings.TrimSpace(request.Name)

	if err := utilities.ValidateDepartmentName(request.Name); err != nil {
		return dto.DepartmentResponse{}, err
	}

	department := model.Department{
		Name: request.Name,
	}

	if err := s.repository.Create(&department); err != nil {
		return dto.DepartmentResponse{}, fmt.Errorf(
			"create department: %w", err,
		)
	}

	return dto.DepartmentResponse{
		ID:   department.ID,
		Name: department.Name,
	}, nil
}

// GetAll retrieves all departments.
func (s *departmentService) GetAll() ([]dto.DepartmentResponse, error) {
	departments, err := s.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get departments: %w", err)
	}

	responses := make([]dto.DepartmentResponse, 0, len(departments))

	for _, department := range departments {
		responses = append(responses, dto.DepartmentResponse{
			ID:   department.ID,
			Name: department.Name,
		})
	}

	return responses, nil
}

// GetByID retrieves a department using its ID.
func (s *departmentService) GetByID(
	id uint,
) (dto.DepartmentResponse, error) {

	department, err := s.repository.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.DepartmentResponse{}, fmt.Errorf("department not found")
		}

		return dto.DepartmentResponse{}, fmt.Errorf(
			"get department: %w", err,
		)
	}

	return dto.DepartmentResponse{
		ID:   department.ID,
		Name: department.Name,
	}, nil
}

// Update updates an existing department.
func (s *departmentService) Update(
	id uint,
	request dto.UpdateDepartmentRequest,
) (dto.DepartmentResponse, error) {

	request.Name = strings.TrimSpace(request.Name)

	if err := utilities.ValidateDepartmentName(request.Name); err != nil {
		return dto.DepartmentResponse{}, err
	}

	department, err := s.repository.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.DepartmentResponse{}, fmt.Errorf("department not found")
		}

		return dto.DepartmentResponse{}, fmt.Errorf(
			"get department: %w", err,
		)
	}

	department.Name = request.Name

	if err := s.repository.Update(department); err != nil {
		return dto.DepartmentResponse{}, fmt.Errorf(
			"update department: %w", err,
		)
	}

	return dto.DepartmentResponse{
		ID:   department.ID,
		Name: department.Name,
	}, nil
}

// Delete deletes a department.
func (s *departmentService) Delete(id uint) error {
	_, err := s.repository.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("department not found")
		}

		return fmt.Errorf("get department: %w", err)
	}

	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("delete department: %w", err)
	}

	return nil
}
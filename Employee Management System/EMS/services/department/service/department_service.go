package service

import "EMS/services/department/dto"

// DepartmentService defines department business operations.
type DepartmentService interface {
	Create(request dto.CreateDepartmentRequest) (dto.DepartmentResponse, error)
	GetAll() ([]dto.DepartmentResponse, error)
	GetByID(id uint) (dto.DepartmentResponse, error)
	Update(id uint, request dto.UpdateDepartmentRequest) (dto.DepartmentResponse, error)
	Delete(id uint) error
}
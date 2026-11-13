package repository

import "EMS/services/department/model"

// DepartmentRepository defines department database operations.
type DepartmentRepository interface {
	Create(department *model.Department) error
	GetAll() ([]model.Department, error)
	GetByID(id uint) (*model.Department, error)
	Update(department *model.Department) error
	Delete(id uint) error
}
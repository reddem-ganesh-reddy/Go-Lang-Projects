package repository

import (
	"EMS/services/department/model"

	"gorm.io/gorm"
)

// departmentRepository implements DepartmentRepository.
type departmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository creates a department repository.
func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

// Create inserts a department into the database.
func (r *departmentRepository) Create(
	department *model.Department,
) error {
	return r.db.Create(department).Error
}

// GetAll retrieves all departments.
func (r *departmentRepository) GetAll() ([]model.Department, error) {
	var departments []model.Department

	err := r.db.Find(&departments).Error

	return departments, err
}

// GetByID retrieves a department using its ID.
func (r *departmentRepository) GetByID(
	id uint,
) (*model.Department, error) {
	var department model.Department

	err := r.db.First(&department, id).Error
	if err != nil {
		return nil, err
	}

	return &department, nil
}

// Update modifies an existing department.
func (r *departmentRepository) Update(
	department *model.Department,
) error {
	return r.db.Save(department).Error
}

// Delete removes a department using its ID.
func (r *departmentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Department{}, id).Error
}
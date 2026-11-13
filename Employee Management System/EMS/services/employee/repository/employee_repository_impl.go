
package repository

import (
	"EMS/services/employee/dto"
	"EMS/services/employee/model"

	"gorm.io/gorm"
)

type employeeRepositoryImpl struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepositoryImpl{
		db: db,
	}
}

func (r *employeeRepositoryImpl) Create(employee *model.Employee) error {
	return r.db.Create(employee).Error
}

func (r *employeeRepositoryImpl) GetAll(
	query dto.EmployeeListQuery,
) ([]model.Employee, int64, error) {

	var employees []model.Employee
	var totalRecords int64

	dbQuery := r.db.Model(&model.Employee{})

	// Filter by employee name.
	if query.Name != "" {
		dbQuery = dbQuery.Where("name ILIKE ?", "%"+query.Name+"%")
	}

	// Filter by employee email.
	if query.Email != "" {
		dbQuery = dbQuery.Where("email ILIKE ?", "%"+query.Email+"%")
	}

	// Filter by department ID.
	if query.DepartmentID > 0 {
		dbQuery = dbQuery.Where("department_id = ?", query.DepartmentID)
	}

	// Count matching records before pagination.
	if err := dbQuery.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Apply default pagination values.
	page := query.Page
	if page < 1 {
		page = 1
	}

	limit := query.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Allow only approved sorting columns.
	allowedSortColumns := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"created_at": true,
	}

	sortBy := query.SortBy
	if !allowedSortColumns[sortBy] {
		sortBy = "id"
	}

	sortOrder := query.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	offset := (page - 1) * limit

	// Apply sorting and pagination.
	err := dbQuery.
		Order(sortBy + " " + sortOrder).
		Limit(limit).
		Offset(offset).
		Find(&employees).Error

	if err != nil {
		return nil, 0, err
	}

	return employees, totalRecords, nil
}

func (r *employeeRepositoryImpl) GetByID(id uint) (*model.Employee, error) {
	var employee model.Employee

	err := r.db.First(&employee, id).Error
	if err != nil {
		return nil, err
	}

	return &employee, nil
}

func (r *employeeRepositoryImpl) Update(employee *model.Employee) error {
	return r.db.Save(employee).Error
}

func (r *employeeRepositoryImpl) Delete(id uint) error {
	result := r.db.Delete(&model.Employee{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *employeeRepositoryImpl) GetByEmail(email string) (*model.Employee, error) {
	var employee model.Employee

	err := r.db.Where("email = ?", email).First(&employee).Error
	if err != nil {
		return nil, err
	}

	return &employee, nil
}

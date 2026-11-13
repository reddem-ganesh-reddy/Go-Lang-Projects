
package repository

import (
	"EMS/services/employee/dto"
	"EMS/services/employee/model"
)

type EmployeeRepository interface {
	Create(employee *model.Employee) error

	GetAll(query dto.EmployeeListQuery) ([]model.Employee, int64, error)

	GetByID(id uint) (*model.Employee, error)

	Update(employee *model.Employee) error

	Delete(id uint) error

	GetByEmail(email string) (*model.Employee, error)
}

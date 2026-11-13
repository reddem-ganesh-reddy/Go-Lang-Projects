
package service

import "EMS/services/employee/dto"

type EmployeeService interface {
	Create(request dto.CreateEmployeeRequest) (dto.EmployeeResponse, error)

	GetAll(query dto.EmployeeListQuery) (dto.EmployeeListResponse, error)

	GetByID(id uint) (dto.EmployeeResponse, error)

	Update(id uint, request dto.UpdateEmployeeRequest) (dto.EmployeeResponse, error)

	Delete(id uint) error
}

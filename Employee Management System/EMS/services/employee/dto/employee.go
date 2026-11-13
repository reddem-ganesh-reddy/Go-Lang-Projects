
package dto

import "time"

type CreateEmployeeRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}

type UpdateEmployeeRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}

type EmployeeResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	DepartmentID uint      `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

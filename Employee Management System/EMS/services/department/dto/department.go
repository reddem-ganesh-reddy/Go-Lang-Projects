package dto

// CreateDepartmentRequest contains department creation details.
type CreateDepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateDepartmentRequest contains department update details.
type UpdateDepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}

// DepartmentResponse contains department information returned by the API.
type DepartmentResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
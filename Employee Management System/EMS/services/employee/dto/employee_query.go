
package dto

type EmployeeListQuery struct {
	Page         int    `form:"page" json:"page"`
	Limit        int    `form:"limit" json:"limit"`
	Name         string `form:"name" json:"name"`
	Email        string `form:"email" json:"email"`
	DepartmentID uint   `form:"department_id" json:"department_id"`
	SortBy       string `form:"sort_by" json:"sort_by"`
	SortOrder    string `form:"sort_order" json:"sort_order"`
}

type PaginationResponse struct {
	TotalRecords int64 `json:"total_records"`
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
}

type EmployeeListResponse struct {
	Employees  []EmployeeResponse `json:"employees"`
	Pagination PaginationResponse `json:"pagination"`
}

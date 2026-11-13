package main

import "time"

// SwaggerLoginRequest represents login credentials.
type SwaggerLoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin@123"`
}

// SwaggerLoginData represents login response data.
type SwaggerLoginData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Username    string `json:"username"`
	Role        string `json:"role"`
}

// SwaggerLoginResponse represents the login API response.
type SwaggerLoginResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    SwaggerLoginData `json:"data"`
}

// SwaggerEmployeeRequest represents employee creation and update data.
type SwaggerEmployeeRequest struct {
	Name         string `json:"name" example:"John"`
	Email        string `json:"email" example:"john@example.com"`
	DepartmentID uint   `json:"department_id" example:"1"`
}

// SwaggerEmployeeResponse represents an employee.
type SwaggerEmployeeResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	DepartmentID uint      `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SwaggerPaginationResponse represents pagination details.
type SwaggerPaginationResponse struct {
	TotalRecords int64 `json:"total_records"`
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
}

// SwaggerEmployeeListResponse represents the employee list.
type SwaggerEmployeeListResponse struct {
	Employees  []SwaggerEmployeeResponse `json:"employees"`
	Pagination SwaggerPaginationResponse `json:"pagination"`
}

// SwaggerEmployeeDeleteResponse represents employee deletion.
type SwaggerEmployeeDeleteResponse struct {
	Message string `json:"message"`
}

// SwaggerDepartmentRequest represents department creation and update data.
type SwaggerDepartmentRequest struct {
	Name string `json:"name" example:"IT"`
}

// SwaggerDepartment represents a department.
type SwaggerDepartment struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// SwaggerDepartmentResponse represents a single department API response.
type SwaggerDepartmentResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    SwaggerDepartment `json:"data"`
}

// SwaggerDepartmentListResponse represents the department list API response.
type SwaggerDepartmentListResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    []SwaggerDepartment `json:"data"`
}

// SwaggerDepartmentDeleteResponse represents department deletion.
type SwaggerDepartmentDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SwaggerErrorResponse represents the standard API error response.
type SwaggerErrorResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	RequestID     string `json:"request_id,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
}

// SwaggerGatewayErrorResponse represents a Gateway proxy error.
type SwaggerGatewayErrorResponse struct {
	Error string `json:"error"`
}

// swaggerAuthLogin godoc
// @Summary Login
// @Description Authenticate using username and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body SwaggerLoginRequest true "Login credentials"
// @Success 200 {object} SwaggerLoginResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 500 {object} SwaggerErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/auth/login [post]
func swaggerAuthLogin() {}

// swaggerGetEmployees godoc
// @Summary Get all employees
// @Description Retrieve employees with pagination, filtering and sorting.
// @Tags Employee
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Records per page (maximum 100)"
// @Param name query string false "Filter by employee name"
// @Param email query string false "Filter by employee email"
// @Param department_id query int false "Filter by department ID"
// @Param sort_by query string false "Sort field: id, name, email, created_at"
// @Param sort_order query string false "Sort order: asc or desc"
// @Success 200 {object} SwaggerEmployeeListResponse
// @Failure 400 {object} SwaggerGatewayErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/employees [get]
func swaggerGetEmployees() {}

// swaggerGetEmployee godoc
// @Summary Get employee by ID
// @Description Retrieve an employee using employee ID.
// @Tags Employee
// @Produce json
// @Security BearerAuth
// @Param id path int true "Employee ID"
// @Success 200 {object} SwaggerEmployeeResponse
// @Failure 400 {object} SwaggerGatewayErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerGatewayErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/employees/{id} [get]
func swaggerGetEmployee() {}

// swaggerCreateEmployee godoc
// @Summary Create employee
// @Description Create a new employee. Admin access required.
// @Tags Employee
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SwaggerEmployeeRequest true "Employee details"
// @Success 201 {object} SwaggerEmployeeResponse
// @Failure 400 {object} SwaggerGatewayErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 409 {object} SwaggerGatewayErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/employees [post]
func swaggerCreateEmployee() {}

// swaggerUpdateEmployee godoc
// @Summary Update employee
// @Description Update an existing employee. Admin access required.
// @Tags Employee
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Employee ID"
// @Param request body SwaggerEmployeeRequest true "Updated employee details"
// @Success 200 {object} SwaggerEmployeeResponse
// @Failure 400 {object} SwaggerGatewayErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerGatewayErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/employees/{id} [put]
func swaggerUpdateEmployee() {}

// swaggerDeleteEmployee godoc
// @Summary Delete employee
// @Description Delete an employee. Admin access required.
// @Tags Employee
// @Produce json
// @Security BearerAuth
// @Param id path int true "Employee ID"
// @Success 200 {object} SwaggerEmployeeDeleteResponse
// @Failure 400 {object} SwaggerGatewayErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerGatewayErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/employees/{id} [delete]
func swaggerDeleteEmployee() {}

// swaggerGetDepartments godoc
// @Summary Get all departments
// @Description Retrieve all departments.
// @Tags Department
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerDepartmentListResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/departments [get]
func swaggerGetDepartments() {}

// swaggerGetDepartment godoc
// @Summary Get department by ID
// @Description Retrieve a department using department ID.
// @Tags Department
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Success 200 {object} SwaggerDepartmentResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/departments/{id} [get]
func swaggerGetDepartment() {}

// swaggerCreateDepartment godoc
// @Summary Create department
// @Description Create a department. Admin access required.
// @Tags Department
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SwaggerDepartmentRequest true "Department details"
// @Success 201 {object} SwaggerDepartmentResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/departments [post]
func swaggerCreateDepartment() {}

// swaggerUpdateDepartment godoc
// @Summary Update department
// @Description Update a department. Admin access required.
// @Tags Department
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Param request body SwaggerDepartmentRequest true "Updated department details"
// @Success 200 {object} SwaggerDepartmentResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/departments/{id} [put]
func swaggerUpdateDepartment() {}

// swaggerDeleteDepartment godoc
// @Summary Delete department
// @Description Delete a department. Admin access required.
// @Tags Department
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Success 200 {object} SwaggerDepartmentDeleteResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 404 {object} SwaggerErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Router /api/departments/{id} [delete]
func swaggerDeleteDepartment() {}

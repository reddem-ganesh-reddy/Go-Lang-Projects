
package controller

import (
	"net/http"
	"strconv"
	"strings"

	"EMS/services/employee/dto"
	"EMS/services/employee/service"

	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	service service.EmployeeService
}

func NewEmployeeController(employeeService service.EmployeeService) *EmployeeController {
	return &EmployeeController{
		service: employeeService,
	}
}

// CreateEmployee godoc
// @Summary Create an employee
// @Description Creates a new employee after validating the department
// @Tags Employees
// @Accept json
// @Produce json
// @Param employee body dto.CreateEmployeeRequest true "Employee details"
// @Success 201 {object} dto.EmployeeResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/employees [post]
func (c *EmployeeController) CreateEmployee(ctx *gin.Context) {
	var request dto.CreateEmployeeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := c.service.Create(request)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, employee)
}

// GetAllEmployees godoc
// @Summary Get all employees
// @Description Gets employees with optional pagination, filtering and sorting
// @Tags Employees
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Records per page (default 10, maximum 100)"
// @Param name query string false "Filter by employee name"
// @Param email query string false "Filter by employee email"
// @Param department_id query int false "Filter by department ID"
// @Param sort_by query string false "Sort field: id, name, email, created_at"
// @Param sort_order query string false "Sort order: asc or desc"
// @Success 200 {object} dto.EmployeeListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/employees [get]
func (c *EmployeeController) GetAllEmployees(ctx *gin.Context) {
	var query dto.EmployeeListQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	result, err := c.service.GetAll(query)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetEmployeeByID godoc
// @Summary Get an employee by ID
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} dto.EmployeeResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/employees/{id} [get]
func (c *EmployeeController) GetEmployeeByID(ctx *gin.Context) {
	id, ok := parseEmployeeID(ctx)
	if !ok {
		return
	}

	employee, err := c.service.GetByID(id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

// UpdateEmployee godoc
// @Summary Update an employee
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param employee body dto.UpdateEmployeeRequest true "Updated employee details"
// @Success 200 {object} dto.EmployeeResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/employees/{id} [put]
func (c *EmployeeController) UpdateEmployee(ctx *gin.Context) {
	id, ok := parseEmployeeID(ctx)
	if !ok {
		return
	}

	var request dto.UpdateEmployeeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := c.service.Update(id, request)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

// DeleteEmployee godoc
// @Summary Delete an employee
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/employees/{id} [delete]
func (c *EmployeeController) DeleteEmployee(ctx *gin.Context) {
	id, ok := parseEmployeeID(ctx)
	if !ok {
		return
	}

	if err := c.service.Delete(id); err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "employee deleted successfully"})
}

func parseEmployeeID(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return 0, false
	}

	return uint(id), true
}

func (c *EmployeeController) handleError(ctx *gin.Context, err error) {
	message := err.Error()
	lowerMessage := strings.ToLower(message)

	switch {
	case strings.Contains(lowerMessage, "not found"):
		ctx.JSON(http.StatusNotFound, gin.H{"error": message})

	case strings.Contains(lowerMessage, "already exists"):
		ctx.JSON(http.StatusConflict, gin.H{"error": message})

	case strings.Contains(lowerMessage, "invalid"),
		strings.Contains(lowerMessage, "required"),
		strings.Contains(lowerMessage, "must be"),
		strings.Contains(lowerMessage, "between"),
		strings.Contains(lowerMessage, "sort_by"),
		strings.Contains(lowerMessage, "sort_order"):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": message})

	case strings.Contains(lowerMessage, "department service"),
		strings.Contains(lowerMessage, "invalid department"):
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": message})

	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": message})
	}
}

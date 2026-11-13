package controller

import (
	"net/http"
	"strconv"
	"strings"

	"EMS/services/department/dto"
	"EMS/services/department/service"
	apierrors "EMS/shared/errors"

	"github.com/gin-gonic/gin"
)

// DepartmentController handles department HTTP requests.
type DepartmentController struct {
	service service.DepartmentService
}

// NewDepartmentController creates a department controller.
func NewDepartmentController(
	departmentService service.DepartmentService,
) *DepartmentController {
	return &DepartmentController{
		service: departmentService,
	}
}

// Create creates a new department.
//
// @Summary Create department
// @Description Creates a new department.
// @Tags Departments
// @Accept json
// @Produce json
// @Param request body dto.CreateDepartmentRequest true "Department details"
// @Success 201 {object} map[string]interface{} "Department created"
// @Failure 400 {object} apierrors.ErrorResponse "Invalid request"
// @Failure 500 {object} apierrors.ErrorResponse "Internal server error"
// @Router /api/departments [post]
func (ctrl *DepartmentController) Create(c *gin.Context) {
	var request dto.CreateDepartmentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		apierrors.BadRequest(c, "Invalid department request")
		return
	}

	response, err := ctrl.service.Create(request)
	if err != nil {
		handleDepartmentError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Department created successfully",
		"data":    response,
	})
}

// GetAll retrieves all departments.
//
// @Summary Get all departments
// @Description Retrieves all departments.
// @Tags Departments
// @Produce json
// @Success 200 {object} map[string]interface{} "Departments retrieved"
// @Failure 500 {object} apierrors.ErrorResponse "Internal server error"
// @Router /api/departments [get]
func (ctrl *DepartmentController) GetAll(c *gin.Context) {
	response, err := ctrl.service.GetAll()
	if err != nil {
		handleDepartmentError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Departments retrieved successfully",
		"data":    response,
	})
}

// GetByID retrieves a department using its ID.
//
// @Summary Get department by ID
// @Description Retrieves a department using its department ID.
// @Tags Departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} map[string]interface{} "Department retrieved"
// @Failure 400 {object} apierrors.ErrorResponse "Invalid department ID"
// @Failure 404 {object} apierrors.ErrorResponse "Department not found"
// @Failure 500 {object} apierrors.ErrorResponse "Internal server error"
// @Router /api/departments/{id} [get]
func (ctrl *DepartmentController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		apierrors.BadRequest(c, "Invalid department ID")
		return
	}

	response, err := ctrl.service.GetByID(uint(id))
	if err != nil {
		handleDepartmentError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Department retrieved successfully",
		"data":    response,
	})
}

// Update updates an existing department.
//
// @Summary Update department
// @Description Updates a department using its department ID.
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param request body dto.UpdateDepartmentRequest true "Updated department details"
// @Success 200 {object} map[string]interface{} "Department updated"
// @Failure 400 {object} apierrors.ErrorResponse "Invalid request"
// @Failure 404 {object} apierrors.ErrorResponse "Department not found"
// @Failure 500 {object} apierrors.ErrorResponse "Internal server error"
// @Router /api/departments/{id} [put]
func (ctrl *DepartmentController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		apierrors.BadRequest(c, "Invalid department ID")
		return
	}

	var request dto.UpdateDepartmentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		apierrors.BadRequest(c, "Invalid department request")
		return
	}

	response, err := ctrl.service.Update(uint(id), request)
	if err != nil {
		handleDepartmentError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Department updated successfully",
		"data":    response,
	})
}

// Delete deletes a department.
//
// @Summary Delete department
// @Description Deletes a department using its department ID.
// @Tags Departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} map[string]interface{} "Department deleted"
// @Failure 400 {object} apierrors.ErrorResponse "Invalid department ID"
// @Failure 404 {object} apierrors.ErrorResponse "Department not found"
// @Failure 500 {object} apierrors.ErrorResponse "Internal server error"
// @Router /api/departments/{id} [delete]
func (ctrl *DepartmentController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		apierrors.BadRequest(c, "Invalid department ID")
		return
	}

	if err := ctrl.service.Delete(uint(id)); err != nil {
		handleDepartmentError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Department deleted successfully",
	})
}

// handleDepartmentError handles department service errors.
func handleDepartmentError(c *gin.Context, err error) {
	message := err.Error()

	switch {
	case strings.Contains(message, "not found"):
		apierrors.NotFound(c, message)

	case strings.Contains(message, "required"),
		strings.Contains(message, "characters"),
		strings.Contains(message, "exceed"):
		apierrors.BadRequest(c, message)

	default:
		apierrors.InternalServerError(c, "Department operation failed")
	}
}
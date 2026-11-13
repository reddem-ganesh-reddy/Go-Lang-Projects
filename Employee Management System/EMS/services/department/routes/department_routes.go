
package routes

import (
	"EMS/services/department/controller"

	"github.com/gin-gonic/gin"
)

// RegisterDepartmentRoutes registers public Department REST endpoints.
// Authentication and authorization are handled by the API Gateway.
func RegisterDepartmentRoutes(
	router *gin.Engine,
	departmentController *controller.DepartmentController,
	jwtSecret string,
) {
	departments := router.Group("/api/departments")
	{
		departments.POST("", departmentController.Create)
		departments.GET("", departmentController.GetAll)
		departments.GET("/:id", departmentController.GetByID)
		departments.PUT("/:id", departmentController.Update)
		departments.DELETE("/:id", departmentController.Delete)
	}
}

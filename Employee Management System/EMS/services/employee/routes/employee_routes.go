
package routes

import (
	"EMS/services/employee/controller"

	"github.com/gin-gonic/gin"
)

func RegisterEmployeeRoutes(
	router *gin.Engine,
	employeeController *controller.EmployeeController,
) {
	employees := router.Group("/api/employees")
	{
		employees.POST("", employeeController.CreateEmployee)

		employees.GET("", employeeController.GetAllEmployees)

		employees.GET("/:id", employeeController.GetEmployeeByID)

		employees.PUT("/:id", employeeController.UpdateEmployee)

		employees.DELETE("/:id", employeeController.DeleteEmployee)
	}
}

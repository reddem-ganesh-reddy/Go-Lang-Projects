package routes

import (
	"EMS/services/auth/controller"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers the public authentication endpoints.
func RegisterAuthRoutes(router *gin.Engine, authController *controller.AuthController) {
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
	}
}

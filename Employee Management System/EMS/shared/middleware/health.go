package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse contains the health status.
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthCheck checks whether the service is running.
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{
			Status: "UP",
		})
	}
}
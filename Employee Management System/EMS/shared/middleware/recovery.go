package middleware

import (
	"net/http"

	"EMS/shared/logger"

	"github.com/gin-gonic/gin"
)

// Recovery catches unexpected application panics.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log := logger.Get()

		log.Error(
			"Unexpected panic recovered",
			"request_id", GetRequestID(c),
			"correlation_id", GetCorrelationID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"error", recovered,
		)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":          "internal server error",
			"request_id":     GetRequestID(c),
			"correlation_id": GetCorrelationID(c),
		})
	})
}
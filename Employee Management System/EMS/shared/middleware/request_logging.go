
package middleware

import (
	"time"

	"EMS/shared/logger"

	"github.com/gin-gonic/gin"
)

// RequestLogging records details about business HTTP requests.
func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Skip automatic health checks.
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)

		log := logger.Get()

		log.Info(
			"HTTP request completed",
			"request_id", GetRequestID(c),
			"correlation_id", GetCorrelationID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse contains the standard API error response.
type ErrorResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	RequestID      string `json:"request_id,omitempty"`
	CorrelationID  string `json:"correlation_id,omitempty"`
}

// Send sends a standard error response.
func Send(
	c *gin.Context,
	statusCode int,
	message string,
) {
	c.JSON(statusCode, ErrorResponse{
		Success:       false,
		Message:       message,
		RequestID:     c.GetHeader("X-Request-ID"),
		CorrelationID: c.GetHeader("X-Correlation-ID"),
	})
}

// BadRequest sends HTTP 400.
func BadRequest(c *gin.Context, message string) {
	Send(c, http.StatusBadRequest, message)
}

// Unauthorized sends HTTP 401.
func Unauthorized(c *gin.Context, message string) {
	Send(c, http.StatusUnauthorized, message)
}

// Forbidden sends HTTP 403.
func Forbidden(c *gin.Context, message string) {
	Send(c, http.StatusForbidden, message)
}

// NotFound sends HTTP 404.
func NotFound(c *gin.Context, message string) {
	Send(c, http.StatusNotFound, message)
}

// Conflict sends HTTP 409.
func Conflict(c *gin.Context, message string) {
	Send(c, http.StatusConflict, message)
}

// InternalServerError sends HTTP 500.
func InternalServerError(c *gin.Context, message string) {
	Send(c, http.StatusInternalServerError, message)
}

// ServiceUnavailable sends HTTP 503.
func ServiceUnavailable(c *gin.Context, message string) {
	Send(c, http.StatusServiceUnavailable, message)
}
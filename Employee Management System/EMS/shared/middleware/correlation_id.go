package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const correlationIDHeader = "X-Correlation-ID"

type correlationIDContextKey struct{}

// CorrelationID preserves or generates a correlation ID.
func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := strings.TrimSpace(
			c.GetHeader(correlationIDHeader),
		)

		if correlationID == "" || len(correlationID) > 128 {
			correlationID = newCorrelationID()
		}

		c.Header(correlationIDHeader, correlationID)
		c.Request.Header.Set(correlationIDHeader, correlationID)

		c.Set(correlationIDHeader, correlationID)

		ctx := context.WithValue(
			c.Request.Context(),
			correlationIDContextKey{},
			correlationID,
		)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetCorrelationID retrieves the ID from a Gin request.
func GetCorrelationID(c *gin.Context) string {
	value, exists := c.Get(correlationIDHeader)
	if !exists {
		return ""
	}

	id, ok := value.(string)
	if !ok {
		return ""
	}

	return id
}

// FromContext retrieves the ID from a standard Go context.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(correlationIDContextKey{}).(string)
	return id
}

// newCorrelationID generates a unique correlation ID.
func newCorrelationID() string {
	var value [16]byte

	if _, err := rand.Read(value[:]); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}

	encoded := hex.EncodeToString(value[:])

	return encoded[:8] + "-" +
		encoded[8:12] + "-" +
		encoded[12:16] + "-" +
		encoded[16:20] + "-" +
		encoded[20:]
}
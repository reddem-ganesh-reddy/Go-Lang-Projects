package middleware

import (
	"net/http"
	"strings"

	apierrors "EMS/shared/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth validates JWT tokens and stores user details in the request context.
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the Authorization header.
		authHeader := c.GetHeader("Authorization")

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			apierrors.Unauthorized(c, "missing or invalid authorization token")
			c.Abort()
			return
		}

		// Parse and validate the JWT token.
		token, err := jwt.Parse(
			parts[1],
			func(token *jwt.Token) (interface{}, error) {
				// Accept only the expected signing algorithm.
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(secret), nil
			},
		)

		if err != nil || !token.Valid {
			apierrors.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// Extract username and role from the token claims.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			apierrors.Unauthorized(c, "invalid token claims")
			c.Abort()
			return
		}

		username, usernameOK := claims["username"].(string)
		role, roleOK := claims["role"].(string)

		if !usernameOK || !roleOK || username == "" || role == "" {
			apierrors.Unauthorized(c, "invalid token claims")
			c.Abort()
			return
		}

		if role != "admin" && role != "employee" {
			apierrors.Unauthorized(c, "invalid user role")
			c.Abort()
			return
		}

		// Make authenticated user details available to handlers.
		c.Set("username", username)
		c.Set("role", role)

		c.Next()
	}
}

// RequireRole allows only users with one of the specified roles.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			apierrors.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			apierrors.Unauthorized(c, "invalid user role")
			c.Abort()
			return
		}

		// Check whether the user's role is permitted.
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		apierrors.Forbidden(c, "you do not have permission to access this resource")
		c.Abort()
	}
}

// GetUsername returns the authenticated username from the Gin context.
func GetUsername(c *gin.Context) string {
	username, _ := c.Get("username")
	value, _ := username.(string)
	return value
}

// GetUserRole returns the authenticated user's role from the Gin context.
func GetUserRole(c *gin.Context) string {
	role, _ := c.Get("role")
	value, _ := role.(string)
	return value
}

// Unauthorized is retained as a helper for HTTP unauthorized responses.
func Unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": "unauthorized",
	})
}

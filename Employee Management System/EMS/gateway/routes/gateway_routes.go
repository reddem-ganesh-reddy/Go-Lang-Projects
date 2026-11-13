package routes

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"EMS/shared/discovery"
	"EMS/shared/httpclient"
	"EMS/shared/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterGatewayRoutes registers all public API Gateway routes.
func RegisterGatewayRoutes(
	router *gin.Engine,
	discoveryClient *discovery.Client,
	httpClient *httpclient.Client,
	jwtSecret string,
) {
	// Public authentication endpoint.
	authProxy := createProxyHandler(
		discoveryClient,
		httpClient,
		"auth-service",
	)

	router.POST("/api/auth/login", authProxy)

	// Protected Gateway GraphQL endpoint.
	router.POST(
		"/graphql",
		middleware.JWTAuth(jwtSecret),
		GraphQLProxyHandler(discoveryClient, httpClient),
	)

	// Employee read operations.
	employeeGroup := router.Group(
		"/api/employees",
		middleware.JWTAuth(jwtSecret),
	)

	employeeProxy := createProxyHandler(
		discoveryClient,
		httpClient,
		"employee-service",
	)

	employeeGroup.GET("", employeeProxy)
	employeeGroup.GET("/*path", employeeProxy)

	// Employee write operations: admin only.
	adminEmployeeGroup := router.Group(
		"/api/employees",
		middleware.JWTAuth(jwtSecret),
		middleware.RequireRole("admin"),
	)

	adminEmployeeGroup.POST("", employeeProxy)
	adminEmployeeGroup.PUT("/*path", employeeProxy)
	adminEmployeeGroup.DELETE("/*path", employeeProxy)

	// Department read operations.
	departmentGroup := router.Group(
		"/api/departments",
		middleware.JWTAuth(jwtSecret),
	)

	departmentProxy := createProxyHandler(
		discoveryClient,
		httpClient,
		"department-service",
	)

	departmentGroup.GET("", departmentProxy)
	departmentGroup.GET("/*path", departmentProxy)

	// Department write operations: admin only.
	adminDepartmentGroup := router.Group(
		"/api/departments",
		middleware.JWTAuth(jwtSecret),
		middleware.RequireRole("admin"),
	)

	adminDepartmentGroup.POST("", departmentProxy)
	adminDepartmentGroup.PUT("/*path", departmentProxy)
	adminDepartmentGroup.DELETE("/*path", departmentProxy)
}

// createProxyHandler forwards requests to the discovered service.
func createProxyHandler(
	discoveryClient *discovery.Client,
	httpClient *httpclient.Client,
	serviceName string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceURL, err := discoveryClient.DiscoverService(serviceName)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Service is currently unavailable",
			})
			return
		}

		requestURL := strings.TrimRight(serviceURL, "/") +
			c.Request.URL.RequestURI()

		var requestBody []byte

		if c.Request.Body != nil {
			requestBody, err = io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Unable to read request body",
				})
				return
			}
		}

		headers := make(map[string]string)

		if contentType := c.GetHeader("Content-Type"); contentType != "" {
			headers["Content-Type"] = contentType
		}

		if accept := c.GetHeader("Accept"); accept != "" {
			headers["Accept"] = accept
		}

		if authorization := c.GetHeader("Authorization"); authorization != "" {
			headers["Authorization"] = authorization
		}

		response, err := httpClient.Do(
			c.Request.Context(),
			c.Request.Method,
			requestURL,
			requestBody,
			headers,
		)

		if err != nil {
			statusCode := http.StatusServiceUnavailable

			if errors.Is(err, context.DeadlineExceeded) {
				statusCode = http.StatusGatewayTimeout
			} else {
				var networkError net.Error
				if errors.As(err, &networkError) && networkError.Timeout() {
					statusCode = http.StatusGatewayTimeout
				}
			}

			c.JSON(statusCode, gin.H{
				"error": "Unable to communicate with the service",
			})
			return
		}

		defer response.Body.Close()

		for key, values := range response.Header {
			switch strings.ToLower(key) {
			case "connection",
				"keep-alive",
				"proxy-authenticate",
				"proxy-authorization",
				"te",
				"trailer",
				"transfer-encoding",
				"upgrade":
				continue
			}

			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		c.Status(response.StatusCode)
		_, _ = io.Copy(c.Writer, response.Body)
	}
}

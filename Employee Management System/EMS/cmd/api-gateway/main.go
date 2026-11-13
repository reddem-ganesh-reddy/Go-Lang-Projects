// @title Employee Management System API Gateway
// @version 1.0
// @description Central API Gateway for Authentication, Employee and Department services.
// @host localhost:8080
// @BasePath /
// @schemes http
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter Bearer followed by your JWT token.

package main

import (
	"fmt"
	"log"

	"EMS/config"
	"EMS/gateway/docs"
	"EMS/gateway/routes"
	"EMS/shared/discovery"
	"EMS/shared/httpclient"
	"EMS/shared/logger"
	"EMS/shared/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize structured logging.
	logger.Init(cfg.LogLevel, "api-gateway")

	// Set Swagger host using the configured Gateway port.
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.APIGatewayPort)

	// Initialize Consul discovery client.
	discoveryClient := discovery.New(cfg)

	// Initialize shared HTTP client.
	sharedHTTPClient := httpclient.New(cfg)

	// Initialize Gin router.
	router := gin.New()

	// Register common middleware.
	router.Use(
		middleware.RequestID(),
		middleware.CorrelationID(),
		middleware.RequestLogging(),
		middleware.Recovery(),
	)

	// Gateway health endpoint.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "api-gateway",
			"status":  "UP",
		})
	})

	// Register Gateway REST routes.
	routes.RegisterGatewayRoutes(
		router,
		discoveryClient,
		sharedHTTPClient,
		cfg.JWTSecret,
	)

	// Register Gateway Swagger UI.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start API Gateway.
	address := fmt.Sprintf(":%d", cfg.APIGatewayPort)

	logger.Get().Info(
		"Starting API Gateway",
		"port", cfg.APIGatewayPort,
	)

	if err := router.Run(address); err != nil {
		logger.Get().Error(
			"API Gateway stopped unexpectedly",
			"error", err,
		)
	}
}

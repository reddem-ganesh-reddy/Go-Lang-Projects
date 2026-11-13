package main

import (
	"fmt"
	"log"

	"EMS/config"
	"EMS/services/auth/controller"
	"EMS/services/auth/routes"
	"EMS/services/auth/service"
	"EMS/shared/discovery"
	"EMS/shared/logger"
	"EMS/shared/middleware"

	"github.com/gin-gonic/gin"
)

// main starts the Auth Service.
func main() {
	// Load application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Validate application configuration.
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	// Initialize structured logging.
	if err := logger.Init(cfg.LogLevel, "auth-service"); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	appLogger := logger.Get()

	// Initialize authentication using in-memory credentials.
	authService, err := service.NewAuthService(cfg)
	if err != nil {
		appLogger.Error("failed to initialize authentication service", "error", err)
		return
	}

	// Initialize the Auth Controller.
	authController := controller.NewAuthController(authService)

	// Initialize Consul service discovery.
	consulClient := discovery.New(cfg)

	// Register Auth Service with Consul.
	serviceID := fmt.Sprintf("auth-service-%d", cfg.AuthServicePort)

	err = consulClient.RegisterService(
		"auth-service",
		serviceID,
		"127.0.0.1",
		cfg.AuthServicePort,
	)
	if err != nil {
		appLogger.Error("failed to register Auth Service", "error", err)
		return
	}

	// Create the Gin router.
	router := gin.New()

	// Register common middleware.
	router.Use(middleware.RequestID())
	router.Use(middleware.CorrelationID())
	router.Use(middleware.RequestLogging())
	router.Use(middleware.Recovery())

	// Register the health endpoint.
	router.GET("/health", middleware.HealthCheck())

	// Register authentication REST endpoints.
	routes.RegisterAuthRoutes(router, authController)

	// Start the Auth Service.
	address := fmt.Sprintf(":%d", cfg.AuthServicePort)

	appLogger.Info("Auth Service started", "address", address)

	if err := router.Run(address); err != nil {
		appLogger.Error("Auth Service stopped unexpectedly", "error", err)

		if deregisterErr := consulClient.DeregisterService(serviceID); deregisterErr != nil {
			appLogger.Error("failed to deregister Auth Service", "error", deregisterErr)
		}
	}
}

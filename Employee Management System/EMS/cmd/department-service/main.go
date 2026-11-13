package main

import (
	"fmt"
	"log"

	"EMS/cmd/department-service/docs"
	"EMS/config"
	"EMS/database"
	"EMS/services/department/controller"
	departmentGraphQL "EMS/services/department/graphql/handler"
	"EMS/services/department/graphql/resolvers"
	"EMS/services/department/model"
	"EMS/services/department/repository"
	"EMS/services/department/routes"
	"EMS/services/department/service"
	"EMS/shared/discovery"
	"EMS/shared/logger"
	"EMS/shared/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Department Service API
// @version 1.0
// @description REST API documentation for the Department Service.
// @BasePath /
// @schemes http

// main starts the Department Service.
func main() {
	// Load application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Configuration error: ", err)
	}

	// Validate application configuration.
	if err := cfg.Validate(); err != nil {
		log.Fatal("Configuration validation error: ", err)
	}

	// Initialize structured logging.
	if err := logger.Init(cfg.LogLevel, "department-service"); err != nil {
		log.Fatal("Logger initialization error: ", err)
	}

	// Connect to the Department database.
	db, err := database.Connect(cfg, cfg.DepartmentDBName)
	if err != nil {
		logger.Get().Error("Department database connection failed", "error", err)
		log.Fatal(err)
	}

	// Create the Department table if it does not exist.
	if err := db.AutoMigrate(&model.Department{}); err != nil {
		log.Fatal("Department migration failed: ", err)
	}

	// Initialize the Department repository.
	departmentRepository := repository.NewDepartmentRepository(db)

	// Initialize the Department service.
	departmentService := service.NewDepartmentService(departmentRepository)

	// Initialize the Department controller.
	departmentController := controller.NewDepartmentController(departmentService)

	// Initialize the GraphQL resolver.
	departmentResolver := resolvers.NewDepartmentResolver(departmentService)

	// Initialize the GraphQL handler.
	graphQLHandler, err := departmentGraphQL.NewDepartmentGraphQLHandler(
		departmentResolver,
	)
	if err != nil {
		log.Fatal("GraphQL initialization failed: ", err)
	}

	// Configure Swagger documentation.
	docs.SwaggerInfo.Title = "Department Service API"
	docs.SwaggerInfo.Description = "Department Service REST API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"

	// Create the Gin router.
	router := gin.New()

	// Register common middleware.
	router.Use(middleware.RequestID())
	router.Use(middleware.CorrelationID())
	router.Use(middleware.RequestLogging())
	router.Use(middleware.Recovery())

	// Register the public health endpoint.
	router.GET("/health", middleware.HealthCheck())

	// Register Department REST endpoints without JWT protection.
	routes.RegisterDepartmentRoutes(
		router,
		departmentController,
		"",
	)

	// Register GraphQL endpoint without JWT protection.
	router.Any("/graphql", graphQLHandler)

	// Register Swagger UI without JWT protection.
	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	// Initialize the Consul client.
	consulClient := discovery.New(cfg)

	// Register Department Service with Consul.
	if err := consulClient.RegisterService(
		"department-service",
		"department-service-8081",
		"127.0.0.1",
		cfg.DepartmentServicePort,
	); err != nil {
		log.Fatal("Department Service Consul registration failed: ", err)
	}

	logger.Get().Info(
		"Department Service registered with Consul",
		"service", "department-service",
		"address", "127.0.0.1",
		"port", cfg.DepartmentServicePort,
	)

	// Start the Department Service.
	address := fmt.Sprintf(":%d", cfg.DepartmentServicePort)

	logger.Get().Info("Department Service starting", "address", address)

	if err := router.Run(address); err != nil {
		_ = consulClient.DeregisterService("department-service-8081")
		log.Fatal("Department Service failed to start: ", err)
	}
}

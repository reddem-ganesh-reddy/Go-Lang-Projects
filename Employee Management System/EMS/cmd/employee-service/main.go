
package main

import (
	"fmt"
	"log"
	"net/http"

	"EMS/config"
	"EMS/database"

	employeeController "EMS/services/employee/controller"
	employeeGraphQLHandler "EMS/services/employee/graphql/handler"
	employeeGraphQLSchema "EMS/services/employee/graphql/schema"
	employeeResolver "EMS/services/employee/graphql/resolvers"
	employeeRepository "EMS/services/employee/repository"
	employeeRoutes "EMS/services/employee/routes"
	employeeService "EMS/services/employee/service"
	employeeUtilities "EMS/services/employee/utilities"

	"EMS/services/employee/model"

	"EMS/shared/discovery"
	"EMS/shared/httpclient"
	"EMS/shared/logger"
	"EMS/shared/middleware"

	employeeDocs "EMS/cmd/employee-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Employee Service API
// @version 1.0
// @description Employee Management REST API
// @BasePath /
// @schemes http

func main() {

	// Load application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Configuration loading failed: ", err)
	}

	// Initialize structured logger.
	if err := logger.Init(cfg.LogLevel, "employee-service"); err != nil {
		log.Fatal("Logger initialization failed: ", err)
	}

	// Connect to Employee database.
	db, err := database.Connect(cfg, cfg.EmployeeDBName)
	if err != nil {
		log.Fatal("Employee database connection failed: ", err)
	}

	// Automatically create or update Employee table.
	if err := db.AutoMigrate(&model.Employee{}); err != nil {
		log.Fatal("Employee database migration failed: ", err)
	}

	// Initialize Employee Repository.
	repo := employeeRepository.NewEmployeeRepository(db)

	// Initialize Consul client.
	consulClient := discovery.New(cfg)

	// Initialize shared HTTP client.
	sharedHTTPClient := httpclient.New(cfg)

	// Initialize Department HTTP client.
	departmentClient := employeeUtilities.NewDepartmentClient(
		consulClient,
		sharedHTTPClient,
	)

	// Initialize Employee Service.
	empService := employeeService.NewEmployeeService(
		repo,
		departmentClient,
	)

	// Initialize Employee Controller.
	empController := employeeController.NewEmployeeController(empService)

	// Initialize Employee GraphQL Resolver.
	resolver := employeeResolver.NewEmployeeResolver(empService)

	// Create Employee GraphQL Schema.
	graphqlSchema, err := employeeGraphQLSchema.NewEmployeeSchema(resolver)
	if err != nil {
		log.Fatal("Employee GraphQL schema creation failed: ", err)
	}

	// Initialize GraphQL Handler.
	graphqlHandler := employeeGraphQLHandler.NewEmployeeGraphQLHandler(
		graphqlSchema,
	)

	// Update Swagger information.
	employeeDocs.SwaggerInfo.Title = "Employee Service API"
	employeeDocs.SwaggerInfo.Description = "Employee Management REST API"
	employeeDocs.SwaggerInfo.Version = "1.0"
	employeeDocs.SwaggerInfo.BasePath = "/"
	employeeDocs.SwaggerInfo.Host = fmt.Sprintf(
		"localhost:%d",
		cfg.EmployeeServicePort,
	)
	employeeDocs.SwaggerInfo.Schemes = []string{"http"}

	// Create Gin router.
	router := gin.New()

	// Register shared middleware.
	router.Use(middleware.RequestID())
	router.Use(middleware.CorrelationID())
	router.Use(middleware.RequestLogging())
	router.Use(middleware.Recovery())

	// Register health check.
	router.GET("/health", middleware.HealthCheck())

	// Register Employee REST routes.
	employeeRoutes.RegisterEmployeeRoutes(router, empController)

	// Register Employee GraphQL route.
	employeeGraphQLHandler.RegisterEmployeeGraphQLRoutes(
		router,
		graphqlHandler,
	)

	// Register Swagger UI.
	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	// Register Employee Service with Consul.
	if err := consulClient.RegisterService(
		"employee-service",
		"employee-service-8082",
		"127.0.0.1",
		cfg.EmployeeServicePort,
	); err != nil {
		log.Fatal("Employee Service Consul registration failed: ", err)
	}

	// Start Employee Service.
	address := fmt.Sprintf(":%d", cfg.EmployeeServicePort)

	log.Printf("Employee Service starting on port %d", cfg.EmployeeServicePort)

	if err := router.Run(address); err != nil {
		_ = consulClient.DeregisterService("employee-service-8082")
		log.Fatal("Employee Service startup failed: ", err)
	}

	_ = http.MethodGet
}

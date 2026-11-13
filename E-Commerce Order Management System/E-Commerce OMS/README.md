# E-Commerce Order Management System

### Microservices-Based Backend Application Using Go (Golang)

A backend application developed using **Go (Golang)** and **Microservices Architecture** to manage user authentication, products, inventory reservations, and customer orders.

The system consists of four main components: API Gateway, Auth Service, Product Service, and Order Service. It uses REST APIs for communication, Consul for service discovery, PostgreSQL for data persistence, and JWT for authentication and authorization.

The application also supports GraphQL, Swagger documentation, distributed request tracking, structured logging, and resilience mechanisms.

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Key Features](#2-key-features)
3. [Technology Stack](#3-technology-stack)
4. [System Architecture](#4-system-architecture)
5. [Microservices](#5-microservices)
6. [Business Functionalities](#6-business-functionalities)
7. [Order Processing Workflow](#7-order-processing-workflow)
8. [API Documentation](#8-api-documentation)
9. [Database Architecture](#9-database-architecture)
10. [Resilience and Observability](#10-resilience-and-observability)
11. [Project Structure](#11-project-structure)
12. [AWS Architecture and CI/CD](#12-aws-architecture-and-cicd)
13. [Installation and Setup](#13-installation-and-setup)
14. [API Testing](#14-api-testing)
15. [Learning Outcomes](#15-learning-outcomes)
16. [Author](#16-author)

---

## 1. Project Overview

The E-Commerce Order Management System is a Go-based backend application designed using microservices architecture.

It provides functionalities for managing products, handling inventory reservations, processing customer orders, and authenticating users.

The application is divided into four independently running components:

- **API Gateway:** Acts as the centralized entry point for client requests.
- **Auth Service:** Handles user authentication and JWT token generation.
- **Product Service:** Manages product information and inventory reservations.
- **Order Service:** Handles order creation, retrieval, confirmation, and cancellation.

The services communicate through HTTP REST APIs and use Consul for service registration and discovery.

PostgreSQL is used for persistent storage, while GORM provides ORM-based database operations.

The application also incorporates GraphQL, Swagger/OpenAPI documentation, role-based authorization, structured logging, request tracing, and fault-tolerance mechanisms.

## 2. Key Features

### Authentication and Authorization
- User login and JWT token generation.
- JWT-based authentication.
- Role-Based Access Control (RBAC).
- Protected API endpoints.
- Middleware-based authentication and authorization.

### Product Management
- Create products.
- Retrieve all products.
- Retrieve products by ID.
- Update product information.
- Delete products.
- Product filtering, sorting, and pagination.

### Inventory Reservation Management
- Reserve product inventory.
- Retrieve inventory reservation details.
- List all reservations.
- Release reserved inventory.
- Track reservation status.

### Order Management
- Create customer orders.
- Retrieve individual orders.
- List all orders.
- Confirm existing orders.
- Cancel existing orders.
- Automatically coordinate inventory reservation during order processing.

### API and Infrastructure
- REST API development using Gin.
- GraphQL queries and mutations.
- Swagger/OpenAPI documentation.
- Centralized API Gateway.
- Consul-based service discovery.
- Inter-service HTTP communication.
- Environment-based configuration.

### Reliability and Monitoring
- HTTP timeouts.
- Bounded retries with exponential backoff.
- Circuit breaker implementation.
- Request ID tracking.
- Correlation ID propagation.
- Structured application logging.
- Recovery middleware.
- Health check endpoints.

---

## 3. Technology Stack

| Category | Technology | Purpose |
|---|---|---|
| Programming Language | Go (Golang) | Backend development |
| Web Framework | Gin | REST API development |
| ORM | GORM | Database operations |
| Database | PostgreSQL | Persistent data storage |
| API Architecture | REST | Client and service communication |
| Query Language | GraphQL | Flexible data queries and mutations |
| Authentication | JWT | User authentication |
| Authorization | RBAC | Role-based access control |
| Service Discovery | Consul | Service registration and discovery |
| API Documentation | Swagger / OpenAPI | REST API documentation |
| API Testing | SoapUI | REST API testing |
| Version Control | Git and GitHub | Source code management |
| Cloud Database | Amazon RDS | Managed PostgreSQL hosting |
| Object Storage | Amazon S3 | Build artifact storage |
| CI/CD | AWS CodeBuild | Application build and testing |
| CI/CD | AWS CodePipeline | Build and deployment automation |
| Application Deployment | AWS Elastic Beanstalk | Application hosting |

---

## 4. System Architecture

The application follows a microservices architecture in which each service performs a specific business responsibility.

The API Gateway acts as the single entry point for external client requests and forwards them to the appropriate backend service.

### Architecture Diagram

```text
                     CLIENT APPLICATION
                             |
                             v
                      API GATEWAY
                        PORT 8080
                             |
                 Authentication / Routing
                             |
          +------------------+------------------+
          |                  |                  |
          v                  v                  v
     AUTH SERVICE       PRODUCT SERVICE     ORDER SERVICE
      PORT 8083           PORT 8082          PORT 8081
          |                  |                  |
          |                  |<-----------------|
          |                  |                  |
          |                  | Inventory        |
          |                  | Reservation      |
          |                  | Communication   |
          |                  |                  |
          v                  v                  v
      AUTH DATA           PRODUCT DB          ORDER DB


                 CONSUL SERVICE DISCOVERY
                            |
                  SERVICE REGISTRATION
                   AND DISCOVERY

                  SHARED COMPONENTS
          +--------------------------------+
          | HTTP Client                    |
          | Circuit Breaker                |
          | Structured Logger              |
          | Request ID Middleware          |
          | Correlation ID Middleware      |
          | Health Check Middleware        |
          | Recovery Middleware            |
          +--------------------------------+
```

---

## 5. Microservices

### 5.1 API Gateway

**Port:** 8080

The API Gateway acts as the centralized entry point for client requests.

Responsibilities:

- Receives incoming client requests.
- Routes requests to the appropriate backend service.
- Handles JWT authentication and authorization.
- Applies request logging middleware.
- Maintains request IDs and correlation IDs.
- Supports centralized REST API access.
- Provides GraphQL Gateway functionality.

### 5.2 Auth Service

**Port:** 8083

The Auth Service manages user authentication and token generation.

Responsibilities:

- Accepts user login requests.
- Validates user credentials.
- Generates JWT access tokens.
- Supports token validation.
- Provides authentication functionality to the API Gateway.

### 5.3 Product Service

**Port:** 8082

The Product Service manages product information and inventory reservation operations.

Responsibilities:

- Create, retrieve, update, and delete products.
- Retrieve products with filtering, sorting, and pagination.
- Reserve available product stock.
- Retrieve inventory reservation details.
- Release reserved inventory.
- Expose REST and GraphQL APIs.

### 5.4 Order Service

**Port:** 8081

The Order Service manages customer orders and their lifecycle.

Responsibilities:

- Create customer orders.
- Retrieve individual orders.
- List existing orders.
- Confirm orders.
- Cancel orders.
- Communicate with the Product Service for inventory reservation and release.
- Expose REST and GraphQL APIs.

### 5.5 Communication Between Microservices

The services communicate using HTTP REST APIs.

The API Gateway forwards client requests to the appropriate backend service. The Order Service communicates directly with the Product Service for inventory-related operations.

The application workflow is:

1. The client sends a login request through the API Gateway.
2. The Gateway forwards the login request to the Auth Service.
3. The Auth Service validates credentials and generates a JWT token.
4. The client includes the JWT token in subsequent protected requests.
5. The Gateway validates authentication and authorization.
6. Product requests are forwarded to the Product Service.
7. Order requests are forwarded to the Order Service.
8. During order creation, the Order Service communicates with the Product Service to reserve the required inventory.
9. When an order is cancelled, the Order Service requests the Product Service to release the reserved inventory.
10. The response is returned to the client through the API Gateway.

---

## 6. Business Functionalities

### 6.1 User Authentication

Users authenticate through the Auth Service.

After successful authentication, the service generates a JWT token. The token is used to access protected APIs through the API Gateway.

### 6.2 Product Management

The Product Service supports:

- Creating products.
- Retrieving all products.
- Retrieving products by ID.
- Updating product information.
- Deleting products.
- Filtering and sorting products.
- Pagination.

### 6.3 Inventory Reservation Management

The Product Service manages inventory reservations.

It supports:

- Creating inventory reservations.
- Retrieving reservation details.
- Listing all reservations.
- Releasing reserved stock.
- Tracking reservation status.

### 6.4 Create Order

The customer provides the required order information, including customer ID, product IDs, and quantities.

The system validates the request, checks product availability, communicates with the Product Service to reserve inventory, and saves the order with the initial status `CREATED`.

### 6.5 Get Order

Retrieves an individual order using its unique order ID.

The response contains available order information, including customer details, product information, quantities, and order status.

### 6.6 List Orders

Retrieves existing orders and their associated order items from the Order Service.

### 6.7 Update Order Status

The Order Service supports the following order statuses:

| Status | Description |
|---|---|
| CREATED | Order has been created. |
| CONFIRMED | Order has been confirmed. |
| CANCELLED | Order has been cancelled. |

When an order is cancelled, the system requests the Product Service to release the previously reserved inventory.

---

## 7. Order Processing Workflow

The following workflow illustrates how an order is processed.

```text
       CUSTOMER
           |
           v
       CREATE ORDER
           |
           v
      API GATEWAY
           |
           v
      ORDER SERVICE
           |
           v
      VALIDATE ORDER
           |
           v
   CHECK PRODUCT AVAILABILITY
           |
           v
     PRODUCT SERVICE
           |
           v
   RESERVE PRODUCT STOCK
           |
           v
      SAVE ORDER
           |
           v
    STATUS: CREATED
           |
           v
    ORDER RESPONSE
           |
           v
        CUSTOMER
```

### Order Cancellation Workflow

```text
       CANCEL ORDER
            |
            v
       ORDER SERVICE
            |
            v
     UPDATE ORDER STATUS
            |
            v
    REQUEST STOCK RELEASE
            |
            v
      PRODUCT SERVICE
            |
            v
    RELEASE RESERVED STOCK
            |
            v
     STATUS: CANCELLED
```

---

## 8. API Documentation

All REST APIs are accessible through the API Gateway at:

`http://localhost:8080`

### 8.1 Auth Service

**User Login**

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/auth/login` | Authenticates a user and returns a JWT token. |

### 8.2 Product Service

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/products` | Retrieves products with optional filtering, sorting, and pagination. |
| GET | `/api/products/{id}` | Retrieves a product by ID. |
| POST | `/api/products` | Creates a new product. |
| PUT | `/api/products/{id}` | Updates an existing product. |
| DELETE | `/api/products/{id}` | Deletes a product. |

### 8.3 Inventory Reservation APIs

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/inventory/reservations` | Retrieves all inventory reservations. |
| GET | `/api/inventory/reservations/{reservation_id}` | Retrieves a reservation by ID. |
| POST | `/api/inventory/reservations` | Creates an inventory reservation. |
| POST | `/api/inventory/reservations/release` | Releases reserved inventory. |

### 8.4 Order Service

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/orders` | Creates a new order. |
| GET | `/api/orders` | Retrieves all orders. |
| GET | `/api/orders/{id}` | Retrieves an order by ID. |
| PATCH | `/api/orders/{id}/confirm` | Confirms an existing order. |
| PATCH | `/api/orders/{id}/cancel` | Cancels an existing order. |

### 8.5 GraphQL APIs

The Product and Order services expose GraphQL endpoints.

| Service | Endpoint | Purpose |
|---|---|---|
| Product Service | `http://localhost:8082/graphql` | Product and inventory queries and mutations. |
| Order Service | `http://localhost:8081/graphql` | Order queries and mutations. |

### 8.6 Health Check APIs

| Service | Method | Endpoint |
|---|---|---|
| API Gateway | GET | `http://localhost:8080/health` |
| Auth Service | GET | `http://localhost:8083/health` |
| Product Service | GET | `http://localhost:8082/health` |
| Order Service | GET | `http://localhost:8081/health` |

These endpoints help verify whether the respective services are running.

---

## 9. Database Architecture

PostgreSQL is used for persistent data storage, and GORM handles database operations.

The application follows database separation for its main business services.

| Database | Service | Stored Information |
|---|---|---|
| `product_db` | Product Service | Product details, stock information, and inventory reservations. |
| `order_db` | Order Service | Orders, order items, and order statuses. |

### Database Independence

- The Product Service manages its own product and inventory data.
- The Order Service manages its own order data.
- Services communicate through APIs instead of directly accessing another service's database.
- GORM is used for ORM-based database operations.
- Environment variables are used for database configuration.

---

## 10. Resilience and Observability

The application incorporates shared components to improve reliability and simplify troubleshooting.

| Component | Purpose |
|---|---|
| Retry Mechanism | Retries eligible failed service requests. |
| Exponential Backoff | Increases the delay between retry attempts. |
| HTTP Timeout | Prevents requests from waiting indefinitely. |
| Circuit Breaker | Helps prevent repeated calls to unhealthy services. |
| Request ID | Identifies individual HTTP requests. |
| Correlation ID | Tracks requests across multiple services. |
| Structured Logging | Records application events in a structured format. |
| Recovery Middleware | Handles unexpected request-processing panics. |
| Health Checks | Monitors service availability. |

The shared components are maintained under the `shared/` directory.

---

## 11. Project Structure

```text
E-Commerce OMS/
│
├── cmd/
│   ├── api-gateway/
│   │   ├── main.go
│   │   ├── graphql_gateway.go
│   │   ├── swagger.go
│   │   └── docs/
│   │
│   ├── auth-service/
│   │   └── main.go
│   │
│   ├── order-service/
│   │   ├── main.go
│   │   ├── order_inventory_reservation.go
│   │   ├── testing_graphql.go
│   │   └── testing_swagger.go
│   │
│   └── product-service/
│       ├── main.go
│       ├── tesing_graphql.go
│       └── testing_swagger.go
│
├── config/
│   └── config.go
│
├── database/
│   └── database.go
│
├── services/
│   ├── auth/
│   ├── order/
│   │   ├── controller/
│   │   ├── docs/
│   │   ├── dto/
│   │   ├── graphql/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── routes/
│   │   ├── service/
│   │   └── utilities/
│   │
│   └── product/
│       ├── controller/
│       ├── docs/
│       ├── dto/
│       ├── graphql/
│       ├── model/
│       ├── repository/
│       ├── routes/
│       ├── service/
│       └── utilities/
│
├── shared/
│   ├── discovery/
│   ├── httpclient/
│   ├── logger/
│   └── middleware/
│
├── test-app/
│   ├── Admin Testing – Order Service and Inventory Flow.MD
│   ├── Admin Testing – Product and Inventory Services.MD
│   └── User – Order Service and Inventory Flow Testing.MD
│
├── .env.example
├── .gitignore
├── buildspec.yml
├── consul-config.ps1
├── go.mod
├── go.sum
├── Procfile
└── start-all.ps1
```

### Layered Architecture

The Product and Order services follow a layered architecture.

| Layer | Responsibility |
|---|---|
| Controller | Handles HTTP requests and responses. |
| Service | Implements business logic. |
| Repository | Performs database operations. |
| Model | Defines database entities. |
| DTO | Defines request and response structures. |
| Routes | Registers API endpoints. |
| Utilities | Handles validation and helper functions. |
| GraphQL | Manages schemas, queries, mutations, and resolvers. |
| Docs | Contains Swagger documentation. |

---

## 12. AWS Architecture and CI/CD

The project includes AWS deployment and CI/CD configuration using AWS services.

The intended cloud architecture uses managed database hosting, object storage, automated builds, and application deployment.

### AWS Services

**12.1 GitHub**

GitHub is used for source code management and version control.

It serves as the source repository for the AWS CodePipeline workflow.

**12.2 Amazon RDS**

Amazon Relational Database Service is used to host PostgreSQL databases in the cloud.

The application uses database connection settings supplied through environment variables.

**12.3 Amazon S3**

Amazon Simple Storage Service is used to store application build artifacts and deployment packages.

**12.4 AWS CodeBuild**

AWS CodeBuild handles the application build process.

The project includes a `buildspec.yml` file that defines build-related commands and configuration.

**12.5 AWS Elastic Beanstalk**

AWS Elastic Beanstalk is the target application hosting environment.

It provides application deployment and environment management capabilities.

**12.6 AWS CodePipeline**

AWS CodePipeline connects the source repository, build process, and deployment stages.

The intended workflow is:

```text
       GITHUB REPOSITORY
               |
               v
        AWS CODEPIPELINE
               |
               v
          AWS CODEBUILD
               |
               v
        BUILD ARTIFACTS
               |
               v
          AMAZON S3
               |
               v
      ELASTIC BEANSTALK
               |
               v
       APPLICATION ENVIRONMENT
               |
               v
          AMAZON RDS
```

This architecture provides a foundation for automated application delivery and cloud deployment.

---

## 13. Installation and Setup

### Prerequisites

Install the following tools:

- Go (Golang)
- PostgreSQL
- Consul
- Git
- SoapUI (for API testing)

### Step 1: Clone the Repository

```powershell
git clone <repository-url>
```

### Step 2: Navigate to the Project Directory

```powershell
cd "E-Commerce OMS"
```

### Step 3: Configure Environment Variables

Create the `.env` file using the provided `.env.example` template.

Configure the required values, including:

- Database connection settings.
- Service ports.
- JWT configuration.
- Consul address.
- Service URLs.
- HTTP timeout and retry settings.
- Circuit breaker settings.
- Logging configuration.

**Security:** Never commit the actual `.env` file, database passwords, or JWT secrets to GitHub.

### Step 4: Download Dependencies

```powershell
go mod download
```

### Step 5: Verify Dependencies

```powershell
go mod verify
```

### Step 6: Start Infrastructure

Ensure PostgreSQL is running and configure Consul according to the project settings.

### Step 7: Start Microservices

Run the provided PowerShell startup script:

```powershell
.\start-all.ps1
```

The script starts the application services according to the configured environment.

### Step 8: Verify Service Health

Check the health endpoints:

```text
http://localhost:8080/health
http://localhost:8083/health
http://localhost:8082/health
http://localhost:8081/health
```

---

## 14. API Testing

SoapUI is used to test the REST APIs and validate application responses.

The `test-app/` directory contains separate testing documentation for different roles and business workflows.

### Testing Documentation

| Document | Description |
|---|---|
| Admin Testing – Order Service and Inventory Flow | Tests administrative order and inventory operations. |
| Admin Testing – Product and Inventory Services | Tests product and inventory management operations. |
| User – Order Service and Inventory Flow Testing | Tests user-level order and inventory-related workflows. |

Testing covers:

- Authentication and JWT validation.
- Role-based authorization.
- Product management.
- Inventory reservation and release.
- Order creation.
- Order retrieval.
- Order confirmation.
- Order cancellation.
- Inventory restoration during cancellation.
- Service health endpoints.

---

## 15. Learning Outcomes

This project provided practical experience in:

- Designing and implementing Go-based microservices.
- Developing REST APIs using Gin.
- Implementing GraphQL APIs.
- Working with PostgreSQL and GORM.
- Implementing JWT authentication and RBAC.
- Designing layered application architecture.
- Implementing inter-service communication.
- Integrating Consul service discovery.
- Implementing inventory reservation workflows.
- Applying retry, timeout, and circuit breaker patterns.
- Implementing request tracking and structured logging.
- Documenting and testing APIs using Swagger and SoapUI.
- Understanding AWS infrastructure and CI/CD architecture.

---

## 16. Author

**Reddem Ganesh Reddy**

---

**Project:** E-Commerce Order Management System

**Domain:** Backend Development / Microservices

**Programming Language:** Go (Golang)

**Architecture:** Microservices Architecture

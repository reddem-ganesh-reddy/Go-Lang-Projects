# Employee Management System (EMS)

A Go-based Employee Management System built using a microservices architecture. The application provides employee and department management through REST APIs and GraphQL, with centralized authentication and authorization through an API Gateway.

---

## 1. Project Overview

The Employee Management System consists of four microservices:

1. **API Gateway** – Central entry point for client requests. Handles request routing, authentication, and authorization.
2. **Auth Service** – Handles user login, password verification, and JWT generation.
3. **Employee Service** – Manages employee information and communicates with the Department Service to validate department details.
4. **Department Service** – Manages department information.

The application uses PostgreSQL for data storage, GORM for database operations, Consul for service discovery, JWT for authentication, and Swagger for API documentation.

## 2. Technologies Used

| Technology | Purpose |
|---|---|
| Go | Backend programming language |
| Gin | HTTP web framework |
| GORM | ORM and database operations |
| PostgreSQL | Relational database |
| REST API | HTTP-based service communication |
| GraphQL | Query and mutation operations |
| Swagger | API documentation and testing |
| JWT | Authentication |
| bcrypt | Password hashing |
| Consul | Service discovery |
| log/slog | Structured logging |
| PowerShell | Windows commands and service startup |

## 3. System Architecture

```text
                         Client
                           |
                           v
                      API Gateway
                       Port 8080
                           |
              +------------+------------+
              |            |            |
              v            v            v
         Auth Service  Employee Service  Department Service
           :8083          :8082              :8081
                            |
                            v
                     Department Service
                            |
                            v
                      department_db

                  Employee Service
                            |
                            v
                       employee_db
```

### Service Details

| Service | Port | Responsibility |
|---|---:|---|
| API Gateway | 8080 | Centralized routing, authentication, authorization |
| Department Service | 8081 | Department management |
| Employee Service | 8082 | Employee management and department validation |
| Auth Service | 8083 | Login and JWT generation |

### Service Communication

- Clients access the system through the API Gateway.
- The Gateway forwards requests to the appropriate service.
- Employee Service communicates with Department Service to validate department information during employee creation and updates.
- Employee Service does not directly access the Department Service database.
- Consul is used for service discovery.
- Shared HTTP client functionality is used for inter-service communication.

## 4. Project Structure

```text
EMS/
│
├── cmd/
│   ├── api-gateway/
│   │   ├── main.go
│   │   ├── swagger.go
│   │   └── graphql_swagger.go
│   │
│   ├── auth-service/
│   │   └── main.go
│   │
│   ├── employee-service/
│   │   └── main.go
│   │
│   └── department-service/
│       └── main.go
│
├── config/
│   └── config.go
│
├── database/
│   └── database.go
│
├── gateway/
│   ├── docs/
│   ├── graphql/
│   │   ├── handler/
│   │   ├── inputs/
│   │   ├── mutations/
│   │   ├── queries/
│   │   ├── resolvers/
│   │   ├── schema/
│   │   └── types/
│   └── routes/
│
├── services/
│   ├── auth/
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
│   ├── department/
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
│   └── employee/
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
│   │   └── consul.go
│   │
│   ├── errors/
│   │   └── response.go
│   │
│   ├── httpclient/
│   │   ├── client.go
│   │   └── circuit_breaker.go
│   │
│   ├── logger/
│   │   └── logger.go
│   │
│   └── middleware/
│       ├── correlation_id.go
│       ├── health.go
│       ├── jwt_auth.go
│       ├── recovery.go
│       ├── request_id.go
│       └── request_logging.go
│
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

### Layered Architecture

Each business service follows a layered architecture:

- **Controller:** Handles incoming HTTP requests.
- **DTO:** Defines request and response structures.
- **Model:** Defines database entities.
- **Repository:** Performs database operations using GORM.
- **Service:** Contains business logic.
- **Routes:** Registers REST API routes.
- **GraphQL:** Contains schemas, handlers, resolvers, queries, mutations, inputs, and types.
- **Utilities:** Contains validation and helper functions.
- **Docs:** Stores generated Swagger documentation.

## 5. Database Design

The project uses two PostgreSQL databases.

| Database | Owner |
|---|---|
| employee_db | Employee Service |
| department_db | Department Service |

### Employee Database

Stores employee information:

- Employee ID
- Employee name
- Email
- Department ID
- Created timestamp
- Updated timestamp

### Department Database

Stores department information:

- Department ID
- Department name

### Database Responsibilities

- Employee Service manages `employee_db`.
- Department Service manages `department_db`.
- Employee Service validates departments by communicating with Department Service.
- Auth Service does not use a separate `auth_db`.
- Database credentials are loaded through environment configuration.

## 6. Prerequisites

Install the following software:

1. Go 1.27.1 or the Go version configured for the project.
2. PostgreSQL.
3. Consul.
4. Git (optional).

Verify Go installation:

```powershell
go version
```

Verify PostgreSQL is running and accessible.

Verify Consul:

```powershell
consul version
```

## 7. Environment Configuration

The application loads environment variables from the root `.env` file using `godotenv`.

Configuration includes:

- Database connection settings
- Service ports
- JWT secret and expiration
- Consul address
- HTTP client timeout
- Maximum retry count
- Retry delay
- Circuit breaker failure threshold
- Circuit breaker reset timeout
- Logging level
- Authentication credentials

Example resilience configuration values:

| Setting | Example |
|---|---:|
| HTTP client timeout | 5 seconds |
| Maximum retries | 3 |
| Circuit breaker failure threshold | 5 |
| Circuit breaker reset timeout | 30 seconds |

Use the exact environment-variable names defined in `config/config.go`.

**Important:**

- Do not commit `.env` to Git.
- Do not hardcode passwords, JWT secrets, database credentials, or service URLs in Go files.
- Configure the required environment variables before starting the services.
- Change development credentials before deployment.

## 8. Installation and Setup

### Step 1: Navigate to the Project

```powershell
cd D:\Tenarai\EMS
```

### Step 2: Download Dependencies

```powershell
go mod tidy
```

### Step 3: Configure PostgreSQL

Start PostgreSQL and ensure these databases are available:

```text
employee_db
department_db
```

Configure their connection details in `.env`.

### Step 4: Start Consul

Run Consul in development mode:

```powershell
consul agent -dev
```

Default Consul address:

```text
127.0.0.1:8500
```

Consul dashboard:

http://localhost:8500/ui/

The services register themselves with Consul when they start.

## 9. Swagger Documentation

Swagger provides interactive documentation for REST APIs.

The project provides Swagger UI for:

- API Gateway
- Employee Service
- Department Service

Auth Service does not have a separate Swagger UI.

### 9.1 Install Swagger CLI

Install the Swagger command-line tool:

```powershell
go install github.com/swaggo/swag/cmd/swag@v1.16.6
```

Check the installation:

```powershell
swag --version
```

If PowerShell cannot find `swag`, add the Go bin directory to the current PATH:

```powershell
$env:Path += ";$(go env GOPATH)\bin"
```

### 9.2 Generate API Gateway Swagger Documentation

Navigate to the project root:

```powershell
cd D:\Tenarai\EMS
```

Run:

```powershell
swag init -g main.go -d .\cmd\api-gateway -o .\gateway\docs
```

This generates:

```text
gateway/
└── docs/
    ├── docs.go
    ├── swagger.json
    └── swagger.yaml
```

### 9.3 Generate Employee Service Swagger Documentation

Use the following command if the Employee Service Swagger annotations are in its controller package and the general API information is in `cmd/employee-service/main.go`:

```powershell
swag init -g main.go -d .\cmd\employee-service,.\services\employee\controller -o .\services\employee\docs
```

Generated documentation is stored in:

```text
services/
└── employee/
    └── docs/
        ├── docs.go
        ├── swagger.json
        └── swagger.yaml
```

### 9.4 Generate Department Service Swagger Documentation

Use the following command if the Department Service Swagger annotations are in its controller package and the general API information is in `cmd/department-service/main.go`:

```powershell
swag init -g main.go -d .\cmd\department-service,.\services\department\controller -o .\services\department\docs
```

Generated documentation is stored in:

```text
services/
└── department/
    └── docs/
        ├── docs.go
        ├── swagger.json
        └── swagger.yaml
```

If your annotations are located in additional packages, include those directories in the `-d` argument.

### 9.5 When Should Swagger Be Regenerated?

Regenerate Swagger documentation whenever you:

- Add a new REST endpoint.
- Modify an existing REST endpoint.
- Change request or response structures.
- Change endpoint descriptions.
- Change API annotations.
- Add or modify authentication documentation.

Run the relevant `swag init` command after making those changes.

### 9.6 Access Swagger UI

| Service | Swagger URL |
|---|---|
| API Gateway | http://localhost:8080/swagger/index.html |
| Department Service | http://localhost:8081/swagger/index.html |
| Employee Service | http://localhost:8082/swagger/index.html |

The Gateway Swagger UI documents centralized Auth, Employee, and Department REST endpoints and the Gateway GraphQL endpoint.

## 10. Running the Services

Run each service from the project root in a separate PowerShell terminal.

### Department Service

```powershell
go run .\cmd\department-service
```

Port: `8081`

### Employee Service

```powershell
go run .\cmd\employee-service
```

Port: `8082`

### Auth Service

```powershell
go run .\cmd\auth-service
```

Port: `8083`

### API Gateway

```powershell
go run .\cmd\api-gateway
```

Port: `8080`

Start PostgreSQL and Consul first. Start Department and Employee Services before the Gateway so service registration and health checks can become available.

## 11. Authentication and Authorization

The application uses JWT-based authentication and role-based access control.

### User Roles

| Role | Permissions |
|---|---|
| Admin | Employee and department management |
| Employee | Read-only access to permitted information |

### Development Login Credentials

| Role | Username | Password |
|---|---|---|
| Admin | admin | admin@123 |
| Employee | emp | emp@123 |

These are development credentials. Change them before using the application outside local development.

### Authentication Flow

1. Client sends login credentials through the API Gateway.
2. Gateway forwards the login request to Auth Service.
3. Auth Service verifies the credentials.
4. Password verification uses bcrypt.
5. Auth Service generates a JWT.
6. Client includes the JWT in the Authorization header for protected requests.
7. Gateway validates the token and applies role-based authorization.
8. Authorized requests are forwarded to the appropriate service.

Authorization header:

```text
Authorization: Bearer <JWT_TOKEN>
```

### Admin Permissions

- Create, view, update, and delete employees.
- Create, view, update, and delete departments.
- Execute permitted GraphQL queries and mutations.

### Employee Permissions

- Read permitted employee information.
- Execute permitted GraphQL queries.
- Cannot create, update, or delete employees.
- Cannot create, update, or delete departments.

## 12. GraphQL

GraphQL is available through the Employee Service, Department Service, and API Gateway.

### Employee GraphQL

Supported operations:

- Get all employees.
- Get employee by ID.
- Pagination.
- Filtering.
- Sorting.
- Create employee.
- Update employee.
- Delete employee.

Example query:

```graphql
query {
  employees {
    employees {
      id
      name
      email
      department_id
    }
    pagination {
      total_records
      current_page
      page_size
      total_pages
    }
  }
}
```

Example query by ID:

```graphql
query {
  employee(id: 1) {
    id
    name
    email
    department_id
  }
}
```

Example mutation:

```graphql
mutation {
  createEmployee(
    input: {
      name: "John"
      email: "john@example.com"
      department_id: 1
    }
  ) {
    id
    name
    email
    department_id
  }
}
```

### Department GraphQL

Supported operations:

- Get all departments.
- Get department by ID.
- Create department.
- Update department.
- Delete department.

Example query:

```graphql
query {
  departments {
    id
    name
  }
}
```

Example query by ID:

```graphql
query {
  department(id: "1") {
    id
    name
  }
}
```

Example mutation:

```graphql
mutation {
  createDepartment(input: { name: "Testing" }) {
    id
    name
  }
}
```

### API Gateway GraphQL

The Gateway exposes a centralized GraphQL endpoint:

```text
POST http://localhost:8080/graphql
```

The Gateway forwards GraphQL operations to the appropriate Employee or Department Service.

Example JSON request:

```json
{
  "query": "query { departments { id name } }"
}
```

Example employee request:

```json
{
  "query": "query { employee(id: 1) { id name email department_id } }"
}
```

**Important:** The current Gateway GraphQL proxy routes an operation to either Employee Service or Department Service. An operation containing root fields from both services is not supported by the current proxy.

## 13. GraphQL Testing Through Swagger

Open:

http://localhost:8080/swagger/index.html

Locate:

```text
POST /graphql
```

Testing steps:

1. Click **Try it out**.
2. Click **Authorize**.
3. Enter `Bearer` followed by a valid JWT.
4. Enter the GraphQL request as a JSON body.
5. Click **Execute**.
6. Verify the response.

Example:

```json
{
  "query": "query { departments { id name } }"
}
```

For mutations, send the corresponding GraphQL mutation inside the JSON `query` field.

**Note:** GraphQL may return HTTP 200 even when the response contains errors. Always inspect the `errors` field.

## 14. Middleware and Shared Components

### Request ID

Generates or accepts a request ID for tracking an individual request.

### Correlation ID

Maintains the same correlation ID across requests between the API Gateway and downstream services.

### JWT Authentication

Validates JWT tokens and extracts authorization information.

### Recovery Middleware

Handles unexpected panics during request processing.

### Health Middleware

Provides service health endpoints.

### Request Logging

Uses structured logging to record request details.

### Shared HTTP Client

Provides:

- Configurable timeout.
- Retry configuration.
- Circuit breaker.
- Request ID forwarding.
- Correlation ID forwarding.
- Error handling.

Retries should be restricted to appropriate operations. Non-idempotent operations such as record creation and GraphQL mutations should not be automatically retried without safeguards.

## 15. Logging and Monitoring

The application uses Go's `log/slog` for structured logging.

Supported logging levels:

- DEBUG
- INFO
- WARN
- ERROR

Log information may include:

- Service name.
- Request ID.
- Correlation ID.
- Operation.
- Error details.

Passwords, JWT tokens, and database credentials must not be logged.

Consul dashboard:

http://localhost:8500/ui/

## 16. Health Checks

Each service exposes a health endpoint:

```text
GET /health
```

Test Gateway:

```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/health"
```

Test Department Service:

```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8081/health"
```

Test Employee Service:

```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8082/health"
```

Test Auth Service:

```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8083/health"
```

Check the returned health information and Consul service status.

## 17. Development Commands

Run these commands from the project root.

### Download and Maintain Dependencies

```powershell
go mod tidy
```

### Format Go Files

Format all Go source files:

```powershell
Get-ChildItem -Recurse -Filter *.go -File | ForEach-Object {
    gofmt -w $_.FullName
}
```

### Build All Packages

```powershell
go build ./...
```

### Run All Tests

```powershell
go test ./...
```

### Run Tests With Verbose Output

```powershell
go test -v ./...
```

### Build Individual Services

API Gateway:

```powershell
go build -o .\bin\api-gateway.exe .\cmd\api-gateway
```

Employee Service:

```powershell
go build -o .\bin\employee-service.exe .\cmd\employee-service
```

Department Service:

```powershell
go build -o .\bin\department-service.exe .\cmd\department-service
```

Auth Service:

```powershell
go build -o .\bin\auth-service.exe .\cmd\auth-service
```

## 18. Testing

### REST API Testing

Use Gateway Swagger UI:

http://localhost:8080/swagger/index.html

Test:

- Login.
- Employee CRUD operations.
- Department CRUD operations.
- Admin authorization.
- Employee authorization.
- Missing or invalid JWT.
- Invalid request data.
- Non-existing records.
- Department validation during employee creation and updates.

### GraphQL Testing

Use Gateway Swagger UI and test:

- Employee queries.
- Department queries.
- Employee queries by ID.
- Department queries by ID.
- Pagination.
- Filtering.
- Sorting.
- Employee mutations.
- Department mutations.
- Invalid input.
- Missing or invalid authentication.
- Role-based access.

### Build and Test Verification

Run:

```powershell
go build ./...
go test ./...
```

Both commands should complete successfully before considering the project ready for use.

## 19. Common HTTP Status Codes

| Status Code | Description |
|---|---|
| 200 | Request successful |
| 201 | Resource created |
| 204 | Request successful with no response body |
| 400 | Invalid request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Resource not found |
| 409 | Conflict |
| 500 | Internal server error |
| 502 | Bad gateway |
| 503 | Service unavailable |
| 504 | Gateway timeout |

Actual responses depend on the endpoint and error-handling implementation.

## 20. Troubleshooting

### Go Dependency Issues

Run:

```powershell
go mod tidy
```

### Swagger Command Not Found

Add Go's binary directory to the current PowerShell session:

```powershell
$env:Path += ";$(go env GOPATH)\bin"
```

Verify:

```powershell
swag --version
```

### Swagger Documentation Not Updating

Regenerate the documentation using the relevant `swag init` command and restart the corresponding service.

### Database Connection Failure

Check:

- PostgreSQL is running.
- Database names are correct.
- Database credentials in `.env` are valid.
- Database connection configuration is correct.

### Consul Service Discovery Failure

Check:

- Consul agent is running.
- Consul address is correct in `.env`.
- Services have registered successfully.
- Health checks are passing.

### Unauthorized Response

Check:

- JWT is present.
- Token is valid.
- Authorization header uses `Bearer <JWT_TOKEN>`.
- User role has permission for the requested operation.

### Downstream Service Unavailable

Check:

- The required microservice is running.
- Its configured port is correct.
- It is registered in Consul.
- Its health endpoint is responding.

## 21. Important Design Decisions

- The project contains exactly four services: API Gateway, Auth Service, Employee Service, and Department Service.
- No separate Auth database is used.
- Employee Service owns employee data.
- Department Service owns department data.
- Employee Service validates departments through Department Service communication.
- API Gateway is the centralized client entry point.
- Employee and Department Services retain their own REST and GraphQL interfaces for independent testing.
- Consul provides service discovery.
- Configuration is loaded through `.env`.
- Shared middleware and HTTP client components provide common functionality.
- No HR Service, Payroll Service, Docker, or Kafka is included in the current project.

## 22. Future Improvements

Potential improvements, if required:

- More comprehensive integration testing.
- GraphQL query depth and complexity limits.
- Consistent pagination limits.
- Improved circuit breaker half-open handling.
- Centralized metrics and distributed tracing.
- Production-ready secrets management.
- Containerized deployment.

## 23. Project Information

**Project Name:** Employee Management System (EMS)  
**Programming Language:** Go  
**Architecture:** Microservices  
**Framework:** Gin  
**ORM:** GORM  
**Database:** PostgreSQL  
**API Styles:** REST and GraphQL  
**Documentation:** Swagger  
**Authentication:** JWT and bcrypt  
**Service Discovery:** Consul  
**Logging:** log/slog
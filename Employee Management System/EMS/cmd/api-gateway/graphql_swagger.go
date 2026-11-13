package main

// SwaggerGraphQLRequest represents a GraphQL HTTP request.
type SwaggerGraphQLRequest struct {
	Query         string                 `json:"query" example:"query { employees { employees { id name email } } }"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

// SwaggerGraphQLError represents a GraphQL error.
type SwaggerGraphQLError struct {
	Message string `json:"message"`
}

// SwaggerGraphQLResponse represents a GraphQL response.
type SwaggerGraphQLResponse struct {
	Data   map[string]interface{} `json:"data,omitempty"`
	Errors []SwaggerGraphQLError  `json:"errors,omitempty"`
}

// swaggerGatewayGraphQL godoc
// @Summary Gateway GraphQL endpoint
// @Description Execute Employee or Department GraphQL queries and mutations. A valid JWT is required. Employees can execute queries only. Administrators can execute queries and mutations. Employee and Department root fields must be requested in separate operations.
// @Tags GraphQL
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SwaggerGraphQLRequest true "GraphQL query and variables"
// @Success 200 {object} SwaggerGraphQLResponse
// @Failure 400 {object} SwaggerErrorResponse
// @Failure 401 {object} SwaggerErrorResponse
// @Failure 403 {object} SwaggerErrorResponse
// @Failure 503 {object} SwaggerGatewayErrorResponse
// @Failure 504 {object} SwaggerGatewayErrorResponse
// @Router /graphql [post]
func swaggerGatewayGraphQL() {}

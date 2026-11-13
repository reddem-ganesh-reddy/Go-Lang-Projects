package routes

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"EMS/shared/discovery"
	"EMS/shared/httpclient"
	"EMS/shared/middleware"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

type graphQLRequest struct {
	Query         string `json:"query"`
	OperationName string `json:"operationName"`
}

type graphQLRoute struct {
	serviceName string
	isMutation  bool
}

// GraphQLProxyHandler authenticates and authorizes GraphQL requests
// before forwarding them to the appropriate service.
func GraphQLProxyHandler(
	discoveryClient *discovery.Client,
	httpClient *httpclient.Client,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Unable to read GraphQL request",
			})
			return
		}

		var request graphQLRequest
		if err := json.Unmarshal(body, &request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid GraphQL request body",
			})
			return
		}

		if strings.TrimSpace(request.Query) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "GraphQL query is required",
			})
			return
		}

		route, err := determineGraphQLRoute(
			request.Query,
			request.OperationName,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		// Only administrators can execute GraphQL mutations.
		role := middleware.GetUserRole(c)

		if route.isMutation && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Only administrators can execute GraphQL mutations",
			})
			return
		}

		// Discover the destination service.
		serviceURL, err := discoveryClient.DiscoverService(route.serviceName)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "GraphQL service is currently unavailable",
			})
			return
		}

		requestURL := strings.TrimRight(serviceURL, "/") + "/graphql"

		headers := map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}

		// Forward the JWT so the destination service can validate it too.
		if authorization := c.GetHeader("Authorization"); authorization != "" {
			headers["Authorization"] = authorization
		}

		if accept := c.GetHeader("Accept"); accept != "" {
			headers["Accept"] = accept
		}

		// Forward the request without retrying GraphQL POST operations.
		response, err := httpClient.Do(
			c.Request.Context(),
			http.MethodPost,
			requestURL,
			body,
			headers,
		)

		if err != nil {
			statusCode := http.StatusServiceUnavailable

			if errors.Is(err, c.Request.Context().Err()) {
				statusCode = http.StatusGatewayTimeout
			} else {
				var networkError net.Error
				if errors.As(err, &networkError) && networkError.Timeout() {
					statusCode = http.StatusGatewayTimeout
				}
			}

			c.JSON(statusCode, gin.H{
				"success": false,
				"message": "Unable to communicate with the GraphQL service",
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

// determineGraphQLRoute identifies the service and operation type
// from the GraphQL document.
func determineGraphQLRoute(
	query string,
	operationName string,
) (graphQLRoute, error) {
	document, err := parser.Parse(parser.ParseParams{
		Source: source.NewSource(&source.Source{
			Body: []byte(query),
			Name: "Gateway GraphQL request",
		}),
	})
	if err != nil {
		return graphQLRoute{}, errors.New("invalid GraphQL query")
	}

	var selected *ast.OperationDefinition

	for _, definition := range document.Definitions {
		operation, ok := definition.(*ast.OperationDefinition)
		if !ok {
			continue
		}

		if operationName != "" {
			if operation.Name != nil && operation.Name.Value == operationName {
				selected = operation
				break
			}
		} else {
			if selected != nil {
				return graphQLRoute{}, errors.New(
					"operationName is required when multiple operations are provided",
				)
			}
			selected = operation
		}
	}

	if selected == nil {
		return graphQLRoute{}, errors.New("GraphQL operation not found")
	}

	if selected.Operation != "query" && selected.Operation != "mutation" {
		return graphQLRoute{}, errors.New("GraphQL subscriptions are not supported")
	}

	employeeFields := map[string]bool{
		"employees":      true,
		"employee":       true,
		"createEmployee": true,
		"updateEmployee": true,
		"deleteEmployee": true,
	}

	departmentFields := map[string]bool{
		"departments":      true,
		"department":       true,
		"createDepartment": true,
		"updateDepartment": true,
		"deleteDepartment": true,
	}

	var serviceName string

	for _, selection := range selected.SelectionSet.Selections {
		field, ok := selection.(*ast.Field)
		if !ok {
			return graphQLRoute{}, errors.New(
				"top-level GraphQL fragments are not supported by the Gateway",
			)
		}

		fieldName := field.Name.Value

		if employeeFields[fieldName] {
			if serviceName != "" && serviceName != "employee-service" {
				return graphQLRoute{}, errors.New(
					"employee and department fields must be requested in separate GraphQL operations",
				)
			}
			serviceName = "employee-service"
			continue
		}

		if departmentFields[fieldName] {
			if serviceName != "" && serviceName != "department-service" {
				return graphQLRoute{}, errors.New(
					"employee and department fields must be requested in separate GraphQL operations",
				)
			}
			serviceName = "department-service"
			continue
		}

		return graphQLRoute{}, errors.New(
			"unsupported GraphQL field: " + fieldName,
		)
	}

	if serviceName == "" {
		return graphQLRoute{}, errors.New("no supported GraphQL fields found")
	}

	return graphQLRoute{
		serviceName: serviceName,
		isMutation:  selected.Operation == "mutation",
	}, nil
}

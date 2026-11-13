package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"order-management/config"

	"github.com/gin-gonic/gin"
)

// graphQLRequest contains the GraphQL request data sent by Swagger or a client.
type graphQLRequest struct {
	Query string `json:"query"`
}

// newGraphQLGatewayHandler creates one GraphQL endpoint for both services.
func newGraphQLGatewayHandler(
	cfg config.Config,
	productProxy gin.HandlerFunc,
	orderProxy gin.HandlerFunc,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		query, err := getGraphQLQuery(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		serviceName, adminRequired, err := identifyGraphQLService(query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Every GraphQL operation requires authentication.
		roleValue, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authentication is required",
			})
			return
		}

		role, ok := roleValue.(string)
		if !ok || (role != "admin" && role != "user") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user role",
			})
			return
		}

		// Product mutations are restricted to Admin.
		if adminRequired && role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "only admin can execute Product mutations",
			})
			return
		}

		// Forward the original request to the correct existing service.
		switch serviceName {
		case "product":
			productProxy(c)
		case "order":
			orderProxy(c)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "unsupported GraphQL service",
			})
		}
	}
}

// getGraphQLQuery reads the query from a GET or POST GraphQL request.
func getGraphQLQuery(c *gin.Context) (string, error) {
	switch c.Request.Method {
	case http.MethodGet:
		query := strings.TrimSpace(c.Query("query"))
		if query == "" {
			return "", fmt.Errorf("the query parameter is required")
		}
		return query, nil

	case http.MethodPost:
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return "", fmt.Errorf("unable to read GraphQL request body")
		}

		// Restore the original body so the downstream service receives it.
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Request.ContentLength = int64(len(body))

		var request graphQLRequest
		if err := json.Unmarshal(body, &request); err != nil {
			return "", fmt.Errorf("request body must contain valid GraphQL JSON")
		}

		query := strings.TrimSpace(request.Query)
		if query == "" {
			return "", fmt.Errorf("the query field is required")
		}

		return query, nil

	default:
		return "", fmt.Errorf("only GET and POST methods are supported")
	}
}

// identifyGraphQLService determines which existing service owns the operation.
//
// Product operations:
//
//	Queries: product, products
//	Mutations: createProduct, updateProduct, deleteProduct
//
// Order operations:
//
//	Queries: orders, order
//	Mutations: createOrder, confirmOrder, cancelOrder
func identifyGraphQLService(query string) (string, bool, error) {
	tokens := graphQLNameTokens(query)

	productFields := map[string]bool{
		"product":       true,
		"products":      true,
		"createProduct": true,
		"updateProduct": true,
		"deleteProduct": true,
	}

	orderFields := map[string]bool{
		"orders":       true,
		"order":        true,
		"createOrder":  true,
		"confirmOrder": true,
		"cancelOrder":  true,
	}

	productFound := false
	orderFound := false
	adminRequired := false

	for _, token := range tokens {
		if productFields[token] {
			productFound = true
		}

		if orderFields[token] {
			orderFound = true
		}

		switch token {
		case "createProduct", "updateProduct", "deleteProduct":
			adminRequired = true
		}
	}

	if productFound && orderFound {
		return "", false, fmt.Errorf(
			"one GraphQL request cannot mix Product and Order operations; send separate requests to /graphql",
		)
	}

	if productFound {
		return "product", adminRequired, nil
	}

	if orderFound {
		return "order", false, nil
	}

	return "", false, fmt.Errorf(
		"unable to identify the service; use a supported Product or Order GraphQL operation",
	)
}

// graphQLNameTokens extracts GraphQL names while ignoring comments and strings.
//
// This prevents words inside comments and string values from being mistaken
// for GraphQL operation names.
func graphQLNameTokens(query string) []string {
	tokens := make([]string, 0)
	length := len(query)

	for i := 0; i < length; {
		current := query[i]

		// Ignore whitespace and commas.
		if isGraphQLWhitespace(current) || current == ',' {
			i++
			continue
		}

		// Ignore comments.
		if current == '#' {
			for i < length && query[i] != '\n' {
				i++
			}
			continue
		}

		// Ignore GraphQL block strings.
		if i+2 < length && query[i:i+3] == `"""` {
			i += 3
			for i+2 < length && query[i:i+3] != `"""` {
				i++
			}
			if i+2 < length {
				i += 3
			}
			continue
		}

		// Ignore normal quoted strings, including escaped characters.
		if current == '"' {
			i++
			for i < length {
				if query[i] == '\\' {
					i += 2
					continue
				}
				if query[i] == '"' {
					i++
					break
				}
				i++
			}
			continue
		}

		// Ignore variable names such as $product.
		if current == '$' {
			i++
			for i < length && isGraphQLNameCharacter(query[i]) {
				i++
			}
			continue
		}

		// Extract GraphQL names.
		if isGraphQLNameStart(current) {
			start := i
			i++

			for i < length && isGraphQLNameCharacter(query[i]) {
				i++
			}

			tokens = append(tokens, query[start:i])
			continue
		}

		i++
	}

	return tokens
}

func isGraphQLWhitespace(value byte) bool {
	return value == ' ' ||
		value == '\t' ||
		value == '\n' ||
		value == '\r'
}

func isGraphQLNameStart(value byte) bool {
	return value == '_' ||
		(value >= 'a' && value <= 'z') ||
		(value >= 'A' && value <= 'Z')
}

func isGraphQLNameCharacter(value byte) bool {
	return isGraphQLNameStart(value) ||
		(value >= '0' && value <= '9')
}

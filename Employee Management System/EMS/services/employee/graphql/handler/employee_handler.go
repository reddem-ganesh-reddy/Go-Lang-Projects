
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	graphqlHandler "github.com/graphql-go/handler"
)

type EmployeeGraphQLHandler struct {
	handler http.Handler
}

func NewEmployeeGraphQLHandler(
	employeeSchema graphql.Schema,
) *EmployeeGraphQLHandler {

	handler := graphqlHandler.New(&graphqlHandler.Config{
		Schema:   &employeeSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	return &EmployeeGraphQLHandler{
		handler: handler,
	}
}

// ServeHTTP handles GraphQL API and GraphiQL UI requests.
func (h *EmployeeGraphQLHandler) ServeHTTP(
	ctx *gin.Context,
) {
	h.handler.ServeHTTP(ctx.Writer, ctx.Request)
}

// RegisterEmployeeGraphQLRoutes registers GraphQL routes.
func RegisterEmployeeGraphQLRoutes(
	router *gin.Engine,
	graphqlHandler *EmployeeGraphQLHandler,
) {
	router.Any("/graphql", graphqlHandler.ServeHTTP)
}

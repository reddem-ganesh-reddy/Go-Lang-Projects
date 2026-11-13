package handler

import (
	"context"

	"EMS/services/department/graphql/resolvers"
	"EMS/services/department/graphql/schema"
	"EMS/shared/middleware"

	"github.com/gin-gonic/gin"
	graphqlhandler "github.com/graphql-go/handler"
)

// NewDepartmentGraphQLHandler creates the Department GraphQL HTTP handler.
func NewDepartmentGraphQLHandler(
	resolver *resolvers.DepartmentResolver,
) (gin.HandlerFunc, error) {

	// Create the Department GraphQL schema.
	departmentSchema, err := schema.NewDepartmentSchema(resolver)
	if err != nil {
		return nil, err
	}

	// Configure the GraphQL HTTP handler.
	graphQLHandler := graphqlhandler.New(&graphqlhandler.Config{
		Schema:   &departmentSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Pass the authenticated user role into the GraphQL request context.
	return func(c *gin.Context) {
		userRole := middleware.GetUserRole(c)

		ctx := context.WithValue(
			c.Request.Context(),
			resolvers.UserRoleContextKey,
			userRole,
		)

		request := c.Request.WithContext(ctx)
		graphQLHandler.ServeHTTP(c.Writer, request)
	}, nil
}

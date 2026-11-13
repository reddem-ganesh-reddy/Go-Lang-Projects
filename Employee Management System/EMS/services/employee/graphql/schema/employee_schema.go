
package schema

import (
	"EMS/services/employee/graphql/mutations"
	"EMS/services/employee/graphql/queries"
	"EMS/services/employee/graphql/resolvers"

	"github.com/graphql-go/graphql"
)

func NewEmployeeSchema(
	resolver *resolvers.EmployeeResolver,
) (graphql.Schema, error) {

	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "EmployeeQuery",
			Fields: queries.GetEmployeeQueries(resolver),
		},
	)

	mutationType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "EmployeeMutation",
			Fields: mutations.GetEmployeeMutations(resolver),
		},
	)

	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)
}

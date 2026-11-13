package schema

import (
	"EMS/services/department/graphql/mutations"
	"EMS/services/department/graphql/queries"
	"EMS/services/department/graphql/resolvers"

	"github.com/graphql-go/graphql"
)

// NewDepartmentSchema creates the department GraphQL schema.
func NewDepartmentSchema(
	resolver *resolvers.DepartmentResolver,
) (graphql.Schema, error) {

	query := queries.NewDepartmentQueries(resolver)
	mutation := mutations.NewDepartmentMutations(resolver)

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
}
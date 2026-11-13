
package queries

import (
	"EMS/services/employee/graphql/resolvers"
	"EMS/services/employee/graphql/inputs"

	"github.com/graphql-go/graphql"
)

func GetEmployeeQueries(resolver *resolvers.EmployeeResolver) graphql.Fields {
	return graphql.Fields{
		"employees": &graphql.Field{
			Type:        resolvers.EmployeeListType,
			Description: "Get employees with optional pagination, filtering and sorting",
			Args: graphql.FieldConfigArgument{
				"query": &graphql.ArgumentConfig{
					Type: inputs.EmployeeListQueryInput,
				},
			},
			Resolve: resolver.GetAll,
		},

		"employee": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.NewObject(graphql.ObjectConfig{Name: "EmployeeQueryResult", Fields: graphql.Fields{"id": &graphql.Field{Type: graphql.Int}, "name": &graphql.Field{Type: graphql.String}, "email": &graphql.Field{Type: graphql.String}, "department_id": &graphql.Field{Type: graphql.Int}, "created_at": &graphql.Field{Type: graphql.String}, "updated_at": &graphql.Field{Type: graphql.String}}})),
			Description: "Get an employee by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: resolver.GetByID,
		},
	}
}

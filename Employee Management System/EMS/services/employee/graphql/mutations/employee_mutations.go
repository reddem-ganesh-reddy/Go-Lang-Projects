
package mutations

import (
	"EMS/services/employee/graphql/inputs"
	"EMS/services/employee/graphql/resolvers"
	"EMS/services/employee/graphql/types"

	"github.com/graphql-go/graphql"
)

func GetEmployeeMutations(
	resolver *resolvers.EmployeeResolver,
) graphql.Fields {
	return graphql.Fields{
		"createEmployee": &graphql.Field{
			Type:        types.EmployeeType,
			Description: "Create a new employee",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(inputs.CreateEmployeeInput),
				},
			},
			Resolve: resolver.Create,
		},

		"updateEmployee": &graphql.Field{
			Type:        types.EmployeeType,
			Description: "Update an existing employee",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(inputs.UpdateEmployeeInput),
				},
			},
			Resolve: resolver.Update,
		},

		"deleteEmployee": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Delete an employee",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: resolver.Delete,
		},
	}
}

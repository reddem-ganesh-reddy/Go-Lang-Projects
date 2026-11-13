package queries

import (
	"EMS/services/department/graphql/resolvers"
	"EMS/services/department/graphql/types"

	"github.com/graphql-go/graphql"
)

// NewDepartmentQueries creates department GraphQL queries.
func NewDepartmentQueries(
	resolver *resolvers.DepartmentResolver,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "DepartmentQuery",
		Fields: graphql.Fields{
			"departments": &graphql.Field{
				Type: graphql.NewList(types.DepartmentType),
				Description: "Get all departments",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.GetAll()
				},
			},

			"department": &graphql.Field{
				Type: types.DepartmentType,
				Description: "Get department by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(string)
					return resolver.GetByID(id)
				},
			},
		},
	})
}
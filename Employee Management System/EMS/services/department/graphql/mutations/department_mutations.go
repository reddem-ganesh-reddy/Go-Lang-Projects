package mutations

import (
	"errors"

	"EMS/services/department/graphql/inputs"
	"EMS/services/department/graphql/resolvers"
	"EMS/services/department/graphql/types"

	"github.com/graphql-go/graphql"
)

// NewDepartmentMutations creates Department GraphQL mutations.
func NewDepartmentMutations(
	resolver *resolvers.DepartmentResolver,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "DepartmentMutation",
		Fields: graphql.Fields{
			"createDepartment": &graphql.Field{
				Type:        types.DepartmentType,
				Description: "Create a department",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(inputs.CreateDepartmentInput),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Read and validate department input.
					input, ok := p.Args["input"].(map[string]interface{})
					if !ok {
						return nil, errors.New("invalid department input")
					}

					name, ok := input["name"].(string)
					if !ok {
						return nil, errors.New("department name is required")
					}

					return resolver.Create(name)
				},
			},

			"updateDepartment": &graphql.Field{
				Type:        types.DepartmentType,
				Description: "Update a department",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(inputs.UpdateDepartmentInput),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Read and validate department ID.
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, errors.New("invalid department ID")
					}

					// Read and validate department input.
					input, ok := p.Args["input"].(map[string]interface{})
					if !ok {
						return nil, errors.New("invalid department input")
					}

					name, ok := input["name"].(string)
					if !ok {
						return nil, errors.New("department name is required")
					}

					return resolver.Update(id, name)
				},
			},

			"deleteDepartment": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "Delete a department",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Read and validate department ID.
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, errors.New("invalid department ID")
					}

					return resolver.Delete(id)
				},
			},
		},
	})
}

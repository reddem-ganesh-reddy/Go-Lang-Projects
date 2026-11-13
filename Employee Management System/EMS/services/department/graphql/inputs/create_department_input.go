package inputs

import "github.com/graphql-go/graphql"

// CreateDepartmentInput defines department creation input.
var CreateDepartmentInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "CreateDepartmentInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)
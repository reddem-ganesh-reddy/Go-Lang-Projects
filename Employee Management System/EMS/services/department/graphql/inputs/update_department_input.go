package inputs

import "github.com/graphql-go/graphql"

// UpdateDepartmentInput defines department update input.
var UpdateDepartmentInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "UpdateDepartmentInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)
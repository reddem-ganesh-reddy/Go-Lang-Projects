
package inputs

import "github.com/graphql-go/graphql"

var CreateEmployeeInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "CreateEmployeeInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"email": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"department_id": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
	},
)

var UpdateEmployeeInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "UpdateEmployeeInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"email": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"department_id": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
	},
)

var EmployeeListQueryInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "EmployeeListQueryInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"page": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"limit": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"email": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"department_id": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"sort_by": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"sort_order": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)

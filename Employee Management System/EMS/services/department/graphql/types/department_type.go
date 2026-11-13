package types

import (
	"strconv"

	"EMS/services/department/dto"

	"github.com/graphql-go/graphql"
)

// DepartmentType defines the GraphQL Department object.
var DepartmentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Department",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				department := p.Source.(dto.DepartmentResponse)
				return strconv.FormatUint(uint64(department.ID), 10), nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				department := p.Source.(dto.DepartmentResponse)
				return department.Name, nil
			},
		},
	},
})
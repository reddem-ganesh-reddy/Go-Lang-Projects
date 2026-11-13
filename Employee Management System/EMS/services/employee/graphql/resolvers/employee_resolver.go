
package resolvers

import (
	"fmt"

	"EMS/services/employee/dto"
	"EMS/services/employee/graphql/types"
	"EMS/services/employee/service"

	"github.com/graphql-go/graphql"
)

type EmployeeResolver struct {
	service service.EmployeeService
}

func NewEmployeeResolver(employeeService service.EmployeeService) *EmployeeResolver {
	return &EmployeeResolver{
		service: employeeService,
	}
}

// GraphQL response type for pagination.
var PaginationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "EmployeePagination",
		Fields: graphql.Fields{
			"total_records": &graphql.Field{
				Type: graphql.Int,
			},
			"current_page": &graphql.Field{
				Type: graphql.Int,
			},
			"page_size": &graphql.Field{
				Type: graphql.Int,
			},
			"total_pages": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

// GraphQL response type for employee list.
var EmployeeListType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "EmployeeList",
		Fields: graphql.Fields{
			"employees": &graphql.Field{
				Type: graphql.NewList(types.EmployeeType),
			},
			"pagination": &graphql.Field{
				Type: PaginationType,
			},
		},
	},
)

// Get all employees with optional filters and pagination.
func (r *EmployeeResolver) GetAll(
	params graphql.ResolveParams,
) (interface{}, error) {

	query := dto.EmployeeListQuery{}

	if value, ok := params.Args["query"].(map[string]interface{}); ok {
		if v, exists := value["page"]; exists {
			query.Page = v.(int)
		}

		if v, exists := value["limit"]; exists {
			query.Limit = v.(int)
		}

		if v, exists := value["name"]; exists {
			query.Name = v.(string)
		}

		if v, exists := value["email"]; exists {
			query.Email = v.(string)
		}

		if v, exists := value["department_id"]; exists {
			query.DepartmentID = uint(v.(int))
		}

		if v, exists := value["sort_by"]; exists {
			query.SortBy = v.(string)
		}

		if v, exists := value["sort_order"]; exists {
			query.SortOrder = v.(string)
		}
	}

	return r.service.GetAll(query)
}

// Get an employee by ID.
func (r *EmployeeResolver) GetByID(
	params graphql.ResolveParams,
) (interface{}, error) {

	id, ok := params.Args["id"].(int)
	if !ok || id <= 0 {
		return nil, fmt.Errorf("invalid employee ID")
	}

	return r.service.GetByID(uint(id))
}

// Create an employee.
func (r *EmployeeResolver) Create(
	params graphql.ResolveParams,
) (interface{}, error) {

	input, ok := params.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("employee input is required")
	}

	request := dto.CreateEmployeeRequest{
		Name:         input["name"].(string),
		Email:        input["email"].(string),
		DepartmentID: uint(input["department_id"].(int)),
	}

	return r.service.Create(request)
}

// Update an employee.
func (r *EmployeeResolver) Update(
	params graphql.ResolveParams,
) (interface{}, error) {

	id, ok := params.Args["id"].(int)
	if !ok || id <= 0 {
		return nil, fmt.Errorf("invalid employee ID")
	}

	input, ok := params.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("employee input is required")
	}

	request := dto.UpdateEmployeeRequest{
		Name:         input["name"].(string),
		Email:        input["email"].(string),
		DepartmentID: uint(input["department_id"].(int)),
	}

	return r.service.Update(uint(id), request)
}

// Delete an employee.
func (r *EmployeeResolver) Delete(
	params graphql.ResolveParams,
) (interface{}, error) {

	id, ok := params.Args["id"].(int)
	if !ok || id <= 0 {
		return nil, fmt.Errorf("invalid employee ID")
	}

	if err := r.service.Delete(uint(id)); err != nil {
		return nil, err
	}

	return true, nil
}

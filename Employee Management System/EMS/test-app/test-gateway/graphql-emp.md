# Employee Service – GraphQL API Testing

## 1. Testing Details

| Item | Details |
|---|---|
| Project | Employee Management System |
| Service | Employee Service |
| Testing Method | Gateway Swagger UI |
| API Gateway URL | `http://localhost:8080` |
| GraphQL Endpoint | `POST /graphql` |
| Content-Type | `application/json` |

### Swagger Testing Procedure

1. Start Consul and all four services.
2. Open `http://localhost:8080/swagger/index.html`.
3. Locate `POST /graphql`.
4. Click **Try it out**.
5. Click **Authorize** and provide the JWT token.
6. Enter the required JSON request body.
7. Click **Execute**.
8. Verify the HTTP status code and GraphQL response.

All GraphQL requests must use JSON request bodies.

### Authentication Requirements

- Admin token: Required for CREATE, UPDATE and DELETE operations.
- Employee token: Permitted for READ queries.
- Missing or invalid token: Expected to be rejected by Gateway authentication.

---

## 2. Employee GraphQL Schema Reference

### Employee Fields

- `id`
- `name`
- `email`
- `department_id`
- `created_at`
- `updated_at`

### Employee Pagination Fields

- `total_records`
- `current_page`
- `page_size`
- `total_pages`

### EmployeeListQueryInput Fields

| Field | Type | Purpose |
|---|---|---|
| page | Int | Page number |
| limit | Int | Records per page |
| name | String | Filter by employee name |
| email | String | Filter by email |
| department_id | Int | Filter by department |
| sort_by | String | Sorting field |
| sort_order | String | ASC or DESC |

### Default Values

| Field | Default |
|---|---|
| page | 1 |
| limit | 10 |
| sort_by | id |
| sort_order | asc |

---

## 3. READ – Basic Employee Queries

### E01 – Get All Employees

```json
{
  "query": "query { employees { employees { id name email department_id created_at updated_at } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Employee records and pagination information.

### E02 – Get Employee by ID

```json
{
  "query": "query { employee(id: 1) { id name email department_id created_at updated_at } }"
}
```

Expected: Details of employee ID 1.

### E03 – Get Employee with Non-Existing ID

```json
{
  "query": "query { employee(id: 99999) { id name email department_id } }"
}
```

Expected: Null or a not-found error, depending on resolver behavior.

### E04 – Get Employee with Invalid ID

```json
{
  "query": "query { employee(id: -1) { id name email department_id } }"
}
```

Expected: Invalid ID error.

---

## 4. PAGINATION TEST CASES

### E05 – First Page

```json
{
  "query": "query { employees(query: { page: 1, limit: 5 }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: First five matching records, if available.

### E06 – Second Page

```json
{
  "query": "query { employees(query: { page: 2, limit: 5 }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Next page of records without repeating the first page.

### E07 – Third Page

```json
{
  "query": "query { employees(query: { page: 3, limit: 5 }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Third page of records, if available.

### E08 – Different Page Size

```json
{
  "query": "query { employees(query: { page: 1, limit: 10 }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Up to ten records.

### E09 – Large Page Size

```json
{
  "query": "query { employees(query: { page: 1, limit: 100 }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Verify whether the service accepts or caps the requested page size.

### E10 – Page Beyond Available Records

```json
{
  "query": "query { employees(query: { page: 999, limit: 5 }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Empty employee list with pagination metadata, or the application's defined validation response.

### E11 – Page Number Zero

```json
{
  "query": "query { employees(query: { page: 0, limit: 5 }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Verify defaulting to page 1 or validation rejection.

### E12 – Negative Page Number

```json
{
  "query": "query { employees(query: { page: -1, limit: 5 }) { employees { id name email } } }"
}
```

Expected: Validation rejection or defined normalization behavior.

### E13 – Zero Page Size

```json
{
  "query": "query { employees(query: { page: 1, limit: 0 }) { employees { id name email } } }"
}
```

Expected: Verify defaulting or validation rejection.

### E14 – Negative Page Size

```json
{
  "query": "query { employees(query: { page: 1, limit: -5 }) { employees { id name email } } }"
}
```

Expected: Validation rejection or defined normalization behavior.

---

## 5. FILTERING TEST CASES

Use values available in the employee database.

### E15 – Filter by Employee Name

```json
{
  "query": "query { employees(query: { name: \"Ganesh\" }) { employees { id name email department_id } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Employees matching the name filter.

### E16 – Filter by Email

```json
{
  "query": "query { employees(query: { email: \"test@example.com\" }) { employees { id name email department_id } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Employees matching the email filter.

### E17 – Filter by Department ID

```json
{
  "query": "query { employees(query: { department_id: 1 }) { employees { id name email department_id } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Employees belonging to department 1.

### E18 – Filter by Non-Existing Name

```json
{
  "query": "query { employees(query: { name: \"NoSuchEmployee\" }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Empty employee list.

### E19 – Filter by Non-Existing Email

```json
{
  "query": "query { employees(query: { email: \"notfound@example.com\" }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Empty employee list.

### E20 – Filter by Non-Existing Department

```json
{
  "query": "query { employees(query: { department_id: 99999 }) { employees { id name email department_id } } }"
}
```

Expected: Empty employee list.

### E21 – Filter by Name and Email

```json
{
  "query": "query { employees(query: { name: \"Ganesh\", email: \"test@example.com\" }) { employees { id name email } } }"
}
```

Expected: Records matching both filters.

### E22 – Filter by Name and Department

```json
{
  "query": "query { employees(query: { name: \"Ganesh\", department_id: 1 }) { employees { id name email department_id } } }"
}
```

Expected: Employees matching the name filter in department 1.

### E23 – Filter by Email and Department

```json
{
  "query": "query { employees(query: { email: \"test@example.com\", department_id: 1 }) { employees { id name email department_id } } }"
}
```

Expected: Employees matching the email filter in department 1.

### E24 – Filter by All Three Fields

```json
{
  "query": "query { employees(query: { name: \"Ganesh\", email: \"test@example.com\", department_id: 1 }) { employees { id name email department_id } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Records matching all supplied filters.

---

## 6. SORTING TEST CASES

Use supported sorting fields and verify the returned order.

### E25 – Sort by ID Ascending

```json
{
  "query": "query { employees(query: { sort_by: \"id\", sort_order: \"asc\" }) { employees { id name email } } }"
}
```

Expected: Employee IDs in ascending order.

### E26 – Sort by ID Descending

```json
{
  "query": "query { employees(query: { sort_by: \"id\", sort_order: \"desc\" }) { employees { id name email } } }"
}
```

Expected: Employee IDs in descending order.

### E27 – Sort by Name Ascending

```json
{
  "query": "query { employees(query: { sort_by: \"name\", sort_order: \"asc\" }) { employees { id name email } } }"
}
```

Expected: Names in ascending order.

### E28 – Sort by Name Descending

```json
{
  "query": "query { employees(query: { sort_by: \"name\", sort_order: \"desc\" }) { employees { id name email } } }"
}
```

Expected: Names in descending order.

### E29 – Sort by Email Ascending

```json
{
  "query": "query { employees(query: { sort_by: \"email\", sort_order: \"asc\" }) { employees { id name email } } }"
}
```

Expected: Emails in ascending order.

### E30 – Sort by Email Descending

```json
{
  "query": "query { employees(query: { sort_by: \"email\", sort_order: \"desc\" }) { employees { id name email } } }"
}
```

Expected: Emails in descending order.

### E31 – Invalid Sorting Field

```json
{
  "query": "query { employees(query: { sort_by: \"invalid_field\", sort_order: \"asc\" }) { employees { id name email } } }"
}
```

Expected: Validation error or defined fallback behavior.

### E32 – Invalid Sorting Order

```json
{
  "query": "query { employees(query: { sort_by: \"name\", sort_order: \"invalid\" }) { employees { id name email } } }"
}
```

Expected: Validation error or defined fallback behavior.

### E33 – Uppercase Sorting Order

```json
{
  "query": "query { employees(query: { sort_by: \"name\", sort_order: \"ASC\" }) { employees { id name email } } }"
}
```

Expected: Verify whether uppercase ASC is accepted.

---

## 7. COMBINATION TEST CASES

These tests verify pagination, filtering and sorting together.

### E34 – Pagination + Name Filter

```json
{
  "query": "query { employees(query: { page: 1, limit: 5, name: \"Ganesh\" }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

### E35 – Pagination + Email Filter

```json
{
  "query": "query { employees(query: { page: 1, limit: 5, email: \"test@example.com\" }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

### E36 – Pagination + Department Filter

```json
{
  "query": "query { employees(query: { page: 1, limit: 5, department_id: 1 }) { employees { id name email department_id } pagination { total_records current_page page_size total_pages } } }"
}
```

### E37 – Pagination + Sorting

```json
{
  "query": "query { employees(query: { page: 1, limit: 5, sort_by: \"name\", sort_order: \"asc\" }) { employees { id name email } pagination { total_records current_page page_size total_pages } } }"
}
```

### E38 – Name Filter + Sorting

```json
{
  "query": "query { employees(query: { name: \"Ganesh\", sort_by: \"name\", sort_order: \"asc\" }) { employees { id name email } } }"
}
```

### E39 – Department Filter + Sorting

```json
{
  "query": "query { employees(query: { department_id: 1, sort_by: \"id\", sort_order: \"desc\" }) { employees { id name email department_id } } }"
}
```

### E40 – Email Filter + Sorting

```json
{
  "query": "query { employees(query: { email: \"test@example.com\", sort_by: \"email\", sort_order: \"asc\" }) { employees { id name email } } }"
}
```

### E41 – Multiple Filters + Pagination

```json
{
  "query": "query { employees(query: { page: 1, limit: 5, name: \"Ganesh\", department_id: 1 }) { employees { id name email department_id } pagination { total_records current_page page_size total_pages } } }"
}
```

### E42 – Multiple Filters + Sorting

```json
{
  "query": "query { employees(query: { name: \"Ganesh\", department_id: 1, sort_by: \"name\", sort_order: \"desc\" }) { employees { id name email department_id } } }"
}
```

### E43 – All Filters + Pagination + Sorting

```json
{
  "query": "query { employees(query: { page: 1, limit: 5, name: \"Ganesh\", email: \"test@example.com\", department_id: 1, sort_by: \"name\", sort_order: \"asc\" }) { employees { id name email department_id } pagination { total_records current_page page_size total_pages } } }"
}
```

### E44 – Second Page with Filters and Sorting

```json
{
  "query": "query { employees(query: { page: 2, limit: 5, department_id: 1, sort_by: \"id\", sort_order: \"desc\" }) { employees { id name email department_id } pagination { total_records current_page page_size total_pages } } }"
}
```

Expected: Apply filters first, sort the matching records, and return the requested page. Verify pagination totals against the filtered result set.

---

## 8. CREATE – Employee Mutation

Use an existing department ID.

### E45 – Create Employee Successfully

```json
{
  "query": "mutation { createEmployee(input: { name: \"Test Employee\", email: \"testemployee@example.com\", department_id: 1 }) { id name email department_id created_at updated_at } }"
}
```

Expected: Employee created successfully.

Record the generated employee ID for subsequent update and delete tests.

### E46 – Create Employee with Missing Name

```json
{
  "query": "mutation { createEmployee(input: { email: \"missingname@example.com\", department_id: 1 }) { id name } }"
}
```

Expected: GraphQL validation error because name is required.

### E47 – Create Employee with Missing Email

```json
{
  "query": "mutation { createEmployee(input: { name: \"Test Employee\", department_id: 1 }) { id name } }"
}
```

Expected: GraphQL validation error.

### E48 – Create Employee with Missing Department

```json
{
  "query": "mutation { createEmployee(input: { name: \"Test Employee\", email: \"missingdept@example.com\" }) { id name } }"
}
```

Expected: GraphQL validation error.

### E49 – Create Employee with Empty Name

```json
{
  "query": "mutation { createEmployee(input: { name: \"\", email: \"emptyname@example.com\", department_id: 1 }) { id name } }"
}
```

Expected: Application validation error.

### E50 – Create Employee with Invalid Email

```json
{
  "query": "mutation { createEmployee(input: { name: \"Test Employee\", email: \"invalid-email\", department_id: 1 }) { id name } }"
}
```

Expected: Email validation error, if email format validation is implemented.

### E51 – Create Employee with Non-Existing Department

```json
{
  "query": "mutation { createEmployee(input: { name: \"Invalid Department Employee\", email: \"invaliddept@example.com\", department_id: 99999 }) { id name } }"
}
```

Expected: Department not found. Employee record must not be created.

### E52 – Create Employee with Duplicate Email

Use an email address that already exists.

```json
{
  "query": "mutation { createEmployee(input: { name: \"Duplicate Employee\", email: \"testemployee@example.com\", department_id: 1 }) { id name email } }"
}
```

Expected: Duplicate email rejected if uniqueness validation or a database unique constraint is configured.

---

## 9. UPDATE – Employee Mutation

Use the employee ID generated by E45.

### E53 – Update Employee Successfully

```json
{
  "query": "mutation { updateEmployee(id: 1, input: { name: \"Updated Employee\", email: \"updatedemployee@example.com\", department_id: 1 }) { id name email department_id updated_at } }"
}
```

Expected: Employee details updated successfully.

### E54 – Update Employee with Non-Existing ID

```json
{
  "query": "mutation { updateEmployee(id: 99999, input: { name: \"Updated Employee\", email: \"notfound@example.com\", department_id: 1 }) { id name } }"
}
```

Expected: Employee not found.

### E55 – Update Employee with Invalid ID

```json
{
  "query": "mutation { updateEmployee(id: -1, input: { name: \"Updated Employee\", email: \"updated@example.com\", department_id: 1 }) { id name } }"
}
```

Expected: Invalid ID error.

### E56 – Update Employee with Empty Name

```json
{
  "query": "mutation { updateEmployee(id: 1, input: { name: \"\", email: \"updated@example.com\", department_id: 1 }) { id name } }"
}
```

Expected: Validation error.

### E57 – Update Employee with Invalid Email

```json
{
  "query": "mutation { updateEmployee(id: 1, input: { name: \"Updated Employee\", email: \"invalid-email\", department_id: 1 }) { id name } }"
}
```

Expected: Email validation error, if implemented.

### E58 – Update Employee with Non-Existing Department

```json
{
  "query": "mutation { updateEmployee(id: 1, input: { name: \"Updated Employee\", email: \"updated@example.com\", department_id: 99999 }) { id name } }"
}
```

Expected: Department not found. Existing employee department must remain unchanged.

---

## 10. DELETE – Employee Mutation

Use test employee records to avoid deleting existing project data.

### E59 – Delete Employee Successfully

```json
{
  "query": "mutation { deleteEmployee(id: 1) }"
}
```

Expected response:

```json
{
  "data": {
    "deleteEmployee": true
  }
}
```

### E60 – Delete Non-Existing Employee

```json
{
  "query": "mutation { deleteEmployee(id: 99999) }"
}
```

Expected: Not-found error or the application's defined deletion response.

### E61 – Delete Employee with Invalid ID

```json
{
  "query": "mutation { deleteEmployee(id: -1) }"
}
```

Expected: Invalid ID error.

### E62 – Verify Deleted Employee

```json
{
  "query": "query { employee(id: 1) { id name email department_id } }"
}
```

Expected: Null or not-found error.

---

## 11. AUTHORIZATION TEST CASES

Execute all authorization tests through Gateway Swagger UI.

| Test ID | Token | Operation | Expected Result |
|---|---|---|---|
| E63 | No token | Get all employees | 401 Unauthorized |
| E64 | No token | Get employee by ID | 401 Unauthorized |
| E65 | No token | Create employee | 401 Unauthorized |
| E66 | Employee token | Get all employees | Allowed |
| E67 | Employee token | Get employee by ID | Allowed |
| E68 | Employee token | Create employee | 403 Forbidden |
| E69 | Employee token | Update employee | 403 Forbidden |
| E70 | Employee token | Delete employee | 403 Forbidden |
| E71 | Admin token | Get all employees | Allowed |
| E72 | Admin token | Get employee by ID | Allowed |
| E73 | Admin token | Create employee | Allowed |
| E74 | Admin token | Update employee | Allowed |
| E75 | Admin token | Delete employee | Allowed |
| E76 | Invalid token | Any operation | 401 Unauthorized |

---

## 12. FINAL TESTING CHECKLIST

### Basic Queries

- [ ] Get all employees
- [ ] Get employee by ID
- [ ] Get non-existing employee
- [ ] Test invalid employee ID

### Pagination

- [ ] First page
- [ ] Second page
- [ ] Third page
- [ ] Different page sizes
- [ ] Large page size
- [ ] Page beyond available records
- [ ] Page number zero
- [ ] Negative page number
- [ ] Zero page size
- [ ] Negative page size

### Filtering

- [ ] Filter by name
- [ ] Filter by email
- [ ] Filter by department
- [ ] Filter with no matching results
- [ ] Name + email
- [ ] Name + department
- [ ] Email + department
- [ ] All filters together

### Sorting

- [ ] ID ascending
- [ ] ID descending
- [ ] Name ascending
- [ ] Name descending
- [ ] Email ascending
- [ ] Email descending
- [ ] Invalid sorting field
- [ ] Invalid sorting order
- [ ] Uppercase sorting order

### Combination Tests

- [ ] Pagination + name filter
- [ ] Pagination + email filter
- [ ] Pagination + department filter
- [ ] Pagination + sorting
- [ ] Filtering + sorting
- [ ] Multiple filters + pagination
- [ ] Multiple filters + sorting
- [ ] All filters + pagination + sorting

### CRUD and Validation

- [ ] Create employee
- [ ] Update employee
- [ ] Delete employee
- [ ] Verify deletion
- [ ] Missing required fields
- [ ] Empty employee name
- [ ] Invalid email
- [ ] Duplicate email
- [ ] Invalid department
- [ ] Non-existing employee ID

### Authorization

- [ ] Requests without a token
- [ ] Requests with an invalid token
- [ ] Employee read access
- [ ] Employee mutation restrictions
- [ ] Admin CRUD access

---

## 13. TESTING OBSERVATIONS

Record the following details for every test case:

1. Test ID
2. Request JSON
3. HTTP status code
4. GraphQL response
5. Expected result
6. Actual result
7. PASS or FAIL

### Important Notes

- GraphQL may return HTTP 200 even when the response contains an error. Check both the HTTP status code and the `errors` field.
- Missing required GraphQL fields may be rejected before resolver execution.
- Empty strings, invalid email addresses, duplicate email addresses and invalid department IDs require application-level validation.
- Verify pagination, filtering and sorting against actual database records.
- Test invalid pagination values and unsupported sorting fields to establish the service's actual behavior.
- Use test employee records and valid department IDs.
- Avoid deleting existing project data.
- Record actual results instead of assuming that every expected validation response is implemented.

---

**End of Employee Service GraphQL API Testing Documentation.**
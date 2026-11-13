# Department Service – GraphQL API Testing

## 1. Testing Details

| Item | Details |
|---|---|
| Project | Employee Management System |
| Service | Department Service |
| Testing Method | Gateway Swagger UI |
| API Gateway | `http://localhost:8080` |
| GraphQL Endpoint | `POST /graphql` |
| HTTP Method | POST |
| Content-Type | `application/json` |

### Swagger Testing Procedure

1. Start Consul and all four services.
2. Open `http://localhost:8080/swagger/index.html`.
3. Locate `POST /graphql`.
4. Click **Try it out**.
5. Click **Authorize** and provide the admin JWT token.
6. Enter the JSON request body.
7. Click **Execute**.
8. Verify the HTTP status code and GraphQL response.

All GraphQL requests must use JSON request bodies.

### Authorization

- Admin token: Required for CREATE, UPDATE and DELETE operations.
- Employee token: Permitted for READ queries.
- Missing or invalid token: Expected to be rejected by Gateway authentication.

---

## 2. Department GraphQL Reference

### Department Fields

- `id`
- `name`

### GraphQL Operations

| Operation | GraphQL Field | Access |
|---|---|---|
| Create | `createDepartment` | Admin |
| Get all | `departments` | Admin and Employee |
| Get by ID | `department` | Admin and Employee |
| Update | `updateDepartment` | Admin |
| Delete | `deleteDepartment` | Admin |

---

## 3. CREATE – Department Mutation

### D01 – Create Department Successfully

**Authorization:** Admin token

**Request:**

```json
{
  "query": "mutation { createDepartment(input: { name: \"Testing\" }) { id name } }"
}
```

**Expected response (example):**

```json
{
  "data": {
    "createDepartment": {
      "id": "1",
      "name": "Testing"
    }
  }
}
```

Expected: Department created successfully.

Record the generated department ID for subsequent tests. The ID shown above is only an example.

---

## 4. READ – Department Queries

### D02 – Get All Departments

**Request:**

```json
{
  "query": "query { departments { id name } }"
}
```

Expected: List of departments, including the department created in D01.

### D03 – Get Department by ID

Replace `1` with the department ID returned by D01.

**Request:**

```json
{
  "query": "query { department(id: \"1\") { id name } }"
}
```

Expected: Details of the selected department.

---

## 5. UPDATE – Department Mutation

### D04 – Update Department Successfully

**Authorization:** Admin token

Replace `1` with the department ID returned by D01.

**Request:**

```json
{
  "query": "mutation { updateDepartment(id: \"1\", input: { name: \"Updated Testing\" }) { id name } }"
}
```

**Expected response (example):**

```json
{
  "data": {
    "updateDepartment": {
      "id": "1",
      "name": "Updated Testing"
    }
  }
}
```

Expected: Department name updated successfully.

---

## 6. DELETE – Department Mutation

### D05 – Delete Department Successfully

**Authorization:** Admin token

Use a test department that is not assigned to any employee.

Replace `1` with the test department ID.

**Request:**

```json
{
  "query": "mutation { deleteDepartment(id: \"1\") }"
}
```

**Expected response:**

```json
{
  "data": {
    "deleteDepartment": true
  }
}
```

Expected: Department deleted successfully.

### D06 – Verify Department Deletion

**Request:**

```json
{
  "query": "query { department(id: \"1\") { id name } }"
}
```

Expected: Null or a not-found error, depending on resolver behavior.

---

## 7. NEGATIVE TEST CASES

### D07 – Create Department with Empty Name

**Authorization:** Admin token

**Request:**

```json
{
  "query": "mutation { createDepartment(input: { name: \"\" }) { id name } }"
}
```

Expected: Application validation error, if empty-name validation is implemented.

### D08 – Create Department with Missing Name

**Authorization:** Admin token

**Request:**

```json
{
  "query": "mutation { createDepartment(input: {}) { id name } }"
}
```

Expected: GraphQL validation error if the name field is required in the input schema.

### D09 – Get Department with Invalid ID

**Request:**

```json
{
  "query": "query { department(id: \"abc\") { id name } }"
}
```

Expected: Invalid ID error or not-found response, depending on ID handling.

### D10 – Get Department with Non-Existing ID

**Request:**

```json
{
  "query": "query { department(id: \"99999\") { id name } }"
}
```

Expected: Null or not-found error.

### D11 – Create Department with Duplicate Name

**Authorization:** Admin token

**Request:**

```json
{
  "query": "mutation { createDepartment(input: { name: \"Testing\" }) { id name } }"
}
```

Expected: Duplicate department name rejected if uniqueness validation or a database constraint is configured. Otherwise, record the actual behavior.

---

## 8. AUTHORIZATION TEST CASES

Execute the following tests through Gateway Swagger UI.

| Test ID | Token | Operation | Expected Result |
|---|---|---|---|
| D12 | No token | Get all departments | 401 Unauthorized |
| D13 | No token | Get department by ID | 401 Unauthorized |
| D14 | No token | Create department | 401 Unauthorized |
| D15 | Employee token | Get all departments | Allowed |
| D16 | Employee token | Get department by ID | Allowed |
| D17 | Employee token | Create department | 403 Forbidden |
| D18 | Employee token | Update department | 403 Forbidden |
| D19 | Employee token | Delete department | 403 Forbidden |
| D20 | Admin token | Get all departments | Allowed |
| D21 | Admin token | Get department by ID | Allowed |
| D22 | Admin token | Create department | Allowed |
| D23 | Admin token | Update department | Allowed |
| D24 | Admin token | Delete department | Allowed |
| D25 | Invalid token | Any operation | 401 Unauthorized |

---

## 9. FINAL TESTING CHECKLIST

### CREATE

- [ ] Create department successfully
- [ ] Create department with empty name
- [ ] Create department with missing name
- [ ] Create department with duplicate name

### READ

- [ ] Get all departments
- [ ] Get department by ID
- [ ] Get department with invalid ID
- [ ] Get department with non-existing ID

### UPDATE

- [ ] Update department successfully
- [ ] Verify updated department name

### DELETE

- [ ] Delete test department
- [ ] Verify department deletion

### AUTHORIZATION

- [ ] Requests without a token
- [ ] Requests with an invalid token
- [ ] Employee read access
- [ ] Employee mutation restrictions
- [ ] Admin CRUD access

---

## 10. TESTING OBSERVATIONS

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
- Missing required fields may be rejected before resolver execution.
- Empty names and duplicate department names require application-level validation or database constraints.
- Use the department ID returned by D01 for subsequent CRUD tests.
- Avoid deleting departments assigned to existing employees.
- Employee Service depends on Department Service to validate department IDs during employee creation and updates.
- Record actual behavior instead of assuming that every expected validation response is implemented.

---

**End of Department Service GraphQL API Testing Documentation.**
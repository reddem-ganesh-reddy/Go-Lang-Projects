# Admin Testing – Product and Inventory Services

## 1. Overview

This document explains how to test the Product Service and Inventory Service using Swagger UI and GraphQL through the API Gateway.

All requests will be executed using an Admin account.

**Base URL:**

```text
http://localhost:8080
```

**GraphQL Endpoint:**

```text
POST http://localhost:8080/graphql
```

**Testing Tool:** Swagger UI / GraphQL Playground

**Authorization:** Admin JWT Token

---

## 2. Start the Application

Open PowerShell in the project root directory and execute:

```powershell
.\start-all.ps1
```

Wait until all services start successfully.

Open API Gateway Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

Verify that the API Gateway documentation loads successfully.

---

## 3. Admin Authentication

### Test 1: Admin Login

**Method:** POST

**Endpoint:**

```text
/api/auth/login
```

**Request Body:**

```json
{
  "username": "admin",
  "password": "admin@123"
}
```

Use the credentials configured in the `.env` file if they differ.

**Expected Result:**

- HTTP Status: `200 OK`
- Admin login should be successful.
- Response should contain a JWT access token.
- User role should be `admin`.

Copy the generated JWT token.

### Authorize Swagger

1. Click the **Authorize** button in Swagger UI.
2. Enter the following value:

```text
Bearer YOUR_ADMIN_JWT_TOKEN
```

3. Click Authorize.
4. Click Close.

All subsequent protected requests will use the Admin token.

---

# 4. Product Service – GraphQL Testing

GraphQL URL:

```text
POST http://localhost:8080/graphql
```

Headers:

```text
Content-Type: application/json
Authorization: Bearer YOUR_ADMIN_JWT_TOKEN
```

### Test 1: Get All Products

**Operation:** Query – `products`

**Request Body:**

```json
{
  "query": "{ products { items { id name description price stock } page limit totalRecords totalPages } }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- All available products should be returned.
- Each product should contain its ID, name, description, price and stock.
- Pagination details should be displayed.

### Test 2: Get Product by ID

**Operation:** Query – `product`

**Request Body:**

```json
{
  "query": "{ product(id: \"1\") { id name description price stock createdAt updatedAt } }"
}
```

Replace `1` with an existing product ID.

**Expected Result:**

- HTTP Status: `200 OK`
- The requested product should be returned.
- All product details should match the database.

### Test 3: Create Product

**Operation:** Mutation – `createProduct`

**Request Body:**

```json
{
  "query": "mutation { createProduct(input: { name: \"Wireless Mouse\", description: \"Wireless USB Mouse\", price: 800, stock: 25 }) { id name description price stock } }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- Product should be created successfully.
- A new product ID should be generated.
- Product name should be Wireless Mouse.
- Price should be 800.
- Initial stock should be 25.

Save the generated Product ID for the next tests.

### Test 4: Update Product

**Operation:** Mutation – `updateProduct`

Use the Product ID generated in Test 3.

**Request Body:**

```json
{
  "query": "mutation { updateProduct(id: \"YOUR_PRODUCT_ID\", input: { name: \"Updated Wireless Mouse\", description: \"Updated Wireless USB Mouse\", price: 1000, stock: 30 }) { id name description price stock updatedAt } }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- Product details should be updated.
- Name should be Updated Wireless Mouse.
- Price should be 1000.
- Stock should be 30.
- Updated timestamp should reflect the modification.

### Test 5: Delete Product

**Operation:** Mutation – `deleteProduct`

Use the same Product ID.

**Request Body:**

```json
{
  "query": "mutation { deleteProduct(id: \"YOUR_PRODUCT_ID\") }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- Delete operation should succeed.
- Product should no longer appear in the product listing.

Perform this test last because the product will no longer be available.

---

# 5. Inventory Service – GraphQL Testing

Inventory testing verifies the available inventory records and their quantities.

### Test 1: Get All Inventory

**Operation:** Query – Use the inventory listing query defined in the Inventory GraphQL schema.

**Request Body:**

```json
{
  "query": "{ inventories { id productId quantity } }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- All inventory records should be returned.
- Product IDs and available quantities should be displayed.

*Note: Use the actual query and field names from your Inventory GraphQL schema if they differ.*

### Test 2: Get Inventory by ID

**Operation:** Query – Use the inventory lookup query defined in the schema.

**Request Body:**

```json
{
  "query": "{ inventory(id: 1) { id productId quantity } }"
}
```

Replace `1` with an existing inventory ID.

**Expected Result:**

- HTTP Status: `200 OK`
- The requested inventory record should be returned.
- Product ID and available quantity should be correct.

### Test 3: Verify Product Stock

**Operation:** Query – `product`

**Request Body:**

```json
{
  "query": "{ product(id: \"1\") { id name stock } }"
}
```

Replace `1` with an existing Product ID.

**Expected Result:**

- HTTP Status: `200 OK`
- Product details should be returned.
- Stock should match the current available product quantity.

---

# 6. Final Testing Checklist

## Admin – Product Service

- [ ] Admin login successful.
- [ ] Swagger authorized using Admin JWT.
- [ ] Get all products.
- [ ] Get product by ID.
- [ ] Create product.
- [ ] Update product.
- [ ] Delete product.

## Admin – Inventory Service

- [ ] Get all inventory.
- [ ] Get inventory by ID.
- [ ] Verify product stock.

---

## 7. Testing Notes

1. Use the Admin JWT token for all protected operations.
2. Execute GraphQL requests through the API Gateway.
3. Use actual Product IDs and Inventory IDs returned by the application.
4. Save generated IDs because they may be required in subsequent requests.
5. Execute the delete operation only after completing all other product tests.
6. HTTP `200 OK` does not always mean a GraphQL operation succeeded. Check the response body for `data` and `errors`.
7. Record actual responses and compare them with the expected results.

**Next Testing Document:** Admin – Order Service and Inventory Flow.

This will cover order creation, automatic stock reduction, order confirmation, order cancellation and inventory restoration.
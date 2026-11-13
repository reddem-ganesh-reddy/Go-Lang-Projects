# Admin Testing – Order Service and Inventory Flow

## 1. Overview

This document explains how to test the Order Service and its interaction with the Inventory Service using Swagger UI and GraphQL through the API Gateway.

The main objective is to verify that inventory is automatically reserved when an order is created and restored when an order is cancelled.

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

## 2. Prerequisites

Before starting the tests:

1. Start all application services using `.\start-all.ps1`.
2. Open API Gateway Swagger UI.
3. Log in using Admin credentials.
4. Authorize Swagger using the Admin JWT token.
5. Ensure that Product Service and Inventory Service are running.
6. Ensure that at least one product is available with sufficient stock.

For this testing, we will use the following sample data.

| Field | Value |
|---|---|
| Customer ID | 1002 |
| Product ID | 1 |
| Product Name | Laptop |
| Unit Price | 55000 |
| Initial Stock | 20 |
| Order Quantity | 2 |

Use actual IDs and product details from your database if they are different.

## 3. Get Product Details Before Creating an Order

### Test 1: Get Product by ID

**Operation:** Query – `product`

**Request Body:**

```json
{
  "query": "{ product(id: \"1\") { id name description price stock } }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- Product details should be returned.
- Product name should be Laptop.
- Price should be 55000.
- Available stock should be 20.

Record the initial stock quantity.

**Initial Stock = 20**

## 4. Order Service – Create Order

### Test 2: Create a New Order

**Operation:** Mutation – `createOrder`

We will create an order for two laptops.

**Request Body:**

```json
{
  "query": "mutation { createOrder(input: { customerId: 1002, items: [{ productId: 1, quantity: 2 }] }) { id customerId status totalAmount items { id productId productName quantity unitPrice subtotal } } }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- Order should be created successfully.
- A new Order ID should be generated.
- Customer ID should be 1002.
- Product ID should be 1.
- Ordered quantity should be 2.
- Unit price should be 55000.
- Subtotal and total amount should be 110000.
- Order should have its initial status according to the application.

Save the generated Order ID.

**Order Calculation:**

```text
Unit Price = 55000
Quantity  = 2

Total Amount = 55000 × 2
             = 110000
```

## 5. Verify Inventory After Order Creation

### Test 3: Get Product by ID

After creating the order, retrieve the same product again.

**Operation:** Query – `product`

**Request Body:**

```json
{
  "query": "{ product(id: \"1\") { id name price stock } }"
}
```

**Expected Result:**

The available stock should decrease automatically.

| Description | Quantity |
|---|---:|
| Initial stock | 20 |
| Ordered quantity | 2 |
| Available stock | 18 |

**Expected Stock = 18**

This verifies that the Order Service automatically communicates with the Inventory Service during order creation.

### Test 4: Get All Inventory

Execute the inventory listing query exposed by your Inventory GraphQL schema.

**Expected Result:**

- HTTP Status: `200 OK`
- Inventory records should be returned.
- The inventory quantity for Product ID 1 should reflect the reserved quantity.

### Test 5: Get Inventory by ID

Use the actual inventory ID associated with Product ID 1.

Execute the inventory lookup query exposed by your GraphQL schema.

**Expected Result:**

- Correct inventory record should be returned.
- Available quantity should reflect the order reservation.

## 6. Get All Orders

### Test 6: Get All Orders

**Operation:** Query – `orders`

**Request Body:**

```json
{
  "query": "{ orders { id customerId status totalAmount items { id productId productName quantity unitPrice subtotal } } }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- All permitted orders should be displayed.
- The newly created order should be present.
- Order details and total amount should be correct.

## 7. Get Order by ID

### Test 7: Get an Existing Order

**Operation:** Query – `order`

Replace `1` with the Order ID generated during order creation.

**Request Body:**

```json
{
  "query": "{ order(id: 1) { id customerId status totalAmount items { id productId productName quantity unitPrice subtotal } } }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- Requested order should be returned.
- Customer ID should be 1002.
- Product quantity should be 2.
- Total amount should be 110000.

## 8. Confirm Order

### Test 8: Confirm the Created Order

**Operation:** Mutation – `confirmOrder`

Replace `1` with the actual Order ID.

**Request Body:**

```json
{
  "query": "mutation { confirmOrder(id: 1) { id customerId status totalAmount } }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- Order should be confirmed successfully.
- Status should change to `CONFIRMED`.

### Test 9: Verify Confirmed Order

Retrieve the same order again.

**Request Body:**

```json
{
  "query": "{ order(id: 1) { id customerId status totalAmount } }"
}
```

**Expected Result:**

The order should have status `CONFIRMED`.

### Test 10: Check Product Stock After Confirmation

**Request Body:**

```json
{
  "query": "{ product(id: \"1\") { id name stock } }"
}
```

**Expected Result:**

Stock should remain 18.

Confirming an order should not reserve the same inventory again.

## 9. Cancel Order

### Test 11: Cancel the Confirmed Order

**Operation:** Mutation – `cancelOrder`

Replace `1` with the actual Order ID.

**Request Body:**

```json
{
  "query": "mutation { cancelOrder(id: 1) { id customerId status totalAmount } }"
}
```

**Expected Result:**

- HTTP Status: `200 OK`
- Order should be cancelled successfully, provided cancellation of confirmed orders is supported.
- Status should change to `CANCELLED`.
- Reserved inventory should be released.

### Test 12: Verify Cancelled Order

**Request Body:**

```json
{
  "query": "{ order(id: 1) { id customerId status totalAmount } }"
}
```

**Expected Result:**

The order should have status `CANCELLED`.

### Test 13: Verify Product Stock After Cancellation

**Request Body:**

```json
{
  "query": "{ product(id: \"1\") { id name stock } }"
}
```

**Expected Result:**

The previously reserved quantity should be restored.

| Operation | Stock |
|---|---:|
| Initial stock | 20 |
| After creating order | 18 |
| After confirming order | 18 |
| After cancelling order | 20 |

**Final Stock = 20**

This confirms that the inventory reservation is released successfully after cancellation.

## 10. Verify Inventory Restoration

### Test 14: Get All Inventory Again

Execute the inventory listing query.

**Expected Result:**

- Inventory records should be returned.
- The available quantity should reflect the restored stock.
- Product ID 1 should have 20 available units, assuming no other transactions occurred.

### Test 15: Get Inventory by ID Again

Retrieve the same inventory record used in the earlier tests.

**Expected Result:**

The inventory quantity should match the restored quantity.

## 11. Additional Order Validation Tests

### Test 16: Create Order with Zero Quantity

**Request Body:**

```json
{
  "query": "mutation { createOrder(input: { customerId: 1002, items: [{ productId: 1, quantity: 0 }] }) { id status } }"
}
```

**Expected Result:**

Order creation should fail with a validation error.

### Test 17: Create Order with Empty Items

**Request Body:**

```json
{
  "query": "mutation { createOrder(input: { customerId: 1002, items: [] }) { id status } }"
}
```

**Expected Result:**

Order creation should fail because at least one item is required.

### Test 18: Create Order with Nonexistent Product

**Request Body:**

```json
{
  "query": "mutation { createOrder(input: { customerId: 1002, items: [{ productId: 999999, quantity: 1 }] }) { id status } }"
}
```

**Expected Result:**

Order creation should fail with a product-not-found error.

### Test 19: Create Order with Insufficient Stock

Assume the product has only 2 units available.

Attempt to order 5 units.

**Request Body:**

```json
{
  "query": "mutation { createOrder(input: { customerId: 1002, items: [{ productId: 1, quantity: 5 }] }) { id status totalAmount } }"
}
```

**Expected Result:**

- Order creation should fail.
- Insufficient stock error should be returned.
- Inventory should not become negative.
- No incomplete order should remain.

### Test 20: Cancel the Same Order Again

Execute the cancellation mutation again using the already cancelled Order ID.

**Request Body:**

```json
{
  "query": "mutation { cancelOrder(id: 1) { id status } }"
}
```

**Expected Result:**

- Duplicate cancellation should be rejected or handled idempotently.
- Inventory must not be restored a second time.
- Stock should remain 20.

---

## 12. Final Testing Checklist

### Admin – Order Service

- [ ] Get product by ID before creating order.
- [ ] Create order successfully.
- [ ] Save generated Order ID.
- [ ] Get all orders.
- [ ] Get order by ID.
- [ ] Confirm order.
- [ ] Verify confirmed order status.
- [ ] Cancel confirmed order.
- [ ] Verify cancelled order status.

### Admin – Inventory Verification

- [ ] Record initial product stock.
- [ ] Verify stock reduction after order creation.
- [ ] Verify stock remains unchanged after confirmation.
- [ ] Verify stock restoration after cancellation.
- [ ] Get all inventory after cancellation.
- [ ] Get inventory by ID after cancellation.
- [ ] Verify duplicate cancellation does not restore stock again.

### Order Validation

- [ ] Zero quantity.
- [ ] Empty order items.
- [ ] Nonexistent product.
- [ ] Insufficient stock.
- [ ] Repeated cancellation.

---

## 13. Final Expected Inventory Result

For this test, the initial inventory is 20 laptops and the order quantity is 2.

The expected flow is:

```text
Initial Stock: 20

Create Order
     |
     v
Reserve 2 Units
     |
     v
Available Stock: 18
     |
     v
Confirm Order
     |
     v
Available Stock: 18
     |
     v
Cancel Order
     |
     v
Release 2 Units
     |
     v
Final Stock: 20
```

**Testing Objective:** Verify that the Admin can manage orders and that the Order Service correctly communicates with the Inventory Service for stock reservation and restoration.

**Next Testing Document:** User – Product and Inventory Testing.
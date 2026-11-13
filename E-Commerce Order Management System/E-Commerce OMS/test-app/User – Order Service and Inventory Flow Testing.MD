# User – Order Service and Inventory Flow Testing

## 1. Login as User

**Gateway URL:** `http://localhost:8080`

**GraphQL Endpoint:** `http://localhost:8080/graphql`

Open Swagger UI:

`http://localhost:8080/swagger/index.html`

Login using the User credentials:

```json
{
  "username": "user",
  "password": "user@123"
}
```

Copy the JWT token from the login response.

In Swagger, open the GraphQL endpoint and pass the token in the Authorization header:

```text
Bearer YOUR_USER_TOKEN
```

Use the User token for all the following test cases.

## 2. Check Product Availability

Before creating an order, verify that the product exists and has sufficient stock.

**GraphQL Query:**

```graphql
query {
  product(id: "1") {
    id
    name
    description
    price
    stock
  }
}
```

**Expected Result:**
- Product details should be displayed.
- Stock should be greater than zero.
- Note the initial stock quantity for later verification.

## 3. Create an Order

Create an order using the customer ID associated with the logged-in User account. Do not use another customer's ID.

**GraphQL Mutation:**

```graphql
mutation {
  createOrder(
    input: {
      customerId: YOUR_CUSTOMER_ID
      items: [
        {
          productId: 1
          quantity: 1
        }
      ]
    }
  ) {
    id
    customerId
    status
    totalAmount
    items {
      id
      productId
      productName
      quantity
      unitPrice
      subtotal
    }
  }
}
```

**Expected Result:**
- Order should be created successfully.
- The order should belong to the logged-in User.
- Order ID, customer ID, status and total amount should be returned.
- Save the generated Order ID for the remaining test cases.
- Inventory should be reserved automatically if that is the configured order-creation flow.

## 4. Verify Inventory After Order Creation

Run the product query again:

```graphql
query {
  product(id: "1") {
    id
    name
    stock
  }
}
```

**Expected Result:**

If the initial stock was 20 and the order quantity was 1, the updated stock should be 19.

## 5. Get All Orders

Run the order-list query available in the GraphQL schema.

**Example:**

```graphql
query {
  orders {
    id
    customerId
    status
    totalAmount
  }
}
```

**Expected Result:**
- The logged-in User should see only their own orders.
- Orders belonging to other customers should not be exposed.
- If the existing `orders` query returns all customers' orders, verify that the service or authorization layer enforces user ownership. It must not expose other users' orders.

## 6. Get Order by ID

Use the Order ID returned during order creation.

```graphql
query {
  order(id: YOUR_ORDER_ID) {
    id
    customerId
    status
    totalAmount
    items {
      id
      productId
      productName
      quantity
      unitPrice
      subtotal
    }
  }
}
```

**Expected Result:**
- The order details should be displayed.
- The order must belong to the logged-in User.

## 7. Confirm the Order

Use the Order ID generated in Step 3.

```graphql
mutation {
  confirmOrder(id: YOUR_ORDER_ID) {
    id
    status
    totalAmount
  }
}
```

**Expected Result:**
- The User should be able to confirm their own order if the order is in a valid state.
- The order status should change according to the application's confirmation flow.
- Inventory should not be deducted a second time if it was already reserved during order creation.

## 8. Cancel the Order

If the application allows cancellation in the current order state, run:

```graphql
mutation {
  cancelOrder(id: YOUR_ORDER_ID) {
    id
    status
    totalAmount
  }
}
```

**Expected Result:**
- The User should be able to cancel their own order if cancellation is allowed by the business rules.
- The order status should change to Cancelled.
- Reserved inventory should be released and stock should be restored.

If cancellation is not allowed after confirmation, create another order and test cancellation while it is in a cancellable state.

## 9. Verify Inventory After Cancellation

Run the product query again:

```graphql
query {
  product(id: "1") {
    id
    name
    stock
  }
}
```

**Expected Result:**

If the initial stock was 20:

- Before order creation: 20
- After creating an order for quantity 1: 19
- After confirmation: 19
- After cancellation: 20

## 10. Negative Test Cases

### Test 1: Create Order with Zero Quantity

Try creating an order with quantity `0`.

**Expected Result:** The application should reject the request with a validation error.

### Test 2: Create Order with Empty Items

Try creating an order with an empty items list.

**Expected Result:** The application should reject the request.

### Test 3: Create Order with Invalid Product ID

Use a product ID that does not exist.

**Expected Result:** The application should return a meaningful product-not-found error.

### Test 4: Create Order with Insufficient Stock

Request a quantity greater than the available stock.

**Expected Result:** The application should reject the order. Inventory should remain unchanged.

### Test 5: Access Another User's Order

Try retrieving another customer's order using their Order ID.

**Expected Result:** Access should be denied or the order should be returned as not found. Other customers' order details must not be exposed.

### Test 6: Confirm Another User's Order

Try confirming an order belonging to another customer.

**Expected Result:** The application should reject the request.

### Test 7: Cancel Another User's Order

Try cancelling an order belonging to another customer.

**Expected Result:** The application should reject the request. The order status and inventory should remain unchanged.

### Test 8: Confirm an Already Confirmed Order

Try confirming the same order again.

**Expected Result:** The application should handle the invalid state without duplicating inventory deductions.

### Test 9: Cancel an Already Cancelled Order

Try cancelling the same order again.

**Expected Result:** The application should return a meaningful error or handle the repeated request safely without releasing inventory twice.

## 11. Final Verification Checklist

- [ ] User login is successful.
- [ ] User can view product details.
- [ ] User can create an order.
- [ ] User can view only their own orders.
- [ ] User can view their own order by ID.
- [ ] User can confirm their own order.
- [ ] User can cancel their own order when permitted.
- [ ] Inventory is reserved during order creation.
- [ ] Inventory is not deducted again during confirmation.
- [ ] Inventory is restored after cancellation.
- [ ] User cannot access another customer's order.
- [ ] User cannot confirm or cancel another customer's order.
- [ ] Invalid order requests are rejected with meaningful errors.

**Important:** GraphQL can return HTTP `200 OK` even when an operation fails. Always check the response's `errors` and `data` fields to verify the actual result.
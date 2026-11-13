package main

/*

Swagger URL: `http://localhost:8081/swagger/index.html`

## Order Service — CRUD testing in Swagger


### 1. CREATE — Create an order

POST

`http://localhost:8081/api/orders`

Test 1: Create order successfully

Request body:

```
{
  "customer_id": 1001,
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    }
  ]
}
```

Expected: `201 Created` with the order ID, customer ID, status, items and total amount.


Test 2: Missing customer ID

```
{
  "items": [
    {
      "product_id": 1,
      "quantity": 1
    }
  ]
}
```

Expected: `400 Bad Request`.

Test 3: Empty items

```
{
  "customer_id": 1001,
  "items": []
}
```

Expected: `400 Bad Request`.

Test 4: Missing product ID

```
{
  "customer_id": 1001,
  "items": [
    {
      "quantity": 1
    }
  ]
}
```

Expected: `400 Bad Request`.

Test 5: Zero quantity

```
{
  "customer_id": 1001,
  "items": [
    {
      "product_id": 1,
      "quantity": 0
    }
  ]
}
```

Expected: `400 Bad Request`.

Test 6: Negative quantity

Use `"quantity": -2` in the request.

Expected: `400 Bad Request`.

Test 7: Product does not exist

Use `"product_id": 999999` and `"quantity": 1`.

Expected: order creation fails; verify the actual error status and message returned by your service.

Test 8: Invalid JSON

```
{
  "customer_id": 1001,
  "items":
}
```

Expected: `400 Bad Request`.

### 2. READ — Get all orders

GET

`http://localhost:8081/api/orders`

Test 9: Get all orders

No request body.

Expected: `200 OK` with an array of orders. If no orders exist, the response may be `[]`.



Test 10: Verify created orders

Run the same GET request after creating an order.

Expected: the response includes the order you created, with the correct customer ID, items, status and total amount.

### 3. READ — Get order by ID

GET

`http://localhost:8081/api/orders/{id}`

Test 11: Existing order

URL: `http://localhost:8081/api/orders/1`

Expected: `200 OK` with order details.

Test 12: Nonexistent order

URL: `http://localhost:8081/api/orders/999999`

Expected: `404 Not Found`.

Test 13: Invalid order ID

URL: `http://localhost:8081/api/orders/abc`

Expected: `400 Bad Request`.

Test 14: Zero order ID

URL: `http://localhost:8081/api/orders/0`

Expected: `400 Bad Request`, according to your current controller validation.

### 4. UPDATE — Confirm an order

PATCH

`http://localhost:8081/api/orders/{id}/confirm`

Test 15: Confirm a created order

First create a new order and use its returned ID.

Example URL: `http://localhost:8081/api/orders/2`

Append `/confirm` to the URL:

`http://localhost:8081/api/orders/2/confirm`

No request body.

Expected: `200 OK` and status changes to `"confirmed"`.

Test 16: Confirm an already confirmed order

Repeat the confirmation request for the same ID.

Expected: an error because the status transition is not allowed.

Test 17: Confirm a nonexistent order

URL: `http://localhost:8081/api/orders/999999/confirm`

Expected: an error response.

Test 18: Confirm using an invalid ID

URL: `http://localhost:8081/api/orders/abc/confirm`

Expected: `400 Bad Request`.

### 5. UPDATE — Cancel an order

PATCH

`http://localhost:8081/api/orders/{id}/cancel`

Test 19: Cancel a created order

Use a separate order that is still in `"created"` status.

Example URL: `http://localhost:8081/api/orders/3/cancel`

No request body.

Expected: `200 OK` and status changes to `"canceled"`.

Test 20: Cancel a confirmed order

Create and confirm another order, then cancel it.

Expected: cancellation succeeds if your service allows the `confirmed → canceled` transition.

Test 21: Cancel an already canceled order

Repeat cancellation for the same order ID.

Expected: an error because the status transition is not allowed.

Test 22: Cancel a nonexistent order

URL: `http://localhost:8081/api/orders/999999/cancel`

Expected: an error response.

Test 23: Cancel using an invalid ID

URL: `http://localhost:8081/api/orders/abc/cancel`

Expected: `400 Bad Request`.

### 6. DELETE — Delete an order

our current Order Service does not expose a `DELETE /api/orders/{id}` endpoint. Cancellation is not the same as deletion: cancellation changes the order's status and may release reserved stock.


*/
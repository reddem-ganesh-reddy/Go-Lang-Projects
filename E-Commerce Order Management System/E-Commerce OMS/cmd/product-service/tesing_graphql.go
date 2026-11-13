package main

/*

Get ALL:

{
  products {
    items {
      id
      name
      price
      stock
    }
    totalRecords
    totalPages
  }
}


Use these test cases in Product Service GraphiQL UI at `http://localhost:8082/graphql`. They cover product listing, filtering, sorting, pagination, individual product lookup, create, update, and delete.

## 1. Create test products
===================================
Run each mutation separately in GraphiQL. These records give us different prices, stock values, and names for testing combinations.

Test 1 — Create a laptop
------------------------------------

```
mutation {
  createProduct(input: {
    name: "Laptop"
    description: "Business laptop"
    price: 55000
    stock: 10
  }) {
    id
    name
    price
    stock
  }
}
```

Test 2 — Create a mouse
------------------------------
```
mutation {
  createProduct(input: {
    name: "Wireless Mouse"
    description: "Wireless USB mouse"
    price: 800
    stock: 25
  }) {
    id
    name
    price
    stock
  }
}
```

Test 3 — Create a keyboard with zero stock
----------------------------------------------

```
mutation {
  createProduct(input: {
    name: "Keyboard"
    description: "Mechanical keyboard"
    price: 1200
    stock: 0
  }) {
    id
    name
    price
    stock
  }
}
```

Test 4 — Create another laptop accessory
----------------------------------------------
```
mutation {
  createProduct(input: {
    name: "Laptop Stand"
    description: "Adjustable aluminum stand"
    price: 1800
    stock: 7
  }) {
    id
    name
    price
    stock
  }
}
```

If we already created some of these products, don't run the mutations again unless we want duplicate test records.

## 2. Test listing and pagination
========================================
Test 5 — Get all products
----------------------------------------------
```
{
  products {
    items {
      id
      name
      price
      stock
    }
    page
    limit
    totalRecords
    totalPages
  }
}
```

Test 6 — First page, two products
----------------------------------------------
```
{
  products(filter: {page: 1, limit: 2}) {
    items {
      id
      name
    }
    page
    limit
    totalRecords
    totalPages
  }
}
```

Test 7 — Second page, two products
----------------------------------------------
```
{
  products(filter: {page: 2, limit: 2}) {
    items {
      id
      name
    }
    page
    limit
    totalRecords
    totalPages
  }
}
```

Test 8 — Page size omitted
----------------------------------------------
```
{
  products(filter: {page: 1}) {
    items {
      id
      name
    }
    page
    limit
  }
}
```

This should use our service's default page size.

## 3. Test filtering
============================================
Test 9 — Filter by name
----------------------------------------------
```
{
  products(filter: {name: "Laptop"}) {
    items {
      id
      name
      price
    }
    totalRecords
  }
}
```

Test 10 — Minimum price
----------------------------------------------
```
{
  products(filter: {minPrice: 1000}) {
    items {
      name
      price
    }
  }
}
```

Test 11 — Maximum price
-----------------------------
```
{
  products(filter: {maxPrice: 2000}) {
    items {
      name
      price
    }
  }
}
```

Test 12 — Price range
-----------------------
```
{
  products(filter: {minPrice: 500, maxPrice: 2000}) {
    items {
      name
      price
    }
  }
}
```

Test 13 — Products in stock
----------------------------------
```
{
  products(filter: {inStock: true}) {
    items {
      name
      stock
    }
  }
}
```

Test 14 — Products out of stock
----------------------------------
```
{
  products(filter: {inStock: false}) {
    items {
      name
      stock
    }
  }
}
```

## 4. Test sorting and combined filters
=======================================
Test 15 — Sort by price, low to high
----------------------------------
```
{
  products(filter: {sortBy: "price", sortOrder: "asc"}) {
    items {
      name
      price
    }
  }
}
```

Test 16 — Sort by price, high to low
--------------------------------------
```
{
  products(filter: {sortBy: "price", sortOrder: "desc"}) {
    items {
      name
      price
    }
  }
}
```

Test 17 — Combine stock filter, price range, sorting, and pagination
------------------------------------------------
```
{
  products(filter: {
    minPrice: 500
    maxPrice: 2000
    inStock: true
    sortBy: "price"
    sortOrder: "desc"
    page: 1
    limit: 2
  }) {
    items {
      id
      name
      price
      stock
    }
    totalRecords
    totalPages
  }
}
```

Test 18 — Filter and sort without pagination values
----------------------------------------------
```
{
  products(filter: {
    inStock: true
    sortBy: "name"
    sortOrder: "asc"
  }) {
    items {
      name
      stock
    }
    page
    limit
  }
}
```

## 5. Test individual product lookup
=========================================
Use an actual product ID returned by one of our create mutations.

```
query {
  product(id: "1") {
    id
    name
    description
    price
    stock
    createdAt
    updatedAt
  }
}
```

If ID `1` isn't present in our database, replace it with an existing ID.

## 6. Test update and delete
====================================
Replace `"1"` with an existing product ID.

Test 20 — Update a product
=============================
```
mutation {
  updateProduct(
    id: "1"
    input: {
      name: "Updated Laptop"
      description: "Updated business laptop"
      price: 60000
      stock: 15
    }
  ) {
    id
    name
    price
    stock
    updatedAt
  }
}
```

Test 21 — Delete a product
=----------------------------------
Use a test product ID that we no longer need, because deletion changes the database.

```
mutation {
  deleteProduct(id: "4")
}
```

If our API returns a Boolean, `true` indicates successful deletion.

## 7. Test validation and error handling
======================================
These are expected to produce errors or validation failures rather than successful results.

| Test                | Input                                                |
| ------------------- | ---------------------------------------------------- |
| Invalid product ID  | `product(id: "abc")`                                 |
| Nonexistent product | `product(id: "999999")`                              |
| Negative price      | Create with `price: -100`                            |
| Negative stock      | Create with `stock: -5`                              |
| Invalid page        | `products(filter: {page: 0})`                        |
| Invalid page size   | `products(filter: {limit: 101})`                     |
| Invalid sort field  | `products(filter: {sortBy: "password"})`             |
| Invalid sort order  | `products(filter: {sortOrder: "sideways"})`          |
| Invalid price range | `products(filter: {minPrice: 5000, maxPrice: 1000})` |


*/
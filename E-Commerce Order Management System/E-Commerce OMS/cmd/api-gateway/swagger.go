package main

// swaggerAuthLogin godoc
import (
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// registerGatewaySwagger registers the Swagger UI and its JSON document.
func registerGatewaySwagger(router *gin.Engine) {
	router.GET("/openapi.json", func(c *gin.Context) {
		c.JSON(200, gatewaySwaggerSpec())
	})

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL("/openapi.json"),
		),
	)
}

func gatewaySwaggerSpec() gin.H {
	return gin.H{
		"swagger": "2.0",
		"info": gin.H{
			"title":       "Order Management API Gateway",
			"version":     "1.0.0",
			"description": "Single API Gateway for Authentication, Product Service, Inventory Reservation, Order Service, and both existing GraphQL services.",
		},
		"host":     "localhost:8080",
		"basePath": "/",
		"schemes":  []string{"http"},
		"consumes": []string{"application/json"},
		"produces": []string{"application/json"},

		"securityDefinitions": gin.H{
			"BearerAuth": gin.H{
				"type":        "apiKey",
				"name":        "Authorization",
				"in":          "header",
				"description": "Enter: Bearer followed by the JWT token returned by POST /api/auth/login.",
			},
		},

		"paths":       gatewaySwaggerPaths(),
		"definitions": gatewaySwaggerDefinitions(),
	}
}

func gatewaySwaggerPaths() gin.H {
	return gin.H{
		"/health": gin.H{
			"get": operation(
				"Gateway health check",
				"Checks whether the API Gateway is running. This operational endpoint does not require a JWT.",
				"Health",
				nil,
				nil,
				map[string]string{"200": "Gateway is healthy"},
				false,
			),
		},

		"/api/auth/login": gin.H{
			"post": operation(
				"Login",
				"Authenticates an Admin or User and returns a JWT. This is the public login endpoint. Use the returned token in Swagger Authorize before accessing protected APIs.",
				"Authentication",
				nil,
				bodyParameter("LoginRequest", "Username and password"),
				map[string]string{
					"200": "Login successful. Returns token, token_type, and role.",
					"400": "Invalid request body.",
					"401": "Invalid username or password.",
					"500": "Token generation failed.",
				},
				false,
			),
		},

		"/api/products": gin.H{
			"get": operation(
				"Get products",
				"Returns a paginated product list. Both Admin and User can access this endpoint.",
				"Products",
				productQueryParameters(),
				nil,
				map[string]string{
					"200": "Returns ProductListResponse.",
					"400": "Invalid query parameters.",
					"500": "Unable to retrieve products.",
				},
				true,
			),
			"post": operation(
				"Create product",
				"Creates a product. Admin role is required.",
				"Products",
				nil,
				bodyParameter("CreateProductRequest", "Product information"),
				map[string]string{
					"201": "Returns the created ProductResponse.",
					"400": "Invalid product data.",
					"500": "Unable to create product.",
				},
				true,
			),
		},

		"/api/products/{id}": gin.H{
			"get": operation(
				"Get product by ID",
				"Returns one product using its numeric ID. Both roles can access this endpoint.",
				"Products",
				[]interface{}{pathParameter("id", "Product ID", "integer")},
				nil,
				map[string]string{
					"200": "Returns ProductResponse.",
					"400": "Invalid product ID.",
					"404": "Product not found.",
					"500": "Unable to retrieve product.",
				},
				true,
			),
			"put": operation(
				"Update product",
				"Updates an existing product. Admin role is required.",
				"Products",
				[]interface{}{pathParameter("id", "Product ID", "integer")},
				bodyParameter("UpdateProductRequest", "Product fields to update"),
				map[string]string{
					"200": "Returns the updated ProductResponse.",
					"400": "Invalid product ID or request body.",
					"404": "Product not found.",
					"500": "Unable to update product.",
				},
				true,
			),
			"delete": operation(
				"Delete product",
				"Deletes an existing product. Admin role is required.",
				"Products",
				[]interface{}{pathParameter("id", "Product ID", "integer")},
				nil,
				map[string]string{
					"200": "Product deleted successfully.",
					"400": "Invalid product ID.",
					"404": "Product not found.",
					"500": "Unable to delete product.",
				},
				true,
			),
		},

		"/api/inventory/reservations": gin.H{
			"get": operation(
				"Get all inventory reservations",
				"Returns all inventory reservations. Admin role is required. Normal Users cannot list all reservations.",
				"Inventory",
				nil,
				nil,
				map[string]string{
					"200": "Returns an array of InventoryReservation.",
					"500": "Unable to retrieve reservations.",
				},
				true,
			),
			"post": operation(
				"Reserve stock",
				"Creates an inventory reservation. Admin role is required. Order creation uses the existing internal Order Service to Product Service reservation flow.",
				"Inventory",
				nil,
				bodyParameter("ReserveStockRequest", "Reservation information"),
				map[string]string{
					"200": "Returns reservation confirmation.",
					"400": "Invalid reservation data.",
					"500": "Unable to reserve stock.",
				},
				true,
			),
		},

		"/api/inventory/reservations/{reservation_id}": gin.H{
			"get": operation(
				"Get reservation by reservation ID",
				"Returns one inventory reservation. Both Admin and User can access this endpoint.",
				"Inventory",
				[]interface{}{pathParameter("reservation_id", "Reservation ID", "string")},
				nil,
				map[string]string{
					"200": "Returns InventoryReservation.",
					"400": "Invalid reservation ID.",
					"404": "Reservation not found.",
					"500": "Unable to retrieve reservation.",
				},
				true,
			),
		},

		"/api/inventory/reservations/release": gin.H{
			"post": operation(
				"Release stock reservation",
				"Releases an existing inventory reservation. Admin role is required.",
				"Inventory",
				nil,
				bodyParameter("ReleaseStockRequest", "Reservation ID to release"),
				map[string]string{
					"200": "Returns release confirmation.",
					"400": "Invalid reservation ID.",
					"404": "Reservation not found.",
					"500": "Unable to release reservation.",
				},
				true,
			),
		},

		"/api/orders": gin.H{
			"get": operation(
				"Get all orders",
				"Returns all orders. Both Admin and User can access this endpoint.",
				"Orders",
				nil,
				nil,
				map[string]string{
					"200": "Returns an array of OrderResponse.",
					"500": "Unable to retrieve orders.",
				},
				true,
			),
			"post": operation(
				"Create order",
				"Creates an order. The Order Service obtains product information and reserves stock internally through the existing Product Service integration. Both Admin and User can create orders.",
				"Orders",
				nil,
				bodyParameter("CreateOrderRequest", "Customer ID and order items"),
				map[string]string{
					"201": "Returns the created OrderResponse.",
					"400": "Invalid order data.",
					"404": "Referenced product or department-related resource not found, where applicable.",
					"500": "Unable to create order.",
				},
				true,
			),
		},

		"/api/orders/{id}": gin.H{
			"get": operation(
				"Get order by ID",
				"Returns one order using its numeric ID. Both roles can access this endpoint.",
				"Orders",
				[]interface{}{pathParameter("id", "Order ID", "integer")},
				nil,
				map[string]string{
					"200": "Returns OrderResponse.",
					"400": "Invalid order ID.",
					"404": "Order not found.",
					"500": "Unable to retrieve order.",
				},
				true,
			),
		},

		"/api/orders/{id}/confirm": gin.H{
			"patch": operation(
				"Confirm order",
				"Confirms an existing order. Both Admin and User can access this endpoint.",
				"Orders",
				[]interface{}{pathParameter("id", "Order ID", "integer")},
				nil,
				map[string]string{
					"200": "Returns the updated OrderResponse.",
					"400": "Invalid order ID or order state.",
					"404": "Order not found.",
					"500": "Unable to confirm order.",
				},
				true,
			),
		},

		"/api/orders/{id}/cancel": gin.H{
			"patch": operation(
				"Cancel order",
				"Cancels an existing order. The Order Service releases its inventory reservations through the existing internal integration. Both Admin and User can access this endpoint.",
				"Orders",
				[]interface{}{pathParameter("id", "Order ID", "integer")},
				nil,
				map[string]string{
					"200": "Returns the updated OrderResponse.",
					"400": "Invalid order ID or order state.",
					"404": "Order not found.",
					"500": "Unable to cancel order.",
				},
				true,
			),
		},

		"/graphql": gin.H{
			"post": operation(
				"Execute Product or Order GraphQL operation",
				"One common GraphQL endpoint for both existing services. Product root fields are product, products, createProduct, updateProduct, and deleteProduct. Order root fields are orders, order, createOrder, confirmOrder, and cancelOrder. Product mutations require Admin. All Order operations allow Admin and User. A single request must target either Product or Order, not both. Send variables as a JSON object.",
				"GraphQL",
				nil,
				bodyParameter("GraphQLRequest", "GraphQL query, optional operationName, and variables"),
				map[string]string{
					"200": "GraphQL response containing data and/or errors. GraphQL execution errors can be returned with HTTP 200.",
					"400": "Invalid request, unsupported operation, or mixed Product and Order operation.",
					"401": "Missing, invalid, or expired JWT.",
					"403": "User role does not have permission for this operation.",
					"502": "Downstream service failed.",
					"503": "Downstream service circuit breaker is open.",
				},
				true,
			),
			// "get": operation(
			// 	"Execute GraphQL query using GET",
			// 	"Runs a GraphQL query using the query URL parameter. Authentication is required. POST is recommended for mutations and for requests containing variables.",
			// 	"GraphQL",
			// 	[]interface{}{
			// 		queryParameter("query", "GraphQL query text", "string", true),
			// 		queryParameter("operationName", "Optional operation name", "string", false),
			// 		queryParameter("variables", "Optional JSON-encoded GraphQL variables", "string", false),
			// 	},
			// 	nil,
			// 	map[string]string{
			// 		"200": "GraphQL response containing data and/or errors.",
			// 		"400": "Invalid or unsupported GraphQL operation.",
			// 		"401": "Missing, invalid, or expired JWT.",
			// 		"403": "User role does not have permission for this operation.",
			// 		"502": "Downstream service failed.",
			// 		"503": "Downstream service circuit breaker is open.",
			// 	},
			// 	true,
			// ),
		},
	}
}

func gatewaySwaggerDefinitions() gin.H {
	return gin.H{
		"LoginRequest": objectDefinition(
			"Login request",
			[]string{"username", "password"},
			gin.H{
				"username": stringProperty("Username configured for the Auth Service."),
				"password": stringProperty("Password configured for the Auth Service."),
			},
			gin.H{
				"username": "admin",
				"password": "your-admin-password",
			},
		),

		"LoginResponse": objectDefinition(
			"JWT login response",
			[]string{"token", "token_type", "role"},
			gin.H{
				"token":      stringProperty("JWT access token."),
				"token_type": stringProperty("Token type. Normally Bearer."),
				"role":       stringProperty("Authenticated role: admin or user."),
			},
			gin.H{
				"token":      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				"token_type": "Bearer",
				"role":       "admin",
			},
		),

		"CreateProductRequest": objectDefinition(
			"Product creation request",
			[]string{"name", "price"},
			gin.H{
				"name":        stringProperty("Product name."),
				"description": stringProperty("Product description."),
				"price":       numberProperty("Product price."),
				"stock":       integerProperty("Available stock quantity."),
			},
			gin.H{
				"name":        "Laptop",
				"description": "Business laptop",
				"price":       55000.0,
				"stock":       20,
			},
		),

		"UpdateProductRequest": objectDefinition(
			"Product update request",
			nil,
			gin.H{
				"name":        stringProperty("Product name."),
				"description": stringProperty("Product description."),
				"price":       numberProperty("Product price."),
				"stock":       integerProperty("Available stock quantity."),
			},
			gin.H{
				"name":        "Updated Laptop",
				"description": "Updated product description",
				"price":       60000.0,
				"stock":       15,
			},
		),

		"ProductResponse": objectDefinition(
			"Product response",
			[]string{"id", "name", "price", "stock", "created_at", "updated_at"},
			gin.H{
				"id":          integerProperty("Product ID."),
				"name":        stringProperty("Product name."),
				"description": stringProperty("Product description."),
				"price":       numberProperty("Product price."),
				"stock":       integerProperty("Available stock quantity."),
				"created_at":  dateProperty("Product creation timestamp."),
				"updated_at":  dateProperty("Last product update timestamp."),
			},
			gin.H{
				"id":          1,
				"name":        "Laptop",
				"description": "Business laptop",
				"price":       55000.0,
				"stock":       20,
				"created_at":  "2026-10-31T10:00:00Z",
				"updated_at":  "2026-10-31T10:00:00Z",
			},
		),

		"ProductListResponse": objectDefinition(
			"Paginated product response",
			[]string{"items", "page", "limit", "total_pages", "total_records"},
			gin.H{
				"items":         arrayProperty("#/definitions/ProductResponse", "Product records."),
				"page":          integerProperty("Current page number."),
				"limit":         integerProperty("Maximum records per page."),
				"total_pages":   integerProperty("Total number of pages."),
				"total_records": integerProperty("Total number of matching products."),
			},
			gin.H{
				"items":         []interface{}{},
				"page":          1,
				"limit":         10,
				"total_pages":   1,
				"total_records": 0,
			},
		),

		"ReserveStockRequest": objectDefinition(
			"Inventory reservation request",
			[]string{"reservation_id", "product_id", "quantity"},
			gin.H{
				"reservation_id": stringProperty("Unique reservation ID."),
				"product_id":     integerProperty("Product ID."),
				"quantity":       integerProperty("Quantity to reserve. Must be greater than zero."),
			},
			gin.H{
				"reservation_id": "RES-1001",
				"product_id":     1,
				"quantity":       2,
			},
		),

		"ReleaseStockRequest": objectDefinition(
			"Inventory release request",
			[]string{"reservation_id"},
			gin.H{
				"reservation_id": stringProperty("Reservation ID to release."),
			},
			gin.H{
				"reservation_id": "RES-1001",
			},
		),

		"InventoryReservation": objectDefinition(
			"Inventory reservation response",
			[]string{"id", "reservation_id", "product_id", "quantity", "status", "created_at", "updated_at"},
			gin.H{
				"id":             integerProperty("Reservation database ID."),
				"reservation_id": stringProperty("Unique reservation ID."),
				"product_id":     integerProperty("Product ID."),
				"quantity":       integerProperty("Reserved quantity."),
				"status":         stringProperty("Reservation status."),
				"created_at":     dateProperty("Reservation creation timestamp."),
				"updated_at":     dateProperty("Last reservation update timestamp."),
			},
			gin.H{
				"id":             1,
				"reservation_id": "RES-1001",
				"product_id":     1,
				"quantity":       2,
				"status":         "reserved",
				"created_at":     "2026-10-31T10:00:00Z",
				"updated_at":     "2026-10-31T10:00:00Z",
			},
		),

		"ReservationResponse": objectDefinition(
			"Reservation action response",
			[]string{"message", "reservation_id"},
			gin.H{
				"message":        stringProperty("Result message."),
				"reservation_id": stringProperty("Reservation ID."),
				"product_id":     integerProperty("Product ID, returned by reserve operation."),
				"quantity":       integerProperty("Reserved quantity, returned by reserve operation."),
			},
			gin.H{
				"message":        "Stock reserved successfully",
				"reservation_id": "RES-1001",
				"product_id":     1,
				"quantity":       2,
			},
		),

		"CreateOrderRequest": objectDefinition(
			"Order creation request",
			[]string{"customer_id", "items"},
			gin.H{
				"customer_id": integerProperty("Customer ID. Must be greater than zero."),
				"items":       arrayProperty("#/definitions/OrderItemRequest", "Items included in the order."),
			},
			gin.H{
				"customer_id": 101,
				"items": []interface{}{
					gin.H{
						"product_id": 1,
						"quantity":   2,
					},
				},
			},
		),

		"OrderItemRequest": objectDefinition(
			"Order item request",
			[]string{"product_id", "quantity"},
			gin.H{
				"product_id": integerProperty("Product ID. Must be greater than zero."),
				"quantity":   integerProperty("Requested quantity. Must be greater than zero."),
			},
			gin.H{
				"product_id": 1,
				"quantity":   2,
			},
		),

		"OrderResponse": objectDefinition(
			"Order response",
			[]string{"id", "customer_id", "status", "total_amount", "items", "created_at", "updated_at"},
			gin.H{
				"id":           integerProperty("Order ID."),
				"customer_id":  integerProperty("Customer ID."),
				"status":       stringProperty("Order status."),
				"total_amount": numberProperty("Total order amount calculated by the Order Service."),
				"items":        arrayProperty("#/definitions/OrderItemResponse", "Order items."),
				"created_at":   dateProperty("Order creation timestamp."),
				"updated_at":   dateProperty("Last order update timestamp."),
			},
			gin.H{
				"id":           1,
				"customer_id":  101,
				"status":       "pending",
				"total_amount": 110000.0,
				"items": []interface{}{
					gin.H{
						"id":           1,
						"product_id":   1,
						"product_name": "Laptop",
						"quantity":     2,
						"unit_price":   55000.0,
						"subtotal":     110000.0,
					},
				},
				"created_at": "2026-10-31T10:00:00Z",
				"updated_at": "2026-10-31T10:00:00Z",
			},
		),

		"OrderItemResponse": objectDefinition(
			"Order item response",
			[]string{"id", "product_id", "product_name", "quantity", "unit_price", "subtotal"},
			gin.H{
				"id":           integerProperty("Order item ID."),
				"product_id":   integerProperty("Product ID."),
				"product_name": stringProperty("Product name."),
				"quantity":     integerProperty("Ordered quantity."),
				"unit_price":   numberProperty("Price per unit."),
				"subtotal":     numberProperty("Item subtotal."),
			},
			gin.H{
				"id":           1,
				"product_id":   1,
				"product_name": "Laptop",
				"quantity":     2,
				"unit_price":   55000.0,
				"subtotal":     110000.0,
			},
		),

		"GraphQLRequest": objectDefinition(
			"GraphQL HTTP request",
			[]string{"query"},
			gin.H{
				"query":         stringProperty("GraphQL query or mutation text."),
				"operationName": stringProperty("Optional operation name when the document contains multiple operations."),
				"variables": gin.H{
					"type":                 "object",
					"description":          "Optional GraphQL variables object.",
					"additionalProperties": gin.H{},
				},
			},
			gin.H{
				"query": "query { products { items { id name price stock } page limit totalRecords totalPages } }",
			},
		),

		"GraphQLResponse": objectDefinition(
			"GraphQL response",
			nil,
			gin.H{
				"data": gin.H{
					"type":                 "object",
					"description":          "GraphQL result data. The returned fields depend on the selected operation.",
					"additionalProperties": gin.H{},
				},
				"errors": gin.H{
					"type":        "array",
					"description": "GraphQL execution or validation errors, when present.",
					"items": gin.H{
						"type": "object",
						"properties": gin.H{
							"message": stringProperty("GraphQL error message."),
						},
					},
				},
			},
			gin.H{
				"data": gin.H{
					"products": gin.H{
						"items":        []interface{}{},
						"page":         1,
						"limit":        10,
						"totalRecords": 0,
						"totalPages":   1,
					},
				},
			},
		),

		"ErrorResponse": objectDefinition(
			"Gateway error response",
			[]string{"error"},
			gin.H{
				"error": stringProperty("Error message."),
			},
			gin.H{
				"error": "invalid or expired token",
			},
		),
	}
}

func operation(
	summary string,
	description string,
	tag string,
	parameters []interface{},
	body interface{},
	responseDescriptions map[string]string,
	secured bool,
) gin.H {
	params := make([]interface{}, 0)

	if parameters != nil {
		params = append(params, parameters...)
	}

	if body != nil {
		params = append(params, body)
	}

	responses := gin.H{}
	for status, description := range responseDescriptions {
		response := gin.H{
			"description": description,
		}

		if status == "200" || status == "201" {
			switch {
			case strings.Contains(description, "ProductListResponse"):
				response["schema"] = refSchema("ProductListResponse")
			case strings.Contains(description, "ProductResponse"):
				response["schema"] = refSchema("ProductResponse")
			case strings.Contains(description, "OrderResponse"):
				response["schema"] = refSchema("OrderResponse")
			case strings.Contains(description, "InventoryReservation"):
				response["schema"] = refSchema("InventoryReservation")
			case strings.Contains(description, "array of OrderResponse"):
				response["schema"] = arraySchema("#/definitions/OrderResponse")
			case strings.Contains(description, "array of InventoryReservation"):
				response["schema"] = arraySchema("#/definitions/InventoryReservation")
			case strings.Contains(description, "token"):
				response["schema"] = refSchema("LoginResponse")
			case strings.Contains(description, "Product"):
				response["schema"] = refSchema("ProductResponse")
			case strings.Contains(description, "reservation confirmation"),
				strings.Contains(description, "release confirmation"):
				response["schema"] = refSchema("ReservationResponse")
			}
		}

		if status == "400" || status == "401" || status == "403" ||
			status == "404" || status == "500" || status == "502" ||
			status == "503" {
			response["schema"] = refSchema("ErrorResponse")
		}

		responses[status] = response
	}

	result := gin.H{
		"tags":        []string{tag},
		"summary":     summary,
		"description": description,
		"parameters":  params,
		"responses":   responses,
	}

	if secured {
		result["security"] = []interface{}{
			gin.H{"BearerAuth": []string{}},
		}
	}

	return result
}

func bodyParameter(definition string, description string) gin.H {
	return gin.H{
		"name":        "body",
		"in":          "body",
		"required":    true,
		"description": description,
		"schema":      refSchema(definition),
	}
}

func pathParameter(name string, description string, dataType string) gin.H {
	return gin.H{
		"name":        name,
		"in":          "path",
		"required":    true,
		"type":        dataType,
		"description": description,
	}
}

func queryParameter(
	name string,
	description string,
	dataType string,
	required bool,
) gin.H {
	return gin.H{
		"name":        name,
		"in":          "query",
		"required":    required,
		"type":        dataType,
		"description": description,
	}
}

func objectDefinition(
	description string,
	required []string,
	properties gin.H,
	example interface{},
) gin.H {
	result := gin.H{
		"type":        "object",
		"description": description,
		"properties":  properties,
		"example":     example,
	}

	if required != nil {
		result["required"] = required
	}

	return result
}

func stringProperty(description string) gin.H {
	return gin.H{
		"type":        "string",
		"description": description,
	}
}

func integerProperty(description string) gin.H {
	return gin.H{
		"type":        "integer",
		"format":      "int64",
		"description": description,
	}
}

func numberProperty(description string) gin.H {
	return gin.H{
		"type":        "number",
		"format":      "double",
		"description": description,
	}
}

func dateProperty(description string) gin.H {
	return gin.H{
		"type":        "string",
		"format":      "date-time",
		"description": description,
	}
}

func arrayProperty(itemsRef string, description string) gin.H {
	return gin.H{
		"type":        "array",
		"description": description,
		"items":       refSchema(strings.TrimPrefix(itemsRef, "#/definitions/")),
	}
}

func refSchema(definition string) gin.H {
	return gin.H{
		"$ref": "#/definitions/" + definition,
	}
}

func arraySchema(itemsRef string) gin.H {
	return gin.H{
		"type":  "array",
		"items": refSchema(strings.TrimPrefix(itemsRef, "#/definitions/")),
	}
}

func productQueryParameters() []interface{} {
	return []interface{}{
		queryParameter("page", "Page number. Defaults to 1.", "integer", false),
		queryParameter("limit", "Number of records per page. Defaults to 10; maximum 100.", "integer", false),
		queryParameter("name", "Optional product name filter.", "string", false),
		queryParameter("min_price", "Optional minimum product price.", "number", false),
		queryParameter("max_price", "Optional maximum product price.", "number", false),
		queryParameter("in_stock", "Optional stock availability filter.", "boolean", false),
		queryParameter("sort_by", "Sort field: id, name, price, stock, created_at, or updated_at.", "string", false),
		queryParameter("sort_order", "Sort order: asc or desc.", "string", false),
	}
}

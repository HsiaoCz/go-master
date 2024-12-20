# Shopping Mall Application

A modern shopping mall application built with Go, featuring user authentication, product management, and order processing.

## Features

- User Authentication (Register/Login)
- Product Management (Create, List, View)
- Order Processing
- Role-based Access Control (Admin/User)
- JWT-based Authentication
- MongoDB Database

## Prerequisites

- Go 1.21 or later
- MongoDB
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/shopping-mall.git
cd shopping-mall
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
export MONGO_URI="mongodb://localhost:27017"
export JWT_SECRET="your-secret-key"
```

## Running the Application

1. Start MongoDB:
```bash
mongod
```

2. Run the application:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Public Endpoints

- `POST /api/register` - Register a new user
- `POST /api/login` - Login and get JWT token
- `GET /api/products` - List all products
- `GET /api/products/:id` - Get product details

### Protected Endpoints (Requires Authentication)

- `POST /api/orders` - Create a new order
- `GET /api/orders` - List user's orders
- `GET /api/orders/:id` - Get order details

### Admin Endpoints (Requires Admin Role)

- `POST /api/products` - Create a new product
- `PUT /api/orders/:id/status` - Update order status

## Authentication

Include the JWT token in the Authorization header for protected endpoints:
```
Authorization: Bearer your-jwt-token
```

## Example Requests

### Register a New User
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john_doe", "email": "john@example.com", "password": "secret123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "secret123"}'
```

### Create a Product (Admin Only)
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{"name": "Product 1", "description": "Description", "price": 99.99, "category": "Electronics", "stock": 100, "image_url": "https://example.com/image.jpg"}'
```

### Create an Order
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{"items": [{"product_id": "product-uuid", "quantity": 1}], "shipping_addr": "123 Main St", "payment_method": "credit_card"}'
```

## Error Handling

The API returns appropriate HTTP status codes and error messages in JSON format:

```json
{
  "error": "error message here"
}
```

## Security

- Passwords are hashed using bcrypt
- JWT tokens are used for authentication
- Role-based access control for admin functions
- Environment variables for sensitive configuration

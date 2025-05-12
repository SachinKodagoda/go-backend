# Go MongoDB API Backend

## Requirements

- Go 1.18 or higher
- MongoDB 4.4 or higher

## Project Structure

```
├── main.go                  # Entry point
├── models/                  # Data models
│   └── models.go
├── db/                      # Database connection and initialization
│   └── db.go
├── handlers/                # API handlers
│   ├── category_handlers.go
│   └── product_handlers.go
├── middleware/              # Middleware functions
│   └── middleware.go
├── routes/                  # API routes
│   └── routes.go
└── .env                     # Environment variables
```

## Setup Instructions

1. Clone the repository:

```bash
git clone https://github.com/yourusername/go-backend.git
cd go-backend
```

2. Install dependencies:

```bash
go mod init go-backend
go get github.com/gorilla/mux
go get go.mongodb.org/mongo-driver/mongo
go get github.com/joho/godotenv
```

3. Create a `.env` file (use the `.env.template` as a guide):

```bash
cp .env.template .env
```

4. Make sure MongoDB is running

5. Run the application:

```bash
go run main.go
```

The server will start on http://localhost:8080 by default (or the port specified in the .env file).

## API Endpoints

### Categories

- `GET /api/categories` - Get all categories

### Products

- `GET /api/products` - Get all products (with pagination, filtering, and sorting)
  - Query parameters:
    - `page`: Page number (default: 1)
    - `page_size`: Number of items per page (default: 10)
    - `category_id`: Filter by category ID
    - `category_group`: Filter by category group
    - `_sort` or `sortField`: Field to sort by
    - `_order` or `sortOrder`: Sort order (asc or desc)
    - `_start`: Starting index
    - `_limit`: Limit number of results
- `GET /api/products/{id}` - Get a product by ID
- `PUT /api/products/{id}` - Update a product

### Health Check

- `GET /api/health` - API health check

## Example API Calls

### Get All Categories

```bash
curl -X GET http://localhost:8080/api/categories
```

### Get Products with Pagination

```bash
curl -X GET "http://localhost:8080/api/products?page=1&page_size=5"
```

### Get Products by Category

```bash
curl -X GET "http://localhost:8080/api/products?category_id=2"
```

### Get Products Sorted by Name

```bash
curl -X GET "http://localhost:8080/api/products?_sort=name&_order=asc"
```

### Get a Product by ID

```bash
curl -X GET http://localhost:8080/api/products/1
```

### Update a Product

```bash
curl -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "name": "iPhone 13 Pro",
    "category_id": "2",
    "category_group": "1",
    "attributes": [
      {
        "code": "price",
        "value": 999,
        "type": "number"
      },
      {
        "code": "color",
        "value": "Gold",
        "type": "text"
      },
      {
        "code": "in_stock",
        "value": true,
        "type": "boolean"
      }
    ]
  }'
```

# 🚀 Go MongoDB API Backend

## 📋 Requirements

- Go 1.18 or higher (1.24.0 is used)
- MongoDB 4.4 or higher

## 📁 Project Structure

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

## 🛠️ Setup Instructions

1. Clone the repository:

```bash
git clone https://github.com/SachinKodagoda/go-backend.git
cd go-backend
```

2. Install dependencies: (Following is already included in the `go.mod` file)

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

Following is an example of the `.env` file:

```bash
# MongoDB connection string
MONGODB_URI=mongodb+srv://duminda:test123@cluster0.gnfih.mongodb.net
# Database name
DB_NAME=mydb
# Server port
PORT=8080
```

4. Make sure MongoDB is running.

5. Run the application:

```bash
go run main.go
```

The server will start on http://localhost:8080 by default (or the port specified in the .env file).

## 🔌 API Endpoints

### 📊 Categories

- `GET /api/categories` - Get all categories

### 🛒 Products

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

### 💓 Health Check

- `GET /api/health` - API health check

## 🔍 Example API Calls

```bash
curl -X GET http://localhost:8080/api/categories
curl -X GET "http://localhost:8080/api/products?page=1&page_size=5"
curl -X GET "http://localhost:8080/api/products?category_id=2"
curl -X GET "http://localhost:8080/api/products?_sort=name&_order=asc"
curl -X GET http://localhost:8080/api/products/1
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

## Other Notes :

- The code is written in Go and uses the MongoDB driver for database operations.
- The API is built using the Gorilla Mux router.
- The API supports pagination, filtering, and sorting for the products endpoint.
- The API includes a health check endpoint to verify if the server is running.
- The code includes error handling for database operations and API requests.
- The code uses JSON for data interchange between the server and clients.
- The code includes a middleware function for logging requests and responses.
- To build

```bash
go build -o go-backend main.go
```

- This is deployed in Render.com and the URL is

```bash
https://go-backend-s2eg.onrender.com/api/health
```

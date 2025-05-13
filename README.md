# 🚀 Go MongoDB API Backend

## 💻 Requirements

- ✅ Go 1.18+ (1.24.0 recommended)
- 🍃 MongoDB 4.4+

## 📁 Project Structure

```
├── main.go                  # Entry point
├── models/                  # Data models
│   └── models.go
├── db/                      # Database connection and initialization
│   └── db.go
├── handlers/                # API handlers
│   ├── category_handlers.go
│   ├── product_handlers.go
│   └── user_handlers.go
├── middleware/              # Middleware functions
│   └── middleware.go
├── routes/                  # API routes
│   └── routes.go
└── .env                     # Environment variables
```

## 🔑 Key Features

- 🔄 RESTful API endpoints for products and categories
- 📊 Advanced filtering, pagination, and sorting
- 🔒 Environment-based configuration
- 🧩 Modular project structure

## 🛠️ Setup Instructions

### 1️⃣ Clone the repository:

```bash
git clone https://github.com/SachinKodagoda/go-backend.git
cd go-backend
```

### 2️⃣ Install dependencies: (Following is already included in the `go.mod` file)

```bash
go mod init go-backend
go get github.com/gorilla/mux
go get go.mongodb.org/mongo-driver/mongo
go get github.com/joho/godotenv
```

### 3️⃣ Create a `.env` file (use the `.env.template` as a guide):

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

### 4️⃣ Make sure MongoDB is running.

### 5️⃣ Run the application:

```bash
go run main.go
```

The server will start on http://localhost:8080 by default (or the port specified in the .env file).

## 🔌 API Reference

### 📊 Categories

| Method | Endpoint          | Description        |
| ------ | ----------------- | ------------------ |
| GET    | `/api/categories` | Get all categories |

### 🛒 Products

| Method | Endpoint             | Description       | Query Parameters                         |
| ------ | -------------------- | ----------------- | ---------------------------------------- |
| GET    | `/api/products`      | Get all products  | `page`, `page_size`, `category_id`, etc. |
| GET    | `/api/products/{id}` | Get product by ID | -                                        |
| PUT    | `/api/products/{id}` | Update product    | -                                        |

- 📄 `page`: Page number (default: 1)
- 🔢 `page_size`: Items per page (default: 10)
- 🏷️ `category_id`: Filter by category ID
- 📊 `_sort`/`sortField`: Field to sort by
- 🔃 `_order`/`sortOrder`: Sort order (`asc` or `desc`)

### 🛒 Users

| Method | Endpoint          | Description    |
| ------ | ----------------- | -------------- |
| GET    | `/api/users`      | Get all users  |
| GET    | `/api/users/{id}` | Get user by ID |

### 💓 Health Check

- `GET /api/health` - API health check

## 🔍 Example API Calls

```bash
curl -X GET http://localhost:8080/api/categories
curl -X GET http://localhost:8080/api/users
curl -X GET http://localhost:8080/api/users?email=user@gmail.com&password=user123
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

## 💻 Development

### Building the Application

```bash
go build -o go-backend main.go
```

- This is deployed in Render.com and the URL is

https://go-backend-s2eg.onrender.com/api/health

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Category represents a product category
type Category struct {
	ID       string  `json:"id" bson:"id"`
	Name     string  `json:"name" bson:"name"`
	ParentID *string `json:"parent_id" bson:"parent_id,omitempty"`
}

// Attribute represents a product attribute with dynamic type
type Attribute struct {
	Code  string      `json:"code" bson:"code"`
	Value interface{} `json:"value" bson:"value"`
	Type  string      `json:"type" bson:"type"`
	Label string      `json:"label" bson:"label"` // Human-readable label for the attribute
}

// Product represents a product
type Product struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	CategoryID    string             `json:"category_id" bson:"category_id"`
	CategoryGroup string             `json:"category_group" bson:"category_group"`
	Attributes    []Attribute        `json:"attributes" bson:"attributes"`
}

// PaginationParams represents parameters for pagination and filtering
type PaginationParams struct {
	Page          int    `json:"page"`
	PageSize      int    `json:"page_size"`
	CategoryID    string `json:"category_id"`
	CategoryGroup string `json:"category_group"`
	SortField     string `json:"sortField"`
	SortOrder     string `json:"sortOrder"` // "asc" or "desc"
	Start         int    `json:"_start"`    // For pagination
	Limit         int    `json:"_limit"`    // For pagination
}

// ProductsResponse represents the response for paginated products
type ProductsResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// User represents a user in the system
type User struct {
	ID       string `json:"id" bson:"id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"` // In production, store hashed passwords
	Name     string `json:"name" bson:"name"`
	Role     string `json:"role" bson:"role"`
}

// UserResponse represents a user response without sensitive data
type UserResponse struct {
	ID    string `json:"id" bson:"id"`
	Email string `json:"email" bson:"email"`
	Name  string `json:"name" bson:"name"`
	Role  string `json:"role" bson:"role"`
}

package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go-backend/db"
	"go-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GET /products endpoint with pagination and filtering
func GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.GetProductsCollection()

	// Parse query parameters
	params := parseProductsQueryParams(r)

	// Build filter
	filter := bson.M{}
	if params.CategoryID != "" {
		filter["category_id"] = params.CategoryID
	}
	if params.CategoryGroup != "" {
		filter["category_group"] = params.CategoryGroup
	}

	// Build options for sorting and pagination
	findOptions := options.Find()
	if params.SortField != "" {
		sortValue := 1 // asc
		if params.SortOrder == "desc" {
			sortValue = -1 // desc
		}

		// Map API field names to MongoDB field names
		sortFieldName := params.SortField
		if params.SortField == "id" {
			sortFieldName = "_id"
		}

		findOptions.SetSort(bson.D{{Key: sortFieldName, Value: sortValue}})

		// Add case-insensitive collation for string fields
		if params.SortField != "id" && params.SortField != "_id" {
			// Use simple collation with case-insensitive comparison
			findOptions.SetCollation(&options.Collation{
				Locale:   "en",
				Strength: 2, // 2 = case-insensitive
			})
		}
	}

	// First get total count
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		http.Error(w, "Error counting products", http.StatusInternalServerError)
		return
	}

	// Apply pagination
	findOptions.SetSkip(int64(params.Start))
	findOptions.SetLimit(int64(params.Limit))

	// Execute query
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Parse results
	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		http.Error(w, "Error parsing products", http.StatusInternalServerError)
		return
	}

	// Build response
	response := models.ProductsResponse{
		Products: products,
		Total:    total,
	}

	// Return products as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// POST /products endpoint
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse request body
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields (except ID which will be generated)
	if product.Name == "" || product.CategoryID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Generate a new ObjectID for the product
	product.ID = primitive.NewObjectID()

	collection := db.GetProductsCollection()

	// Insert the product
	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}

	// Return the created product with the generated ID
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GET /products/{id} endpoint
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	collection := db.GetProductsCollection()

	// Try to convert the string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ObjectID format", http.StatusBadRequest)
		return
	}

	// Find product by ObjectID using _id field
	var product models.Product
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching product", http.StatusInternalServerError)
		}
		return
	}

	// Return product as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// PUT /products/{id} endpoint
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse request body
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Try to convert the string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ObjectID format", http.StatusBadRequest)
		return
	}

	// Ensure we use the ID from the URL
	product.ID = objectID

	collection := db.GetProductsCollection()

	// Update product
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"name":           product.Name,
		"category_id":    product.CategoryID,
		"category_group": product.CategoryGroup,
		"attributes":     product.Attributes,
	}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, "Error updating product", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Return updated product
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// Helper function to parse query parameters
func parseProductsQueryParams(r *http.Request) models.PaginationParams {
	params := models.PaginationParams{
		Page:      1,
		PageSize:  10,
		SortOrder: "asc",
	}

	// Parse page
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	// Parse page_size
	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			params.PageSize = pageSize
		}
	}

	// Parse category_id
	params.CategoryID = r.URL.Query().Get("category_id")

	// Parse category_group
	params.CategoryGroup = r.URL.Query().Get("category_group")

	// Parse sort field
	params.SortField = r.URL.Query().Get("_sort")
	if params.SortField == "" {
		params.SortField = r.URL.Query().Get("sortField")
	}

	// Parse sort order
	sortOrder := r.URL.Query().Get("_order")
	if sortOrder == "" {
		sortOrder = r.URL.Query().Get("sortOrder")
	}
	if sortOrder == "desc" {
		params.SortOrder = "desc"
	}

	// Parse pagination parameters
	if startStr := r.URL.Query().Get("_start"); startStr != "" {
		if start, err := strconv.Atoi(startStr); err == nil && start >= 0 {
			params.Start = start
		}
	} else {
		// Calculate start based on page and page_size
		params.Start = (params.Page - 1) * params.PageSize
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("_limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
		}
	} else {
		params.Limit = params.PageSize
	}

	return params
}

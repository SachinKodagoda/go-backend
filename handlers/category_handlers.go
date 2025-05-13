package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-backend/db"
	"go-backend/models"

	"go.mongodb.org/mongo-driver/bson"
)

// GET /categories endpoint
func GetCategories(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.GetCategoriesCollection()

	// Find all categories
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Parse results
	var categories []models.Category
	if err = cursor.All(ctx, &categories); err != nil {
		http.Error(w, "Error parsing categories", http.StatusInternalServerError)
		return
	}

	// Return categories as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GET /categories/{id} endpoint
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Neet to implement this
}

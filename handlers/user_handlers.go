package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go-backend/db"
	"go-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// This is developed for the interview only.
// In production level applications, we should have strong features like
// Encription/Decription, JWT accessToken/JWT refreshToken Authentication, Bearer <token>, Authorization

// GetUsers handles requests to get users, also used for authentication
func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.GetUsersCollection()

	// Build filter based on query parameters
	filter := bson.M{}

	// For authentication - check email and password
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	if email != "" {
		filter["email"] = email
	}

	// Execute query
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Parse results
	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		http.Error(w, "Error parsing users", http.StatusInternalServerError)
		return
	}

	// For authentication - if password was provided, we need to check it
	if password != "" && len(users) > 0 {
		authenticated := false
		var authenticatedUser models.User

		for _, user := range users {
			// In production, use password hashing and proper comparison
			if user.Password == password {
				authenticated = true
				authenticatedUser = user
				break
			}
		}

		if !authenticated {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Return only the authenticated user without password
		userResp := models.UserResponse{
			ID:    authenticatedUser.ID,
			Email: authenticatedUser.Email,
			Name:  authenticatedUser.Name,
			Role:  authenticatedUser.Role,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userResp)
		return
	}

	// Convert to response objects (without passwords)
	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		})
	}

	// Return users as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponses)
}

// GetUserByID retrieves a single user by ID
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	collection := db.GetUsersCollection()

	// Find user by ID
	var user models.User
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching user", http.StatusInternalServerError)
		}
		return
	}

	// Return user without password
	userResp := models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}

	// Return user as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userResp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

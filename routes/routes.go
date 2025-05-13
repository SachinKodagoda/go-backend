// routes/routes.go
package routes

import (
	"net/http"

	"go-backend/handlers"
	"go-backend/middleware"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up the API routes
func RegisterRoutes(router *mux.Router) {
	// Apply global middleware
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.JSONContentTypeMiddleware)

	// Create API subrouter
	api := router.PathPrefix("/api").Subrouter()

	// Categories endpoints
	api.HandleFunc("/categories", handlers.GetCategories).Methods("GET", "OPTIONS")

	// Products endpoints
	api.HandleFunc("/products", handlers.GetProducts).Methods("GET", "OPTIONS")
	api.HandleFunc("/products", handlers.CreateProduct).Methods("POST", "OPTIONS")
	api.HandleFunc("/products/{id}", handlers.GetProductByID).Methods("GET", "OPTIONS")
	api.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT", "OPTIONS")

	// Users endpoints
	api.HandleFunc("/users", handlers.GetUsers).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/{id}", handlers.GetUserByID).Methods("GET", "OPTIONS")

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Handle 404
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"endpoint not found"}`))
	})
}

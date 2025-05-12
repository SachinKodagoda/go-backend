package main

import (
	"log"
	"net/http"
	"os"

	"go-backend/db"
	"go-backend/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to MongoDB
	err = db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Disconnect()

	// Initialize database with sample data if needed
	err = db.InitializeDatabase()
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
	}

	// Create router
	router := mux.NewRouter()

	// Register routes
	routes.RegisterRoutes(router)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

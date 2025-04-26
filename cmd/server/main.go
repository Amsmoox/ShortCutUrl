package main

import (
	"log"
	"net/http"
	"os"
	"shorturl/internal/db"
	"shorturl/internal/handlers"
	"shorturl/internal/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file specifically for main, db.InitDB also loads it but good practice here too
	err := godotenv.Load() // Load .env from current working dir (project root)
	if err != nil {
		log.Println("Warning: .env file not found in project root, relying on environment variables or db init load.")
	}

	// Initialize Database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.DB.Close()

	// Get Database connection
	dbConn := db.GetDB()

	// Initialize Handlers with DB connection
	handlers.NewHandlers(dbConn)

	// Create Router
	r := mux.NewRouter()

	// Setup Routes
	routes.SetupRoutes(r)

	// Get Port from environment or default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Default port if not specified
		log.Printf("Defaulting to port %s", port)
	}

	// Start Server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 
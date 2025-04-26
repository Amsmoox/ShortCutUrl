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

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func main() {
	// Load .env file specifically for main, db.InitDB also loads it but good practice here too
	err := godotenv.Load() // Load .env from current working dir (project root)
	if err != nil {
		log.Println(ColorYellow + "Warning: .env file not found in project root, relying on environment variables or db init load." + ColorReset)
	}

	// Initialize Database
	if err := db.InitDB(); err != nil {
		log.Fatalf(ColorRed+"Failed to initialize database: %v"+ColorReset, err)
	}
	defer db.DB.Close()

	// Initialize Database Tables (creates if not exists)
	if err := db.InitTables(); err != nil {
		log.Fatalf(ColorRed+"Failed to initialize database tables: %v"+ColorReset, err)
	}

	// Get Database connection
	dbConn := db.GetDB()

	// Initialize Handlers with DB connection
	handlers.NewHandlers(dbConn)

	// Create Router
	r := mux.NewRouter()

	// Setup static file serving before specific routes
	// This serves files from the "static" directory directly.
	// For example, /style.css would serve static/style.css
	fs := http.FileServer(http.Dir("./static/"))
	// Use PathPrefix to serve static files; StripPrefix removes /static/ from the request path
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Setup Application Routes (including the root route "/")
	routes.SetupRoutes(r)

	// Get Port from environment or default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Default port if not specified
		log.Printf(ColorYellow+"Warning: APP_PORT not set in .env. Defaulting to port %s"+ColorReset, port)
	}

	// Start Server
	log.Printf(ColorGreen+"Server starting on port %s"+ColorReset, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf(ColorRed+"Failed to start server: %v"+ColorReset, err)
	}
} 
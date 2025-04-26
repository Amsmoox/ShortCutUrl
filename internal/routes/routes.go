package routes

import (
	"net/http"
	"shorturl/internal/handlers"

	"github.com/gorilla/mux"
)

// SetupRoutes configures the application routes.
func SetupRoutes(router *mux.Router) {
	// Serve static files (like index.html)
	// We'll configure the actual file server in main.go
	// This handler specifically serves index.html for the root path
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	}).Methods("GET")

	// Route for creating a short URL (POST /api/shorten)
	router.HandleFunc("/api/shorten", handlers.CreateShortURL).Methods("POST")

	// Route for redirecting using the short code (GET /{shortCode})
	router.HandleFunc("/{shortCode:[a-zA-Z0-9]{6}}", handlers.RedirectURL).Methods("GET")

	// Optional: Add a simple health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
} 
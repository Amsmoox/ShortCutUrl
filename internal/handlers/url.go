package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"shorturl/internal/models"
	"time"

	"github.com/gorilla/mux"
)

const shortCodeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortCodeLength = 6

var db *sql.DB

// NewHandlers sets up the database connection for handlers.
func NewHandlers(database *sql.DB) {
	db = database
	rand.Seed(time.Now().UnixNano())
}

// generateShortCode creates a random short code.
func generateShortCode() string {
	b := make([]byte, shortCodeLength)
	for i := range b {
		b[i] = shortCodeChars[rand.Intn(len(shortCodeChars))]
	}
	return string(b)
}

// CreateShortURL handles requests to create a new short URL.
func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	// TODO: Add validation for URL format

	shortCode := generateShortCode()
	// TODO: Check if shortCode already exists and regenerate if necessary

	url := models.URL{
		OriginalURL: req.URL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
	}

	// TODO: Implement database insertion
	// For now, just log and return the generated short code
	log.Printf("Generated short code '%s' for URL '%s'", url.ShortCode, url.OriginalURL)

	// Simulate successful creation for now
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"short_code": url.ShortCode})
}

// RedirectURL handles requests to redirect to the original URL.
func RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	if shortCode == "" {
		http.Error(w, "Short code cannot be empty", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to redirect using short code: %s", shortCode)

	// TODO: Implement database lookup for the original URL based on shortCode
	// For now, simulate a lookup and redirect
	originalURL := "https://example.com" // Placeholder

	if originalURL == "" { // Simulate not found
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound) // Use 302 Found for temporary redirect
} 
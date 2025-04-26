package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"shorturl/internal/models"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Import necessary for side-effects (driver registration)
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

// respondWithError sends a JSON error response.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sends a JSON response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error marshalling response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// shortCodeExists checks if a short code already exists in the database.
// TODO: Implement the actual database query when the schema is ready.
func shortCodeExists(code string) (bool, error) {
	log.Printf("Checking existence of short code: %s (DB check not implemented yet)", code)
	// Placeholder: Query the database, e.g.:
	// SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)
	// Handle sql.ErrNoRows for non-existence vs. other errors.
	return false, nil // Simulate non-existence for now
}

// CreateShortURL handles requests to create a new short URL.
func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.URL == "" {
		respondWithError(w, http.StatusBadRequest, "URL cannot be empty")
		return
	}

	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		respondWithError(w, http.StatusBadRequest, "Invalid URL format (must start with http:// or https://)")
		return
	}

	// TODO: Add more robust URL validation

	var shortCode string
	var exists bool
	maxAttempts := 5 // Prevent infinite loops

	for i := 0; i < maxAttempts; i++ {
		shortCode = generateShortCode()
		var checkErr error
		exists, checkErr = shortCodeExists(shortCode)
		if checkErr != nil {
			log.Printf("Error checking short code existence: %v", checkErr)
			respondWithError(w, http.StatusInternalServerError, "Internal server error checking short code")
			return
		}
		if !exists {
			break // Found unique code
		}
		log.Printf("Short code %s already exists, generating a new one...", shortCode)
	}

	if exists {
		log.Println("Failed to generate a unique short code after several attempts")
		respondWithError(w, http.StatusInternalServerError, "Failed to generate unique short code after several attempts")
		return
	}

	url := models.URL{
		OriginalURL: req.URL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
	}

	// TODO: Implement database insertion
	log.Printf("Generated unique short code '%s' for URL '%s'", url.ShortCode, url.OriginalURL)

	// Simulate successful creation for now
	respondWithJSON(w, http.StatusCreated, map[string]string{"short_code": url.ShortCode})
}

// RedirectURL handles requests to redirect to the original URL.
func RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	if shortCode == "" {
		respondWithError(w, http.StatusBadRequest, "Short code cannot be empty")
		return
	}

	log.Printf("Attempting to redirect using short code: %s", shortCode)

	// TODO: Implement database lookup for the original URL based on shortCode
	originalURL := "https://example.com" // Placeholder
	found := true                         // Simulate finding the URL

	if !found { // Simulate not found
		respondWithError(w, http.StatusNotFound, "Short code not found")
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
} 
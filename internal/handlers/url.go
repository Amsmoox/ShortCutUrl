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
	// Placeholder: Query the database to see if 'code' exists in the 'urls' table.
	// Example Query: SELECT 1 FROM urls WHERE short_code = $1 LIMIT 1
	// var exists bool // Removed unused variable
	// err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)", code).Scan(&exists)
	// if err != nil {
	// 	 if err == sql.ErrNoRows {
	// 		 return false, nil // Doesn't exist
	// 	 }
	// 	 return false, fmt.Errorf("database error checking short code: %w", err)
	// }
	// For now, assume it doesn't exist to allow generation.
	// In a real scenario, you'd handle potential DB errors here.
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
	
	// Basic check if URL seems valid (this should be more robust)
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		respondWithError(w, http.StatusBadRequest, "Invalid URL format (must start with http:// or https://)")
		return
	}

	// TODO: Add more robust URL validation

	var shortCode string
	var exists bool
	// var err error // Removed unused variable
	maxAttempts := 5 // Prevent infinite loops in case of high collision rate

	for i := 0; i < maxAttempts; i++ {
		shortCode = generateShortCode()
		var checkErr error // Use a separate variable for the error in this scope
		exists, checkErr = shortCodeExists(shortCode) 
		if checkErr != nil {
			log.Printf("Error checking short code existence: %v", checkErr)
			respondWithError(w, http.StatusInternalServerError, "Internal server error checking short code")
			return
		}
		if !exists {
			break // Found a unique code
		}
		log.Printf("Short code %s already exists, generating a new one...", shortCode)
	}

	// Need to check the error from the *last* call to shortCodeExists if the loop finished
	// However, the logic guarantees 'exists' is false if we break, and true if the loop completes.
	// The error check inside the loop already handles DB errors during the search.

	if exists {
		log.Println("Failed to generate a unique short code after several attempts")
		respondWithError(w, http.StatusInternalServerError, "Failed to generate unique short code after several attempts")
		return
	}

	url := models.URL{
		OriginalURL: req.URL,
		ShortCode:   shortCode, // Use the validated unique short code
		CreatedAt:   time.Now(),
	}

	// TODO: Implement database insertion
	// For now, just log and return the generated short code
	log.Printf("Generated unique short code '%s' for URL '%s'", url.ShortCode, url.OriginalURL)

	// Simulate successful creation for now
	// w.Header().Set("Content-Type", "application/json") // Set by respondWithJSON
	// w.WriteHeader(http.StatusCreated) // Set by respondWithJSON
	// json.NewEncoder(w).Encode(map[string]string{"short_code": url.ShortCode}) // Use respondWithJSON
	respondWithJSON(w, http.StatusCreated, map[string]string{"short_code": url.ShortCode})
}

// RedirectURL handles requests to redirect to the original URL.
func RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	if shortCode == "" {
		// Although the route regex should prevent this, good to have a check
		respondWithError(w, http.StatusBadRequest, "Short code cannot be empty")
		return
	}

	log.Printf("Attempting to redirect using short code: %s", shortCode)

	// TODO: Implement database lookup for the original URL based on shortCode
	// For now, simulate a lookup and redirect
	originalURL := "https://example.com" // Placeholder
	found := true // Simulate finding the URL

	if !found { // Simulate not found
		respondWithError(w, http.StatusNotFound, "Short code not found")
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound) // Use 302 Found for temporary redirect
} 
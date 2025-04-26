package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"shorturl/internal/models"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Import necessary for side-effects (driver registration)
)

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
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
		log.Printf(ColorRed+"Error marshalling JSON response: %v"+ColorReset, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error marshalling response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// shortCodeExists checks if a short code already exists in the database.
func shortCodeExists(code string) (bool, error) {
	log.Printf(ColorBlue+"Checking database for existence of short code: %s"+ColorReset, code)
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)"
	err := db.QueryRow(query, code).Scan(&exists)
	if err != nil {
		// An error here is a server error, not just non-existence
		return false, fmt.Errorf("database error checking short code: %w", err)
	}
	return exists, nil
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
			log.Printf(ColorRed+"Error checking short code existence: %v"+ColorReset, checkErr)
			respondWithError(w, http.StatusInternalServerError, "Internal server error checking short code")
			return
		}
		if !exists {
			break // Found unique code
		}
		log.Printf(ColorYellow+"Short code %s already exists, generating a new one..."+ColorReset, shortCode)
	}

	if exists {
		log.Println(ColorRed + "Failed to generate a unique short code after several attempts" + ColorReset)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate unique short code after several attempts")
		return
	}

	url := models.URL{
		OriginalURL: req.URL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(), // Handled by DB default now, but good to have
	}

	// Insert into database
	query := `INSERT INTO urls (original_url, short_code) VALUES ($1, $2) RETURNING id, created_at`
	err := db.QueryRow(query, url.OriginalURL, url.ShortCode).Scan(&url.ID, &url.CreatedAt)
	if err != nil {
		log.Printf(ColorRed+"Failed to insert URL into database: %v"+ColorReset, err)
		// Check for duplicate short_code violation (specific error might depend on DB driver)
		if strings.Contains(err.Error(), "unique constraint") {
			respondWithError(w, http.StatusConflict, "Failed to generate unique short code, please try again.")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to save URL")
		}
		return
	}

	log.Printf(ColorGreen+"Successfully inserted URL '%s' with short code '%s' (ID: %d)"+ColorReset, url.OriginalURL, url.ShortCode, url.ID)

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

	log.Printf(ColorBlue+"Attempting to redirect using short code: %s"+ColorReset, shortCode)

	var originalURL string
	query := "SELECT original_url FROM urls WHERE short_code = $1"
	err := db.QueryRow(query, shortCode).Scan(&originalURL)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf(ColorYellow+"Short code not found in database: %s"+ColorReset, shortCode)
			respondWithError(w, http.StatusNotFound, "Short code not found")
		} else {
			log.Printf(ColorRed+"Database error looking up short code %s: %v"+ColorReset, shortCode, err)
			respondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	// Successfully found the original URL
	log.Printf(ColorGreen+"Redirecting short code %s to %s"+ColorReset, shortCode, originalURL)
	http.Redirect(w, r, originalURL, http.StatusFound)
} 
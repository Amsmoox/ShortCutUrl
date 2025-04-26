package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

var DB *sql.DB

// InitDB initializes the database connection pool.
func InitDB() error {
	// Load .env file from the root directory
	err := godotenv.Load() // Assumes .env is in the current working directory when the app starts
	if err != nil {
		// If .env is not found, try loading from one level up (common when running tests or from cmd/)
		err = godotenv.Load("../.env")
		if err != nil {
			log.Println(ColorYellow + "Warning: .env file not found, relying on environment variables." + ColorReset)
			// Continue without error if .env is optional or vars are set externally
		}
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close() // Close the connection if ping fails
		return fmt.Errorf(ColorRed+"failed to ping database: %w"+ColorReset, err)
	}

	DB = db
	log.Println(ColorGreen + "Database connection established successfully." + ColorReset)
	return nil
}

// InitTables creates required database tables if they don't exist.
func InitTables() error {
	// SQL statement to create the urls table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		original_url TEXT NOT NULL,
		short_code VARCHAR(6) UNIQUE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	// Execute the SQL
	_, err := DB.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf(ColorRed+"failed to create tables: %w"+ColorReset, err)
	}

	log.Println(ColorGreen + "Database tables initialized successfully." + ColorReset)
	return nil
}

// GetDB returns the database connection pool.
func GetDB() *sql.DB {
	return DB
} 
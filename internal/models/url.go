package models

import (
	"time"
)

// URL represents the structure for a shortened URL in the database.
type URL struct {
	ID          int       `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
} 
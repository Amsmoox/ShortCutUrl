# ⚡ QuickLink - Go URL Shortener

A fast, minimalist URL shortening service built with Go. Features a clean web interface and a REST API.

## ✨ Features

*   **Shorten URLs:** Convert long URLs into compact short codes
*   **Clean UI:** User-friendly frontend with copy to clipboard functionality
*   **API Integration:** REST API for programmatic access
*   **Redirection:** Fast redirects from short codes to original URLs
*   **Smart Generation:** Automatic unique short code generation with database checks
*   **Colored Logs:** Color-coded logging for easy debugging
*   **Auto-setup:** Automatic database table creation on startup
*   **Environment Config:** Simple configuration via `.env` file

## 📂 Project Structure

```
url-shortener/
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── internal/
│   ├── db/
│   │   ├── db.go           # Database connection logic
│   │   └── migrations.sql  # SQL migrations documentation
│   ├── handlers/
│   │   └── url.go          # HTTP request handlers
│   ├── models/
│   │   └── url.go          # Data structures (URL model)
│   └── routes/
│       └── routes.go       # Route definitions
├── static/
│   └── index.html          # Web frontend
├── .env                    # Environment variables
├── .gitignore              # Git ignore rules
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # This file
```

## 🚀 Getting Started

### Prerequisites

*   **Go:** Version 1.18 or higher recommended
*   **PostgreSQL:** Running instance with database created
*   **Git:** For cloning the repository

### Setup

1.  **Clone the repository:**
    ```bash
    git clone <your-repository-url>
    cd url-shortener
    ```

2.  **Create Environment File:**
    Create a `.env` file in the project root with your database configuration:
    ```
    DB_NAME=shorturl
    DB_USER=postgres
    DB_PASSWORD=your-password
    DB_HOST=localhost
    DB_PORT=5432
    APP_PORT=8080 
    ```

3.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

4.  **Run the Application:**
    ```bash
    go run cmd/server/main.go
    ```
    
The server will start on `http://localhost:8080` (or your configured port). The application will automatically:
- Connect to your PostgreSQL database
- Create the necessary `urls` table if it doesn't exist
- Start serving the web interface and API

## 📝 Usage

### Web Interface

1. Open your browser and navigate to `http://localhost:8080`
2. Enter a URL in the input field (must start with http:// or https://)
3. Click "Shorten URL"
4. Copy the generated short URL using the "Copy" button

### API

**Create a Short URL:**
```bash
curl -X POST -H "Content-Type: application/json" -d '{"url":"https://example.com/long/url"}' http://localhost:8080/api/shorten
```

Response:
```json
{
  "short_code": "abcdef"
}
```

### Redirection

To use a short URL, simply access `http://localhost:8080/{short_code}` in your browser or follow the link from the web interface.

## ⚡ Performance Notes

- The application uses connection pooling for database efficiency
- Short code lookups are optimized with database indexing
- Color-coded logging helps identify potential issues quickly

## 📋 Future Improvements

- Add analytics for tracking link usage
- Implement custom short code selection
- Add user authentication for managing links
- Create expiration dates for links
- Add rate limiting for API protection

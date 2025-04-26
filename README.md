# ⚡ QuickLink - Go URL Shortener

A simple URL shortening service built with Go, featuring a web interface and a REST API.

## ✨ Features

*   **Shorten URLs:** Convert long URLs into compact short codes.
*   **Web Interface:** Easy-to-use frontend to paste URLs and get shortened links.
*   **API Endpoint:** `/api/shorten` for programmatic URL shortening.
*   **Redirection:** Short codes redirect users to the original URL.
*   **Unique Code Generation:** Attempts to generate unique short codes (database check logic is planned).
*   **Copy to Clipboard:** Quickly copy the generated short URL from the web UI.
*   **Environment-based Configuration:** Database and app settings configured via a `.env` file.

## 📂 Project Structure

```
url-shortener/
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── internal/
│   ├── db/
│   │   └── db.go           # Database connection logic
│   ├── handlers/
│   │   └── url.go          # HTTP request handlers
│   ├── models/
│   │   └── url.go          # Data structures (URL model)
│   └── routes/
│       └── routes.go       # Route definitions (using gorilla/mux)
├── static/
│   └── index.html        # Simple HTML frontend
├── .env                    # Environment variables (DB config, ports - MUST BE CREATED)
├── .gitignore              # Git ignore rules
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # This file
```

## 🚀 Getting Started

### Prerequisites

*   **Go:** Version 1.18 or higher recommended.
*   **PostgreSQL:** A running instance is required for database storage.
*   **Git:** For cloning the repository.

### Setup

1.  **Clone the repository:**
    ```bash
    git clone <your-repository-url>
    cd url-shortener
    ```

2.  **Create Environment File:**
    Copy the example or create a `.env` file in the project root:
    ```bash
    cp .env.example .env # If you create an example file
    # OR create .env manually
    ```
    Edit `.env` with your database credentials and desired application port:
    ```dotenv
    DB_NAME=shorturl
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_HOST=localhost
    DB_PORT=5432
    APP_PORT=8080 
    ```

3.  **Set up Database:**
    *   Ensure your PostgreSQL server is running.
    *   Connect to PostgreSQL and create the database specified in `DB_NAME` (`shorturl` by default).
    *   **(TODO)** Create the necessary `urls` table. The schema definition needs to be added.
      ```sql
      -- Example table schema (adjust as needed)
      CREATE TABLE IF NOT EXISTS urls (
          id SERIAL PRIMARY KEY,
          original_url TEXT NOT NULL,
          short_code VARCHAR(6) UNIQUE NOT NULL, -- Ensure uniqueness
          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
      );
      ```

4.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

### Running the Application

```bash
go run cmd/server/main.go
```

The server will start, typically on `http://localhost:8080` (or the `APP_PORT` you specified).

##  kullanım

*   **Web Interface:** Open your browser and navigate to `http://localhost:8080`.
    *   Paste a long URL into the input field.
    *   Click "Shorten URL".
    *   The shortened URL will appear, along with a button to copy it.
*   **API:** Send a POST request to `/api/shorten` with a JSON body:
    ```json
    {
      "url": "https://your-long-url.com/example"
    }
    ```
    The response will be:
    ```json
    {
      "short_code": "abcdef"
    }
    ```
*   **Redirection:** Accessing `http://localhost:8080/{short_code}` (e.g., `http://localhost:8080/abcdef`) will redirect to the original URL (once database interaction is fully implemented).

## 🚧 TODO

*   Implement database insertion logic in `CreateShortURL` handler.
*   Implement database lookup logic in `RedirectURL` handler.
*   Implement the actual database check in the `shortCodeExists` function.
*   Add database schema creation script/migration.
*   Add input validation for URLs.
*   Add more robust error handling.
*   Consider adding rate limiting.

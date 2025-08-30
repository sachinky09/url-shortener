package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"shortUrl"`
	Code     string `json:"code"`
}

var db *pgxpool.Pool

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func generateCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	code := base64.URLEncoding.EncodeToString(bytes)
	return strings.TrimRight(code, "=")[:8]
}

func sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, message string, status int) {
	sendJSON(w, map[string]string{"error": message}, status)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		sendError(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Add protocol if missing
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		req.URL = "https://" + req.URL
	}

	// Check if URL already exists
	var existingCode string
	err := db.QueryRow(context.Background(),
		"SELECT code FROM urls WHERE long_url = $1 LIMIT 1", req.URL).Scan(&existingCode)
	
	if err == nil {
		resp := ShortenResponse{
			ShortURL: "http://localhost:8080/" + existingCode,
			Code:     existingCode,
		}
		sendJSON(w, resp, http.StatusOK)
		return
	}

	// Generate new code
	code := generateCode()
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	_, err = db.Exec(context.Background(),
		"INSERT INTO urls (code, long_url, created_at) VALUES ($1, $2, NOW())",
		code, req.URL)

	if err != nil {
		log.Printf("Database error: %v", err)
		sendError(w, "Database error", http.StatusInternalServerError)
		return
	}

	resp := ShortenResponse{
		ShortURL: baseURL + "/" + code,
		Code:     code,
	}
	sendJSON(w, resp, http.StatusOK)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" || code == "favicon.ico" {
		sendError(w, "Short code required", http.StatusBadRequest)
		return
	}

	var longURL string
	err := db.QueryRow(context.Background(),
		"SELECT long_url FROM urls WHERE code = $1", code).Scan(&longURL)

	if err != nil {
		sendError(w, "URL not found", http.StatusNotFound)
		return
	}

	// Increment clicks async
	go func() {
		db.Exec(context.Background(),
			"UPDATE urls SET clicks = clicks + 1 WHERE code = $1", code)
	}()

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	var err error
	db, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE NOT NULL,
			long_url TEXT NOT NULL,
			clicks INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW()
		)`)
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", shortenHandler)
	mux.HandleFunc("/", redirectHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(mux)))
}
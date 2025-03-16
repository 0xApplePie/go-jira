package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/0xApplePie/go-jira/internal/server/graphql"
	"github.com/0xApplePie/go-jira/internal/store"
	"github.com/graphql-go/handler"
)

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func waitForDB(dbURL string) error {
	var db *sql.DB
	var err error
	
	// Try to connect to the database with retries
	for i := 0; i < 30; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Printf("Failed to open DB connection: %v", err)
			time.Sleep(time.Second)
			continue
		}
		
		err = db.Ping()
		if err == nil {
			db.Close()
			return nil
		}
		
		log.Printf("Failed to connect to DB: %v", err)
		db.Close()
		time.Sleep(time.Second)
	}
	
	return err
}

func main() {
	// Get database URL from environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Wait for database to be ready
	log.Println("Waiting for database connection...")
	if err := waitForDB(dbURL); err != nil {
		log.Fatalf("Database connection failed after retries: %v", err)
	}
	log.Println("Database connection established")

	// Initialize store
	store, err := store.NewPostgresStore(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}
	defer store.Close()

	// Create GraphQL schema
	schema, err := graphql.NewSchema(store)
	if err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Create GraphQL handler
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Wrap the GraphQL handler with CORS middleware
	http.Handle("/graphql", corsMiddleware(h))

	// Start server
	log.Printf("Server starting on http://localhost:8080/graphql")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
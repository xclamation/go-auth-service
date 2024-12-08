package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/xclamation/go-auth-service/internal/auth"
	"github.com/xclamation/go-auth-service/internal/database"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database pool
	dbURL := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// Initialize the database
	db := database.New(pool)

	// Initialize the auth service
	authService := auth.NewAuthService(db)

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/token", authService.GenerateTokenPair).Methods("POST")
	router.HandleFunc("/refresh", authService.RefreshTokenPair).Methods("POST")

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

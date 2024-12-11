package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/xclamation/go-auth-service/internal/auth"
	"github.com/xclamation/go-auth-service/internal/database"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Error loading .env file")
	}

	// Initialize the logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Initialize the database pool
	dbURL := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to create connection pool")
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
	logrus.Info("Starting server on :8080")
	logrus.Fatal(http.ListenAndServe(":8080", router))
}

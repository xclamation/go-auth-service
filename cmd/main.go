package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

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
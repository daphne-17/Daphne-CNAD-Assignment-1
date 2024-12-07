package main

import (
	"fmt"
	"log"
	"net/http"
	"user-service/database"
	"user-service/handlers"
	"user-service/middleware"

	"github.com/rs/cors"
)

func main() {
	// Initialize database connections
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer database.CloseDB() // Ensure all databases are closed on exit

	// Create a router
	mux := http.NewServeMux()

	// Register routes
	mux.Handle("/users/profile", middleware.Authenticate(http.HandlerFunc(handlers.GetUserProfileHandler)))
	mux.Handle("/users/update", middleware.Authenticate(http.HandlerFunc(handlers.UpdateUserDetailsHandler)))
	mux.Handle("/users/membership", middleware.Authenticate(http.HandlerFunc(handlers.ViewMembershipHandler)))
	mux.Handle("/users/rentalhistory", middleware.Authenticate(http.HandlerFunc(handlers.RentalHistoryHandler)))

	// Enable CORS using rs/cors library
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)

	// Start the server
	fmt.Println("User-Service is running on port 8081")
	if err := http.ListenAndServe(":8081", corsHandler); err != nil {
		log.Fatalf("Failed to start User-Service: %v", err)
	}
}

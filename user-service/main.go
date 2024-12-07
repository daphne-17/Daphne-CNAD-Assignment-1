package main

import (
	"fmt"
	"log"
	"net/http"
	"user-service/database"
	"user-service/handlers"
	"user-service/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize database connections
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer database.CloseDB() // Ensure all databases are closed on exit

	// Create a Gorilla Mux router
	router := mux.NewRouter()

	// Register routes with middleware
	router.Handle("/users/profile", middleware.Authenticate(http.HandlerFunc(handlers.GetUserProfileHandler))).Methods("GET")
	router.Handle("/users/update", middleware.Authenticate(http.HandlerFunc(handlers.UpdateUserDetailsHandler))).Methods("PUT")
	router.Handle("/users/membership", middleware.Authenticate(http.HandlerFunc(handlers.ViewMembershipHandler))).Methods("GET")
	router.Handle("/users/rentalhistory", middleware.Authenticate(http.HandlerFunc(handlers.RentalHistoryHandler))).Methods("GET")

	// Enable CORS using rs/cors library
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins, restrict in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	// Start the server
	fmt.Println("User-Service is running on port 8081")
	if err := http.ListenAndServe(":8081", corsHandler); err != nil {
		log.Fatalf("Failed to start User-Service: %v", err)
	}
}

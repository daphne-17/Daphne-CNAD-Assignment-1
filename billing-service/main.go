package main

import (
	"billing-service/database"
	"billing-service/handlers"
	"billing-service/middleware"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize database connections
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer database.CloseDB() // Ensure databases are closed on exit

	// Create a new router
	router := mux.NewRouter()

	// Define routes
	router.Handle("/billing/calculate", middleware.Authenticate(http.HandlerFunc(handlers.CalculateCostHandler))).Methods("POST")
	router.Handle("/billing/confirm", middleware.Authenticate(http.HandlerFunc(handlers.ConfirmReservationHandler))).Methods("POST")

	// Add CORS middleware
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins; restrict in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	handler := corsOptions.Handler(router)

	// Start the server
	port := ":8083"
	fmt.Printf("Billing service is running on port %s\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

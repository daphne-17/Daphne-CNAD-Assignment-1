package main

import (
	"fmt"
	"log"
	"net/http"
	"vehicle-service/database"
	"vehicle-service/handlers"
	"vehicle-service/middleware"

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

	// Add middleware
	router.Use(middleware.Authenticate)

	// Define routes
	router.HandleFunc("/vehicles/available", handlers.FetchAvailableVehicles).Methods("GET")
	router.HandleFunc("/reservations/user", handlers.GetUserReservationsHandler).Methods("GET")
	router.HandleFunc("/reservations", handlers.CreateReservationHandler).Methods("POST")
	router.HandleFunc("/reservations/modify", handlers.ModifyReservationHandler).Methods("PUT")
	router.HandleFunc("/reservations/cancel", handlers.CancelReservationHandler).Methods("DELETE")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Update with your frontend's URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	// Start the server
	port := ":8082"
	fmt.Printf("Vehicle service is running on port %s\n", port)
	if err := http.ListenAndServe(port, corsHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

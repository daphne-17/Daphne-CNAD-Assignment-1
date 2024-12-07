package main

import (
	"auth-service/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/signup", handlers.SignupHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins (use specific domains in production)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	// Start the server
	log.Println("Auth-Service is running on port 8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("Failed to start Auth-Service: %v", err)
	}
}

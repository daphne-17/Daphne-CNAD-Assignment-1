package main

import (
	"auth-service/handlers"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/signup", handlers.SignupHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)

	// Wrap the mux with CORS middleware
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	http.ListenAndServe(":8080", handler)
}

package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"user-service/utils"
)

type contextKey string

const userKey contextKey = "user_id"

// Authenticate middleware validates the JWT and extracts the user ID
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the token from the Authorization header
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Println("Authorization token is missing")
			http.Error(w, "Authorization token missing", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		token = strings.TrimPrefix(token, "Bearer ")
		log.Println("Token received:", token)

		// Validate the token and extract the user ID
		userID, err := utils.ValidateToken(token)
		if err != nil {
			log.Printf("Token validation failed: %v\n", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		log.Println("Extracted UserID from token:", userID)

		// Add UserID to the request context
		ctx := context.WithValue(r.Context(), userKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID retrieves the authenticated UserID from the request context
func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(userKey).(string)
	log.Printf("GetUserID: Retrieved UserID = %v, Present = %v\n", userID, ok)
	if ok {
		return userID
	}
	log.Println("GetUserID: UserID not found in context")
	return ""
}

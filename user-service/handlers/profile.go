package handlers

import (
	"encoding/json"
	"net/http"
	"user-service/database"
	"user-service/middleware"
)

// GetUserProfileHandler retrieves the user's profile details
func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
        SELECT Name, Email, PhoneNumber
        FROM Users
        WHERE UserID = ?`

	var name, email, phone string
	err := database.GetUserDB().QueryRow(query, userID).Scan(&name, &email, &phone)
	if err != nil {
		http.Error(w, "Failed to retrieve profile details", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"name":  name,
		"email": email,
		"phone": phone,
	})
}

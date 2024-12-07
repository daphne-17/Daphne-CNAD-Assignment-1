package handlers

import (
	"encoding/json"
	"net/http"
	"user-service/database"
	"user-service/middleware"
)

func UpdateUserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := "UPDATE Users SET Name = ?, Email = ?, PhoneNumber = ? WHERE UserID = ?"
	_, err := database.GetUserDB().Exec(query, req.Name, req.Email, req.Phone, userID)
	if err != nil {
		http.Error(w, "Failed to update user details", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User details updated successfully"))
}

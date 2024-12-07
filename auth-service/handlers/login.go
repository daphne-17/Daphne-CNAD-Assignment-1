package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"auth-service/database"
	"auth-service/models"
	"auth-service/utils"

	"golang.org/x/crypto/bcrypt"
)

// LoginHandler handles user login and returns a JWT token on successful authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch user data from the database
	query := "SELECT UserID, PasswordHash FROM Users WHERE Email = ? OR PhoneNumber = ?"
	var user models.User
	err := database.GetDB().QueryRow(query, req.Identifier, req.Identifier).Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send response with UserID and token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"user_id": user.ID,
		"token":   token,
		"expires": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	})
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"auth-service/database"
	"auth-service/models"

	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Generate the next UserID
	nextUserID, err := getNextUserID(database.GetDB())
	if err != nil {
		http.Error(w, "Error generating user ID", http.StatusInternalServerError)
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Insert user into the database with default TierID "T1"
	query := "INSERT INTO Users (UserID, Name, Email, PhoneNumber, PasswordHash, TierID) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = database.GetDB().Exec(query, nextUserID, req.Username, req.Email, req.Phone, hashedPassword, "T1")
	if err != nil {
		http.Error(w, "Error saving user to the database", http.StatusInternalServerError)
		fmt.Printf("Database error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Signup successful"))
}

// getNextUserID generates the next UserID based on the highest existing UserID.
func getNextUserID(db *sql.DB) (string, error) {
	var lastID string
	err := db.QueryRow("SELECT UserID FROM Users ORDER BY CAST(SUBSTRING(UserID, 2) AS UNSIGNED) DESC LIMIT 1").Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No users exist, start at U6
			return "U6", nil
		}
		return "", err
	}

	// Extract the numeric part of the UserID and increment it
	num, err := strconv.Atoi(lastID[1:]) // Skip the "U" prefix
	if err != nil {
		return "", err
	}

	nextID := fmt.Sprintf("U%d", num+1) // Generate the next ID without leading zeroes
	return nextID, nil
}

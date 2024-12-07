package models

// SignupRequest represents the data required for user signup.
type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// LoginRequest represents the data required for user login.
type LoginRequest struct {
	Identifier string `json:"identifier"` // Email or phone
	Password   string `json:"password"`
}

// User represents the structure of a user in the database.
type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"password_hash"`
}

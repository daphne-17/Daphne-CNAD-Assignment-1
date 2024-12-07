package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Define a custom struct for the claims
type Claims struct {
	UserID string `json:"sub"` // 'sub' is the standard field for the subject (userID)
	jwt.RegisteredClaims
}

var jwtKey = []byte("your_secret_key") // Ensure this matches in all services

// GenerateToken creates a new JWT for the authenticated user with userID as the subject.
func GenerateToken(userID string) (string, error) {
	// Use the custom claims with the userID
	claims := &Claims{
		UserID: userID, // Store the userID in 'sub'
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 1 day
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // Token issued at the current time
		},
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return the token
	return token.SignedString(jwtKey)
}

// ValidateToken parses and validates a JWT token and extracts the user ID (sub).
func ValidateToken(tokenStr string) (string, error) {
	// Parse the token with the custom claims structure
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil // Extract the userID from 'sub'
	}
	return "", errors.New("invalid token")
}

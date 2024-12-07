package handlers

import (
	"encoding/json"
	"net/http"
	"user-service/database"
	"user-service/middleware"
)

// ViewMembershipHandler retrieves the membership tier for a given user
func ViewMembershipHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
        SELECT MembershipTiers.TierName, MembershipTiers.HourlyRate, MembershipTiers.BookingLimit
        FROM Users
        INNER JOIN MembershipTiers ON Users.TierID = MembershipTiers.TierID
        WHERE Users.UserID = ?`

	var tierName string
	var hourlyRate float64
	var bookingLimit int

	err := database.GetUserDB().QueryRow(query, userID).Scan(&tierName, &hourlyRate, &bookingLimit)
	if err != nil {
		http.Error(w, "Failed to retrieve membership details", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"tier_name":     tierName,
		"hourly_rate":   hourlyRate,
		"booking_limit": bookingLimit,
	})
}

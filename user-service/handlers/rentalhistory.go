package handlers

import (
	"encoding/json"
	"net/http"
	"user-service/database"
	"user-service/middleware"
)

func RentalHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
        SELECT ReservationID, VehicleID, StartTime, EndTime, Status
        FROM Reservations
        WHERE UserID = ?
        ORDER BY StartTime DESC`

	rows, err := database.GetReservationDB().Query(query, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve rental history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rentals []map[string]interface{}
	for rows.Next() {
		var reservationID, vehicleID, status string
		var startTime, endTime string
		if err := rows.Scan(&reservationID, &vehicleID, &startTime, &endTime, &status); err != nil {
			http.Error(w, "Error reading rental history", http.StatusInternalServerError)
			return
		}
		rentals = append(rentals, map[string]interface{}{
			"reservation_id": reservationID,
			"vehicle_id":     vehicleID,
			"start_time":     startTime,
			"end_time":       endTime,
			"status":         status,
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":        userID,
		"rental_history": rentals,
	})
}

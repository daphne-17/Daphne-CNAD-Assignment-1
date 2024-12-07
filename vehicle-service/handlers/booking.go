package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"vehicle-service/database"
	"vehicle-service/middleware"

	_ "github.com/go-sql-driver/mysql"
)

type ReservationRequest struct {
	VehicleID string    `json:"vehicle_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func CreateReservationHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, _ := io.ReadAll(r.Body)
	log.Printf("Raw request body: %s", string(body))

	var req struct {
		VehicleID string `json:"vehicle_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Parse the time fields
	startTime, err := time.Parse("2006-01-02T15:04:05", req.StartTime)
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse("2006-01-02T15:04:05", req.EndTime)
	if err != nil {
		log.Printf("Error parsing end time: %v", err)
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	if endTime.Before(startTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	// Check vehicle availability
	var status string
	query := "SELECT Status FROM Vehicles WHERE VehicleID = ? AND Status = 'Available'"
	err = database.GetVehicleDB().QueryRow(query, req.VehicleID).Scan(&status)
	if err != nil {
		log.Printf("Error checking vehicle availability: %v", err)
		http.Error(w, "Vehicle is not available", http.StatusInternalServerError)
		return
	}

	// Generate next reservation ID
	reservationID, err := getNextReservationID()
	if err != nil {
		log.Printf("Error generating reservation ID: %v", err)
		http.Error(w, "Error generating reservation ID", http.StatusInternalServerError)
		return
	}

	// Insert reservation
	insertQuery := `INSERT INTO Reservations (ReservationID, UserID, VehicleID, StartTime, EndTime, Status) 
                    VALUES (?, ?, ?, ?, ?, 'Active')`
	_, err = database.GetReservationDB().Exec(insertQuery, reservationID, userID, req.VehicleID, startTime, endTime)
	if err != nil {
		log.Printf("Error creating reservation: %v", err)
		http.Error(w, "Error creating reservation", http.StatusInternalServerError)
		return
	}

	// Update vehicle status
	updateQuery := "UPDATE Vehicles SET Status = 'Reserved' WHERE VehicleID = ?"
	_, err = database.GetVehicleDB().Exec(updateQuery, req.VehicleID)
	if err != nil {
		log.Printf("Error updating vehicle status: %v", err)
		http.Error(w, "Error updating vehicle status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Reservation created successfully with ID: %s", reservationID)))
}

// ModifyReservationHandler allows users to modify their reservations
func ModifyReservationHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Unauthorized: UserID not found in token")
		return
	}

	// Parse request payload
	var req struct {
		ReservationID string `json:"reservation_id"`
		NewStartTime  string `json:"new_start_time"`
		NewEndTime    string `json:"new_end_time"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		log.Printf("Error reading request body: %v\n", err)
		return
	}

	log.Printf("Raw request body: %s", string(body))

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request payload: %v\n", err)
		return
	}

	// Validate required fields
	if req.ReservationID == "" || req.NewStartTime == "" || req.NewEndTime == "" {
		http.Error(w, "Missing required fields in request payload", http.StatusBadRequest)
		log.Println("Missing required fields: ReservationID, NewStartTime, or NewEndTime")
		return
	}

	// Parse datetime strings into `time.Time`
	layout := "2006-01-02T15:04:05"
	newStartTime, err := time.Parse(layout, req.NewStartTime)
	if err != nil {
		http.Error(w, "Invalid new start time format", http.StatusBadRequest)
		log.Printf("Error parsing new start time: %v\n", err)
		return
	}

	newEndTime, err := time.Parse(layout, req.NewEndTime)
	if err != nil {
		http.Error(w, "Invalid new end time format", http.StatusBadRequest)
		log.Printf("Error parsing new end time: %v\n", err)
		return
	}

	if newEndTime.Before(newStartTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		log.Println("Validation failed: End time is before start time")
		return
	}

	log.Printf("Parsed NewStartTime=%v, NewEndTime=%v\n", newStartTime, newEndTime)

	// Check if the reservation exists and is associated with the user
	var rawStartTime string
	query := "SELECT StartTime FROM Reservations WHERE ReservationID = ? AND UserID = ?"
	err = database.GetReservationDB().QueryRow(query, req.ReservationID, userID).Scan(&rawStartTime)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Reservation not found or unauthorized", http.StatusNotFound)
			log.Println("No rows found for the given ReservationID and UserID")
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Database error: %v\n", err)
		}
		return
	}

	// Parse `rawStartTime` into a `time.Time` object
	currentStartTime, err := time.Parse("2006-01-02 15:04:05", rawStartTime)
	if err != nil {
		http.Error(w, "Error parsing reservation start time", http.StatusInternalServerError)
		log.Printf("Error parsing start time from DB: %v\n", err)
		return
	}
	log.Printf("Parsed reservation start time: %v\n", currentStartTime)

	// Check if the modification is within 48 hours of the reservation start time
	if time.Until(currentStartTime).Hours() < 48 {
		http.Error(w, "Modification not allowed within 48 hours of reservation start", http.StatusForbidden)
		log.Printf("Modification denied: ReservationID=%s, UserID=%s, HoursUntilStart=%.2f\n", req.ReservationID, userID, time.Until(currentStartTime).Hours())
		return
	}

	// Update the reservation in the database
	updateQuery := `
		UPDATE Reservations
		SET StartTime = ?, EndTime = ?, Status = 'Active'
		WHERE ReservationID = ? AND UserID = ?`
	_, err = database.GetReservationDB().Exec(updateQuery, newStartTime.Format("2006-01-02 15:04:05"), newEndTime.Format("2006-01-02 15:04:05"), req.ReservationID, userID)
	if err != nil {
		http.Error(w, "Error modifying reservation", http.StatusInternalServerError)
		log.Printf("Error updating reservation: %v\n", err)
		return
	}

	log.Printf("Reservation modified successfully: ReservationID=%s, UserID=%s\n", req.ReservationID, userID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reservation modified successfully"))
}

// CancelReservationHandler allows users to cancel their reservations
func CancelReservationHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("CancelReservationHandler: UserID from token = %s\n", userID)

	var req struct {
		ReservationID string `json:"reservation_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("CancelReservationHandler: ReservationID from request = %s\n", req.ReservationID)

	// Validate the reservation
	var rawStartTime string
	query := "SELECT StartTime FROM Reservations WHERE ReservationID = ? AND UserID = ?"
	err := database.GetReservationDB().QueryRow(query, req.ReservationID, userID).Scan(&rawStartTime)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		http.Error(w, "Reservation not found or unauthorized", http.StatusNotFound)
		return
	}

	// Parse StartTime from string to time.Time
	currentStartTime, err := time.Parse("2006-01-02 15:04:05", rawStartTime)
	if err != nil {
		http.Error(w, "Error parsing reservation start time", http.StatusInternalServerError)
		log.Printf("Error parsing start time: %v\n", err)
		return
	}

	log.Printf("Parsed reservation start time: %v\n", currentStartTime)

	// Validate cancellation time
	timeUntilStart := time.Until(currentStartTime).Hours()
	if timeUntilStart < 48 {
		http.Error(w, "Cancellation not allowed within 48 hours of reservation start", http.StatusForbidden)
		log.Printf("Cancellation not allowed. Time until start: %.2f hours\n", timeUntilStart)
		return
	}

	// Update reservation status to 'Cancelled'
	cancelQuery := `
		UPDATE Reservations
		SET Status = 'Cancelled'
		WHERE ReservationID = ? AND UserID = ?`
	_, err = database.GetReservationDB().Exec(cancelQuery, req.ReservationID, userID)
	if err != nil {
		log.Printf("Error canceling reservation: %v\n", err)
		http.Error(w, "Error canceling reservation", http.StatusInternalServerError)
		return
	}

	// Free the vehicle
	var vehicleID string
	fetchVehicleQuery := "SELECT VehicleID FROM Reservations WHERE ReservationID = ?"
	err = database.GetReservationDB().QueryRow(fetchVehicleQuery, req.ReservationID).Scan(&vehicleID)
	if err != nil {
		log.Printf("Error fetching vehicle ID from reservation: %v\n", err)
		http.Error(w, "Error fetching vehicle ID", http.StatusInternalServerError)
		return
	}

	updateQuery := "UPDATE Vehicles SET Status = 'Available' WHERE VehicleID = ?"
	_, err = database.GetVehicleDB().Exec(updateQuery, vehicleID)
	if err != nil {
		log.Printf("Error updating vehicle status: %v\n", err)
		http.Error(w, "Error updating vehicle status", http.StatusInternalServerError)
		return
	}

	// Delete the reservation
	deleteReservationQuery := "DELETE FROM Reservations WHERE ReservationID = ? AND UserID = ?"
	_, err = database.GetReservationDB().Exec(deleteReservationQuery, req.ReservationID, userID)
	if err != nil {
		log.Printf("Error deleting reservation: %v\n", err)
		http.Error(w, "Error deleting reservation", http.StatusInternalServerError)
		return
	}

	log.Printf("Reservation %s successfully canceled by user %s\n", req.ReservationID, userID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reservation canceled successfully"))
}

func GetUserReservationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Fetching reservations for UserID: %s\n", userID)

	query := `
		SELECT ReservationID, VehicleID, StartTime, EndTime, Status 
		FROM Reservations 
		WHERE UserID = ? 
		ORDER BY StartTime ASC`

	rows, err := database.GetReservationDB().Query(query, userID)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		http.Error(w, "Failed to fetch reservations", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reservations []map[string]interface{}
	for rows.Next() {
		var reservationID, vehicleID, status string
		var rawStartTime, rawEndTime []byte

		// Read raw time as []byte
		if err := rows.Scan(&reservationID, &vehicleID, &rawStartTime, &rawEndTime, &status); err != nil {
			log.Printf("Error scanning reservation data: %v\n", err)
			http.Error(w, "Error reading reservations", http.StatusInternalServerError)
			return
		}

		// Convert []byte to string
		startTime := string(rawStartTime)
		endTime := string(rawEndTime)

		reservations = append(reservations, map[string]interface{}{
			"reservation_id": reservationID,
			"vehicle_id":     vehicleID,
			"start_time":     startTime,
			"end_time":       endTime,
			"status":         status,
		})
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v\n", err)
		http.Error(w, "Error processing reservations", http.StatusInternalServerError)
		return
	}

	// Check if no reservations are available
	if len(reservations) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "No bookings available"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

// getNextReservationID generates the next reservation ID
func getNextReservationID() (string, error) {
	var lastID string
	query := "SELECT ReservationID FROM Reservations ORDER BY ReservationID DESC LIMIT 1"
	err := database.GetReservationDB().QueryRow(query).Scan(&lastID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "R1", nil // First reservation
		}
		return "", err
	}

	// Increment the last reservation ID
	var number int
	_, err = fmt.Sscanf(lastID, "R%d", &number)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("R%d", number+1), nil
}

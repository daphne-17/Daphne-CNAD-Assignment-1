package models

import "time"

// Reservation represents a reservation made by a user for a vehicle.
type Reservation struct {
	ReservationID string    `json:"reservation_id"` // Unique ID for the reservation
	UserID        string    `json:"user_id"`        // User ID who made the reservation
	VehicleID     string    `json:"vehicle_id"`     // Vehicle ID reserved
	StartTime     time.Time `json:"start_time"`     // Start time of the reservation
	EndTime       time.Time `json:"end_time"`       // End time of the reservation
	Status        string    `json:"status"`         // Reservation status (Active, Cancelled, Completed)
}

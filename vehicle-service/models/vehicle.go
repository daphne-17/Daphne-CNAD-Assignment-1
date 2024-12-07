package models

// Vehicle represents a car in the vehicle database.
type Vehicle struct {
	VehicleID     string  `json:"vehicle_id"`     // Unique ID of the vehicle
	Model         string  `json:"model"`          // Vehicle model (e.g., Tesla Model 3)
	VehicleNumber string  `json:"vehicle_number"` // Vehicle number or plate
	ChargeLevel   int     `json:"charge_level"`   // Charge level (for electric vehicles)
	Location      string  `json:"location"`       // The location where the car is
	Status        string  `json:"status"`         // Current status (Available, Reserved, Maintenance)
	RatePerHour   float64 `json:"rate_per_hour"`  // The hourly rate of the vehicle
}

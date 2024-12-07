package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"vehicle-service/database"
)

// FetchAvailableVehicles fetches all available vehicles from the database
func FetchAvailableVehicles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT VehicleID, Model, RatePerHour, Status 
		FROM Vehicles 
		WHERE Status = 'Available'`

	rows, err := database.GetVehicleDB().Query(query)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		http.Error(w, "Failed to retrieve vehicles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vehicles []map[string]interface{}
	for rows.Next() {
		var vehicleID, model, status string
		var ratePerHour float64
		if err := rows.Scan(&vehicleID, &model, &ratePerHour, &status); err != nil {
			log.Printf("Error scanning row: %v\n", err)
			http.Error(w, "Error reading vehicle data", http.StatusInternalServerError)
			return
		}
		vehicles = append(vehicles, map[string]interface{}{
			"vehicle_id":    vehicleID,
			"model":         model,
			"rate_per_hour": ratePerHour,
			"status":        status,
		})
	}

	if len(vehicles) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No available vehicles found"})
		return
	}

	json.NewEncoder(w).Encode(vehicles)
}

package handlers

import (
	"billing-service/database"
	"billing-service/middleware"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CostRequest struct {
	VehicleID string `json:"vehicle_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	PromoCode string `json:"promo_code,omitempty"`
}

type CostResponse struct {
	TotalCost      float64 `json:"total_cost"`
	DiscountAmount float64 `json:"discount_amount,omitempty"`
	FinalCost      float64 `json:"final_cost"`
}

type ReservationRequest struct {
	VehicleID string  `json:"vehicle_id"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	PromoCode string  `json:"promo_code,omitempty"`
	TotalCost float64 `json:"total_cost"`
}

// CalculateCostHandler calculates the estimated cost and applies promo code if provided
func CalculateCostHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Error: Missing UserID in JWT token")
		return
	}

	var req CostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request: %v\n", err)
		return
	}

	if req.VehicleID == "" || req.StartTime == "" || req.EndTime == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		log.Println("Error: Missing required fields (VehicleID, StartTime, EndTime)")
		return
	}

	startTime, err := time.Parse("2006-01-02T15:04:05", req.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		log.Printf("Error parsing start time: %v\n", err)
		return
	}

	endTime, err := time.Parse("2006-01-02T15:04:05", req.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		log.Printf("Error parsing end time: %v\n", err)
		return
	}

	duration := endTime.Sub(startTime).Hours()
	if duration <= 0 {
		http.Error(w, "Invalid time range", http.StatusBadRequest)
		log.Println("Error: End time must be after start time")
		return
	}

	var ratePerHour float64
	query := "SELECT RatePerHour FROM Vehicles WHERE VehicleID = ?"
	err = database.GetVehicleDB().QueryRow(query, req.VehicleID).Scan(&ratePerHour)
	if err != nil {
		http.Error(w, "Vehicle not found or unavailable", http.StatusNotFound)
		log.Printf("Error fetching vehicle rate: %v\n", err)
		return
	}

	totalCost := ratePerHour * duration
	finalCost := totalCost
	discountAmount := 0.0

	if req.PromoCode != "" {
		var discountRate float64
		promoQuery := "SELECT DiscountPercent FROM Promotions WHERE PromoCode = ? AND EndDate >= CURDATE()"
		err = database.GetBillingDB().QueryRow(promoQuery, req.PromoCode).Scan(&discountRate)
		if err != nil {
			http.Error(w, "Invalid or expired promo code", http.StatusBadRequest)
			log.Printf("Invalid promo code: %s\n", req.PromoCode)
			return
		}
		finalCost = totalCost - discountRate
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CostResponse{
		TotalCost:      totalCost,
		DiscountAmount: discountAmount,
		FinalCost:      finalCost,
	})
}

// ConfirmReservationHandler confirms the reservation and creates a billing entry
func ConfirmReservationHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Error: Missing UserID in JWT token")
		return
	}

	var req ReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request: %v\n", err)
		return
	}

	if req.VehicleID == "" || req.StartTime == "" || req.EndTime == "" || req.TotalCost <= 0 {
		http.Error(w, "Missing required fields or invalid cost", http.StatusBadRequest)
		log.Println("Error: Missing required fields or invalid cost")
		return
	}

	billingID, err := getNextBillingID()
	if err != nil {
		http.Error(w, "Error generating billing ID", http.StatusInternalServerError)
		log.Printf("Error generating billing ID: %v\n", err)
		return
	}

	billingQuery := `
        INSERT INTO Billing (BillingID, ReservationID, UserID, TotalAmount, PromoID)
        VALUES (?, ?, ?, ?, 'Paid', ?)`
	_, err = database.GetBillingDB().Exec(billingQuery, billingID, userID, req.VehicleID, req.TotalCost, req.PromoCode)
	if err != nil {
		http.Error(w, "Error saving billing information", http.StatusInternalServerError)
		log.Printf("Error saving billing: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"billing_id": billingID,
		"message":    "Reservation and payment confirmed successfully",
	})
}

func getNextBillingID() (string, error) {
	var lastID string
	query := "SELECT BillingID FROM Billing ORDER BY BillingID DESC LIMIT 1"
	err := database.GetBillingDB().QueryRow(query).Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "B1", nil
		}
		return "", fmt.Errorf("error retrieving last billing ID: %w", err)
	}

	var number int
	_, err = fmt.Sscanf(lastID, "B%d", &number)
	if err != nil {
		return "", fmt.Errorf("error parsing last billing ID: %w", err)
	}

	return fmt.Sprintf("B%d", number+1), nil
}

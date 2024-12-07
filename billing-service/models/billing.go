package models

// Billing represents a billing record for a reservation
type Billing struct {
	BillingID     string  `json:"billing_id"`     // e.g., B1
	ReservationID string  `json:"reservation_id"` // e.g., R1
	UserID        string  `json:"user_id"`        // e.g., U1
	TotalAmount   float64 `json:"total_amount"`   // e.g., 40.00
	PromoID       *string `json:"promo_id"`       // e.g., P1 (nullable)
}

package models

// Membership represents a user’s membership details.
type Membership struct {
	TierID       string  `json:"tier_id"`       // Unique ID for the tier
	TierName     string  `json:"tier_name"`     // Name of the membership tier (Basic, Premium, etc.)
	HourlyRate   float64 `json:"hourly_rate"`   // Discounted hourly rate based on membership
	BookingLimit int     `json:"booking_limit"` // Max number of bookings allowed
}

package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var (
	BillingDB     *sql.DB // Connection for billing_db
	ReservationDB *sql.DB // Connection for reservation_db
	VehicleDB     *sql.DB // Connection for vehicle_db
)

// InitDB initializes the database connections
func InitDB() error {
	var err error

	// Connect to billing_db
	billingDSN := "root:seventeen@tcp(localhost:3306)/billing_db"
	BillingDB, err = sql.Open("mysql", billingDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to billing_db: %w", err)
	}
	if err = BillingDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping billing_db: %w", err)
	}
	log.Println("Connected to billing_db successfully!")

	// Connect to reservation_db
	reservationDSN := "root:seventeen@tcp(localhost:3306)/reservation_db"
	ReservationDB, err = sql.Open("mysql", reservationDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to reservation_db: %w", err)
	}
	if err = ReservationDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping reservation_db: %w", err)
	}
	log.Println("Connected to reservation_db successfully!")

	// Connect to vehicle_db
	vehicleDSN := "root:seventeen@tcp(localhost:3306)/vehicle_db"
	VehicleDB, err = sql.Open("mysql", vehicleDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to vehicle_db: %w", err)
	}
	if err = VehicleDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping vehicle_db: %w", err)
	}
	log.Println("Connected to vehicle_db successfully!")

	return nil
}

// CloseDB closes all database connections
func CloseDB() {
	if BillingDB != nil {
		if err := BillingDB.Close(); err != nil {
			log.Printf("Error closing billing_db: %v", err)
		} else {
			log.Println("Closed connection to billing_db")
		}
	}

	if ReservationDB != nil {
		if err := ReservationDB.Close(); err != nil {
			log.Printf("Error closing reservation_db: %v", err)
		} else {
			log.Println("Closed connection to reservation_db")
		}
	}

	if VehicleDB != nil {
		if err := VehicleDB.Close(); err != nil {
			log.Printf("Error closing vehicle_db: %v", err)
		} else {
			log.Println("Closed connection to vehicle_db")
		}
	}
}

// GetBillingDB returns the billing_db connection
func GetBillingDB() *sql.DB {
	if BillingDB == nil {
		log.Fatal("billing_db is not initialized. Call InitDB first.")
	}
	return BillingDB
}

// GetReservationDB returns the reservation_db connection
func GetReservationDB() *sql.DB {
	if ReservationDB == nil {
		log.Fatal("reservation_db is not initialized. Call InitDB first.")
	}
	return ReservationDB
}

// GetVehicleDB returns the vehicle_db connection
func GetVehicleDB() *sql.DB {
	if VehicleDB == nil {
		log.Fatal("vehicle_db is not initialized. Call InitDB first.")
	}
	return VehicleDB
}

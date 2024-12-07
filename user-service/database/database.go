package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var (
	UserDB        *sql.DB // Connection for user_db
	ReservationDB *sql.DB // Connection for reservation_db
)

// InitDB initializes database connections
func InitDB() error {
	var err error

	// Connect to user_db
	userDSN := "root:seventeen@tcp(localhost:3306)/user_db"
	UserDB, err = sql.Open("mysql", userDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to user_db: %w", err)
	}

	if err = UserDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping user_db: %w", err)
	}
	log.Println("Connected to user_db successfully!")

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

	return nil
}

// CloseDB closes all database connections
func CloseDB() {
	if UserDB != nil {
		if err := UserDB.Close(); err != nil {
			log.Printf("Error closing user_db: %v", err)
		} else {
			log.Println("Closed connection to user_db")
		}
	}

	if ReservationDB != nil {
		if err := ReservationDB.Close(); err != nil {
			log.Printf("Error closing reservation_db: %v", err)
		} else {
			log.Println("Closed connection to reservation_db")
		}
	}
}

// GetUserDB returns the user_db connection
func GetUserDB() *sql.DB {
	if UserDB == nil {
		log.Fatal("user_db is not initialized. Call InitDB first.")
	}
	return UserDB
}

// GetReservationDB returns the reservation_db connection
func GetReservationDB() *sql.DB {
	if ReservationDB == nil {
		log.Fatal("reservation_db is not initialized. Call InitDB first.")
	}
	return ReservationDB
}

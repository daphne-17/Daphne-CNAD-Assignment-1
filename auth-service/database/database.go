package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var db *sql.DB

func init() {
	var err error
	dsn := "root:seventeen@tcp(localhost:3306)/user_db"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Database connection error: %v", err))
	}
	fmt.Println("Connected to database successfully!")
}

func GetDB() *sql.DB {
	return db
}

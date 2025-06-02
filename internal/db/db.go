package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"

	"gohst/internal/config"
)

// DBManager manages the database connection
type DBManager struct {
	DB *sql.DB
}

var (
	Database *DBManager // Global database manager
	once     sync.Once  // Ensure InitDB runs only once
)

// InitDB initializes the database connection
func InitDB() {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.DB.Host,
			config.DB.Port,
			config.DB.User,
			config.DB.Password,
			config.DB.DBName,
		)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Error connecting to Postgres: %v", err)
		}

		// Ping to verify the connection is working
		if err = db.Ping(); err != nil {
			log.Fatalf("Postgres ping failed: %v", err)
		}

		log.Println("Connected to Postgres")

		Database = &DBManager{DB: db}
	})
}

// CloseDB closes the database connection
func CloseDB() {
	if Database != nil && Database.DB != nil {
		Database.DB.Close()
		log.Println("Database connection closed.")
	}
}

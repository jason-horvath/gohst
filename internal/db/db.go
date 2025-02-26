package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"

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
	once.Do(func() { // Ensures this runs only once
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			config.DB.User,
			config.DB.Password,
			config.DB.Host,
			strconv.Itoa(config.DB.Port),
			config.DB.DBName,
		)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error connecting to MySQL: %v", err)
		}

		// Ping to verify the connection is working
		if err = db.Ping(); err != nil {
			log.Fatalf("MySQL ping failed: %v", err)
		}

		log.Println("Connected to MySQL")

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

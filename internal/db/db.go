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
const PRIMARY_DB_NAME = "primary"
type DBManager struct {
	DB *sql.DB
}

var (
	Database  *DBManager                 // Global database manager
	Databases map[string]*DBManager      // Multiple database connections
	once      sync.Once                  // Ensure InitDB runs only once
	multiOnce sync.Once                  // Ensure InitMultiDB runs only once
)

// InitDBPool initializes multiple database connections from app configs
func InitDBPool(pool *config.DatabaseConfigPool) {
	multiOnce.Do(func() {
		Databases = make(map[string]*DBManager)

		for name, dbConfig := range pool.GetConfigs() {
			sslMode := dbConfig.SSLMode
			if sslMode == "" {
				sslMode = "disable" // Default for local development
			}

			dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				dbConfig.Host,
				dbConfig.Port,
				dbConfig.User,
				dbConfig.Password,
				dbConfig.DBName,
				sslMode,
			)

			db, err := sql.Open("postgres", dsn)
			if err != nil {
				log.Fatalf("Error connecting to database %s: %v", name, err)
			}

			if err = db.Ping(); err != nil {
				log.Fatalf("Database %s ping failed: %v", name, err)
			}

			log.Printf("Connected to database: %s", name)
			Databases[name] = &DBManager{DB: db}
		}
	})
}

// GetDB returns a specific database connection by name
func GetDB(name string) *DBManager {
	if db, exists := Databases[name]; exists {
		return db
	}

	return nil
}

// GetPrimaryDB returns the primary database connection
func GetPrimaryDB() *DBManager {
	if primary := GetDB("primary"); primary != nil {
		return primary
	}
	return Database // Fallback to legacy Database
}

// CloseDB closes all database connections
func CloseDB() {
	// Close primary database connection
	if Database != nil && Database.DB != nil {
		Database.DB.Close()
		log.Println("Primary database connection closed.")
	}

	// Close all multi-database connections
	for name, dbManager := range Databases {
		if dbManager != nil && dbManager.DB != nil {
			dbManager.DB.Close()
			log.Printf("Database connection '%s' closed.", name)
		}
	}
}

// CloseDBPool closes all database connections and clears the Databases map
func CloseDBPool() {
	// Close all multi-database connections
	for name, dbManager := range Databases {
		if dbManager != nil && dbManager.DB != nil {
			dbManager.DB.Close()
			log.Printf("Database connection '%s' closed.", name)
		}
	}

	Databases = make(map[string]*DBManager) // Clear the map after closing
}

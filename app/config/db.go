package config

import (
	"gohst/internal/config"
)

// InitDBConfigs initializes the database configurations for the app
func CreateDBConfigs() *config.DatabaseConfigPool {
	dbConfigPool := config.NewDatabaseConfigPool()

	// Add primary database
	primaryDB := &config.DatabaseConfig{
		Host:     config.GetEnv("DB_HOST", "localhost").(string),
		Port:     config.GetEnv("DB_PORT", config.DB_DEFAULT_PORT).(int),
		User:     config.GetEnv("DB_USER", "gohst").(string),
		Password: config.GetEnv("DB_PASSWORD", "password").(string),
		DBName:   config.GetEnv("DB_NAME", "gohst").(string),
	}

	dbConfigPool.Add("primary", primaryDB)

	// Example: Additional Anaylitics database configuration
	// Add analytics database if configured
	// if config.GetEnv("DB_ANALYTICS_HOST", "") != "" {
	// 	analyticsDB := &config.DatabaseConfig{
	// 		Host:     config.GetEnv("DB_ANALYTICS_HOST", "localhost").(string),
	// 		Port:     config.GetEnv("DB_ANALYTICS_PORT", 5432).(int),
	// 		User:     config.GetEnv("DB_ANALYTICS_USER", primaryDB.User).(string),
	// 		Password: config.GetEnv("DB_ANALYTICS_PASSWORD", primaryDB.Password).(string),
	// 		DBName:   config.GetEnv("DB_ANALYTICS_NAME", "gohst_analytics").(string),
	// 	}
	// 	dbConfigPool.Add("analytics", analyticsDB)
	// }

	return dbConfigPool
}


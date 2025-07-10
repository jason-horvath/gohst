package config

import (
	"gohst/internal/config"
)

var DBConfigMap *config.DatabaseConfigPool

// InitDBConfigs initializes the database configurations for the app
func InitDBConfigs() {
	DBConfigMap = config.NewDatabaseConfigPool()

	// Add primary database
	primaryDB := &config.DatabaseConfig{
		Host:     config.GetEnv("DB_HOST", "localhost").(string),
		Port:     config.GetEnv("DB_PORT", 5432).(int),
		User:     config.GetEnv("DB_USER", "gohst").(string),
		Password: config.GetEnv("DB_PASSWORD", "password").(string),
		DBName:   config.GetEnv("DB_NAME", "gohst").(string),
	}
	DBConfigMap.Add("primary", primaryDB)

	// Add analytics database if configured
	if config.GetEnv("DB_ANALYTICS_HOST", "") != "" {
		analyticsDB := &config.DatabaseConfig{
			Host:     config.GetEnv("DB_ANALYTICS_HOST", "localhost").(string),
			Port:     config.GetEnv("DB_ANALYTICS_PORT", 5432).(int),
			User:     config.GetEnv("DB_ANALYTICS_USER", primaryDB.User).(string),
			Password: config.GetEnv("DB_ANALYTICS_PASSWORD", primaryDB.Password).(string),
			DBName:   config.GetEnv("DB_ANALYTICS_NAME", "gohst_analytics").(string),
		}
		DBConfigMap.Add("analytics", analyticsDB)
	}
}

// GetPrimaryDB returns the primary database config
func GetPrimaryDB() *config.DatabaseConfig {
	if config, exists := DBConfigMap.Get("primary"); exists {
		return config
	}
	return nil
}

// GetAnalyticsDB returns the analytics database config, fallback to primary
func GetAnalyticsDB() *config.DatabaseConfig {
	primary, _ := DBConfigMap.Get("primary")
	return DBConfigMap.GetOrDefault("analytics", primary)
}

// GetDBConfig returns a database config by name
func GetDBConfig(name string) (*config.DatabaseConfig, bool) {
	return DBConfigMap.Get(name)
}

// GetAvailableDatabases returns all configured database names
func GetAvailableDatabases() []string {
	return DBConfigMap.Names()
}


package config

import (
	"gohst/internal/config"
)

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string
	Version     string
	Environment string
	Debug       bool
	URL         string

	// Feature flags
	Features FeatureFlags

	// Business logic settings
	Pagination PaginationConfig
	Upload     UploadConfig
}

// FeatureFlags controls optional application features
type FeatureFlags struct {
	EnableRegistration bool
	EnableUserProfiles bool
	EnableNotifications bool
	MaintenanceMode    bool
}

// PaginationConfig controls pagination behavior
type PaginationConfig struct {
	DefaultLimit int
	MaxLimit     int
}

// UploadConfig controls file upload settings
type UploadConfig struct {
	MaxFileSize   int64  // in bytes
	AllowedTypes  []string
	UploadPath    string
}

// App holds the global application configuration
var App *AppConfig

// Initialize sets up the application configuration
func Initialize() {
	App = &AppConfig{
		Name:        config.GetEnv("APP_NAME", "Gohst Application").(string),
		Version:     config.GetEnv("APP_VERSION", "1.0.0").(string),
		Environment: config.GetEnv("APP_ENV", "development").(string),
		Debug:       config.GetEnv("APP_DEBUG", false).(bool),
		URL:         config.GetEnv("APP_URL", "http://localhost:8080").(string),

		Features: FeatureFlags{
			EnableRegistration:  config.GetEnv("FEATURE_REGISTRATION", true).(bool),
			EnableUserProfiles:  config.GetEnv("FEATURE_USER_PROFILES", true).(bool),
			EnableNotifications: config.GetEnv("FEATURE_NOTIFICATIONS", false).(bool),
			MaintenanceMode:     config.GetEnv("MAINTENANCE_MODE", false).(bool),
		},

		Pagination: PaginationConfig{
			DefaultLimit: config.GetEnv("PAGINATION_DEFAULT_LIMIT", 20).(int),
			MaxLimit:     config.GetEnv("PAGINATION_MAX_LIMIT", 100).(int),
		},

		Upload: UploadConfig{
			MaxFileSize:  config.GetEnv("UPLOAD_MAX_FILE_SIZE", int64(10<<20)).(int64), // 10MB
			AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
			UploadPath:   config.GetEnv("UPLOAD_PATH", "static/uploads").(string),
		},
	}
}

// IsProduction returns true if the app is running in production
func IsProduction() bool {
	return App.Environment == "production"
}

// IsDevelopment returns true if the app is running in development
func IsDevelopment() bool {
	return App.Environment == "development"
}

// IsMaintenanceMode returns true if maintenance mode is enabled
func IsMaintenanceMode() bool {
	return App.Features.MaintenanceMode
}

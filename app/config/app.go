package config

import (
	"gohst/internal/config"
)

// AppConfig holds application-specific configuration
type AppConfig struct {
    EnvKey       string // The application environment (e.g., "development", "production").
    URL          string // The application URL.
	DistPath	 string // The path to the distribution directory.
    Port         int    // The port on which the application listens.
	Name        string
	Version     string
	Environment string
	Debug       bool

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
func InitAppConfig() *AppConfig {
	App = &AppConfig{
		Name:        config.GetEnv("APP_NAME", "Gohst Application").(string),
		Version:     config.GetEnv("APP_VERSION", "1.0.0").(string),
		EnvKey: 	 config.GetEnv("APP_ENV_KEY", "development").(string),
		Debug:       config.GetEnv("APP_DEBUG", false).(bool),
		URL:         config.GetEnv("APP_URL", "http://localhost:3030").(string),
		Port:        config.GetEnv("APP_PORT", config.APP_DEFAULT_PORT).(int),
		DistPath:    config.GetEnv("APP_DIST_PATH", "static/dist").(string),

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
			MaxFileSize:  int64(config.GetEnv("UPLOAD_MAX_FILE_SIZE", 10485760).(int)), // 10MB as int, then convert to int64
			AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
			UploadPath:   config.GetEnv("UPLOAD_PATH", "static/uploads").(string),
		},
	}

	return App
}

// GetAppConfig returns the global application configuration
func (ac *AppConfig) GetURL() string {
	return ac.URL
}

// DistPath returns the port on which the application listens
func (ac *AppConfig) GetDistPath() string {
	return ac.DistPath
}

// IsProduction returns true if the app is running in production
func (ac *AppConfig) IsProduction() bool {
	return ac.EnvKey == "production"
}

// IsDevelopment returns true if the app is running in development
func (ac *AppConfig) IsDevelopment() bool {
	return ac.EnvKey == "development"
}

// IsMaintenanceMode returns true if maintenance mode is enabled
func (ac *AppConfig) IsMaintenanceMode() bool {
	return ac.Features.MaintenanceMode
}

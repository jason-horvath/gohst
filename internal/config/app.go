package config

const APP_DEFAULT_PORT = 3030

const APP_CSRF_KEY string = "csrf_token"

// AppConfigProvider - Simple, single interface
type AppConfigProvider interface {
	GetURL() string
	GetDistPath() string
    IsProduction() bool
    IsDevelopment() bool
    IsMaintenanceMode() bool
}

var appProvider AppConfigProvider

// Framework registration function
func RegisterAppConfig(provider AppConfigProvider) {
    appProvider = provider
}

// GetAppConfigProvider retrieves the registered AppConfigProvider.
func GetAppConfig() AppConfigProvider {
	if appProvider == nil {
		panic("AppConfigProvider is not registered")
	}
	return appProvider
}


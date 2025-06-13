package config

import (
	"strconv"
	"strings"
)

const APP_DEFAULT_PORT = 3030

const APP_CSRF_KEY string = "csrf_token"

// AppConfig holds the application's configuration values.
type AppConfig struct {
    EnvKey       string // The application environment (e.g., "development", "production").
    URL          string // The application URL.
	DistPath	 string // The path to the distribution directory.
    Port         int    // The port on which the application listens.
	CSRFName	 string // The key used for CSRF protection.
}

// App is the global application configuration variable.
var App *AppConfig

// InitApp initializes the application configuration by loading values from environment variables.
// It returns a pointer to the initialized AppConfig struct.
func initApp() {
	usePort := GetEnv("APP_PORT", APP_DEFAULT_PORT).(int)
	App = &AppConfig{
		EnvKey: GetEnv("APP_ENV_KEY", "development").(string),
		URL: GetEnv("APP_URL", "http://localhost:" + strconv.Itoa(usePort)).(string),
		DistPath: GetEnv("APP_DIST_PATH", "static/dist").(string),
		Port: GetEnv("APP_PORT", APP_DEFAULT_PORT).(int),
		CSRFName: APP_CSRF_KEY,
	}
}

func (a *AppConfig) IsDevelopment() bool {
	return a.EnvKey == "development"
}

func (a *AppConfig) IsProduction() bool {
	return a.EnvKey == "production"

}

func (a *AppConfig) PortStr() string {
	return strconv.Itoa(a.Port)
}

func (a * AppConfig) normalizeRel(rel string) string {
	if !strings.HasPrefix(rel, "/") {
		rel = "/" + rel
	}

	return rel
}

func (a *AppConfig) FullURL(rel string) string {
	rel = a.normalizeRel(rel)

	return a.URL + rel
}

func (a *AppConfig) DistURL(rel ...string) string {
	var useRel string
	if len(rel) == 0 {
		useRel = ""
	} else {
		useRel = rel[0]
	}

	dist := a.normalizeRel(a.DistPath)
	distRel := a.normalizeRel(useRel)

	return a.FullURL(dist + distRel)
}

func (a *AppConfig) IsHTTPS() bool {
	return strings.HasPrefix(a.URL, "https")
}

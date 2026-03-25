package render

import (
	"strings"

	"gohst/internal/config"
)

// AppURL returns the configured application base URL without trailing slash.
// Use this in templ components wherever the old {{ url }} template function was used.
func AppURL() string {
	return strings.TrimRight(config.GetEnv("APP_URL", "http://localhost:3030").(string), "/")
}


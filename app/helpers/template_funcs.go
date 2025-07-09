package helpers

import (
	"html/template"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	appConfig "gohst/app/config"
)

// AppTemplateFuncs returns application-specific template functions
func AppTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// App-specific functions
		"appName":        func() string { return appConfig.App.Name },
		"appVersion":     func() string { return appConfig.App.Version },
		"isMaintenanceMode": func() bool { return appConfig.IsMaintenanceMode() },

		// Feature flag helpers
		"canRegister":    func() bool { return appConfig.App.Features.EnableRegistration },
		"hasProfiles":    func() bool { return appConfig.App.Features.EnableUserProfiles },
		"hasNotifications": func() bool { return appConfig.App.Features.EnableNotifications },

		// Utility functions
		"formatDate":     func(t time.Time) string { return t.Format("January 2, 2006") },
		"truncate":       func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
		"upper":          strings.ToUpper,
		"lower":          strings.ToLower,
		"title":          func(s string) string { return cases.Title(language.English).String(s) },

		// Upload helpers
		"maxUploadSize":  func() int64 { return appConfig.App.Upload.MaxFileSize },
		"uploadPath":     func() string { return appConfig.App.Upload.UploadPath },

		// Pagination helpers
		"defaultPageSize": func() int { return appConfig.App.Pagination.DefaultLimit },
		"maxPageSize":     func() int { return appConfig.App.Pagination.MaxLimit },
	}
}

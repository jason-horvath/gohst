package helpers

import (
	"html/template"

	"gohst/app/config"
)

// AppTemplateFuncs returns application-specific template functions
func AppTemplateFuncs() template.FuncMap {
	app := config.App

	return template.FuncMap{
		// App identity
		"appName":        func() string { return app.Name },
		"appVersion":     func() string { return app.Version },
		"isMaintenanceMode": func() bool { return app.IsMaintenanceMode() },

		// App environment
		"isDevelopment": func() bool { return app.IsDevelopment() },
        "isProduction":  func() bool { return app.IsProduction() },
        "url":           func() string { return app.URL },

		// Feature flag helpers
		"canRegister":    func() bool { return app.Features.EnableRegistration },
		"hasProfiles":    func() bool { return app.Features.EnableUserProfiles },
		"hasNotifications": func() bool { return app.Features.EnableNotifications },

		// Upload configuration helpers
		"maxUploadSize":  func() int64 { return app.Upload.MaxFileSize },
		"uploadPath":     func() string { return app.Upload.UploadPath },

		// Pagination configuration helpers
		"defaultPageSize": func() int { return app.Pagination.DefaultLimit },
		"maxPageSize":     func() int { return app.Pagination.MaxLimit },
	}
}

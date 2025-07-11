package helpers

import (
	"html/template"

	appConfig "gohst/app/config"
)

// AppTemplateFuncs returns application-specific template functions
func AppTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// App identity
		"appName":        func() string { return appConfig.App.Name },
		"appVersion":     func() string { return appConfig.App.Version },
		"isMaintenanceMode": func() bool { return appConfig.IsMaintenanceMode() },

		// App environment
		"isDevelopment": func() bool { return appConfig.IsDevelopment() },
        "isProduction":  func() bool { return appConfig.IsProduction() },
        "url":           func() string { return appConfig.App.URL },

		// Feature flag helpers
		"canRegister":    func() bool { return appConfig.App.Features.EnableRegistration },
		"hasProfiles":    func() bool { return appConfig.App.Features.EnableUserProfiles },
		"hasNotifications": func() bool { return appConfig.App.Features.EnableNotifications },

		// Upload configuration helpers
		"maxUploadSize":  func() int64 { return appConfig.App.Upload.MaxFileSize },
		"uploadPath":     func() string { return appConfig.App.Upload.UploadPath },

		// Pagination configuration helpers
		"defaultPageSize": func() int { return appConfig.App.Pagination.DefaultLimit },
		"maxPageSize":     func() int { return appConfig.App.Pagination.MaxLimit },
	}
}

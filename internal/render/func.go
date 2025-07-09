package render

import (
	"html/template"

	"gohst/internal/config"
)

// appFuncs holds additional template functions registered by the application
var appFuncs = make(template.FuncMap)

// RegisterTemplateFuncs allows external packages to register additional template functions
func RegisterTemplateFuncs(funcs template.FuncMap) {
	for name, fn := range funcs {
		appFuncs[name] = fn
	}
}

// TemplateFuncs returns the core framework template functions merged with any registered app functions
func TemplateFuncs() template.FuncMap {
	coreFuncs := template.FuncMap{
		"assetsHead":    AssetsHead,
		"isDevelopment": func() bool { return config.App.IsDevelopment() },
		"isProduction":  func() bool { return config.App.IsProduction() },
		"url":           func() string { return config.App.URL },
		"icon":          Icon,
	}

	// Merge app functions (they can override core functions if needed)
	for name, fn := range appFuncs {
		coreFuncs[name] = fn
	}

	return coreFuncs
}

package render

import (
	"html/template"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
		// Framework infrastructure
		"assetsHead": AssetsHead,  // Framework asset management
		"icon":       Icon,        // Framework icon helper

		// Generic utility functions (useful for any app)
		"formatDate": func(t time.Time) string { return t.Format("January 2, 2006") },
		"truncate":   func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"title":      func(s string) string { return cases.Title(language.English).String(s) },
	}

	// Merge app functions (they can override core functions if needed)
	for name, fn := range appFuncs {
		coreFuncs[name] = fn
	}

	return coreFuncs
}

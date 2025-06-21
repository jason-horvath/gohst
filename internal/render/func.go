package render

import (
	"html/template"

	"gohst/internal/config"
)

func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"assetsHead": 	 AssetsHead,
		"isDevelopment": func() bool { return config.App.IsDevelopment() },
		"isProduction":	 func() bool { return config.App.IsProduction() },
		"url":		 	 func() string { return config.App.URL },
		"icon":		 	 Icon,
	}
}

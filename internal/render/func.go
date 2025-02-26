package render

import (
	"html/template"

	"gohst/internal/config"
)

func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"fn_assets_head": 	 AssetsHead,
		"fn_is_development": func() bool { return config.App.IsDevelopment() },
		"fn_is_production":	 func() bool { return config.App.IsProduction() },
	}
}

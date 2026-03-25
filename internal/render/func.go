package render

import (
	"encoding/json"
	"html/template"
	"sort"
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
		"assetsHead": AssetsHead, // Framework asset management
		"icon":       Icon,       // Framework icon helper

		// Template data construction
		"dict": func(pairs ...interface{}) map[string]interface{} {
			m := make(map[string]interface{}, len(pairs)/2)
			for i := 0; i < len(pairs)-1; i += 2 {
				key, _ := pairs[i].(string)
				m[key] = pairs[i+1]
			}
			return m
		},

		// HTML attribute rendering
		"attrs": Attrs, // Render a map as HTML attributes; empty values become boolean attrs

		// Generic utility functions (useful for any app)
		"formatDate": func(t time.Time, format string) string {
			if format == "" {
				format = "January 2, 2006"
			}
			return t.Format(format)
		},
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
		"metaDesc": func(desc string) string {
			desc = strings.TrimSpace(desc)
			if desc == "" {
				return "A marketplace connecting organizations with expert HR practitioners for strategic consulting, assessments, and workforce development."
			}
			return desc
		},
		"jsonLD": func(schema any) template.HTML {
			if schema == nil {
				return ""
			}
			switch v := schema.(type) {
			case string:
				return template.HTML(v)
			default:
				b, err := json.Marshal(v)
				if err != nil {
					return ""
				}
				return template.HTML(string(b))
			}
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": func(s string) string { return cases.Title(language.English).String(s) },
	}

	// Merge app functions (they can override core functions if needed)
	for name, fn := range appFuncs {
		coreFuncs[name] = fn
	}

	return coreFuncs
}

// Attrs renders a map of HTML attributes as a template.HTMLAttr string.
// If a value is empty, only the attribute name is rendered (boolean attribute).
// If a value is non-empty, it renders as key="value".
// Keys are sorted for deterministic output.
func Attrs(attrs map[string]string) template.HTMLAttr {
	if len(attrs) == 0 {
		return ""
	}

	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		b.WriteByte(' ')
		b.WriteString(k)
		if v := attrs[k]; v != "" {
			b.WriteString(`="`)
			b.WriteString(template.HTMLEscapeString(v))
			b.WriteByte('"')
		}
	}
	return template.HTMLAttr(b.String())
}

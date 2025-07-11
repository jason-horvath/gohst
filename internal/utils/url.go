package utils

import (
	"gohst/internal/config"
	"strings"
)

type URLBuilder struct {
	URL string

}
// NewURLBuilder creates a new URLBuilder instance with the provided base URL and distribution path.
func NewURLBuilder(baseURL string) *URLBuilder {
	return &URLBuilder{URL: baseURL}
}

func (u * URLBuilder) normalizeRel(rel string) string {
	if !strings.HasPrefix(rel, "/") {
		rel = "/" + rel
	}

	return rel
}

func (u *URLBuilder) FullURL(rel string) string {
	rel = u.normalizeRel(rel)

	return u.URL + rel
}

func (u *URLBuilder) IsHTTPS() bool {
	return strings.HasPrefix(u.URL, "https")
}


// In internal/utils/url.go
func BuildDistURL(file string) string {
    app := config.GetAppConfig()
    builder := NewURLBuilder(app.GetURL())
    return builder.FullURL(app.GetDistPath() + "/" + file)
}

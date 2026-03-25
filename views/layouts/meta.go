package layouts

import (
	"context"
	"encoding/json"

	"gohst/internal/render"
)

const (
	defaultSiteTitle       = "Gohst"
	defaultSiteDescription = "A lightweight Go web application starter kit with modern frontend capabilities."
	defaultOGType          = "website"
	defaultTwitterCard     = "summary_large_image"
)

// resolvedMeta holds all computed meta values with defaults applied, ready for the layout template.
type resolvedMeta struct {
	Title       string
	Description string
	Canonical   string
	OGImage     string
	OGType      string
	TwitterCard string
	NoIndex     bool
	Schema      string // JSON-encoded, empty if none
}

// resolvePageMeta merges page-level PageMeta with site-wide defaults.
// title is the page title from render.Page; meta may be nil for defaults-only.
func resolvePageMeta(ctx context.Context, title string, meta *render.PageMeta) resolvedMeta {
	appURL := render.AppURL()
	req := render.GetRequestFromCtx(ctx)

	m := resolvedMeta{
		Title:       orDefault(title, defaultSiteTitle),
		Description: defaultSiteDescription,
		Canonical:   req.URL,
		OGImage:     appURL + "/static/images/social/gohst-og-image-1200x630.png",
		OGType:      defaultOGType,
		TwitterCard: defaultTwitterCard,
	}

	if meta == nil {
		return m
	}
	if meta.Title != "" {
		m.Title = meta.Title
	}
	if meta.Description != "" {
		m.Description = meta.Description
	}
	if meta.Canonical != "" {
		m.Canonical = meta.Canonical
	}
	if meta.OGImage != "" {
		m.OGImage = meta.OGImage
	}
	if meta.OGType != "" {
		m.OGType = meta.OGType
	}
	if meta.TwitterCard != "" {
		m.TwitterCard = meta.TwitterCard
	}
	m.NoIndex = meta.NoIndex
	if meta.Schema != nil {
		if b, err := json.Marshal(meta.Schema); err == nil {
			m.Schema = string(b)
		}
	}
	return m
}

func orDefault(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

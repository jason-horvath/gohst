package middleware

import (
	"net/http"

	"gohst/internal/render"
	"gohst/views/pages"
)

type responseWriter404 struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter404) WriteHeader(status int) {
	w.status = status
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *responseWriter404) Write(b []byte) (int, error) {
	if w.status == http.StatusNotFound {
		return len(b), nil // Suppress the default "404 page not found" text
	}
	return w.ResponseWriter.Write(b)
}

// NotFound returns a middleware that intercepts 404 responses and renders a custom 404 page.
func NotFound() func(http.Handler) http.Handler {
	view := render.NewView()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter404{ResponseWriter: w}
			next.ServeHTTP(rw, r)

			if rw.status == http.StatusNotFound {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				view.Render(w, r, pages.NotFoundPage()) //nolint:errcheck
			}
		})
	}
}

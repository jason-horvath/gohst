package middleware

import (
	"io"
	"log"
	"net/http"

	"gohst/internal/config"
)

// SecurityHeaders adds standard security response headers to every request.
// See .ai/security/http-security-headers.md for the full policy spec.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()

		// Prevent MIME type sniffing
		h.Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		h.Set("X-Frame-Options", "DENY")

		// Control referrer information
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Restrict browser features
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")

		// Content Security Policy
		app := config.GetAppConfig()
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-eval'; " +
			"style-src 'self' https://fonts.googleapis.com 'unsafe-inline'; " +
			"font-src 'self' https://fonts.gstatic.com; " +
			"img-src 'self' data: https://images.unsplash.com; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'; " +
			"report-uri /csp-report"

		if app.IsDevelopment() {
			// Allow Vite dev server assets and HMR websocket in development
			csp = "default-src 'self'; " +
				"script-src 'self' 'unsafe-inline' 'unsafe-eval' http://localhost:*; " +
				"style-src 'self' https://fonts.googleapis.com 'unsafe-inline' http://localhost:*; " +
				"font-src 'self' https://fonts.gstatic.com; " +
				"img-src 'self' data: https://images.unsplash.com http://localhost:*; " +
				"connect-src 'self' ws://localhost:* http://localhost:*; " +
				"frame-ancestors 'none'; " +
				"base-uri 'self'; " +
				"form-action 'self'; " +
				"report-uri /csp-report"
		}
		h.Set("Content-Security-Policy", csp)

		// HSTS — only in production (browsers will reject non-HTTPS if sent in dev)
		if app.IsProduction() {
			h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		next.ServeHTTP(w, r)
	})
}

// NoCacheHeaders adds cache-control headers to prevent caching of authenticated pages.
// Wire this into authenticated/dashboard route groups.
func NoCacheHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		h.Set("Pragma", "no-cache")
		h.Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// CSPReportHandler receives Content-Security-Policy violation reports from browsers
// and logs them for monitoring. Returns 204 No Content.
func CSPReportHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(io.LimitReader(r.Body, 10*1024)) // 10 KB max
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("[CSP Violation] %s", string(body))
		w.WriteHeader(http.StatusNoContent)
	})
}

package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"gohst/internal/config"
	"gohst/internal/render"
)

// Recover is a middleware that recovers from panics in HTTP handler goroutines.
// It ensures a single bad request cannot crash the server process.
//
// What this protects: any panic in the handler chain for a specific request.
// What this does NOT protect: panics in background goroutines spawned with `go func()`.
// For those, wrap the goroutine body with a deferred recover() at the call site.
//
// In development: logs the full stack trace and renders it in the browser.
// In production: logs the full stack trace server-side, returns a generic message to the client.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()

				// Always log the full panic + stack trace server-side regardless of environment
				log.Printf("[PANIC RECOVERED] %s %s — %v\n%s", r.Method, r.URL.Path, err, stack)

				app := config.GetAppConfig()
				var errorMsg string
				if app.IsDevelopment() {
					errorMsg = fmt.Sprintf("Panic: %v\n\nStack Trace:\n%s", err, stack)
				} else {
					errorMsg = "An unexpected error occurred. Please try again later."
				}

				render.RenderError(w, "Application Error", errorMsg)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RecoverGoroutine wraps a background goroutine function with panic recovery.
// Use this whenever spawning a goroutine outside the request chain to prevent
// an unhandled panic from crashing the server process.
//
// Usage:
//
//	go middleware.RecoverGoroutine("invoice processor", func() {
//	    processInvoices()
//	})
func RecoverGoroutine(name string, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			log.Printf("[PANIC RECOVERED] goroutine %q — %v\n%s", name, err, stack)
		}
	}()
	fn()
}

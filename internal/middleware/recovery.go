package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"gohst/internal/render"
)

// Recover is a middleware that recovers from panics, logs the error, and renders a nice error page.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				stack := debug.Stack()
				log.Printf("PANIC: %v\n%s", err, stack)

				// Format the error message for the view
				errorMsg := fmt.Sprintf("Panic: %v\n\nStack Trace:\n%s", err, stack)

				// Render the error page
				render.RenderError(w, "Application Panic", errorMsg)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

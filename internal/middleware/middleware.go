package middleware

import "net/http"

// Chain applies multiple middleware in order
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- { // Apply from last to first
		handler = middlewares[i](handler)
	}
	return handler
}

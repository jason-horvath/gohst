package routes

import (
	"net/http"

	"gohst/internal/middleware"
)

// Router defines the contract that application routers must implement
type Router interface {
    SetupRoutes() http.Handler
}

// RegisterRouter allows the application to register its router implementation.
// It wraps the application router with outer framework middleware that applies globally
// to every request before any route group or controller middleware runs:
//
//   - Recover: catches panics so a single bad request cannot crash the server.
//   - SecurityHeaders: sets CSP, frame-options, HSTS, and other hardening headers.
//   - NotFound: intercepts 404 responses and renders the framework not-found page.
func RegisterRouter(r Router) http.Handler {
	return middleware.Chain(
		r.SetupRoutes(),
		middleware.Recover,
		middleware.SecurityHeaders,
		middleware.NotFound(),
	)
}

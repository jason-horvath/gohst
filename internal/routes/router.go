package routes

import (
	"net/http"

	"gohst/internal/middleware"
)

// Router defines the contract that application routers must implement
type Router interface {
    SetupRoutes() http.Handler
}

// RegisterRouter allows the application to register its router implementation
func RegisterRouter(r Router) http.Handler {
    // Wrap the application router with core framework middleware
    return middleware.Recover(r.SetupRoutes())
}

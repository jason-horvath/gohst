package routes

import "net/http"

// Router defines the contract that application routers must implement
type Router interface {
    SetupRoutes() http.Handler
}

// RegisterRouter allows the application to register its router implementation
func RegisterRouter(r Router) http.Handler {
    return r.SetupRoutes()
}

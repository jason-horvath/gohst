package routes

import (
	"net/http"

	"gohst/internal/controllers"
	"gohst/internal/middleware"
	"gohst/internal/session"
)

type RouteConfig struct {
	SessionManager *session.SessionManager
}

func SetupRoutes(rc RouteConfig) http.Handler {
	mux := http.NewServeMux()

	// Single handler for all static files
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	pages := controllers.NewPagesController()
	auth := controllers.NewAuthController()
	mux.HandleFunc("GET /", pages.Index)
	mux.HandleFunc("GET /about", pages.About)
	mux.HandleFunc("GET /post/{id}", pages.Post)
	mux.HandleFunc("GET /login", auth.Login)
	mux.HandleFunc("POST /login", auth.HandleLogin)

	return middleware.Chain(
		mux,
		session.SM.SessionMiddleware,
		middleware.CSRF,
		middleware.Logger,
		middleware.Auth,
	)
}

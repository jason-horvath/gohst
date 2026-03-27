package routes

import (
	"net/http"

	"gohst/app/controllers"
)

type AppRouter struct{}

func NewAppRouter() *AppRouter {
	return &AppRouter{}
}

// Set up all routes and return the mux for the application.
func (r *AppRouter) SetupRoutes() http.Handler {
	mainMux := http.NewServeMux()

	auth := controllers.NewAuthController()
	pages := controllers.NewPagesController()

	fileServer := http.FileServer(http.Dir("static"))
	mainMux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mainMux.Handle("/auth/", http.StripPrefix("/auth", auth.RegisterRoutes()))
	mainMux.Handle("/", pages.RegisterRoutes())

	return mainMux
}

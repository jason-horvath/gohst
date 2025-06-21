package routes

import (
	"net/http"

	"gohst/internal/controllers"
	"gohst/internal/middleware"
	"gohst/internal/session"
)

func SetupRoutes() http.Handler {
	mainMux := http.NewServeMux()

	// Single handler for all static files
	fileServer := http.FileServer(http.Dir("static"))
	mainMux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	publicRoutes := setupPublicRoutes()
	authRoutes := setupAuthRoutes()

	mainMux.Handle("/auth/", http.StripPrefix("/auth", authRoutes)) // Auth routes
	mainMux.Handle("/", publicRoutes) // Public routes



	return mainMux
}

func setupAuthRoutes() http.Handler {
	mux := http.NewServeMux()
	auth := controllers.NewAuthController()

	// Authentication routes
	mux.HandleFunc("GET /login", auth.Login)                // Login page
	mux.HandleFunc("POST /login", auth.HandleLogin)        // Handle login form submission
	mux.HandleFunc("GET /register", auth.Register)         // Registration page
	mux.HandleFunc("POST /register", auth.HandleRegister)  // Handle registration form submission
	mux.HandleFunc("POST /logout", auth.HandleLogout)      // Logout action

	// Apply middleware
	return middleware.Chain(
		mux,
		session.SM.SessionMiddleware,
		middleware.CSRF,
		middleware.Logger,
		middleware.Guest, // Only allow non-authenticated users
	)
}

func setupPublicRoutes() http.Handler {
	mux := http.NewServeMux()
    pages := controllers.NewPagesController()

    // Public informational pages
    mux.HandleFunc("GET /{$}", pages.Index)
	mux.HandleFunc("GET /", pages.NotFound)
	mux.HandleFunc("GET /about", pages.About)
	mux.HandleFunc("GET /post/{id}", pages.Post)

    // Apply middleware
    return middleware.Chain(
        mux,
        session.SM.SessionMiddleware,
        middleware.CSRF,
        middleware.Logger,
    )
}

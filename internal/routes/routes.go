package routes

import (
	"net/http"

	"gohst/internal/controllers"
	"gohst/internal/middleware"
	"gohst/internal/session"
)

// Set up all routes and return the mux for the application.
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

// auth related routes
func setupAuthRoutes() http.Handler {
	mux := http.NewServeMux()
	auth := controllers.NewAuthController()

	// Authentication routes (for guests)
	mux.HandleFunc("GET /login", auth.Login)               // Login page
	mux.HandleFunc("POST /login", auth.HandleLogin)        // Handle login form submission
	mux.HandleFunc("GET /register", auth.Register)         // Registration page
	mux.HandleFunc("POST /register", auth.HandleRegister)  // Handle registration form submission

	guestRoutes := middleware.Chain(
		mux,
		session.SM.SessionMiddleware,
		middleware.CSRF,
		middleware.Logger,
		middleware.Guest, // Only allow non-authenticated users
	)

	// Auth related routes that are require authenticated-only
	authMux := http.NewServeMux()
	authMux.HandleFunc("POST /logout", auth.HandleLogout) // Logout action

	authRoutes := middleware.Chain(
		authMux,
		session.SM.SessionMiddleware,
		middleware.CSRF,
		middleware.Logger,
		middleware.Auth, // Only allow authenticated users
	)

	// Combine both guest and auth routes under a parent mux
	parentMux := http.NewServeMux()
	parentMux.Handle("/", guestRoutes)
	parentMux.Handle("/logout", authRoutes)

	return parentMux
}

// Public routes meant for all users, authenticated or not
func setupPublicRoutes() http.Handler {
	mux := http.NewServeMux()
    pages := controllers.NewPagesController()

    // Public informational pages
    mux.HandleFunc("GET /{$}", pages.Index)
	mux.HandleFunc("GET /", pages.NotFound)
	mux.HandleFunc("GET /post/{id}", pages.Post)

    // Apply middleware
    return middleware.Chain(
        mux,
        session.SM.SessionMiddleware,
        middleware.CSRF,
        middleware.Logger,
    )
}

package middleware

import (
	"net/http"

	"gohst/internal/auth"
	"gohst/internal/session"
)

// Auth middleware ensures that requests are from authenticated users
// If not authenticated, redirects to login page
func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sess := session.FromContext(r.Context())

        // Check if user is authenticated using existing auth function
        if !auth.IsAuthenticated(sess) {
            // Store intended destination for post-login redirect
            sess.SetFlash("error", "Please log in to access this page")

            // Redirect to login page
            http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
            return
        }

        // User is authenticated, proceed to the next handler
        next.ServeHTTP(w, r)
    })
}

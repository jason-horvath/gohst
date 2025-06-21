package middleware

import (
	"net/http"

	"gohst/internal/auth"
	"gohst/internal/session"
)

// Guest middleware - only allows non-authenticated users
func Guest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sess := session.FromContext(r.Context())

        // If user is authenticated, redirect to home page
        if auth.IsAuthenticated(sess) {
            // Optional: add a friendly message
            sess.SetFlash("info", "You are already logged in")

            // Redirect to dashboard or home
            http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
            return
        }

        // User is not authenticated, proceed to pages meant for guests
        next.ServeHTTP(w, r)
    })
}

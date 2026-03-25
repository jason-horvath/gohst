package middleware

import (
	"net/http"

	"gohst/internal/auth"
	"gohst/internal/session"
)

// Role enforces that the authenticated user has one of the allowed roles.
func Role(allowed ...string) func(http.Handler) http.Handler {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, role := range allowed {
		allowedSet[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess := session.FromContext(r.Context())
			authData := auth.GetAuthData(sess)
			if authData == nil {
				sess.SetFlash("error", "Please log in to access this page")
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			rp, ok := authData.(auth.RoleProvider)
			if !ok || rp.RoleName() == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if _, ok := allowedSet[rp.RoleName()]; !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

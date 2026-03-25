package middleware

import (
	"crypto/subtle"
	"gohst/internal/config"
	"gohst/internal/session"
	"gohst/internal/utils"
	"log"
	"net/http"
)

// CSRF protects against Cross-Site Request Forgery attacks.
// Generates a token per session and validates it on state-changing requests.
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.FromContext(r.Context())
		if sess == nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Ensure a CSRF token exists in the session
		token, ok := sess.GetCSRF()
		if !ok || token == "" {
			newToken, err := utils.GenerateCSRF()
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			sess.Set("csrfToken", newToken)
			token = newToken
		}

		// Validate the CSRF token on state-changing requests
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete || r.Method == http.MethodPatch {
			// Check form value first (HTML forms), then X-CSRF-Token header (AJAX/fetch)
			requestToken := r.FormValue(config.APP_CSRF_KEY)
			if requestToken == "" {
				requestToken = r.Header.Get("X-CSRF-Token")
			}
			if requestToken == "" {
				http.Error(w, "CSRF token missing", http.StatusForbidden)
				return
			}

			// Constant-time comparison to prevent timing attacks
			sessionToken, _ := token.(string)
			if subtle.ConstantTimeCompare([]byte(sessionToken), []byte(requestToken)) != 1 {
				app := config.GetAppConfig()
				if app.IsDevelopment() {
					log.Println("CSRF token mismatch")
				}
				http.Error(w, "CSRF token invalid", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

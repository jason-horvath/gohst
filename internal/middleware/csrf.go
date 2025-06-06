package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"gohst/internal/config"
	"gohst/internal/session"
	"log"
	"net/http"
)

// GenerateToken generates a new CSRF token.
func GenerateToken() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(b), nil
}

// CSRFMiddleware is a middleware that protects against CSRF attacks.
func CSRF(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get session from context (same way your render code does)
        sess := session.FromContext(r.Context())
        if sess == nil {
            http.Error(w, "No session found", http.StatusInternalServerError)
            return
        }

        // Check for existing token using the same method your render code uses

        token, ok := sess.Get("csrfToken")
        if !ok || token == "" {
            newToken, err := GenerateToken()
            if err != nil {
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }
            sess.Set("csrfToken", newToken)
            token = newToken
        }

        // Validate the CSRF token on POST requests
        if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
            // Get the CSRF token from the request
            requestToken := r.FormValue(config.App.CSRFName)
            if requestToken == "" {
                http.Error(w, "CSRF token missing", http.StatusBadRequest)
                return
            }

            // Compare the session token with the request token
            if token != requestToken {
                log.Println("CSRF token invalid")
                http.Error(w, "CSRF token invalid", http.StatusBadRequest)
                return
            }

            // Generate a new CSRF token after successful validation
            newToken, err := GenerateToken()
            if err != nil {
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }
            token = newToken
            sess.Set("csrfToken", token)
        }

        // Add the CSRF token to the response context for use in templates
        w.Header().Set("X-CSRF-Token", token.(string))
        next.ServeHTTP(w, r)
    })
}

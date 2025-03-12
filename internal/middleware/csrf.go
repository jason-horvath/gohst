package middleware

import (
	"crypto/rand"
	"encoding/base64"
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
        // Generate a new CSRF token if one doesn't exist in the session
		var sessionID string
		var sessionData *session.SessionData
        sm := session.SM
        sessionData, sessionID = sm.GetSession(r)
        if sessionData == nil {
            sessionManager := session.NewSessionManager()
            sessionID = sessionManager.StartSession(w, r)
        }

        token, ok := sm.GetValue(sessionID, "csrfToken")
        if !ok || token == "" {
            newToken, err := GenerateToken()
            if err != nil {
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }
            token = newToken
            sm.SetValue(sessionID, "csrfToken", token)
        }

        // Validate the CSRF token on POST requests
        if r.Method == http.MethodPost {
            // Get the CSRF token from the request
            requestToken := r.FormValue("csrfToken")
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
            sm.SetValue(sessionID, "csrfToken", token)
        }

        // Add the CSRF token to the response context for use in templates
        w.Header().Set("X-CSRF-Token", token.(string))
        next.ServeHTTP(w, r)
    })
}

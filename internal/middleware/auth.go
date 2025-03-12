package middleware

import (
	"log"
	"net/http"

	"gohst/internal/session"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sm := session.SM
		_, sessionId := sm.GetSession(r)
		log.Println("AUTHORIZATION")
		log.Println("Session ID in Auth middleware:", sessionId)
		sm.SetValue(sessionId, "Authorized", "Is Jason")
		// Add authentication logic here
		// For now, just pass the request
		next.ServeHTTP(w, r)
	})
}

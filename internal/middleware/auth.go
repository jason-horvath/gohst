package middleware

import (
	"log"
	"net/http"

	"gohst/internal/session"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.FromContext(r.Context())
		log.Println("AUTHORIZATION")
		log.Println("Session ID in Auth middleware using context session wrapper:", sess.ID())
		sess.Set("Username", "Jason From Session Set update") // Example of setting a user in the session
		// sm.SetValue(sessionId, "Authorized", "Is Jason")
		// Add authentication logic here
		// For now, just pass the request
		next.ServeHTTP(w, r)
	})
}

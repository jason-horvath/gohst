package session

import (
	"context"
	"log"
	"net/http"
)

// Middleware to attach session to request context
func (sm *SessionManager) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionData, sessionID := sm.GetSession(r)

		if sessionData == nil {
			sessionID = sm.StartSession(w, r)
			sessionData = &SessionData{}
			log.Println("Started a new session with ID:", sessionID)
		} else {
		    sm.SetValue(sessionID, "Name", "Jason")
			name, _ := sm.GetValue(sessionID, "Name")
			log.Println("Name in session:", name)
			log.Println("Session store type:", sm.StoreType)
			log.Println("Session exists with ID:", sessionID)
		}

		// Store session data in request context
		ctx := context.WithValue(r.Context(), sessionIDKey, sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

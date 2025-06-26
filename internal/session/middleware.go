package session

import (
	"context"
	"log"
	"net/http"

	"gohst/internal/config"
)

// Middleware to attach session to request context
func (sm *SessionManager) SessionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //Load or create raw session data
        sessionData, sid := sm.store.GetSession(r)
        if sessionData == nil {
            sessionData, sid = sm.store.StartSession(w, r)

			if config.App.IsDevelopment() {
				log.Println("Started a new session with ID:", sid)
			}
        }

        // Wrap it in Session type
        sess := &Session{
            id:      sid,
            data:    sessionData,
            manager: sm,
            w:       w,
        }

        // Put the *Session into context (so handlers can grab it)
        ctx := context.WithValue(r.Context(), sessionKey, sess)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

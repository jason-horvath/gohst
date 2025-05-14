package session

import (
	"context"
	"log"
	"net/http"
)

// Middleware to attach session to request context
func (sm *SessionManager) SessionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1) Load or create raw session data
        rawData, sid := sm.store.GetSession(r)
        if rawData == nil {
            sid = sm.store.StartSession(w, r)
            rawData, _ = sm.store.GetSession(r)
            log.Println("Started a new session with ID:", sid)
        } else {
            log.Println("Session exists with ID:", sid)
        }

        // 2) Wrap it in our rich Session type
        sess := &Session{
            id:      sid,
            data:    rawData,
            manager: sm,
            w:       w,
        }

        // 3) Put the *Session into context (so handlers can grab it)
        ctx := context.WithValue(r.Context(), sessionKey, sess)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

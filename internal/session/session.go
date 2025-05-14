package session

import (
	"context"
	"net/http"
	"time"
)

type sessKey string

const (
    sessionKey sessKey = "gohst-session" // avoids colliding with SessionIDKey/CSRFKey
)

// Session is what your handlers will actually use
type Session struct {
    id      string
    data    *SessionData
    manager *SessionManager
    w       http.ResponseWriter
}

// FromContext pulls the *Session out of the context (or nil)
func FromContext(ctx context.Context) *Session {
    if sess, _ := ctx.Value(sessionKey).(*Session); sess != nil {
        return sess
    }
    return nil
}

// ID returns the session ID
func (s *Session) ID() string {
    return s.id
}

// Get returns a value (and whether it was present)
func (s *Session) Get(key string) (interface{}, bool)  {
    if s.data == nil {
        return nil, false
    }
    val, ok := s.data.Values[key]
    return val, ok
}

// Set writes a value, persists to the store, and re-sets the cookie
func (s *Session) Set(key string, val interface{}) {
    s.data.Values[key] = val
    s.manager.store.SetValue(s.id, key, val)
    // refresh cookie so client sees it
    http.SetCookie(s.w, &http.Cookie{
        Name:     SESSION_NAME,
        Value:    s.id,
        Path:     "/",
        HttpOnly: true,
        Expires:  time.Now().Add(GetSessionLength()),
    })
}

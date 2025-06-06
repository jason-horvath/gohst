package session

import (
	"context"
	"net/http"
	"time"
)

type sessKey string
type flashMessageKey string
type oldDataKey string

const (
    sessionKey sessKey = "gohst-session" // avoids colliding with SessionIDKey/CSRFKey
	flashKey flashMessageKey = "_gohst_flash_" // for flash messages
	oldKey oldDataKey = "_gohst_old_" // for old values in forms/data
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

func (s *Session) GetCSRF() (interface{}, bool) {
	if s.data == nil {
        return nil, false
    }
    val, ok := s.data.Values[string(CSRFKey)]
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

// SetOld stores a form value for repopulation after a redirect
func (s *Session) SetOld(key string, val interface{}) {
    s.Set("_old_"+key, val)
}

// GetOld retrieves a form value and removes it from the session
func (s *Session) GetOld(key string) interface{} {
    oldKey := "_old_" + key
    val, exists := s.Get(oldKey)
    if !exists {
        return nil
    }

    // Remove after retrieving
    s.Remove(oldKey)
    return val
}

// Remove removes a key from the session
func (s *Session) Remove(key string) {
    if s.data == nil {
        return
    }
    delete(s.data.Values, key)

    // Persist the change to the storage backend
    s.manager.store.Remove(s.id, key)
}

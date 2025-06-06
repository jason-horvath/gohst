package session

import (
	"context"
	"net/http"
)

type SessionStore interface {
	StartSession(w http.ResponseWriter, r *http.Request) (*SessionData, string)
	GetSession(r *http.Request) (*SessionData, string)
	SetValue(sessionID string, key string, value interface{})
	GetValue(sessionID string, key string) (interface{}, bool)
	GetSessionByID(ctx context.Context, sessionID string) (*SessionData, error)
	Remove(sessionID string, key string) error
}

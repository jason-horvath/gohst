package session

import (
	"context"
	"gohst/internal/config"
	"time"
)

// SessionData stores session values
type SessionData struct {
	ID 		string					`json:"id"`
	Values  map[string]interface{} 	`json:"values"`
	Expires time.Time				`json:"expires"`
	manager SessionStore			`json:"-"`
}

// GetSessionLength retrieves the session duration from ENV and returns it as time.Duration
func GetSessionLength() time.Duration {
	// Read from config
	session := config.Session
	sessionLength := config.Session.Length

	if session != nil {
		if session.Length > 0 {
			sessionLength = session.Length
		}
	}

	// Convert minutes to time.Duration
	return time.Duration(sessionLength) * time.Minute
}

// Set a session value through the session manager to maintain one source of truth
func (sd *SessionData) Set(key string, value any)  {
    sd.Values[key] = value
    sd.manager.SetValue(sd.ID, key, value);
	updated, _ := sd.manager.GetSessionByID(context.Background(), sd.ID)
	sd.Values = updated.Values
}

// Get Session
func (sd *SessionData) Get(key string) (any, bool) {
    val, ok := sd.manager.GetValue(sd.ID, key)
    return val, ok
}

func (sd *SessionData) CSRF() (string, bool) {
	val, ok := sd.manager.GetValue(sd.ID, string(CSRFKey))
	if !ok {
		return "", false
	}
	return val.(string), true
}

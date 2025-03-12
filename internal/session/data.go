package session

import (
	"gohst/internal/config"
	"time"
)

// SessionData stores session values
type SessionData struct {
	Values  map[string]interface{} 	`json:"values"`
	Expires time.Time				`json:"expires"`
}

// GetSessionLength retrieves the session duration from ENV and returns it as time.Duration
func GetSessionLength() time.Duration {
	// Read from config
	session := config.Session
	sessionLength := SESSION_LENGTH_DEFAULT

	if session != nil {
		if session.Length > 0 {
			sessionLength = session.Length
		}
	}

	// Convert minutes to time.Duration
	return time.Duration(sessionLength) * time.Minute
}

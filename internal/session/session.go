package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"gohst/internal/config"
)

const SESSION_LENGTH_DEFAULT = 60

const SESSION_NAME = "session_id"

const SESSION_STORE_DEFAULT = "file"

const (
	SESSION_TYPE_FILE  = "file"
	SESSION_TYPE_REDIS = "redis"
)

var SESSION_VALID_TYPES = []string{
	SESSION_TYPE_FILE,
	SESSION_TYPE_REDIS,
}

// SessionData stores session values
type SessionData struct {
	Values  map[string]interface{} 	`json:"values"`
	Expires time.Time				`json:"expires"`
}

// SessionManager using Redis
type SessionManager struct {
	StoreType string
	store SessionStore
}

// NewSessionManager initializes Redis connection
func NewSessionManager() *SessionManager {
	var store SessionStore
	var storeType string

	sessionStore := config.GetEnv("SESSION_STORE", SESSION_STORE_DEFAULT).(string)

	if !IsValidSessionType(sessionStore) {
		sessionStore = SESSION_STORE_DEFAULT
	}

	if sessionStore == "redis" {
		store, storeType = NewRedisSessionManager() // Redis session manager
	} else {
		sessionFilePath := config.GetEnv("SESSION_FILE_PATH", SESSION_FILE_PATH_DEFAULT).(string)
		store, storeType = NewFileSessionManager(sessionFilePath) // File-based session manager
	}

	return &SessionManager{StoreType: storeType, store: store}
}

// GenerateSessionID creates a unique session ID
func GenerateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// StartSession creates a session and stores it in Redis
func (sm *SessionManager) StartSession(w http.ResponseWriter, r *http.Request) string {
	return sm.store.StartSession(w, r)
}

// GetSession retrieves session data from Redis
func (sm *SessionManager) GetSession(r *http.Request) (*SessionData, string) {
	return sm.store.GetSession(r)
}

// SetValue stores a value in the session
func (sm *SessionManager) SetValue(sessionID string, key string, value interface{}) {
	sm.store.SetValue(sessionID, key, value)
}

// GetValue retrieves a session value
func (sm *SessionManager) GetValue(sessionID string, key string) (interface{}, bool) {
	return sm.store.GetValue(sessionID, key)
}

// GetSessionLength retrieves the session duration from ENV and returns it as time.Duration
func GetSessionLength() time.Duration {
	// Read from config
	session := config.Session
	sessionLength := SESSION_LENGTH_DEFAULT

	if session == nil {
		if session.Length > 0 {
			sessionLength = session.Length
		}
	}

	// Convert minutes to time.Duration
	return time.Duration(sessionLength) * time.Minute
}

// Function to check if a session type is valid
func IsValidSessionType(value string) bool {
	for _, v := range SESSION_VALID_TYPES {
		if v == value {
			return true
		}
	}
	return false
}

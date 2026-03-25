package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"gohst/internal/config"
)

type contextKey string

type csrfKey string

const SessionIDKey contextKey = "sessionID"

const CSRFKey csrfKey = "csrfToken"

// SessionName returns the configured session cookie name.
func SessionName() string {
	if config.Session != nil && config.Session.Name != "" {
		return config.Session.Name
	}
	return "_gohst_session"
}

const SESSION_STORE_DEFAULT = "file"

const (
	SESSION_TYPE_FILE  = "file"
	SESSION_TYPE_REDIS = "redis"
)

var SESSION_VALID_TYPES = []string{
	SESSION_TYPE_FILE,
	SESSION_TYPE_REDIS,
}

var SM *SessionManager
var SMAdmin *SessionManager

// SessionManager using Redis
type SessionManager struct {
	StoreType  string
	store      SessionStore
	cookieName string
}

// Initialize the session setup
func Init() {
	InitSessionManager()
}

// Initialize the session manager
func InitSessionManager() {
	SM = NewSessionManager()
	SMAdmin = NewSessionManagerWithName(config.Session.AdminName)
}

// NewSessionManager initializes Redis connection
func NewSessionManager() *SessionManager {
	return NewSessionManagerWithName(SessionName())
}

// NewSessionManagerWithName initializes session manager with a specific cookie name.
func NewSessionManagerWithName(cookieName string) *SessionManager {
	var store SessionStore
	var storeType string

	if cookieName == "" {
		cookieName = SessionName()
	}

	sessionStore := config.GetEnv("SESSION_STORE", SESSION_STORE_DEFAULT).(string)

	if !IsValidSessionType(sessionStore) {
		sessionStore = SESSION_STORE_DEFAULT
	}

	if sessionStore == "redis" {
		store, storeType = NewRedisSessionManager(cookieName) // Redis session manager
	} else {
		sessionFilePath := config.GetEnv("SESSION_FILE_PATH", SESSION_FILE_PATH_DEFAULT).(string)
		store, storeType = NewFileSessionManager(sessionFilePath, cookieName) // File-based session manager
	}

	return &SessionManager{
		StoreType:  storeType,
		store:      store,
		cookieName: cookieName,
	}
}

// GenerateSessionID creates a unique session ID
func GenerateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// StartSession creates a session and stores it in Redis
func (sm *SessionManager) StartSession(w http.ResponseWriter, r *http.Request) (*SessionData, string) {
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
func (sm *SessionManager) GetCSRF(sessionID string) (string, bool) {
	value, ok := sm.GetValue(sessionID, string(CSRFKey))
	if !ok {
		return "", false
	}
	return value.(string), true
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

// Save saves the entire session data
func (sm *SessionManager) Save(sessionID string, session *SessionData) error {
	return sm.store.Save(sessionID, session)
}

// Delete removes an entire session
func (sm *SessionManager) Delete(sessionID string) error {
	return sm.store.Delete(sessionID)
}

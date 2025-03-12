package session

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)


const SESSION_FILE_PATH_DEFAULT = "tmp/sessions"

const SESSION_FILE_EXT = ".session"

// FileSessionManager manages file-based sessions
type FileSessionManager struct {
	sessions map[string]*SessionData
	mu       sync.Mutex
	dir      string
}

// NewFileSessionManager initializes a file-based session manager
func NewFileSessionManager(storageDir string) (*FileSessionManager, string) {
	os.MkdirAll(storageDir, 0755) // Ensure session directory exists
	return &FileSessionManager{
		sessions: make(map[string]*SessionData),
		dir:      storageDir,
	}, "file"
}

// StartSession creates a session and writes it to a file
func (fsm *FileSessionManager) StartSession(w http.ResponseWriter, r *http.Request) string {
	sessionID := GenerateSessionID()
	sessionLength := GetSessionLength()
	sessionData := &SessionData{Values: make(map[string]interface{}), Expires: time.Now().Add(sessionLength)}

	fsm.mu.Lock()
	fsm.sessions[sessionID] = sessionData
	fsm.mu.Unlock()

	fsm.saveSession(sessionID, sessionData)

	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_NAME,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
	})

	return sessionID
}

// GetSession retrieves session data from file
func (fsm *FileSessionManager) GetSession(r *http.Request) (*SessionData, string) {
	var sessionID string
	cookie, err := r.Cookie(SESSION_NAME)
	if err != nil {
		ctxSessionID, ok := r.Context().Value(SessionIDKey).(string)
		if !ok {
			log.Println("Session ID not found in context")
			return nil, ""
		}
		sessionID = ctxSessionID
	} else {
		sessionID = cookie.Value
	}

	fsm.mu.Lock()
	session, exists := fsm.sessions[sessionID]
	fsm.mu.Unlock()

	// If session is not found in memory, try loading from file
	if !exists {
		session = fsm.loadSession(sessionID)
		if session != nil {
			fsm.mu.Lock()
			fsm.sessions[sessionID] = session // Cache it in memory
			fsm.mu.Unlock()
		}
	}

	// Expired session check
	if session == nil || time.Now().After(session.Expires) {
		return nil, ""
	}

	return session, sessionID
}

// SetValue stores a value in a session file
func (fsm *FileSessionManager) SetValue(sessionID string, key string, value interface{}) {
	fsm.mu.Lock()
	if session, exists := fsm.sessions[sessionID]; exists {
		session.Values[key] = value
		fsm.saveSession(sessionID, session)
	}
	fsm.mu.Unlock()
}

// GetValue retrieves a value from a session file
func (fsm *FileSessionManager) GetValue(sessionID string, key string) (interface{}, bool) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	if session, exists := fsm.sessions[sessionID]; exists {
		val, ok := session.Values[key]
		return val, ok
	}
	return nil, false
}

// Save session to file
func (fsm *FileSessionManager) saveSession(sessionID string, session *SessionData) {
	filePath := filepath.Join(fsm.dir, sessionID+SESSION_FILE_EXT)
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error saving session:", err)
		return
	}
	defer file.Close()
	gob.NewEncoder(file).Encode(session)
}

// Load session from file
func (fsm *FileSessionManager) loadSession(sessionID string) *SessionData {
	filePath := filepath.Join(fsm.dir, sessionID+SESSION_FILE_EXT)
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()
	var session SessionData
	gob.NewDecoder(file).Decode(&session)
	return &session
}

// GetSessionByID fetches session data directly using session ID
func (fsm *FileSessionManager) GetSessionByID(ctx context.Context, sessionID string) (*SessionData, error) {
	filePath := filepath.Join(fsm.dir, sessionID+SESSION_FILE_EXT)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var session SessionData
	gob.NewDecoder(file).Decode(&session)
	return &session, nil
}

// CleanupExpiredSessions removes old session files
func (fsm *FileSessionManager) CleanupExpiredSessions() {
	files, err := os.ReadDir(fsm.dir)
	if err != nil {
		log.Println("Error reading session directory:", err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(fsm.dir, file.Name())

		f, err := os.Open(filePath)
		if err != nil {
			continue
		}
		defer f.Close()

		var session SessionData
		if err := gob.NewDecoder(f).Decode(&session); err != nil {
			os.Remove(filePath) // Corrupted session file
			continue
		}

		if time.Now().After(session.Expires) {
			os.Remove(filePath) // Remove expired session
		}
	}
}

// Cleanup expired sessions periodically
func (fsm *FileSessionManager) StartSessionCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			fsm.CleanupExpiredSessions()
		}
	}()
}



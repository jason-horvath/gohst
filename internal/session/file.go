package session

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gohst/internal/config"
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
	app := config.GetAppConfig()
    // Convert to absolute path if relative
    if !filepath.IsAbs(storageDir) {
        // Get current working directory
        cwd, err := os.Getwd()
        if err == nil {
            storageDir = filepath.Join(cwd, storageDir)

			if app.IsDevelopment() {
				log.Println("Using absolute session path:", storageDir)
			}
        }
    }

    os.MkdirAll(storageDir, 0755) // Ensure session directory exists
    return &FileSessionManager{
        sessions: make(map[string]*SessionData),
        dir:      storageDir,
    }, "file"
}

func (fsm *FileSessionManager) StartSession(w http.ResponseWriter, r *http.Request) (*SessionData, string) {
    // First check if a valid session already exists
    if r != nil {
        cookie, err := r.Cookie(SESSION_NAME)
        if err == nil && cookie.Value != "" {
            // Try to load existing session from file
            existingSession := fsm.loadSession(cookie.Value)
			log.Println("Checking for existing session:", existingSession)
            if existingSession != nil && time.Now().Before(existingSession.Expires) {
                // Valid session found, use it
                fsm.mu.Lock()
                fsm.sessions[cookie.Value] = existingSession // Update memory cache
                fsm.mu.Unlock()
                log.Printf("Reusing existing session: %s", cookie.Value)
                return existingSession, cookie.Value
            }
        }
    }

    // No valid session found, create a new one
    sessionID := GenerateSessionID()
    sessionLength := GetSessionLength()
    sessionData := &SessionData{
        ID:      sessionID,
        Values:  make(map[string]any),
        Expires: time.Now().Add(sessionLength),
    }

    fsm.mu.Lock()
    fsm.sessions[sessionID] = sessionData
    fsm.mu.Unlock()

    fsm.saveSession(sessionID, sessionData)

    http.SetCookie(w, &http.Cookie{
        Name:     SESSION_NAME,
        Value:    sessionID,
        Path:     "/",
        HttpOnly: true,
        Expires:  time.Now().Add(sessionLength),
    })

    log.Printf("Created new session: %s", sessionID)
    return sessionData, sessionID
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
    defer fsm.mu.Unlock()

    session, exists := fsm.sessions[sessionID]
    if !exists {
        // Try to load from file if not in memory
        session = fsm.loadSession(sessionID)
        if session != nil {
            fsm.sessions[sessionID] = session
        } else {
            // Create new session if doesn't exist
            session = &SessionData{
                ID:      sessionID,
                Values:  make(map[string]any),
                Expires: time.Now().Add(GetSessionLength()),
            }
            fsm.sessions[sessionID] = session
        }
    }

    // Now we have a valid session
    session.Values[key] = value
    fsm.saveSession(sessionID, session)
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
func (fsm *FileSessionManager) saveSession(sessionID string, session *SessionData) error {
    // Create a temporary file first
    tempPath := filepath.Join(fsm.dir, sessionID+".tmp")
    file, err := os.Create(tempPath)
    if err != nil {
        return err
    }

    // Encode to temporary file
    if err := gob.NewEncoder(file).Encode(session); err != nil {
        file.Close()
        os.Remove(tempPath)
        return err
    }

    // Ensure all data is written to disk
    if err := file.Sync(); err != nil {
        file.Close()
        os.Remove(tempPath)
        return err
    }

    // Close file
    if err := file.Close(); err != nil {
        os.Remove(tempPath)
        return err
    }

    // Rename to final destination (atomic operation on most filesystems)
    finalPath := filepath.Join(fsm.dir, sessionID+SESSION_FILE_EXT)
    return os.Rename(tempPath, finalPath)
}

// Load session from file
func (fsm *FileSessionManager) loadSession(sessionID string) *SessionData {
	app := config.GetAppConfig()
    filePath := filepath.Join(fsm.dir, sessionID+SESSION_FILE_EXT)

    file, err := os.Open(filePath)
    if err != nil {
        log.Printf("Error opening session file: %v", err)
        return nil
    }
    defer file.Close()

    var session SessionData
    if err := gob.NewDecoder(file).Decode(&session); err != nil {
        log.Printf("Error decoding session: %v", err)
        os.Remove(filePath)
        return nil
    }

	if app.IsDevelopment() {
		log.Printf("Successfully loaded session: %s", sessionID)
	}

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

// RemoveValue deletes a key from the session data
func (fsm *FileSessionManager) Remove(sessionID string, key string) error {
    fsm.mu.Lock()
    defer fsm.mu.Unlock()

    // Check if session exists in memory
    session, exists := fsm.sessions[sessionID]
    if !exists {
        session = fsm.loadSession(sessionID)
        if session == nil {
            return fmt.Errorf("session not found: %s", sessionID)
        }
        fsm.sessions[sessionID] = session
    }

    delete(session.Values, key)

    fsm.saveSession(sessionID, session)

    return nil
}

// Save saves the entire session
func (fsm *FileSessionManager) Save(sessionID string, session *SessionData) error {
    fsm.mu.Lock()
    defer fsm.mu.Unlock()

    fsm.sessions[sessionID] = session
    return fsm.saveSession(sessionID, session)
}

// Delete removes the entire session
func (fsm *FileSessionManager) Delete(sessionID string) error {
    fsm.mu.Lock()
    defer fsm.mu.Unlock()

    delete(fsm.sessions, sessionID)
    filePath := filepath.Join(fsm.dir, sessionID+SESSION_FILE_EXT)
    return os.Remove(filePath)
}

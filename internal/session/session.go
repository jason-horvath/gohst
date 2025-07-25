package session

import (
	"context"
	"log"
	"net/http"
	"time"

	"gohst/internal/config"
	"gohst/internal/utils"
)

type sessKey string
type flashMessageKey string
type oldDataKey string
type fieldErrorsKey string

const (
    sessionKey sessKey = "gohst_session" // avoids colliding with SessionIDKey/CSRFKey
    flashKey flashMessageKey = "_gohst_flash_" // for flash messages
    oldKey oldDataKey = "_gohst_old_" // for old values in forms/data
    fieldErrorsPrefix fieldErrorsKey = "_gohst_field_errors_" // for field-specific errors
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
func (s *Session) Get(key string) (any, bool)  {
    if s.data == nil {
        return nil, false
    }
    val, ok := s.data.Values[key]
    return val, ok
}

// Get CSRF Token from the session
func (s *Session) GetCSRF() (any, bool) {
    if s.data == nil {
        return nil, false
    }
    val, ok := s.data.Values[string(CSRFKey)]
    return val, ok
}

// Set the CSRF Token in the session
func (s *Session) SetCSRF(csrf string) *Session {
    if csrf == "" {
        log.Println("CSRF token is nil, generating a new one")
        csrf, _ = utils.GenerateCSRF() // Generate a new CSRF token
    }
    s.data.Values[string(CSRFKey)] = csrf
    return s
}

// RemoveCSRF removes the CSRF token from the session
func (s *Session) RemoveCSRF() {
    if s.data == nil {
        return
    }
    s.Remove(string(CSRFKey))
}

// Set writes a value, persists to the store, and re-sets the cookie
func (s *Session) Set(key string, val any) {
    s.data.Values[key] = val
    s.manager.store.SetValue(s.id, key, val)
    // refresh cookie so client sees it

    s.setSessionCookie()
}

// setSessionCookie is a helper to standardize cookie settings
func (s *Session) setSessionCookie() {
    // Get environment from config
    app := config.GetAppConfig()
    isProduction := app.IsProduction()

    cookie := &http.Cookie{
        Name:     SESSION_NAME,
        Value:    s.id,
        Path:     "/",
        HttpOnly: true,
        Secure:   isProduction,
        SameSite: http.SameSiteLaxMode, // Set to Lax by default
        Expires:  time.Now().Add(GetSessionLength()),
    }

    // Set to strict for production
    if isProduction {
        cookie.SameSite = http.SameSiteStrictMode
    }

    http.SetCookie(s.w, cookie)
}

// SetFlash stores a flash message that will be displayed once
func (s *Session) SetFlash(key string, val any) {
    flashKey := string(flashKey) + key
    s.Set(flashKey, val)
}

// GetFlash retrieves a flash message and removes it from the session
func (s *Session) GetFlash(key string) any {
    fullKey := string(flashKey) + key
    val, exists := s.Get(fullKey)
    if !exists {
        return nil
    }

    // Remove after retrieving
    s.Remove(fullKey)
    return val
}

// GetAllFlash retrieves all flash messages and removes them
func (s *Session) GetAllFlash() map[string]any {
    if s.data == nil {
        return nil
    }

    flashMessages := make(map[string]any)
    prefix := string(flashKey)
    prefixLen := len(prefix)

    // Find all keys with the flash prefix
    for key, val := range s.data.Values {
        if len(key) > prefixLen && key[:prefixLen] == prefix {
            // Extract the actual key name without prefix
            actualKey := key[prefixLen:]
            flashMessages[actualKey] = val
            // Remove after retrieving
            s.Remove(key)
        }
    }

    return flashMessages
}

// SetOld stores a form value for repopulation after a redirect
func (s *Session) SetOld(key string, val any) {
    s.Set(string(oldKey) + key, val)
}

// GetOld retrieves a form value and removes it from the session
func (s *Session) GetOld(key string) any {
    oldKey := string(oldKey) + key
    val, exists := s.Get(oldKey)
    if !exists {
        return nil
    }

    // Remove after retrieving
    s.Remove(oldKey)
    return val
}

// GetAllOld retrieves all old form values and removes them from the session
func (s *Session) GetAllOld() map[string]any {
    if s.data == nil {
        return nil
    }

    oldValues := make(map[string]any)
    prefix := string(oldKey)
    prefixLen := len(prefix)

    // Find all keys with the old prefix
    for key, val := range s.data.Values {
        if len(key) > prefixLen && key[:prefixLen] == prefix {
            // Extract the actual key name without prefix
            actualKey := key[prefixLen:]
            oldValues[actualKey] = val
            // Remove after retrieving
            s.Remove(key)
        }
    }

    return oldValues
}

// Add these methods that don't clear data
func (s *Session) PeekOld(key string) (any, bool) {
    return s.Get(string(oldKey) + key)
}

// Add these methods that don't clear data
func (s *Session) PeekFlash(key string) (any, bool) {
    return s.Get(string(flashKey) + key)
}

// Similar to GetAllFlash but without clearing
func (s *Session) PeekAllFlash() map[string]any {
    return s.getKeysByPrefix(string(flashKey))
}

// Similarly for old data
func (s *Session) PeekAllOld() map[string]any {
    return s.getKeysByPrefix(string(oldKey))
}

// Add these methods to your session package

// SetFieldErrors stores a slice of error messages for a field
func (s *Session) SetFieldErrors(field string, errors []string) {
    s.Set(string(fieldErrorsPrefix)+field, errors)
}

// AddFieldError appends an error message to the slice for a field, creating the slice if needed
func (s *Session) AddFieldError(field string, errMsg string) {
    key := string(fieldErrorsPrefix) + field
    val, ok := s.Get(key)
    var errs []string
    if ok {
        if existing, ok := val.([]string); ok {
            errs = existing
        }
    }
    errs = append(errs, errMsg)
    s.SetFieldErrors(field, errs)
}

// GetFieldErrors retrieves a slice of error messages for a field and removes it from the session
func (s *Session) GetFieldErrors(field string) []string {
    val, ok := s.Get(string(fieldErrorsPrefix) + field)
    s.Remove(string(fieldErrorsPrefix) + field)
    if ok {
        if errs, ok := val.([]string); ok {
            return errs
        }
    }
    return []string{}
}

// PeekFieldErrors retrieves a slice of error messages for a field without removing it from the session
// PeekFieldErrors retrieves a slice of error messages for a field without removing it from the session
func (s *Session) PeekFieldErrors(field string) []string {
    val, ok := s.Get(string(fieldErrorsPrefix) + field)
    if ok {
        if errs, ok := val.([]string); ok {
            return errs
        }
    }
    return []string{}
}

// PeekAllFieldErrors retrieves all field errors as a map[string][]string without removing them from the session
func (s *Session) PeekAllFieldErrors() map[string][]string {
    raw := s.getKeysByPrefix(string(fieldErrorsPrefix))
    result := make(map[string][]string)
    for k, v := range raw {
        if errs, ok := v.([]string); ok {
            result[k] = errs
        }
    }
    return result
}

// GetAllFieldErrors retrieves all field errors as a map[string][]string and removes them from the session
func (s *Session) GetAllFieldErrors() map[string][]string {
    if s.data == nil {
        return nil
    }
    fieldErrors := make(map[string][]string)
    prefix := string(fieldErrorsPrefix)
    prefixLen := len(prefix)
    for key, val := range s.data.Values {
        if len(key) > prefixLen && key[:prefixLen] == prefix {
            fieldName := key[prefixLen:]
            if errs, ok := val.([]string); ok {
                fieldErrors[fieldName] = errs
            }
            s.Remove(key)
        }
    }
    return fieldErrors
}

// GetFieldError retrieves a field-specific error and removes it from the session
func (s *Session) GetFieldError(field string) (string, bool) {
    val, ok := s.Get(string(fieldErrorsPrefix) + field)
    if !ok {
        return "", false
    }
    // Clear after retrieving
    s.Remove(string(fieldErrorsPrefix) + field)
    return val.(string), true
}

// PeekFieldError retrieves a field-specific error without removing it from the session
func (s *Session) PeekFieldError(field string) (string, bool) {
    val, ok := s.Get(string(fieldErrorsPrefix) + field)
    if !ok {
        return "", false
    }
    return val.(string), true
}

// PeekAllFieldError retrieves all field errors without removing them from the session
func (s *Session) PeekAllFieldError() map[string]any {
    return s.getKeysByPrefix(string(fieldErrorsPrefix))
}

// GetAllFieldError retrieves all field errors and removes them from the session
func (s *Session) GetAllFieldError() map[string]string {
    if s.data == nil {
        return nil
    }

    fieldErrors := make(map[string]string)
    prefix := string(fieldErrorsPrefix)
    prefixLen := len(prefix)

    // Find all keys with the field_error prefix
    for key, val := range s.data.Values {
        if len(key) > prefixLen && key[:prefixLen] == prefix {
            // Extract the actual field name without prefix
            fieldName := key[prefixLen:]
            if strVal, ok := val.(string); ok {
                fieldErrors[fieldName] = strVal
            }
            // Remove after retrieving
            s.Remove(key)
        }
    }

    return fieldErrors
}

// getKeysByPrefix retrieves all values with a specific prefix without removing them
func (s *Session) getKeysByPrefix(prefix string) map[string]any {
    if s.data == nil {
        return nil
    }

    result := make(map[string]any)
    prefixLen := len(prefix)

    // Find all keys with the given prefix
    for key, val := range s.data.Values {
        if len(key) > prefixLen && key[:prefixLen] == prefix {
            // Extract the actual key name without prefix
            actualKey := key[prefixLen:]
            result[actualKey] = val
            // No removal - this is the key difference from GetAllFlash/GetAllOld
        }
    }

    return result
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

// Regenerate creates a new session ID while preserving important data
func (s *Session) Regenerate() {
   // Save essential data before regeneration
    csrfToken, hasCSRF := s.Get("csrfToken")

    // Generate new session ID
    oldID := s.id
    s.id = GenerateSessionID() // Your existing method

    // Create new session data
    oldValues := s.data.Values
    s.data = &SessionData{
        Values: make(map[string]interface{}),
    }

    // Copy ALL old values to new session
    for k, v := range oldValues {
        s.data.Values[k] = v
    }

    // Ensure CSRF token is preserved (belt and suspenders)
    if hasCSRF {
        s.data.Values["csrfToken"] = csrfToken
    }

    // Save and delete
    s.manager.Save(s.id, s.data)
    s.manager.Delete(oldID)

    // Update cookie
    s.setSessionCookie()
}


// RegenerateNew creates a completely new session with no preserved values
func (s *Session) RegenerateNew() {
    // Generate new session ID
    oldID := s.id
    s.id = GenerateSessionID()

    // Create completely empty session data
    s.data = &SessionData{
        Values: make(map[string]interface{}),
    }

    // Generate fresh CSRF token for security
    s.SetCSRF("")

    // Save the new empty session
    s.manager.Save(s.id, s.data)

    // Delete the old session completely from storage
    s.manager.Delete(oldID)

    // Update cookie with new session ID
    s.setSessionCookie()
}

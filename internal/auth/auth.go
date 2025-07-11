package auth

import (
	"gohst/internal/session"
)

const AuthKey = "_gohst_auth_"


// AuthDataProvider is an interface for providing authentication data
type AuthDataProvider interface {
	Data() any
}

// GetAuthData retrieves auth data from the session
// Auth data must implement AuthDataProvider interface
func GetAuthData(sess *session.Session) any {
    auth, ok := sess.Get(AuthKey)
    if !ok || auth == nil {
        return nil
    }

    // Enforce that auth data must implement AuthDataProvider
    provider, ok := auth.(AuthDataProvider)
    if !ok {
        panic("auth data must implement AuthDataProvider interface")
    }

    return provider.Data()
}

// IsAuthenticated checks if a user is authenticated
func IsAuthenticated(sess *session.Session) bool {
    return GetAuthData(sess) != nil
}

// Logout completely clears the session for security
func Logout(sess *session.Session) {
    sess.RegenerateNew()
}

package auth

import (
	"errors"
	"time"

	"gohst/internal/models"
	"gohst/internal/session"
	"gohst/internal/utils"
)

const authKey = "_gohst_auth_"

// Login attempts to authenticate a user with email and password
// Returns the authenticated user and any error that occurred
func Login(sess *session.Session, email, password string) (*models.User, error) {
    // Find user in database
    userModel := models.NewUserModel()
    user, err := userModel.FindByEmail(email)
    if err != nil {
        return nil, err
    }

    // Verify password
    passwordOk, _ := utils.CheckPassword(password, user.PasswordHash)
    if !passwordOk {
        return nil, errors.New("invalid credentials")
    }

    // Store authentication data in session
    authData := &AuthData{
        UserID:     user.ID,
        Email:      user.Email,
        Name:       user.FirstName,
        IsAdmin:    user.RoleID == 1,
        LoggedInAt: time.Now(),
    }

    sess.Set(authKey, authData)
    return user, nil
}

// GetAuthData retrieves auth data from the session
func GetAuthData(sess *session.Session) *AuthData {
    auth, ok := sess.Get(authKey)
    if !ok || auth == nil {
        return nil
    }
    return auth.(*AuthData)
}

// IsAuthenticated checks if a user is authenticated
func IsAuthenticated(sess *session.Session) bool {
    return GetAuthData(sess) != nil
}

// Logout removes auth data from the session
func Logout(sess *session.Session) {
    sess.Remove(authKey)
}

package auth

import (
	"encoding/gob"
	"errors"
	"time"

	"gohst/internal/models"
	"gohst/internal/session"
	"gohst/internal/utils"
)

const authKey = "_gohst_auth_"
type AuthData struct {
    UserID     uint64
    Email      string
    Name       string
    IsAdmin    bool
    LoggedInAt time.Time
}

func init() {
    // Initialization code here
    gob.Register(&AuthData{})
}

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

	roleModel := models.NewRoleModel()

	role, err := roleModel.FindByID(user.RoleID)

	if err != nil {
		return nil, errors.New("role not found")
	}

	isAdmin := role.Name == "admin"

    // Store authentication data in session
    authData := &AuthData{
        UserID:     user.ID,
        Email:      user.Email,
        Name:       user.FirstName,
        IsAdmin:    isAdmin,
        LoggedInAt: time.Now(),
    }

    sess.Set(authKey, authData)
    return user, nil
}

// Register creates a new user account
func Register(email, firstName, lastName, password string) error {
    // Check if email already exists
    userModel := models.NewUserModel()
    existingUser, err := userModel.FindByEmail(email)
    if err == nil && existingUser != nil {
        return errors.New("email already in use")
    }

    // Hash the password
    passwordHash, err := utils.HashPassword(password)
    if err != nil {
        return errors.New("error processing password")
    }

    // Get the default user role (assuming "user" role exists)
    roleModel := models.NewRoleModel()
    role, err := roleModel.FindByName("user")
    if err != nil {
        return errors.New("default role not found")
    }

    // Create the user
    user := &models.User{
        FirstName:    firstName,
        LastName:     lastName,
        Email:        email,
        PasswordHash: passwordHash,
        RoleID:       role.ID,
        Active:       true,
    }

    // Save to database
    err = userModel.Create(user)
    if err != nil {
        return errors.New("failed to create user: " + err.Error())
    }

    return nil
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

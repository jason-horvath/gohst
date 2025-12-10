package models

import (
	"time"
)

type User struct {
	ID		   	 uint64 `db:"id"`
	FirstName    string `db:"firstname"`
	LastName     string `db:"lastname"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	RoleID       uint64 `db:"role_id"`
	Active       bool   `db:"active"`
	Timestamps
}

type UserModel struct {
	*AppModel[User]
}

func NewUserModel() *UserModel {
	return &UserModel{
		AppModel: NewAppModel[User]("users"),
	}
}

// FindByEmail finds a user by email
func (m *UserModel) FindByEmail(email string) (*User, error) {
	user, err := m.FindOneByField("email" , email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// Create inserts a new user
func (m *UserModel) Create(user *User) (int64, error) {
    // Set timestamps
    now := time.Now()
    user.CreatedAt = now
    user.UpdatedAt = now

    // Let the generic Insert handle all the fields
    return m.Insert(user)
}

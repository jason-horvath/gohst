package models

import (
	"fmt"
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
	*Model[User]
}

func NewUserModel() *UserModel {
	return &UserModel{
		Model: NewModel[User]("users"),
	}
}

// FindByEmail finds a user by email
func (m *UserModel) FindByEmail(email string) (*User, error) {
	user := &User{}
	query := fmt.Sprintf(/*sql*/ `SELECT * FROM %s WHERE email = $1`, m.GetTableName())

	user, err := m.FirstOf(query, email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// Create inserts a new user
func (m *UserModel) Create(user *User) error {
    // Set timestamps
    now := time.Now()
    user.CreatedAt = now
    user.UpdatedAt = now

    // Let the generic Insert handle all the fields
    return m.Insert(user)
}

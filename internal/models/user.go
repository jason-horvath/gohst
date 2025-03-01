package models

type User struct {
	Model
	FirstName    string `db:"firstname"`
	LastName     string `db:"lastname"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	RoleID       uint64 `db:"role_id"`
	Active       bool   `db:"active"`
}

type UserModel struct {
	*Model
}

func NewUserModel() *UserModel {
	return &UserModel{
		Model: NewModel("users"),
	}
}

// FindByEmail finds a user by email
func (m *UserModel) FindByEmail(email string) (*User, error) {
	user := &User{}
	query := "SELECT * FROM " + m.GetTableName() + " WHERE email = ?"
	err := m.GetDB().QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Create inserts a new user
func (m *UserModel) Create(user *User) error {
	query := `INSERT INTO users
		(firstname, lastname, email, password_hash, role_id, active)
		VALUES (?, ?, ?, ?, ?, ?)`

	result, err := m.GetDB().Exec(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.RoleID,
		user.Active,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint64(id)
	return nil
}

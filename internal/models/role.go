package models

type Role struct {
	Model
	Name        string `db:"name"`
	Description string `db:"description"`
}

type RoleModel struct {
	*Model
}

func NewRoleModel() *RoleModel {
	return &RoleModel{
		Model: NewModel("roles"),
	}
}

// FindByID finds a role by ID
func (m *RoleModel) FindByID(id uint64) (*Role, error) {
	role := &Role{}
	query := "SELECT * FROM " + m.GetTableName() + " WHERE id = ?"
	err := m.GetDB().QueryRow(query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// FindByName finds a role by name
func (m *RoleModel) FindByName(name string) (*Role, error) {
	role := &Role{}
	query := "SELECT * FROM " + m.GetTableName() + " WHERE name = ?"
	err := m.GetDB().QueryRow(query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return role, nil
}

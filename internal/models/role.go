package models

type Role struct {
	ID		  	uint64 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Timestamps
}

type RoleModel struct {
	*Model[Role]
}

func NewRoleModel() *RoleModel {
	return &RoleModel{
		Model: NewModel[Role]("roles"),
	}
}

// FindByID finds a role by ID
func (m *RoleModel) FindByID(id uint64) (*Role, error) {
	role := &Role{}
	query := "SELECT * FROM " + m.GetTableName() + " WHERE id = $1"
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
	query := "SELECT * FROM " + m.GetTableName() + " WHERE name = $1"
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

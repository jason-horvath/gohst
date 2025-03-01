package models

import (
	"database/sql"
	"gohst/internal/db"
	"time"
)

// Model is the base structure that all models inherit from.
type Model struct {
	ID        uint64     `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	db        *sql.DB
	tableName string
}

// NewModel initializes a new model instance for a given table.
func NewModel(tableName string) *Model {
	return &Model{
		db:        db.Database.DB,
		tableName: tableName,
	}
}

// GetDB returns the database connection
func (m *Model) GetDB() *sql.DB {
	return m.db
}

// GetTableName returns the table name
func (m *Model) GetTableName() string {
	return m.tableName
}

// WithTransaction handles database transactions.
func (m *Model) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// FindByID retrieves a record by ID.
func (m *Model) FindByID(id uint64, dest interface{}) error {
	query := "SELECT * FROM " + m.tableName + " WHERE id = ?"
	return m.db.QueryRow(query, id).Scan(dest)
}

// Delete removes a record by ID.
func (m *Model) Delete(id uint64) error {
	query := "DELETE FROM " + m.tableName + " WHERE id = ?"
	_, err := m.db.Exec(query, id)
	return err
}

// Count returns the total number of records.
func (m *Model) Count() (int, error) {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM " + m.tableName).Scan(&count)
	return count, err
}

// Exists checks if a record exists.
func (m *Model) Exists(id uint64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM " + m.tableName + " WHERE id = ?)"
	err := m.db.QueryRow(query, id).Scan(&exists)
	return exists, err
}

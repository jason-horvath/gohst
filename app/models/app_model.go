package models

import (
	"database/sql"
	"log"
	"time"

	"gohst/internal/models"
)

// AppModel provides basic app-specific functionality that all app models can inherit
type AppModel[T any] struct {
	*models.Model[T]
}

// Timestamps provides common timestamp fields for database models
type Timestamps struct {
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

// NewAppModel creates a new app model with shared app-level functionality
func NewAppModel[T any](tableName string) *AppModel[T] {
	return &AppModel[T]{
		Model: models.NewModel[T](tableName),
	}
}

// LogActivity logs model activity (app-specific logging)
func (a *AppModel[T]) LogActivity(action string, recordID uint64) {
	log.Printf("Model Activity: %s on table %s, record %d", action, a.GetTableName(), recordID)
}

// WithAppTransaction wraps database transactions with app-specific logging
func (a *AppModel[T]) WithAppTransaction(fn func() error) error {
	a.LogActivity("TRANSACTION_START", 0)

	err := a.WithTransaction(func(tx *sql.Tx) error {
		return fn()
	})

	if err != nil {
		a.LogActivity("TRANSACTION_FAILED", 0)
	} else {
		a.LogActivity("TRANSACTION_SUCCESS", 0)
	}

	return err
}

// ValidateAndInsert performs app-level validation before insert
func (a *AppModel[T]) ValidateAndInsert(record *T) error {
	// Add app-specific validation logic here
	a.LogActivity("INSERT_ATTEMPT", 0)

	err := a.Insert(record)
	if err != nil {
		a.LogActivity("INSERT_FAILED", 0)
		return err
	}

	a.LogActivity("INSERT_SUCCESS", 0)
	return nil
}



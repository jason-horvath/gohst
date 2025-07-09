package models

import (
	"time"
)

// ============================================================================
// SoftDeleteModel - For models that support soft deletes
// ============================================================================

// SoftDeleteModel extends AppModel with soft delete functionality
type SoftDeleteModel[T any] struct {
	*AppModel[T]
}

// NewSoftDeleteModel creates a new model with soft delete capabilities
func NewSoftDeleteModel[T any](tableName string) *SoftDeleteModel[T] {
	return &SoftDeleteModel[T]{
		AppModel: NewAppModel[T](tableName),
	}
}

// SoftDelete marks a record as deleted instead of actually deleting it
func (s *SoftDeleteModel[T]) SoftDelete(id uint64) error {
	query := "UPDATE " + s.GetTableName() + " SET deleted_at = $1 WHERE id = $2"
	_, err := s.GetDB().Exec(query, time.Now(), id)
	if err != nil {
		s.LogActivity("SOFT_DELETE_FAILED", id)
		return err
	}
	s.LogActivity("SOFT_DELETE_SUCCESS", id)
	return err
}

// FindActive returns only non-soft-deleted records
func (s *SoftDeleteModel[T]) FindActive(query string, args ...interface{}) ([]T, error) {
	// Modify query to exclude soft-deleted records
	modifiedQuery := query + " AND deleted_at IS NULL"
	return s.AllOf(modifiedQuery, args...)
}

// FindActiveByField returns active records by field
func (s *SoftDeleteModel[T]) FindActiveByField(fieldName string, value interface{}) ([]T, error) {
	query := "SELECT * FROM " + s.GetTableName() + " WHERE " + fieldName + " = $1 AND deleted_at IS NULL"
	return s.AllOf(query, value)
}

// FindActiveByID retrieves an active (non-deleted) record by ID
func (s *SoftDeleteModel[T]) FindActiveByID(id uint64) (*T, error) {
	query := "SELECT * FROM " + s.GetTableName() + " WHERE id = $1 AND deleted_at IS NULL"
	return s.FirstOf(query, id)
}

// Restore brings back a soft-deleted record
func (s *SoftDeleteModel[T]) Restore(id uint64) error {
	query := "UPDATE " + s.GetTableName() + " SET deleted_at = NULL WHERE id = $1"
	_, err := s.GetDB().Exec(query, id)
	if err != nil {
		s.LogActivity("RESTORE_FAILED", id)
		return err
	}
	s.LogActivity("RESTORE_SUCCESS", id)
	return nil
}

// RecentlyDeleted returns records deleted within the specified hours
func (s *SoftDeleteModel[T]) RecentlyDeleted(hours int) ([]T, error) {
	since := time.Now().Add(-time.Duration(hours) * time.Hour)
	query := "SELECT * FROM " + s.GetTableName() + " WHERE deleted_at > $1"
	return s.AllOf(query, since)
}

// IsDeleted checks if a specific record is soft-deleted
func (s *SoftDeleteModel[T]) IsDeleted(id uint64) (bool, error) {
	var deletedAt *time.Time
	query := "SELECT deleted_at FROM " + s.GetTableName() + " WHERE id = $1"
	err := s.GetDB().QueryRow(query, id).Scan(&deletedAt)
	if err != nil {
		return false, err
	}
	return deletedAt != nil, nil
}

// CountActive returns the count of active (non-deleted) records
func (s *SoftDeleteModel[T]) CountActive() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM " + s.GetTableName() + " WHERE deleted_at IS NULL"
	err := s.GetDB().QueryRow(query).Scan(&count)
	return count, err
}

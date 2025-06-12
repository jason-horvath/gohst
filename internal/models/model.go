package models

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"gohst/internal/db"
)

// Model is the base structure that all models inherit from.
type Model[T any] struct {
    db        *sql.DB
    tableName string
}

// NewModel initializes a new model instance for a given table.
func NewModel[T any](tableName string) *Model[T] {
	return &Model[T]{
		db:        db.Database.DB,
		tableName: tableName,
	}
}

// GetDB returns the database connection
func (m *Model[T]) GetDB() *sql.DB {
	return m.db
}

// GetTableName returns the table name
func (m *Model[T]) GetTableName() string {
	return m.tableName
}

// WithTransaction handles database transactions.
func (m *Model[T]) WithTransaction(fn func(*sql.Tx) error) error {
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
func (m *Model[T]) FindByID(id uint64) (*T, error) {
	query := "SELECT * FROM " + m.tableName + " WHERE id = $1"
	return m.FirstOf(query, id)
}

// FindByField returns records that match the specified field value.
// Example: FindByField("email", "user@example.com")
func (m *Model[T]) FindByField(fieldName string, value interface{}) ([]T, error) {
	queryTpl := /*sql*/ `SELECT * FROM %s WHERE %s = $1`
    query := fmt.Sprintf(queryTpl, m.tableName, fieldName)
    return m.AllOf(query, value)
}

// FindOneByField returns a single record that matches the specified field value.
// Example: FindOneByField("email", "user@example.com")
func (m *Model[T]) FindOneByField(fieldName string, value interface{}) (*T, error) {
	queryTpl := /*sql*/ `SELECT * FROM %s WHERE %s = $1 LIMIT 1`
    query := fmt.Sprintf(queryTpl, m.tableName, fieldName)
    return m.FirstOf(query, value)
}

// Delete removes a record by ID.
func (m *Model[T]) Delete(id uint64) error {
	query := "DELETE FROM " + m.tableName + " WHERE id = $1"
	_, err := m.db.Exec(query, id)
	return err
}

// Count returns the total number of records.
func (m *Model[T]) Count() (int, error) {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM " + m.tableName).Scan(&count)
	return count, err
}

// Exists checks if a record exists.
func (m *Model[T]) Exists(id uint64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM " + m.tableName + " WHERE id = $1)"
	err := m.db.QueryRow(query, id).Scan(&exists)
	return exists, err
}

// Recursive function to collect fields with db tags
func collectFieldsWithTags(v reflect.Value, t reflect.Type, values *[]interface{}) {
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)

        // Check if this is an embedded struct
        if field.Anonymous {
            // Recursively collect fields from the embedded struct
            collectFieldsWithTags(v.Field(i), field.Type, values)
            continue
        }

        // Process regular field
        dbTag := field.Tag.Get("db")
        if dbTag != "" && dbTag != "-" {
            *values = append(*values, v.Field(i).Addr().Interface())
        }
    }
}

// Query executes a query and scans a single result into a struct
func (m *Model[T]) First(dest interface{}, query string, args ...interface{}) error {
    // Validate the destination is a pointer to a struct
    v := reflect.ValueOf(dest)
    if v.Kind() != reflect.Ptr || v.IsNil() {
        return fmt.Errorf("destination must be a non-nil pointer to a struct")
    }

    // Dereference pointer to get the struct
    v = v.Elem()
    if v.Kind() != reflect.Struct {
        return fmt.Errorf("destination must point to a struct")
    }

    // Collect fields with db tags, including from embedded structs
    var values []interface{}
    collectFieldsWithTags(v, v.Type(), &values)

    if len(values) == 0 {
        return fmt.Errorf("struct has no fields with db tags")
    }

    // Execute query and scan
    return m.db.QueryRow(query, args...).Scan(values...)
}

// FirstOf returns a single record of type T
func (m *Model[T]) FirstOf(query string, args ...interface{}) (*T, error) {
    dest := new(T)
	log.Println("Destination type:", dest)
    err := m.First(dest, query, args...)
    if err != nil {
        return nil, err
    }
    return dest, nil
}

// QueryAll executes a query and scans multiple results into a slice of structs
func (m *Model[T]) All(dest interface{}, query string, args ...interface{}) error {
    // Validate destination is a pointer to a slice
    sliceValue := reflect.ValueOf(dest)
    if sliceValue.Kind() != reflect.Ptr || sliceValue.IsNil() {
        return fmt.Errorf("destination must be a non-nil pointer to a slice")
    }

    // Get the slice
    sliceValue = sliceValue.Elem()
    if sliceValue.Kind() != reflect.Slice {
        return fmt.Errorf("destination must point to a slice")
    }

    // Get the type of slice elements
    elemType := sliceValue.Type().Elem()
    if elemType.Kind() != reflect.Struct {
        return fmt.Errorf("slice elements must be structs")
    }

    // Execute the query
    rows, err := m.db.Query(query, args...)
    if err != nil {
        return err
    }
    defer rows.Close()

    // Process each row
    for rows.Next() {
        // Create a new struct instance
        newElem := reflect.New(elemType).Elem()

        // Collect field pointers for scanning
        var values []interface{}
        for i := 0; i < elemType.NumField(); i++ {
            field := elemType.Field(i)
            dbTag := field.Tag.Get("db")
            if dbTag != "" && dbTag != "-" {
                values = append(values, newElem.Field(i).Addr().Interface())
            }
        }

        // Scan the row into the new struct
        if err := rows.Scan(values...); err != nil {
            return err
        }

        // Append to the result slice
        sliceValue.Set(reflect.Append(sliceValue, newElem))
    }

    // Check for errors from iteration
    if err = rows.Err(); err != nil {
        return err
    }

    return nil
}

// AllOf returns a slice of records of type T
func (m *Model[T]) AllOf(query string, args ...interface{}) ([]T, error) {
    var dest []T
    err := m.All(&dest, query, args...)
    if err != nil {
        return nil, err
    }
    return dest, nil
}

// Add this to model.go
func (m *Model[T]) Insert(record *T) error {
    // Use reflection to get struct fields
    v := reflect.ValueOf(record).Elem()
    t := v.Type()

    var fields []string
    var placeholders []string
    var values []interface{}

    // Loop through fields with db tags
    paramCount := 1
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)

        // Handle embedded structs
        if field.Anonymous {
            embedType := field.Type
            embedValue := v.Field(i)

            for j := 0; j < embedType.NumField(); j++ {
                embedField := embedType.Field(j)
                dbTag := embedField.Tag.Get("db")

                // Skip ID field for inserts
                if dbTag != "" && dbTag != "-" && dbTag != "id" {
                    fields = append(fields, dbTag)
                    placeholders = append(placeholders, fmt.Sprintf("$%d", paramCount))
                    paramCount++
                    values = append(values, embedValue.Field(j).Interface())
                }
            }
            continue
        }

        // Regular fields
        dbTag := field.Tag.Get("db")
        if dbTag != "" && dbTag != "-" && dbTag != "id" {
            fields = append(fields, dbTag)
            placeholders = append(placeholders, fmt.Sprintf("$%d", paramCount))
            paramCount++
            values = append(values, v.Field(i).Interface())
        }
    }

    // Build query
    query := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s) RETURNING id",
        m.tableName,
        strings.Join(fields, ", "),
        strings.Join(placeholders, ", "),
    )

    // Execute query
    var id uint64
    err := m.db.QueryRow(query, values...).Scan(&id)
    if err != nil {
        return err
    }

    // Set ID in the struct
    for i := 0; i < t.NumField(); i++ {
        if t.Field(i).Tag.Get("db") == "id" {
            v.Field(i).SetUint(id)
            break
        }
    }

    return nil
}

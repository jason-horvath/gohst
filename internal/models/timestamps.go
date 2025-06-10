package models

import "time"

// Timestamps provides common timestamp fields for database models
type Timestamps struct {
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

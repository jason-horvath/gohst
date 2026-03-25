package storage

import (
	"context"
	"io"
	"time"
)

// FileStore is the unified interface for all storage backends.
// Switching from local → S3 → CDN is a config change, not a code change.
type FileStore interface {
	Store(ctx context.Context, path string, content io.Reader, opts StoreOptions) (StoredFile, error)
	Delete(ctx context.Context, path string) error
	URL(ctx context.Context, path string) (string, error)
}

// StoreOptions configures how a file is stored.
type StoreOptions struct {
	ContentType string
	MaxSize     int64
	Public      bool // if true, stored with public-read ACL (CDN-eligible)
}

// StoredFile is returned after a successful store operation.
type StoredFile struct {
	Path        string
	Size        int64
	ContentType string
	StoredAt    time.Time
}

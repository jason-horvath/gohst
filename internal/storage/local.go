package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalConfig configures the local filesystem store.
type LocalConfig struct {
	// PublicRoot is the directory for public files served directly by the web server.
	// e.g. "static" — files here are accessible at /static/<path>
	PublicRoot string
	// BaseURL is the application base URL, used to construct absolute URLs.
	// e.g. "http://localhost:3030"
	BaseURL string
}

// LocalFileStore stores files on the local filesystem.
// Public files go under PublicRoot and are served via the /static/ handler.
// Use this in development; swap for S3FileStore in production via STORAGE_DRIVER.
type LocalFileStore struct {
	cfg LocalConfig
}

func NewLocalFileStore(cfg LocalConfig) *LocalFileStore {
	return &LocalFileStore{cfg: cfg}
}

func (s *LocalFileStore) Store(ctx context.Context, path string, content io.Reader, opts StoreOptions) (StoredFile, error) {
	root := s.cfg.PublicRoot
	fullPath := filepath.Join(root, filepath.FromSlash(path))

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return StoredFile{}, fmt.Errorf("storage: failed to create directory: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return StoredFile{}, fmt.Errorf("storage: failed to create file: %w", err)
	}
	defer f.Close()

	n, err := io.Copy(f, content)
	if err != nil {
		return StoredFile{}, fmt.Errorf("storage: failed to write file: %w", err)
	}

	return StoredFile{
		Path:        path,
		Size:        n,
		ContentType: opts.ContentType,
		StoredAt:    time.Now(),
	}, nil
}

func (s *LocalFileStore) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.cfg.PublicRoot, filepath.FromSlash(path))
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("storage: failed to delete file: %w", err)
	}
	return nil
}

// URL returns the absolute URL for a locally stored public file.
// Path is appended after /static/ per the app's static file handler.
func (s *LocalFileStore) URL(ctx context.Context, path string) (string, error) {
	return s.cfg.BaseURL + "/static/" + path, nil
}

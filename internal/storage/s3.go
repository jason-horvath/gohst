package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Config configures the S3-compatible store.
// Works with AWS S3, Linode Object Storage, MinIO, and DigitalOcean Spaces.
type S3Config struct {
	Endpoint  string // custom endpoint for S3-compatible providers (e.g. Linode)
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
	Prefix    string // path prefix applied to all stored objects
	CDNURL    string // if set, URL() returns CDN URLs instead of direct S3 URLs
}

// S3FileStore stores files in an S3-compatible object store.
type S3FileStore struct {
	cfg    S3Config
	client *s3.Client
}

func NewS3FileStore(cfg S3Config) *S3FileStore {
	creds := credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")

	client := s3.New(s3.Options{
		BaseEndpoint:       aws.String(cfg.Endpoint),
		Region:             cfg.Region,
		Credentials:        creds,
		UsePathStyle:       true, // required for Linode and most S3-compatible providers
	})

	return &S3FileStore{cfg: cfg, client: client}
}

func (s *S3FileStore) Store(ctx context.Context, path string, content io.Reader, opts StoreOptions) (StoredFile, error) {
	key := s.cfg.Prefix + path

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.Bucket),
		Key:         aws.String(key),
		Body:        content,
		ContentType: aws.String(opts.ContentType),
	}

	if opts.Public {
		input.ACL = types.ObjectCannedACLPublicRead
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return StoredFile{}, fmt.Errorf("storage: S3 put failed: %w", err)
	}

	return StoredFile{
		Path:        path,
		ContentType: opts.ContentType,
		StoredAt:    time.Now(),
	}, nil
}

func (s *S3FileStore) Delete(ctx context.Context, path string) error {
	key := s.cfg.Prefix + path
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("storage: S3 delete failed: %w", err)
	}
	return nil
}

// URL returns the public URL for the stored object.
// If CDN_URL is configured, returns a CDN URL (Cloudflare edge URL).
// Otherwise falls back to the direct S3/Linode endpoint URL.
func (s *S3FileStore) URL(ctx context.Context, path string) (string, error) {
	key := s.cfg.Prefix + path
	if s.cfg.CDNURL != "" {
		return s.cfg.CDNURL + "/" + key, nil
	}
	// Direct Linode / S3-compatible URL
	return fmt.Sprintf("%s/%s/%s", s.cfg.Endpoint, s.cfg.Bucket, key), nil
}

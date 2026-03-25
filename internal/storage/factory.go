package storage

import "gohst/internal/config"

// NewFileStore creates the appropriate FileStore based on the STORAGE_DRIVER
// environment variable. Switching between backends is a config change only.
//
//	STORAGE_DRIVER=local  → LocalFileStore  (development)
//	STORAGE_DRIVER=s3     → S3FileStore     (Linode / AWS / production)
func NewFileStore() FileStore {
	driver := config.GetEnv("STORAGE_DRIVER", "local").(string)

	switch driver {
	case "s3":
		return NewS3FileStore(S3Config{
			Endpoint:  config.GetEnv("S3_ENDPOINT", "").(string),
			Bucket:    config.GetEnv("S3_BUCKET", "").(string),
			Region:    config.GetEnv("S3_REGION", "us-east-1").(string),
			AccessKey: config.GetEnv("S3_ACCESS_KEY", "").(string),
			SecretKey: config.GetEnv("S3_SECRET_KEY", "").(string),
			Prefix:    config.GetEnv("S3_PREFIX", "").(string),
			CDNURL:    config.GetEnv("CDN_URL", "").(string),
		})
	default:
		return NewLocalFileStore(LocalConfig{
			PublicRoot: config.GetEnv("STORAGE_PUBLIC_ROOT", "static").(string),
			BaseURL:    config.GetEnv("APP_URL", "http://localhost:3030").(string),
		})
	}
}

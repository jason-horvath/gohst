# Application Configuration

This directory contains application-specific configuration that extends the framework's core configuration.

## Structure

- `app.go` - Main application configuration with feature flags, pagination, uploads, etc.

## Usage

### In Controllers

```go
import appConfig "gohst/app/config"

func (c *SomeController) SomeAction(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "AppName":       appConfig.App.Name,
        "AppVersion":    appConfig.App.Version,
        "IsProduction":  appConfig.IsProduction(),
        "Features":      appConfig.App.Features,
        "MaxUploadSize": appConfig.App.Upload.MaxFileSize,
    }

    c.Render(w, r, "template", data)
}
```

### In Templates

```html
<h1>{{.AppName}} v{{.AppVersion}}</h1>

{{if .Features.EnableRegistration}}
<a href="/auth/register">Register</a>
{{end}} {{if .Features.MaintenanceMode}}
<div class="alert">Site is under maintenance</div>
{{end}}
```

## Configuration

All configuration values can be set via environment variables:

```bash
# Application settings
APP_NAME="My Gohst App"
APP_VERSION="1.0.0"
APP_ENV="production"
APP_DEBUG=false

# Feature flags
FEATURE_REGISTRATION=true
FEATURE_USER_PROFILES=true
FEATURE_NOTIFICATIONS=false
MAINTENANCE_MODE=false

# Pagination
PAGINATION_DEFAULT_LIMIT=20
PAGINATION_MAX_LIMIT=100

# Upload settings
UPLOAD_MAX_FILE_SIZE=10485760  # 10MB
UPLOAD_PATH="static/uploads"
```

## Adding New Configuration

To add new configuration:

1. Add the field to the appropriate struct in `app.go`
2. Add the environment variable handling in `Initialize()`
3. Add the environment variable to `.env.example`
4. Update this README with usage examples

## Framework vs App Config

- **Framework config** (`internal/config/`) - Core functionality (database, sessions, etc.)
- **App config** (`app/config/`) - Application-specific settings (features, business logic, etc.)

Both are initialized automatically when the application starts.

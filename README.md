# Gohst 👻

A modern Go web framework with clean architecture and powerful development tools. Gohst provides a solid foundation for building scalable web applications with a clear separation between framework and application code.

## Architecture

Gohst follows a **layered architecture** with clean separation of concerns:

- **Framework Layer** (`internal/`) - Core framework functionality that can be reused across projects
- **Application Layer** (`app/`) - Your specific business logic and application configuration
- **Interface-based Dependency Injection** - Framework depends on app through well-defined interfaces

This design allows the framework to evolve independently while keeping your application code clean and testable.

## Features

- 🏗️ **Clean Architecture** - Separation between framework and application layers
- � **Interface-based DI** - Framework depends on app through interfaces, not direct coupling
- �🚀 **Hot-reloading** - Go server using Air for rapid development
- 🎨 **Modern Frontend** - Vite for asset building with TypeScript, Tailwind CSS
- 🐳 **Docker Development** - Postgres and PgAdmin in containers
- 📦 **Flexible Sessions** - File or Redis-based session storage
- 🗄️ **Advanced Database** - Multi-database support with generic models and relationships
- 🔄 **Robust Migrations** - Database migrations and seeding with batch tracking
- 🛠️ **Template System** - HTML rendering with layouts, partials, and custom functions
- ⚙️ **Rich Configuration** - Environment-based config with feature flags and validation
- 🔐 **Authentication** - Built-in auth with role-based permissions
- 📝 **Form Handling** - Type-safe forms with validation and error handling

## Directory Structure

```
gohst/
├── app/                        # 🏢 APPLICATION LAYER
│   ├── config/                 # App-specific configuration
│   │   ├── app.go             # Feature flags, pagination, uploads
│   │   ├── db.go              # Database connections setup
│   │   └── README.md          # Configuration documentation
│   ├── controllers/           # Application controllers
│   ├── helpers/               # App-specific template functions
│   ├── models/                # Business domain models
│   ├── routes/                # Application route definitions
│   └── services/              # Business logic services
├── internal/                   # 🔧 FRAMEWORK LAYER
│   ├── auth/                  # Authentication framework
│   ├── config/                # Framework configuration & interfaces
│   ├── controllers/           # Base controller functionality
│   ├── db/                    # Database connection management
│   ├── forms/                 # Form handling and validation
│   ├── middleware/            # HTTP middleware (auth, CSRF, logging)
│   ├── migration/             # Database migration engine
│   ├── models/                # Generic model base with relationships
│   ├── render/                # Template rendering and asset management
│   ├── routes/                # Route registration and handling
│   ├── session/               # Session management (file/Redis)
│   ├── utils/                 # Framework utilities
│   └── validation/            # Input validation framework
├── cmd/                        # 🚀 COMMANDS
│   ├── migrate/               # Database migration CLI
│   ├── web/                   # Main web application
│   └── dev/                   # Development tools
│       ├── gohst_server       # Development server control
│       ├── docker_sql_build   # Database setup
│       └── docker_sql_clear   # Database cleanup
├── database/                   # 📊 DATABASE
│   ├── migrations/            # SQL migration files
│   └── seeds/                 # SQL seed files
├── templates/                  # 🎨 TEMPLATES
│   ├── layouts/               # Layout templates
│   ├── components/            # Reusable components
│   ├── partials/              # Partial templates
│   └── views/                 # Page templates
├── assets/                     # 🎨 FRONTEND ASSETS
│   ├── css/                   # Stylesheet sources
│   ├── js/                    # JavaScript/TypeScript sources
│   └── icons/                 # SVG icons
├── static/                     # 📁 STATIC FILES
│   ├── dist/                  # Compiled frontend assets
│   ├── images/                # Static images
│   └── uploads/               # User uploads
├── docker/                     # 🐳 DOCKER
│   ├── postgres/              # Postgres container setup
│   └── pgadmin/               # PgAdmin configuration
├── tmp/                        # 🗂️ TEMPORARY
│   ├── sessions/              # File-based sessions
│   └── build-errors.log       # Development logs
├── .air.toml                   # Air hot-reload configuration
├── .env                        # Environment variables
├── docker-compose.yml          # Docker services
├── gohst                       # 👻 CLI tool
├── go.mod                      # Go module
├── package.json               # Frontend dependencies
├── tailwind.config.js         # Tailwind CSS configuration
└── vite.config.js             # Vite build configuration
```

## Getting Started

### 1. Clone and Setup

```bash
git clone https://github.com/jason-horvath/gohst.git
cd gohst
cp .env.example .env
```

### 2. Configure Your Application

Edit `.env` with your specific settings:

```bash
APP_NAME="My CRM App"
DB_NAME=my_crm_db
DB_USER=my_user
DB_PASSWORD=my_password
```

### 3. Build Development Environment

```bash
./gohst build
```

This will:

- Start Docker containers (Postgres, PgAdmin)
- Install frontend dependencies
- Set up database connections
- Run initial migrations and seeds

### 4. Start Development

```bash
./gohst dev
```

Your application will be available at:

- **App**: http://localhost:3030
- **PgAdmin**: http://localhost:5050 (gohst@gohst.dev / password)

### 5. Create Your First Feature

```bash
# Create a migration
./gohst migrate:create create_companies_table

# Create app models
# app/models/company.go

# Create controllers
# app/controllers/company_controller.go

# Add routes
# app/routes/routes.go
```

## Architecture Benefits

### For Framework Development

- **Clean separation** - Framework code isolated in `internal/`
- **Interface-driven** - Framework depends on app through interfaces
- **Testable** - Each layer can be tested independently
- **Reusable** - Framework can be extracted for other projects

### For Application Development

- **Focused business logic** - App code stays in `app/`
- **Configuration flexibility** - Environment-driven config with sensible defaults
- **Type safety** - Generic models and form validation
- **Developer experience** - Hot reloading, error handling, debugging tools

## CLI Commands

The `gohst` CLI tool is your primary interface for development, database management, and deployment tasks.

```bash
./gohst <command> [arguments]
```

### Development Environment

- `build` - Build the complete development environment (Docker, NPM, Migrations)
- `dev` - Start the development environment (Docker, Vite, Air)
- `dev:down` - Stop the development environment
- `destroy` - Destroy the environment (removes containers, volumes, node_modules)
- `server:start` - Start the Go server (Air) independently
- `server:stop` - Stop the Go server
- `docker:rebuild` - Rebuild Docker containers
- `docker:sql:build` - Set up database initialization files
- `docker:sql:clear` - Clear database initialization files
- `storage:link` - Link storage assets to the static directory

### Database Migrations

- `migrate:run` - Run all pending migrations
- `migrate:status` - Show migration status
- `migrate:rollback` - Rollback the last batch of migrations
- `migrate:create <name>` - Create a new migration file
- `migrate:full` - Run migrations and seeds together
- `migrate:fresh` - Drop all tables and re-run all migrations
- `migrate:fresh:full` - Drop all tables, re-run migrations, and run seeds

### Database Seeding

- `migrate:seed` - Run all pending seeds
- `migrate:seed:status` - Show seed status
- `migrate:seed:fresh` - Clear all seed records and re-run all seeds
- `migrate:seed:rollback` - Rollback the last batch of seeds
- `migrate:seed:create <name>` - Create a new seed file

### Examples

```bash
# Start working
./gohst dev

# Create a new migration
./gohst migrate:create create_users_table

# Run migrations
./gohst migrate:run

# Reset database completely
./gohst migrate:fresh:full
```

## Optional CLI Alias

Add an alias to your shell configuration:

```bash
# For zsh
echo 'alias gohst="./gohst"' >> ~/.zshrc
source ~/.zshrc

# For bash
echo 'alias gohst="./gohst"' >> ~/.bashrc
source ~/.bashrc
```

Once the alias is added only `gohst <command>` is needed to run commands.

## Development Workflow

### Day-to-Day Development

1. **Start the environment:**

   ```bash
   ./gohst dev
   ```

2. **Make changes** - Air automatically reloads Go code, Vite handles frontend assets

3. **Create database changes:**

   ```bash
   ./gohst migrate:create add_user_preferences
   # Edit the generated migration file
   ./gohst migrate:run
   ```

4. **Add new features:**

   - Create models in `app/models/`
   - Add controllers in `app/controllers/`
   - Define routes in `app/routes/`
   - Create templates in `templates/`

5. **Update configuration:**
   - Modify `app/config/app.go` for business settings
   - Update `.env` for environment-specific values

### Framework Improvements

When building your app, you may need to enhance the framework:

1. **Identify the need** - Missing validation, model feature, etc.
2. **Implement in `internal/`** - Keep framework code separate
3. **Test with your app** - Ensure it works for your use case
4. **Consider extraction** - Framework improvements can benefit other projects

### Production Deployment

```bash
# Build for production
NODE_ENV=production npm run build
CGO_ENABLED=0 GOOS=linux go build -o app cmd/web/main.go

# Set production environment
APP_ENV_KEY=production
SESSION_STORE=redis
DB_HOST=production-db.example.com

# Run migrations
./gohst migrate:run
```

## Framework Patterns
```

````

## Framework Patterns

### Generic Models with Relationships

```go
// Framework provides generic base models
type Company struct {
    gohst.AppModel[Company]
    Name     string `db:"name"`
    Industry string `db:"industry"`
}

// Automatic CRUD operations
company := &models.Company{}
err := company.Create(map[string]interface{}{
    "name":     "Acme Corp",
    "industry": "Technology",
})

// Built-in soft deletes
err = company.SoftDelete(companyID)
````

### Multi-Database Support

```go
// Configure multiple databases
func CreateDBConfigs() *config.DatabaseConfigPool {
    pool := config.NewDatabaseConfigPool()

    // Primary database
    pool.AddDatabase("primary", config.DatabaseConfig{
        Host: "localhost",
        Port: 5432,
        User: "app",
        Password: "secret",
        DBName: "myapp",
        SSLMode: "disable",
    })

    // Analytics database
    pool.AddDatabase("analytics", config.DatabaseConfig{
        Host: "analytics.example.com",
        Port: 5432,
        User: "readonly",
        Password: "secret",
        DBName: "analytics",
        SSLMode: "require",
    })

    return pool
}

// Use specific database
analyticsDB := db.GetDB("analytics")
```

### Template Functions and Helpers

```go
// App-specific template functions
func AppTemplateFuncs() template.FuncMap {
    return template.FuncMap{
        "formatCurrency": func(amount float64) string {
            return fmt.Sprintf("$%.2f", amount)
        },
        "userCan": func(user User, permission string) bool {
            return user.HasPermission(permission)
        },
        "appVersion": func() string {
            return appConfig.App.Version
        },
    }
}

// Register with framework
render.RegisterTemplateFuncs(appHelpers.AppTemplateFuncs())
```

### Form Handling and Validation

```go
// Type-safe form handling
type CompanyForm struct {
    Name     string `form:"name" validate:"required,min=2"`
    Industry string `form:"industry" validate:"required"`
    Website  string `form:"website" validate:"url"`
}

func (c *CompanyController) Create(w http.ResponseWriter, r *http.Request) {
    form := &CompanyForm{}

    if err := forms.ParseAndValidate(r, form); err != nil {
        // Handle validation errors
        c.RenderWithErrors(w, r, "companies/new", form, err)
        return
    }

    // Form is valid, proceed with business logic
}
```

## Session Management

Gohst provides two session storage options:

### 1. File-Based Sessions

- Sessions are stored as `.session` files in the `tmp/sessions/` directory
- Uses Go's `gob` encoding format for serialization
- Each session is stored in a separate file with the pattern `{session-id}.session`
- Good for development and simple applications
- Configure in `.env`:

```bash
SESSION_STORE=file
SESSION_FILE_PATH=tmp/sessions
```

### 2. Redis Sessions

- Sessions are stored in Redis for better performance and scalability
- Recommended for production environments
- Configure in `.env`:

```bash
SESSION_STORE=redis
SESSION_REDIS_HOST=localhost
SESSION_REDIS_PORT=6379
SESSION_REDIS_PASSWORD=
SESSION_REDIS_DB=0
```

Common session configuration:

```bash
SESSION_NAME=session_id        # Cookie name
SESSION_LENGTH=60             # Session duration in minutes
SESSION_CONTEXT_KEY=session   # Context key for session data
```

## Configuration

Gohst uses a **layered configuration approach** with framework and application layers:

### Framework Configuration

The framework provides core functionality configuration:

- Database connections and pooling
- Session management (file/Redis)
- Template rendering and assets
- Middleware and routing
- Migration system

### Application Configuration

Your application extends the framework with business-specific config:

```go
// app/config/app.go
type AppConfig struct {
    Name        string
    Version     string
    Environment string
    Debug       bool

    // Feature flags
    Features FeatureFlags

    // Business settings
    Pagination PaginationConfig
    Upload     UploadConfig
}
```

### Environment Variables

All configuration can be controlled via environment variables:

```bash
# Application
APP_NAME="My Gohst App"
APP_VERSION="1.0.0"
APP_ENV_KEY="production"
APP_DEBUG=false
APP_URL="https://myapp.com"
APP_PORT=3030

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=myapp
DB_PASSWORD=secret
DB_NAME=myapp_production
DB_SSL_MODE=disable              # disable, require, verify-ca, verify-full

# Session Management
SESSION_STORE=redis              # or 'file'
SESSION_NAME=session_id
SESSION_LENGTH=60               # minutes
SESSION_REDIS_HOST=localhost
SESSION_REDIS_PORT=6379

# Feature Flags
FEATURE_REGISTRATION=true
FEATURE_USER_PROFILES=true
FEATURE_NOTIFICATIONS=false
MAINTENANCE_MODE=false

# Business Settings
PAGINATION_DEFAULT_LIMIT=20
PAGINATION_MAX_LIMIT=100
UPLOAD_MAX_FILE_SIZE=10485760   # 10MB in bytes
UPLOAD_PATH="static/uploads"

# Development
VITE_MANIFEST_PATH=static/dist/.vite/manifest.json
VITE_PORT=5174
```

### Interface-Based Dependency Injection

The framework depends on your application through clean interfaces:

```go
// Framework defines what it needs
type AppConfigProvider interface {
    GetURL() string
    GetDistPath() string
    IsProduction() bool
    IsDevelopment() bool
}

// Your app implements the interface
func (ac *AppConfig) IsProduction() bool {
    return ac.Environment == "production"
}

// Framework uses the interface
app := config.GetAppConfig()
if app.IsProduction() {
    // Production-specific logic
}
```

```

```

## What's New

### Recent Framework Improvements

- **🏗️ App Layer Architecture** - Clean separation between framework (`internal/`) and application (`app/`) code
- **🔌 Interface-Based DI** - Framework depends on app through interfaces, enabling true decoupling
- **🗄️ Multi-Database Support** - Configure multiple database connections with connection pooling
- **📊 Advanced Migrations** - Batch tracking, rollbacks, and seeding with comprehensive CLI
- **🎨 Generic Models** - Type-safe models with built-in CRUD operations and relationship support
- **📝 Enhanced Forms** - Type-safe form parsing with validation and error handling
- **⚙️ Rich Configuration** - Feature flags, pagination settings, upload controls, and environment-based config
- **🛠️ Template Functions** - App-specific template helpers with framework registration
- **🔐 Improved Auth** - Role-based authentication with session management
- **🚀 Developer Experience** - Better debugging, error handling, and development workflow

### Framework vs Application Code

The framework now clearly distinguishes between:

**Framework Code** (`internal/`):

- Reusable across projects
- Database connections and models
- HTTP routing and middleware
- Template rendering and assets
- Session management
- Form handling and validation

**Application Code** (`app/`):

- Business-specific logic
- Domain models and controllers
- Application configuration
- Custom template functions
- Business services and helpers

This architecture makes Gohst suitable for:

- **Rapid prototyping** - Start building immediately with solid foundations
- **Production applications** - Scalable architecture with clean separation
- **Framework development** - Easily extract and reuse framework components
- **Team development** - Clear boundaries between framework and application concerns

## Creating New Projects

### Using Gohst as a Template

To create a new project using Gohst as your framework:

```bash
# Clone with fresh git history
git clone https://github.com/jason-horvath/gohst.git my-new-project
cd my-new-project
rm -rf .git
git init
git add .
git commit -m "Initial commit: New project using Gohst framework"

# Configure for your project
# 1. Update .env with your settings
# 2. Modify app/config/app.go with your business logic
# 3. Clear out example code and start building
```

### Framework Evolution

As you build your application, you'll likely improve the framework:

```bash
# Make framework improvements (in internal/)
git add internal/
git commit -m "framework: add file upload utilities"

# Make application features (in app/)
git add app/ templates/ static/
git commit -m "app: add customer management"
```

Later, you can extract framework improvements to benefit other projects:

```bash
# Extract framework commits
git log --oneline --grep="framework:"

# Port improvements back to main framework repo
git format-patch HEAD~5 --grep="framework:" --stdout > framework-improvements.patch
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Author

Jason Horvath - [jason.horvath@larzilla.com](mailto:jason.horvath@larzilla.com)

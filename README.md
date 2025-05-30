# Gohst 👻

A Go web application starter kit with built-in development tools and modern frontend capabilities. Gohst is designed to provide a solid foundation for building web applications with Go and Vite. The project includes a development environment with Docker, Air, and MySQL, as well as a CLI tool for managing the development workflow. Gohst also includes session management, database migrations, and HTML template rendering.

The project is structured to provide a clean separation of concerns and a clear path for extending functionality. Gohst is a great starting point for building web applications with Go and modern frontend tools like Tailwind CSS, Alpine.js, Typescript, and Vite.

NOTE: The project is ongoing and will continue to evolve with new features and improvements as needed. Future updates will include additional basic authentication, frontend tools, more configuration options, and enhanced development workflows.

## Features

- 🚀 Hot-reloading using gohst bash scripts
- 🎨 Vite for frontend assets
- 🐳 Docker-based MySQL development environment
- 📦 Session management (File/Redis support)
- 🔄 Database migrations and seeding
- 🛠️ HTML template rendering with layouts and partials
- 🔧 Environment-based configuration

## Directory Structure

```
gohst/
├── assets/                     # Build assets
│   ├── css/                    # Entry point for build css
│   ├── js/                     # Entry point for build js
│   ├── storage/                # Stroage for media/user assets
│   │   ├── images              # Images for builkd
│   │   └── uploads             # User generated content
├── cmd/
│   ├── dev/                    # Development scripts
│   │   ├── docker_sql_build    # Build docker sql
│   │   ├── docker_sql_clear    # Clear docker sql
│   │   ├── gohst_server        # Ghost hot-reload control
│   │   ├── storage_mgr         # Storage linking and management
│   │   └── vite_process        # Handles vite process stop/start
│   └── web/                    # Main application
├── database/
│   ├── migrations/             # SQL migrations
│   └── seeds/                  # SQL seed files
├── docker/
│   └── mysql/                  # MySQL container setup
├── internal/
│   ├── config/                 # Application configuration
│   ├── controllers/            # HTTP request handlers
│   ├── db/                     # Database connection
│   ├── middleware/             # HTTP middleware
│   ├── render/                 # Template rendering
│   ├── routes/                 # Route definitions
│   └── session/                # Session management
├── static/                     # Static assets
│   └── dist/                   # Compiled assets
├── templates/                  # HTML templates
│   ├── layouts/                # Layout templates
│   ├── pages/                  # Page templates
│   └── partials/               # Partial templates
│── tmp                         # Tmp files and logs
├── .env                        # Environment variables
├── .env.example                # Environment template
├── docker-compose.yml          # Docker services
├── gohst                       # CLI tool
└── package.json                # Frontend dependencies
```

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/jason-horvath/gohst.git
cd gohst
```

2. Copy environment file:

```bash
cp .env.example .env
```

3. Build the development environment:

```bash
./gohst build
```

## CLI Commands

```bash
./gohst <command>
```

Available commands:

- `build` - Set up complete development environment
- `dev` - Bring up the development environment: npm, docker, vite, gohst
- `dev:down` - Bring down the development environment: npm, docker, vite, gohst
- `server:start` - Start the gohst server
- `server:stop` - Stop the ghost server
- `docker:sql:build` - Set up database files
- `docker:sql:clear` - Clear database files
- `docker:rebuild` - Rebuild Docker containers
- `storage:link` - Link assets to the static directory

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

1. Start the environment:

```bash
./gohst up
```

2. Make changes to Go files - Air will automatically reload
3. Edit frontend assets - Vite will handle hot reloading
4. Add database migrations in `database/migrations`
5. Add database seeds in `database/seeds`

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

## Environment Configuration

Key environment variables:

- `APP_ENV_KEY` - Application environment
- `APP_URL` - Application URL
- `APP_PORT` - Server port
- `DB_*` - Database configuration
- `SESSION_*` - Session configuration
- `VITE_*` - Frontend build configuration

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

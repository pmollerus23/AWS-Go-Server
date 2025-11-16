# Project Structure

This document explains the production-grade folder structure of the AWS Go Server.

## Directory Layout

```
AWS-Go-Server/
├── cmd/                        # Application entry points
│   └── server/                 # Main server application
│       └── main.go            # Application entry point
│
├── internal/                   # Private application code (cannot be imported by other projects)
│   ├── aws/                   # AWS-specific code
│   │   ├── client.go         # AWS client initialization
│   │   └── auth.go           # IAM authentication middleware
│   │
│   ├── config/                # Configuration management
│   │   └── config.go         # Configuration structs and loading
│   │
│   ├── handlers/              # HTTP request handlers
│   │   ├── health.go         # Health check handler
│   │   ├── items.go          # Item CRUD handlers
│   │   └── aws.go            # AWS service handlers
│   │
│   ├── middleware/            # HTTP middleware
│   │   ├── logging.go        # Request logging
│   │   ├── recovery.go       # Panic recovery
│   │   └── sizelimit.go      # Request size limiting
│   │
│   ├── models/                # Domain models (empty for now, ready for future use)
│   │
│   └── server/                # HTTP server setup
│       ├── server.go         # Server initialization and lifecycle
│       └── routes.go         # Route definitions
│
├── pkg/                        # Public libraries (can be imported by other projects)
│                               # Currently empty, add reusable packages here
│
├── api/                        # API definitions
│   └── openapi/               # OpenAPI/Swagger specifications
│
├── deployments/               # Deployment configurations
│   └── docker/
│       ├── Dockerfile        # Multi-stage Docker build
│       └── docker-compose.yml # Docker Compose configuration
│
├── scripts/                   # Build and deployment scripts
│
├── .github/                   # GitHub-specific files
│   └── workflows/            # CI/CD workflows (ready for GitHub Actions)
│
├── Makefile                   # Common development tasks
├── .air.toml                  # Hot reload configuration
├── .gitignore                 # Git ignore rules
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
│
└── Documentation files
    ├── README.md
    ├── PROJECT_STRUCTURE.md
    ├── AWS_INTEGRATION.md
    ├── QUICKSTART_AWS.md
    └── RESILIENCE_IMPROVEMENTS.md
```

## Why This Structure?

### `cmd/`
**Purpose**: Application entry points

- Each subdirectory is a separate binary (e.g., `cmd/server`, `cmd/worker`, `cmd/cli`)
- Keeps main packages small and focused
- Easy to build multiple executables from one repository
- Example: `go build ./cmd/server`

### `internal/`
**Purpose**: Private application code

- **Cannot be imported** by other projects (Go enforces this)
- Contains business logic specific to this application
- Protects internal APIs from external use
- Safe place for rapid iteration without breaking external consumers

#### `internal/aws/`
AWS-specific functionality:
- Client initialization with credential management
- IAM authentication middleware
- AWS service integrations

#### `internal/config/`
Configuration management:
- Environment variable loading
- Configuration validation
- Default values
- Type-safe config structures

#### `internal/handlers/`
HTTP request handlers:
- One file per resource or domain (items, health, aws)
- Handles HTTP request/response logic
- Thin layer - delegates to services for business logic
- Returns appropriate HTTP status codes

#### `internal/middleware/`
HTTP middleware:
- Logging, authentication, recovery, etc.
- Reusable across all routes
- Applied in `server/server.go`

#### `internal/models/`
Domain models:
- Business entities (User, Product, Order, etc.)
- Shared across handlers, services, repositories
- Plain Go structs with validation methods

#### `internal/server/`
HTTP server setup:
- Server initialization and configuration
- Route registration
- Middleware application
- Graceful shutdown handling

### `pkg/`
**Purpose**: Public libraries

- **Can be imported** by other projects
- Reusable, well-tested code
- No dependencies on `internal/`
- Examples: utility functions, common types, shared algorithms

### `api/`
**Purpose**: API definitions

- OpenAPI/Swagger specifications
- Protocol buffer definitions
- GraphQL schemas
- API documentation

### `deployments/`
**Purpose**: Deployment configurations

- Docker files
- Kubernetes manifests
- Terraform configurations
- Helm charts

### `scripts/`
**Purpose**: Development and deployment scripts

- Build scripts
- Database migration scripts
- Test data generation
- Deployment automation

## Design Principles

### 1. **Separation of Concerns**
Each package has a single responsibility:
- `handlers`: HTTP I/O
- `services`: Business logic (future)
- `repositories`: Data access (future)
- `models`: Data structures

### 2. **Dependency Direction**
```
main → server → handlers → services → repositories → models
                              ↓
                             aws
```
- Dependencies flow inward (toward models)
- Inner layers don't know about outer layers
- Easy to test in isolation

### 3. **Package Organization**
- Group by **domain** or **feature**, not by type
- Related code lives together
- Easy to find and modify

### 4. **Testability**
- Small, focused packages
- Dependencies injected (logger, config, AWS clients)
- Easy to mock and test

## Common Patterns

### Adding a New Feature

1. **Define the model** (if needed)
   ```go
   // internal/models/user.go
   type User struct {
       ID    int64
       Name  string
       Email string
   }
   ```

2. **Create the handler**
   ```go
   // internal/handlers/users.go
   func HandleUsersGet(logger *slog.Logger) http.Handler {
       // ...
   }
   ```

3. **Register the route**
   ```go
   // internal/server/routes.go
   mux.Handle("GET /api/v1/users", handlers.HandleUsersGet(s.logger))
   ```

4. **Add tests**
   ```go
   // internal/handlers/users_test.go
   func TestHandleUsersGet(t *testing.T) {
       // ...
   }
   ```

### Adding a New Service Layer

When business logic grows:

1. Create `internal/service/users.go`:
   ```go
   type UserService struct {
       logger *slog.Logger
       repo   *repository.UserRepository
   }

   func (s *UserService) GetUser(id int64) (*models.User, error) {
       // Business logic here
   }
   ```

2. Update handlers to use services:
   ```go
   func HandleUsersGet(logger *slog.Logger, svc *service.UserService) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           user, err := svc.GetUser(id)
           // ...
       })
   }
   ```

### Adding a Database Layer

When you need persistence:

1. Create `internal/repository/users.go`:
   ```go
   type UserRepository struct {
       db *sql.DB
   }

   func (r *UserRepository) FindByID(id int64) (*models.User, error) {
       // Database query here
   }
   ```

2. Services use repositories
3. Handlers use services
4. Clean separation of concerns

## Development Workflow

### Local Development
```bash
# Install dev tools
make install-tools

# Run with hot reload
make dev

# Or run directly
make run
```

### Testing
```bash
# Run tests
make test

# With coverage
make test-coverage
```

### Docker
```bash
# Build image
make docker-build

# Run with docker-compose
make docker-up

# View logs
make docker-logs

# Stop
make docker-down
```

### Deployment
```bash
# Build production binary
make build

# The binary is in bin/server
./bin/server
```

## Migration from Old Structure

The old structure:
```
├── main.go
├── server.go
├── routes.go
├── middleware.go
├── aws.go
├── iam_auth.go
├── helpers.go
└── handlers/
    ├── health.go
    ├── items.go
    └── aws.go
```

Has been reorganized into:
```
├── cmd/server/main.go              (was: main.go + parts of server.go)
├── internal/
│   ├── server/
│   │   ├── server.go               (was: server.go)
│   │   └── routes.go               (was: routes.go)
│   ├── middleware/                 (was: middleware.go, split by concern)
│   │   ├── logging.go
│   │   ├── recovery.go
│   │   └── sizelimit.go
│   ├── aws/                        (was: aws.go, iam_auth.go)
│   │   ├── client.go
│   │   └── auth.go
│   ├── config/                     (new: extracted from server.go)
│   │   └── config.go
│   └── handlers/                   (was: handlers/)
│       ├── health.go
│       ├── items.go
│       └── aws.go
```

### Benefits of New Structure

1. **Scalability**: Easy to add new features without cluttering
2. **Testability**: Each package can be tested independently
3. **Team Collaboration**: Clear boundaries reduce merge conflicts
4. **Maintainability**: Easy to find and modify code
5. **Reusability**: `pkg/` can be imported by other projects
6. **Production-Ready**: Follows industry best practices
7. **Docker-Ready**: Dockerfile optimized for production
8. **CI/CD-Ready**: Structure supports automated pipelines

## Next Steps

1. **Add service layer** when business logic grows
2. **Add repository layer** when adding database
3. **Add OpenAPI spec** in `api/openapi/`
4. **Add CI/CD** in `.github/workflows/`
5. **Add database migrations** in `migrations/`
6. **Add monitoring** (Prometheus metrics, health checks)
7. **Add distributed tracing** (OpenTelemetry)

## References

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

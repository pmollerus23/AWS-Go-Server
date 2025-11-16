# AWS Go Server

A production-grade Go web server with AWS integration, featuring resilient error handling, IAM authentication, and clean, scalable architecture following Go best practices and Mat Ryer's battle-tested patterns.

## Features

- ✅ **Production-Ready Structure** - Following [Go standard project layout](https://github.com/golang-standards/project-layout)
- ✅ **AWS Integration** - S3, DynamoDB with IAM authentication using AWS SDK v2
- ✅ **Resilient** - Panic recovery, graceful shutdown, request timeouts, race condition protection
- ✅ **Secure** - IAM auth, request size limits, input validation, non-root Docker user
- ✅ **Observable** - Structured logging with slog, comprehensive error tracking
- ✅ **Docker Ready** - Multi-stage builds, docker-compose, health checks
- ✅ **Developer Friendly** - Makefile, hot reload, extensive documentation

## Quick Start

### Prerequisites

- Go 1.21+
- Docker (optional)
- AWS credentials configured (for AWS features)

### Run Locally

```bash
# Install dependencies
go mod download

# Run the server
make run

# Or with hot reload (installs air if needed)
make dev

# Or build and run
make build
./bin/server
```

The server will start on `http://localhost:8080`

### Test It

```bash
# Health check
curl http://localhost:8080/healthz

# Create an item
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Item","description":"This is a test"}'

# Get all items
curl http://localhost:8080/api/v1/items
```

### AWS Features

```bash
# Configure AWS credentials
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret

# Or use AWS CLI
aws configure

# Test AWS endpoints
curl http://localhost:8080/api/v1/aws/s3/buckets
curl http://localhost:8080/api/v1/aws/dynamodb/tables
```

See [QUICKSTART_AWS.md](./QUICKSTART_AWS.md) for detailed AWS setup.

## Project Structure

```
AWS-Go-Server/
├── cmd/server/              # Application entry point
│   └── main.go             # Minimal main function
├── internal/                # Private application code
│   ├── aws/                # AWS integration
│   │   ├── client.go      # AWS client initialization
│   │   └── auth.go        # IAM authentication middleware
│   ├── config/             # Configuration management
│   │   └── config.go      # Config loading from environment
│   ├── handlers/           # HTTP request handlers
│   │   ├── health.go      # Health check
│   │   ├── items.go       # Item CRUD operations
│   │   └── aws.go         # AWS service endpoints
│   ├── middleware/         # HTTP middleware
│   │   ├── logging.go     # Request logging
│   │   ├── recovery.go    # Panic recovery
│   │   └── sizelimit.go   # Request size limiting
│   ├── models/             # Domain models (ready for use)
│   └── server/             # HTTP server setup
│       ├── server.go      # Server initialization & lifecycle
│       └── routes.go      # Route definitions
├── pkg/                     # Public libraries (ready for use)
├── deployments/docker/      # Deployment configurations
│   ├── Dockerfile          # Multi-stage production build
│   └── docker-compose.yml  # Docker Compose setup
├── scripts/                 # Build and deployment scripts
├── .air.toml               # Hot reload configuration
├── Makefile                # Common development tasks
└── go.mod                  # Go module definition
```

**Why this structure?** See [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) for detailed architecture documentation and design principles.

## Available Commands

```bash
make help              # Show all available commands
make build             # Build the binary
make run               # Run the server
make dev               # Run with hot reload (auto-restarts on changes)
make test              # Run tests
make test-coverage     # Run tests with HTML coverage report
make lint              # Run linter
make fmt               # Format code
make vet               # Run go vet
make check             # Run all checks (fmt, vet, lint, test)
make docker-build      # Build Docker image
make docker-up         # Start with docker-compose
make docker-down       # Stop docker-compose
make docker-logs       # View container logs
make install-tools     # Install dev tools (air, golangci-lint)
make clean             # Clean build artifacts
```

## Configuration

Configuration is loaded from environment variables with sensible defaults:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `localhost` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `AWS_REGION` | `us-east-1` | AWS region |
| `AWS_PROFILE` | (empty) | AWS profile name |

Example `.env` file:
```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
AWS_REGION=us-west-2
AWS_PROFILE=dev
```

## API Endpoints

### Health & Status
- `GET /healthz` - Health check

### Items (CRUD)
- `GET /api/v1/items` - List all items
- `POST /api/v1/items` - Create a new item
  - Request body: `{"name":"string","description":"string"}`
  - Validation: name required, max 100 chars; description max 500 chars

### AWS Services
- `GET /api/v1/aws/s3/buckets` - List S3 buckets
- `GET /api/v1/aws/dynamodb/tables` - List DynamoDB tables

Full API documentation: [AWS_INTEGRATION.md](./AWS_INTEGRATION.md)

## Architecture Principles

This server follows Mat Ryer's [battle-tested architecture patterns](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/):

1. **No Global State** - All dependencies explicitly injected
2. **Testability First** - OS dependencies (stdout, args) injectable for testing
3. **Graceful Shutdown** - Proper context propagation and signal handling
4. **Explicit Dependencies** - Every component declares what it needs
5. **Separation of Concerns** - Clear package boundaries

### Key Patterns

**Handler Maker Functions** - Functions that return handlers:
```go
func HandleItemsGet(logger *slog.Logger) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Handler logic with access to logger via closure
    })
}
```

**Middleware with Dependencies**:
```go
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(h http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger.Info("request started", "path", r.URL.Path)
            h.ServeHTTP(w, r)
        })
    }
}
```

## Resilience Features

- **Panic Recovery** - Server stays running even if handlers panic (middleware/recovery.go:11)
- **Graceful Shutdown** - Clean shutdown on SIGINT/SIGTERM with 10s timeout (server/server.go:37)
- **Request Timeouts** - 15s read/write, 60s idle timeout (server/server.go:41-43)
- **Request Size Limits** - 10MB max request size (middleware/sizelimit.go:8)
- **Race Condition Protection** - Thread-safe concurrent access with RWMutex (handlers/items.go:22)
- **Comprehensive Logging** - Structured logging for all operations

See [RESILIENCE_IMPROVEMENTS.md](./RESILIENCE_IMPROVEMENTS.md) for details.

## Security Features

- **IAM Authentication** - AWS SigV4 signature verification (aws/auth.go:19)
- **Input Validation** - Request data validation with detailed error messages
- **Request Size Limits** - Prevents memory exhaustion attacks
- **Timeout Protection** - Prevents slowloris attacks
- **Non-root Docker User** - Container runs as non-privileged user
- **No Credentials in Code** - Uses AWS credential chain

## Development

### Hot Reload

```bash
# Install development tools (including air)
make install-tools

# Run with auto-reload on file changes
make dev
```

Air watches `.go` files and automatically rebuilds/restarts on changes.

### Testing

```bash
# Run tests
make test

# Run tests with coverage report
make test-coverage
open coverage.html
```

### Code Quality

```bash
# Format code
make fmt

# Run static analysis
make vet

# Run linter
make lint

# Run all checks (fmt + vet + lint + test)
make check
```

## Docker

### Using Docker Compose (Recommended)

```bash
# Start services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down

# Rebuild and restart
make docker-rebuild
```

### Manual Docker Commands

```bash
# Build image
docker build -f deployments/docker/Dockerfile -t aws-go-server:latest .

# Run container
docker run -p 8080:8080 \
  -e AWS_REGION=us-east-1 \
  -e AWS_ACCESS_KEY_ID=your_key \
  -e AWS_SECRET_ACCESS_KEY=your_secret \
  aws-go-server:latest
```

The Docker image:
- Multi-stage build for minimal size (~20MB)
- Runs as non-root user for security
- Includes health checks
- Optimized for production

## Deployment

### Build for Production

```bash
# Build optimized binary
make build

# Binary is in bin/server
./bin/server
```

### Deploy to AWS

Options for deployment:
- **EC2** - Run binary directly or use Docker
- **ECS/Fargate** - Use provided Dockerfile
- **Lambda** - Package as Lambda function (requires adapter)
- **Kubernetes** - Add K8s manifests to `deployments/k8s/`

## Documentation

- [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) - Detailed architecture and design patterns
- [AWS_INTEGRATION.md](./AWS_INTEGRATION.md) - Complete AWS integration guide
- [QUICKSTART_AWS.md](./QUICKSTART_AWS.md) - Quick AWS setup guide
- [RESILIENCE_IMPROVEMENTS.md](./RESILIENCE_IMPROVEMENTS.md) - Resilience features explained

## Adding Features

### Adding a New Endpoint

1. Create handler in `internal/handlers/`:
   ```go
   func HandleNewFeature(logger *slog.Logger) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           // Implementation
       })
   }
   ```

2. Register route in `internal/server/routes.go`:
   ```go
   mux.Handle("GET /api/v1/feature", handlers.HandleNewFeature(s.logger))
   ```

3. Add tests

### Adding a Database

1. Create `internal/repository/` package for data access
2. Add models to `internal/models/`
3. Create service layer in `internal/service/`
4. Inject repositories into handlers via server

See [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) for patterns.

## Performance

- **Concurrent Request Handling** - Leverages Go's goroutines
- **Connection Pooling** - HTTP keep-alive enabled
- **Optimized Docker Image** - Multi-stage build, minimal size
- **Fast Startup** - < 1 second cold start
- **Race-Free Concurrency** - Mutex protection on shared state

## Roadmap

- [ ] Database integration (PostgreSQL)
- [ ] Redis caching layer
- [ ] Rate limiting middleware
- [ ] OpenAPI/Swagger documentation
- [ ] Prometheus metrics
- [ ] Distributed tracing (OpenTelemetry)
- [ ] CI/CD pipelines (GitHub Actions)
- [ ] Kubernetes manifests
- [ ] More AWS service integrations (SQS, SNS, Lambda)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make check` to ensure quality
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- Architecture inspired by [Mat Ryer's 13-year evolution](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/)
- Built with [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)
- Follows [Go standard project layout](https://github.com/golang-standards/project-layout)
- Structured logging with Go's [slog](https://pkg.go.dev/log/slog)

## Support

- Review documentation in the project root
- Check [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) for architecture details
- Open an issue for bugs or feature requests

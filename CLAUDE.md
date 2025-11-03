# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-hello-devops is a simple Go web application designed for DevOps engineers learning software development. It demonstrates professional patterns like HTTP routing, middleware, JSON APIs, testing, and containerization. The application is intentionally simple (serves "Hello World" pages and basic APIs) but structured to teach modern development practices.

## Architecture

### Single-File Application Structure

The entire application is contained in `main.go` for simplicity. Key components:

- **HTTP Handlers**: Functions that process requests (`handleRoot`, `handleHealth`, `handleMessage`)
- **Middleware Pattern**: `loggingMiddleware` wraps handlers to add logging behavior
- **Response Types**: Structs with JSON tags (`HealthResponse`, `MessageResponse`) control JSON serialization
- **Server Configuration**: Uses standard library `http.ServeMux` for routing with proper timeouts

### Development Environment

This project uses a containerized development environment with two services:

- **app**: The Go web application (port 8000)
- **devbox**: Code-server IDE running devops-coderbox (port 8080)

Both run via Docker Compose and share the workspace directory. Changes to Go files require restarting the app container to recompile.

## Common Commands

### Development Workflow

```bash
# Start the full environment (app + IDE)
docker compose up
# Or with Podman:
podman-compose --pod-args '--userns keep-id' up

# Restart the app after code changes (REQUIRED - Go is compiled)
docker compose restart app

# Stop everything
docker compose down
```

### Building and Testing

```bash
# Build the application binary
make build
# Or: go build -o bin/server .

# Run all tests
make test
# Or: go test -v ./...

# Run tests with coverage
make test-coverage
# Or: go test -cover ./...

# Run a single test
go test -v -run TestHandleHealth ./...

# Run benchmarks
make bench
# Or: go test -bench=. ./...

# Format code (always run before committing)
make fmt
# Or: go fmt ./...
```

### Direct Execution (without Docker)

```bash
# Run the app locally
make run
# Or: go run .

# Clean build artifacts
make clean
```

## Testing Patterns

Tests in `main_test.go` demonstrate:

- **httptest Package**: Create fake requests (`httptest.NewRequest`) and record responses (`httptest.NewRecorder`)
- **Handler Testing**: Call handlers directly with test request/response objects
- **JSON Validation**: Unmarshal responses and verify structure
- **Middleware Testing**: Verify middleware calls wrapped handlers correctly
- **Benchmarking**: Functions starting with `Benchmark` measure performance

All tests must be updated when changing response structures.

## Important Development Notes

1. **Go is compiled**: After editing Go files, you MUST restart the app container with `docker compose restart app` to see changes. This is different from interpreted languages.

2. **Environment Variables**: Required variables are in `.env` (created from `.env.example`). Docker Compose will fail with clear errors if required vars are missing:
   - `GITHUB_USERNAME`: For pulling devops-coderbox image
   - `GIT_USER_NAME` and `GIT_USER_EMAIL`: For git commits
   - `CODE_SERVER_PASSWORD`: IDE authentication

3. **Port Configuration**: The app respects the `PORT` environment variable but defaults to 8000.

4. **Health Checks**: The `/health` endpoint returns JSON with status, timestamp, and version. Used by Docker healthchecks and monitoring systems.

5. **Middleware Pattern**: To add functionality to all routes (auth, rate limiting, etc.), wrap handlers with middleware functions following the `loggingMiddleware` pattern.

## Adding New Features

When adding new endpoints:

1. Define response struct with json tags
2. Implement handler function
3. Register route in `main()` with `mux.HandleFunc()`
4. Write tests in `main_test.go`
5. Restart app container to see changes

Example flow is documented extensively in README.md "Adding Your First Feature" section.

## Code Style

- All code is heavily commented to teach concepts
- Standard library preferred over external dependencies
- Handlers log using `log.Printf()`
- JSON encoding errors are logged but don't stop execution (since headers are already written)
- All Go code must be formatted with `go fmt`

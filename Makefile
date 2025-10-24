# Makefile for go-hello-devops
# Makefiles are a classic Unix tool for defining common tasks and their dependencies.
# While many modern languages have their own task runners, Make is still widely used
# because it's universal and extremely powerful.

# The .PHONY target tells Make that these aren't actual files, they're commands
# Without this, if you had a file named "test" in your directory, "make test" would
# get confused
.PHONY: help build test run clean docker-build docker-run dev stop

# The default target runs when you just type "make" with no arguments
# We make it show the help message so people can see what commands are available
help:
	@echo "Available targets:"
	@echo "  make build        - Build the Go binary"
	@echo "  make test         - Run tests"
	@echo "  make run          - Run the application locally"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make docker-build - Build the Docker image"
	@echo "  make docker-run   - Run the application in Docker"
	@echo "  make dev          - Start the full development environment"
	@echo "  make stop         - Stop all containers"

# Build the Go binary
# This compiles your application into an executable file
build:
	@echo "Building application..."
	go build -o bin/server .
	@echo "Build complete! Binary is at bin/server"

# Run all tests
# The -v flag gives verbose output so you can see which tests are running
# The -race flag enables the race detector which finds concurrency bugs
# The -cover flag shows test coverage percentage
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# Run tests with coverage report
# This creates an HTML file showing which lines of code are covered by tests
test-coverage:
	@echo "Running tests with coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run benchmarks
# Benchmarks measure performance. They're useful when you're optimizing code.
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Run the application locally (not in Docker)
# This is useful for quick iteration when you don't need the full Docker environment
run:
	@echo "Starting server on port 8000..."
	go run .

# Clean up build artifacts
# It's good practice to have a clean target that removes generated files
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Build the Docker image for the application
docker-build:
	@echo "Building Docker image..."
	docker build -f Dockerfile.app -t go-hello-devops:latest .
	@echo "Docker image built successfully!"

# Run the application in Docker
docker-run:
	@echo "Starting application in Docker..."
	docker run --rm -p 8000:8000 go-hello-devops:latest

# Start the full development environment (app + IDE)
# This is your main command for day-to-day development
dev:
	@echo "Starting development environment..."
	@echo "Application will be available at http://localhost:8000"
	@echo "IDE will be available at http://localhost:8080"
	@echo "Press Ctrl+C to stop"
	docker-compose up

# Start the development environment in the background
dev-detached:
	@echo "Starting development environment in background..."
	docker-compose up -d
	@echo "Application running at http://localhost:8000"
	@echo "IDE running at http://localhost:8080"
	@echo "Use 'make stop' to stop the environment"

# Stop all containers
stop:
	@echo "Stopping all containers..."
	docker-compose down
	@echo "All containers stopped!"

# View logs from running containers
logs:
	docker-compose logs -f

# Format Go code according to standard formatting rules
# gofmt is the standard Go formatter. All Go code should be formatted with it.
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run the Go linter
# golangci-lint is a meta-linter that runs many different linters
# If you don't have it installed, run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
lint:
	@echo "Running linter..."
	golangci-lint run

# Install development dependencies
# This installs tools that are useful for development but not needed to run the app
deps:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Dependencies installed!"

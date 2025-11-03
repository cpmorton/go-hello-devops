# go-hello-devops

A simple Go web application designed for DevOps engineers learning software development. The app demonstrates professional patterns like HTTP routing, middleware, JSON APIs, testing, and containerization.

## Table of Contents

- [What Is This?](#what-is-this)
- [Quick Start](#quick-start)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Adding Your First Feature](#adding-your-first-feature)
- [Understanding the Code](#understanding-the-code)
- [Common Commands](#common-commands)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

## What Is This?

This project teaches you modern development practices through a simple but professionally-structured Go web application. The app itself just serves "Hello World" pages and basic APIs, but the infrastructure around it demonstrates:

- **HTTP routing and handlers** - How web servers work in Go
- **Middleware patterns** - Adding behavior to handlers
- **JSON APIs** - Building and testing API endpoints
- **Test-driven development** - Writing tests alongside code
- **Containerization** - Docker multi-stage builds
- **Docker Compose** - Orchestrating multiple containers
- **Development environments** - Browser-based IDE with your project

The code is heavily commented to explain not just what things do, but why they're structured that way.

## Quick Start

**Requirements:** Docker Engine installed and running (see [Prerequisites](#prerequisites))

1. Clone this repository:
```bash
git clone https://github.com/yourusername/go-hello-devops.git
cd go-hello-devops
```

2. Configure environment variables:
```bash
# Copy the example file
cp .env.example .env

# Edit .env with your actual values
# At minimum, you MUST set:
#   - GITHUB_USERNAME (your GitHub username)
#   - GIT_USER_NAME (your name for git commits)
#   - GIT_USER_EMAIL (your email for git commits)
#   - CODE_SERVER_PASSWORD (optional, defaults to devops-coderbox)
nano .env
```

3. Start everything:
```bash
docker compose up
# Or, podman:
podman-compose --pod-args '--userns keep-id' up
```

The containers will fail to start if required environment variables aren't set.

4. Open two browser tabs:
   - http://localhost:8000 - Your application
   - http://localhost:8080 - Your IDE (use the password from your .env file)

5. Make changes in the IDE, save files, and restart to see changes:
```bash
# In another terminal
docker compose restart app
```

## Prerequisites

### Docker:
You need Docker Engine (preferred to Docker Desktop) to run this project. Docker Desktop has licensing restrictions for large companies.
Optionally, you may wish to use podman with podman-compose...

#### Installing Docker Engine on WSL2 (Windows)

If you're on Windows, install WSL2 first:

```powershell
wsl --install -d Ubuntu-24.04
```

Restart, then open your Ubuntu terminal and install Docker Engine:

```bash
# Update and install prerequisites
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg lsb-release

# Add Docker's GPG key
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Add Docker repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Start Docker and add yourself to docker group
sudo service docker start
sudo usermod -aG docker $USER
```

Close and reopen your terminal, then verify:
```bash
docker run hello-world
```

**Note:** Each time you open a new WSL2 terminal, run `sudo service docker start`.

#### Other Platforms

**Linux:** Install Docker Engine  or podman and podman-compose from your package manager (see Docker's docs for your distro)

**macOS:** Install Docker Desktop from docker.com (free for small companies, educational use, and personal projects)

## Project Structure

```
go-hello-devops/
├── main.go              # Application code - read this first
├── main_test.go         # Tests - demonstrates testing patterns
├── go.mod              # Go module definition
├── Dockerfile.app      # How to containerize the app
├── docker-compose.yml  # Orchestrates app + IDE
├── .env.example        # Template for environment variables
├── .env               # Your actual environment variables (DO NOT COMMIT)
├── Makefile           # Convenient command shortcuts
└── README.md          # This file
```

**Start by reading `main.go`** - it's heavily commented to teach you Go web development patterns.

## Environment Variables

This project requires several environment variables. Create a `.env` file based on `.env.example`:

```bash
cp .env.example .env
```

**Required variables:**

- `GITHUB_USERNAME` - Your GitHub username (for pulling the devops-coderbox image)
- `GIT_USER_NAME` - Your name for git commits (e.g., "John Smith")
- `GIT_USER_EMAIL` - Your email for git commits (e.g., "john@example.com")

**Optional variables:**

- `CODE_SERVER_PASSWORD` - IDE password (defaults to "devops-coderbox")
- `ANTHROPIC_API_KEY` - Your Anthropic API key from https://console.anthropic.com (do not use if you have a pro or max subscription)

Docker Compose will fail with a clear error message if required variables are missing.

**Security note:** The `.env` file is in `.gitignore` and should NEVER be committed to git. It contains your API key and other secrets.

## Development Workflow

Your typical workflow looks like this:

1. **Start the environment:**
```bash
docker compose up
# Or, if you use podman:
podman-compose --pod-args '--userns keep-id' up
```

2. **Open the IDE** at http://localhost:8080 (password: `devops-coderbox`)

3. **Edit code** in the IDE and save files

4. **Run tests** in the IDE's integrated terminal (Ctrl+` to open):
```bash
go test -v ./...
```

5. **See your changes** by restarting the app:
```bash
docker compose restart app
```

6. **Commit your changes:**
```bash
git add .
git commit -m "Describe what you changed and why"
git push
```

## Adding Your First Feature

Let's add a new API endpoint that returns the current time. This demonstrates the full development cycle.

### Step 1: Define the Response Structure

Open `main.go` in the IDE. After the existing response types (around line 20), add:

```go
type TimeResponse struct {
    UTC       string `json:"utc"`
    LocalTime string `json:"local"`
    Timestamp int64  `json:"timestamp"`
}
```

### Step 2: Implement the Handler

After the existing handlers (around line 120), add:

```go
func handleTime(w http.ResponseWriter, r *http.Request) {
    // This endpoint returns the current time in multiple formats
    now := time.Now()
    
    response := TimeResponse{
        UTC:       now.UTC().Format(time.RFC3339),
        LocalTime: now.Format(time.RFC3339),
        Timestamp: now.Unix(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding time response: %v", err)
    }
}
```

### Step 3: Register the Route

In the `main()` function, where other routes are registered (around line 170), add:

```go
mux.HandleFunc("/api/time", loggingMiddleware(handleTime))
```

### Step 4: Write Tests

Open `main_test.go` and add:

```go
func TestHandleTime(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/api/time", nil)
    rec := httptest.NewRecorder()
    
    handleTime(rec, req)
    
    if rec.Code != http.StatusOK {
        t.Fatalf("Expected status 200, got %d", rec.Code)
    }
    
    var response TimeResponse
    if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse JSON response: %v", err)
    }
    
    if response.UTC == "" {
        t.Error("Expected UTC time to be set")
    }
    
    if response.Timestamp == 0 {
        t.Error("Expected timestamp to be set")
    }
}
```

### Step 5: Test and Run

In the IDE terminal:

```bash
# Run tests
go test -v ./...

# All tests should pass
```

Restart the app:

```bash
docker compose restart app
```

Visit http://localhost:8000/api/time in your browser to see your new endpoint in action.

Congratulations! You just added a feature using test-driven development.

## Understanding the Code

### HTTP Handlers

Handlers are functions that process HTTP requests:

```go
func handleRoot(w http.ResponseWriter, r *http.Request) {
    // w is where you write your response
    // r contains information about the incoming request
}
```

Every endpoint in the application is a handler function.

### Middleware

Middleware wraps handlers to add behavior. The `loggingMiddleware` logs information about every request:

```go
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Code here runs before the handler
        next(w, r)  // Call the actual handler
        // Code here runs after the handler
    }
}
```

This pattern is how you implement authentication, rate limiting, or any cross-cutting concern.

### JSON APIs

To return JSON, create a struct with json tags:

```go
type Response struct {
    Message string `json:"message"`
}

response := Response{Message: "Hello"}
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
```

### Testing

Tests use the `httptest` package to simulate HTTP requests:

```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()
    
    handleRoot(rec, req)
    
    if rec.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", rec.Code)
    }
}
```

Run tests frequently as you develop to catch issues early.

## Common Commands

Using Make (convenient shortcuts):

```bash
make help          # Show all available commands
make dev           # Start development environment
make test          # Run tests
make test-coverage # Generate coverage report
make stop          # Stop all containers
```

Using Docker Compose directly:

```bash
docker compose up              # Start everything
docker compose down            # Stop and remove containers
docker compose restart app     # Restart just the app
docker compose logs -f         # Watch logs
```

Or, using Podman Compose directly:

```bash
podman-compose  --pod-args '--userns keep-id' up    # Start everything
podman-compose down                                 # Stop and remove containers
podman-compose restart app                          # Restart just the app
podman-compose logs -f                              # Watch logs
```

Using Go directly (inside the IDE terminal):

```bash
go run .                       # Run without Docker
go test -v ./...              # Run all tests
go test -cover ./...          # Run tests with coverage
go fmt ./...                  # Format code
go build -o bin/server .      # Build binary
```

## Testing

This project emphasizes test-driven development. The test file (`main_test.go`) demonstrates:

- **Unit testing handlers** - Verify each endpoint works correctly
- **Testing JSON APIs** - Parse and validate JSON responses
- **Testing middleware** - Ensure middleware calls handlers correctly
- **Benchmarking** - Measure handler performance

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
# Open coverage.html in your browser

# Run benchmarks
go test -bench=. ./...
```

### Writing Tests

Follow this pattern:

```go
func TestHandleThing(t *testing.T) {
    // Create fake request
    req := httptest.NewRequest(http.MethodGet, "/thing", nil)
    rec := httptest.NewRecorder()
    
    // Call handler
    handleThing(rec, req)
    
    // Verify response
    if rec.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", rec.Code)
    }
}
```

## Troubleshooting

**Changes don't appear when I refresh the browser:**
- Go is a compiled language. You must rebuild the app container.
- Run: `docker compose restart app`
- Wait for the container to restart, then refresh your browser.

**"Port already in use" error:**
- Something else is using ports 8000 or 8080
- Stop other services or change ports in `docker-compose.yml`

**Tests fail after making changes:**
- Read the test output carefully - it tells you what's wrong
- Common issue: changed response structure but didn't update tests
- Fix: Update tests to match your new code

**IDE won't load at localhost:8080:**
- Is Docker running? Check with `docker ps`
- Is the devbox container running? Should see `hello-devops-devbox` in `docker ps`
- Try http://127.0.0.1:8080 instead of localhost:8080

**"go: command not found" in IDE terminal:**
- The devbox container includes Go. Make sure you're in the container's terminal.
- In code-server, open the integrated terminal with Ctrl+`

**Docker says "permission denied":**
- On Linux/WSL2, did you add yourself to the docker group?
- Run: `sudo usermod -aG docker $USER`
- Then log out and back in

**WSL2: "Cannot connect to Docker daemon":**
- Docker isn't running. Start it: `sudo service docker start`
- You need to run this command each time you open a new WSL2 terminal

## Next Steps

Once you're comfortable with the basics:

1. **Add more features** - Implement new endpoints, add a database, handle file uploads
2. **Learn about authentication** - Add login/logout endpoints with session management
3. **Explore concurrency** - Use goroutines and channels for background tasks
4. **Set up CI/CD** - Add GitHub Actions to automatically test and deploy
5. **Deploy to the cloud** - Run your app in Azure, AWS, or another cloud provider

The goal isn't to build a complex application immediately. The goal is to learn patterns and practices that scale from simple learning projects to production systems.

## Resources

- [Go Tour](https://go.dev/tour/) - Interactive introduction to Go
- [Go by Example](https://gobyexample.com/) - Learn Go through annotated example programs
- [Effective Go](https://go.dev/doc/effective_go) - Go best practices
- [Go standard library](https://pkg.go.dev/std) - Everything included with Go

---

## License

MIT License - see LICENSE file for details.

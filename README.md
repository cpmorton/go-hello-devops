# Go Hello DevOps - Your First Web Application

This is a simple Go web application designed as a learning project for DevOps engineers and systems administrators who want to learn modern software development practices. The application itself is intentionally simple (it serves a "Hello World" page), but the infrastructure around it demonstrates professional development patterns you'll use in real projects.

## What You'll Learn

This project teaches you how to build, test, containerize, and develop web applications using modern tools and patterns. By working through this project, you'll gain hands-on experience with Go programming, Docker containers, test-driven development, and browser-based development environments. The code is heavily commented to explain not just what each piece does, but why it's structured that way and what patterns it demonstrates.

## Quick Start

The fastest way to get started is using Docker Compose, which orchestrates both your application and your development environment in a single command. This requires that you have Docker Desktop (or Docker and Docker Compose) installed on your system.

First, clone this repository to your local machine. Navigate to wherever you keep your projects and run these commands:

```bash
git clone https://github.com/cpmorton/go-hello-devops.git
cd go-hello-devops
```

Now start the development environment with a single command:

```bash
docker-compose up
```

This command starts two containers. The first container runs your actual Go web application, which you can access at http://localhost:8000 in your browser. The second container runs devops-coderbox, which is a browser-based IDE (essentially VS Code running in a container) that you can access at http://localhost:8080. The default password for the IDE is `devops-coderbox`.

When you open the IDE, you'll see the project files already loaded. You can edit the code, and because the project directory is mounted into both containers, your changes are immediately visible. Try editing the message in main.go and save the file, then restart the app container to see your changes.

To stop the environment, press Ctrl+C in the terminal where docker-compose is running, or run `docker-compose down` from another terminal in the same directory.

## What's Included

This project includes everything you need to start learning web development with Go. The main application file demonstrates how to build HTTP servers with proper routing, middleware, and JSON API endpoints. The test file shows you how to write unit tests and benchmarks for your HTTP handlers, which is an essential skill for building reliable software. The Dockerfile shows you how to containerize Go applications using multi-stage builds, which results in small, efficient container images. The docker-compose configuration demonstrates how to orchestrate multiple containers and manage their interactions.

The Makefile provides convenient commands for common development tasks like building, testing, and running your application. You can see all available commands by running `make help` from the project directory.

## Project Structure

Let me walk you through the important files in this project and explain what each one does.

**main.go** is your application code. This file defines the HTTP server, the routes it responds to, and the handlers that process those requests. The code is heavily commented to explain Go's patterns for web development. You'll find examples of serving HTML, returning JSON from APIs, implementing middleware, and configuring HTTP servers with proper timeouts.

**main_test.go** contains tests for your application. These tests demonstrate how to write unit tests for HTTP handlers using Go's httptest package. The tests verify that your handlers return the correct status codes, headers, and response bodies. There are also benchmark functions that measure how fast your handlers execute, which is useful when you're optimizing performance.

**go.mod** defines this project as a Go module and lists its dependencies. Since this application only uses Go's standard library, there are no external dependencies listed yet. As you add third-party packages, they'll automatically be recorded here when you run `go mod tidy`.

**Dockerfile.app** is a multi-stage Dockerfile that builds your Go application into a container image. The first stage compiles your code, and the second stage creates a minimal runtime image with just your binary and the certificates needed for HTTPS. This pattern results in container images that are only about 10-15 MB instead of hundreds of megabytes.

**docker-compose.yml** orchestrates multiple containers together. In this project, it runs both your application and the devops-coderbox development environment, with the appropriate networking and volume mounts configured. This is how you'd typically run complex applications that need multiple services (databases, caches, background workers, etc).

**Makefile** defines common development tasks as simple commands. Instead of remembering long docker commands or complicated test flags, you can just run `make test` or `make build`. The Makefile documents all available commands if you run `make help`.

## Development Workflow

The typical development workflow looks like this. You start the development environment with `docker-compose up`, which launches both your application and the IDE. You open the IDE in your browser at http://localhost:8080 and authenticate with the password `devops-coderbox`. You open the main.go file in the IDE and make some changes to the code. You save the file, and your changes are immediately written to your local filesystem because the project directory is mounted into the container.

To see your changes in action, you need to restart the application container. You can do this by pressing Ctrl+C in the terminal where docker-compose is running, then running `docker-compose up` again. Alternatively, in another terminal, you can run `docker-compose restart app` to restart just the application container without restarting the IDE.

As you're making changes, you should frequently run the tests to verify that your changes don't break existing functionality. You can run tests from inside the IDE's terminal (open Terminal from the menu) by typing `go test -v ./...`, or you can use the Makefile with `make test`. The tests will show you immediately if something is broken, which saves you from discovering problems much later.

## Writing Your First Feature

Let's walk through adding a new feature to the application. This will help you understand the development cycle and get comfortable with the tools.

Open the IDE at http://localhost:8080 and navigate to main.go. Let's add a new API endpoint that returns the current time in different time zones. First, add a new struct type after the existing response types, around line 20:

```go
type TimeResponse struct {
    UTC       string `json:"utc"`
    LocalTime string `json:"local"`
    Timestamp int64  `json:"timestamp"`
}
```

This defines the structure of our JSON response. Now add a new handler function after the existing handlers, around line 120:

```go
func handleTime(w http.ResponseWriter, r *http.Request) {
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

This handler creates a TimeResponse with the current time in different formats and returns it as JSON. Now register this handler with the router by adding this line in the main function where the other handlers are registered, around line 170:

```go
mux.HandleFunc("/api/time", loggingMiddleware(handleTime))
```

Save the file. Now you need to write a test for your new handler. Open main_test.go and add this test function:

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

Save the test file. Now run the tests from the IDE's terminal with `go test -v ./...` and verify that all tests pass, including your new one. Restart the application container with `docker-compose restart app`, then visit http://localhost:8000/api/time in your browser to see your new endpoint in action.

Congratulations! You just added a feature using proper development practices. You wrote the code, added tests to verify it works, and confirmed it works in the running application. This is the cycle you'll follow for every feature you build.

## Understanding the Code

Let me explain some key concepts that appear in this codebase, because understanding these patterns will help you as you start modifying and extending the application.

**HTTP Handlers** are functions that process HTTP requests. In Go, a handler has the signature `func(w http.ResponseWriter, r *http.Request)`. The `w` parameter is where you write your response (headers and body), and the `r` parameter contains information about the incoming request (URL, headers, form data, etc). Every endpoint in your application is implemented as a handler function.

**Middleware** is a pattern for wrapping handlers with additional functionality. The loggingMiddleware in this application wraps every handler to log information about requests. Middleware is how you implement cross-cutting concerns like authentication, rate limiting, or request logging without duplicating code in every handler. Middleware functions take a handler as input and return a new handler that adds some behavior before or after calling the original handler.

**JSON Encoding** is how you convert Go structs into JSON for API responses. The `json.NewEncoder(w).Encode(response)` pattern takes a Go struct and writes it to the response writer as JSON. The struct tags like `json:"status"` control what the fields are named in the JSON output. This is the standard pattern for building JSON APIs in Go.

**Test-Driven Development** is the practice of writing tests for your code before or immediately after writing the code itself. The tests in main_test.go verify that your handlers work correctly by simulating HTTP requests and checking the responses. When you modify your code, running the tests immediately tells you if you broke something. This gives you confidence to make changes and refactor code without fear.

**Struct Tags** are the backtick strings you see after struct fields, like `json:"status"`. These tags provide metadata about the fields that other packages can read. The json package uses these tags to control how structs are encoded to and decoded from JSON. Tags are Go's way of adding annotations to struct fields.

## Common Development Tasks

Here are the commands you'll use most frequently during development, along with explanations of what they do and when you'd use them.

**make dev** starts the full development environment with both the application and the IDE. This is typically the first command you run when starting your work session. It keeps the terminal attached so you can see logs from both containers.

**make dev-detached** starts the environment in the background, which frees up your terminal for other commands. You'd use this when you want the services running but don't need to watch the logs. Use `make stop` to stop services started this way.

**make test** runs all your tests. You should run this frequently as you're making changes to verify that everything still works. The tests run quickly, so there's no excuse not to run them often.

**make test-coverage** generates an HTML report showing which lines of your code are covered by tests. This helps you identify areas that need more testing. Open coverage.html in your browser after running this command.

**make logs** shows the logs from your running containers. This is useful when you've started the environment in detached mode but want to see what's happening.

**make stop** stops all running containers. You'd use this when you're done working or when you need to free up the ports for something else.

## Deploying to Azure (Future Enhancement)

This section is a placeholder for future work. Eventually, you'll add CI/CD pipelines that automatically build your container images and deploy them to Azure Container Apps or Azure Kubernetes Service. The deployment process will include automated testing, security scanning, and blue-green deployments to minimize downtime.

For now, focus on learning Go, writing tests, and building features. The deployment infrastructure will come later as you become more comfortable with the development workflow.

## Troubleshooting

If you encounter issues, here are some common problems and their solutions.

**The IDE won't load at localhost:8080.** Verify that the devbox container is running by executing `docker ps` and looking for a container named hello-devops-devbox. If it's not running, check the docker-compose logs with `make logs` to see if there are any error messages. Make sure no other service is using port 8080 on your machine.

**The application returns 404 for all requests.** This usually means the app container isn't running or didn't start correctly. Run `docker ps` to verify the hello-devops-app container is running. Check the logs with `make logs` and look for error messages during startup.

**Changes to the code don't appear when I refresh the browser.** Remember that Go is a compiled language, so you need to rebuild and restart the application for changes to take effect. The easiest way is to restart the app container with `docker-compose restart app`. You don't need to restart the IDE container.

**Tests fail with unexpected errors.** First, make sure your code compiles by running `go build .` from the IDE terminal. If there are compilation errors, fix those first. If the code compiles but tests fail, read the test output carefully. It will tell you which test failed and what the actual vs expected values were.

**Docker commands fail with permission errors.** On Linux, you need to add your user to the docker group. Run `sudo usermod -aG docker $USER`, then log out and back in. On Windows and Mac with Docker Desktop, make sure Docker Desktop is running.

## Next Steps

Once you're comfortable with this basic application, here are some ideas for extending it and continuing your learning.

Add a database to store data. You'd add a Postgres container to docker-compose.yml and modify your Go code to connect to it and store data. This teaches you about database integration and connection management.

Implement user authentication. Add login and registration endpoints, store user credentials securely (hashed passwords), and protect certain endpoints so they require authentication. This teaches you about security and session management.

Add a frontend framework. Create a React or Vue frontend that calls your API endpoints. This teaches you about building full-stack applications and dealing with CORS.

Set up CI/CD pipelines. Add GitHub Actions workflows that run your tests automatically when you push code, build container images, and deploy to Azure. This teaches you about automated deployment pipelines.

Add monitoring and observability. Integrate Prometheus metrics and structured logging to monitor your application's health and performance. This teaches you about production operations.

The goal isn't to build a real product right now. The goal is to learn the patterns and practices that professional developers use, so that when you do need to build something real, you'll know how to do it properly.

## Resources

Go's official documentation is excellent and should be your first stop for learning about language features. The Go Tour at tour.golang.org provides an interactive introduction to the language. The standard library documentation at pkg.go.dev shows you what's available and how to use it.

For learning about web development in Go specifically, check out Let's Go by Alex Edwards, which is a practical guide to building web applications. The book covers routing, middleware, databases, authentication, and testing in detail.

For Docker and containers, Docker's official documentation is comprehensive. Pay particular attention to the best practices guide for writing Dockerfiles and the Compose documentation for orchestrating multiple containers.

## Contributing

If you find issues with this project or have suggestions for improvements, please open an issue or submit a pull request on GitHub. This is a teaching project, so clarity and educational value are more important than complexity or completeness.

## License

MIT License - see LICENSE file for details.

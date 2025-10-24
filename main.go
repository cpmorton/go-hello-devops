package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// This is a simple HTTP server that demonstrates basic Go web development patterns.
// It's designed to be extended and modified as you learn, so the structure is
// intentionally simple and well-commented.

// HealthResponse represents the JSON structure we send for health check endpoints.
// In Go, we use struct tags to control how fields are serialized to JSON.
// The json:"fieldname" tag tells the JSON encoder what to call this field.
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// MessageResponse represents a simple message response.
// This demonstrates how to structure data for API responses.
type MessageResponse struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

// handleRoot handles requests to the root path "/"
// This is our main page that displays the hello world message.
func handleRoot(w http.ResponseWriter, r *http.Request) {
	// In a real application, you'd probably render an HTML template here.
	// For this simple example, we're just sending plain HTML.
	
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Hello DevOps!</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            text-align: center;
        }
        .container {
            background: rgba(255, 255, 255, 0.1);
            border-radius: 10px;
            padding: 40px;
            backdrop-filter: blur(10px);
        }
        h1 {
            font-size: 3em;
            margin: 0;
        }
        p {
            font-size: 1.2em;
            margin: 20px 0;
        }
        .info {
            margin-top: 30px;
            font-size: 0.9em;
            opacity: 0.8;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>👋 Hello DevOps!</h1>
        <p>Welcome to your first Go web application running in Coderbox.</p>
        <p>This is where your journey begins. Start editing and watch the changes happen!</p>
        <div class="info">
            <p>Try these endpoints:</p>
            <p>GET /health - Check if the service is running</p>
            <p>GET /api/message - Get a JSON response</p>
        </div>
    </div>
</body>
</html>
`
	
	// Set the content type header to tell the browser we're sending HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// Write the HTTP status code. 200 means OK.
	w.WriteHeader(http.StatusOK)
	
	// Write the HTML response
	fmt.Fprint(w, html)
	
	// Log that we served a request. In production, you'd use structured logging.
	log.Printf("Served request to %s from %s", r.URL.Path, r.RemoteAddr)
}

// handleHealth provides a health check endpoint for monitoring and orchestration.
// This is a standard pattern in cloud-native applications. Kubernetes, Docker,
// and cloud platforms use health endpoints to determine if your app is running correctly.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	// Create our health response with current information
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}
	
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Encode the response struct as JSON and write it to the response writer.
	// If encoding fails, we'll get an error, but at that point we've already
	// written the status code, so we just log the error.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding health response: %v", err)
	}
}

// handleMessage provides a simple API endpoint that returns a JSON message.
// This demonstrates the pattern for building JSON APIs in Go.
func handleMessage(w http.ResponseWriter, r *http.Request) {
	response := MessageResponse{
		Message: "This is your first API endpoint! Try modifying this message.",
		Time:    time.Now().Format(time.RFC3339),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding message response: %v", err)
	}
}

// loggingMiddleware wraps HTTP handlers to log requests.
// Middleware is a pattern in web development where you wrap handlers with
// additional functionality. This is how you implement cross-cutting concerns
// like logging, authentication, rate limiting, etc.
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Call the actual handler
		next(w, r)
		
		// Log information about the request after it's been handled
		duration := time.Since(start)
		log.Printf("%s %s completed in %v", r.Method, r.URL.Path, duration)
	}
}

func main() {
	// Get the port from an environment variable, defaulting to 8000 if not set.
	// This is a common pattern for configuring applications in containers.
	// Different environments can set different ports without changing the code.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	
	// Set up our HTTP routes using the standard library's http.ServeMux.
	// ServeMux is a request router that matches incoming requests to handlers.
	mux := http.NewServeMux()
	
	// Register our handlers with the router.
	// We wrap each handler with our logging middleware to get request logs.
	mux.HandleFunc("/", loggingMiddleware(handleRoot))
	mux.HandleFunc("/health", loggingMiddleware(handleHealth))
	mux.HandleFunc("/api/message", loggingMiddleware(handleMessage))
	
	// Configure the HTTP server.
	// In production, you'd want to set timeouts to prevent resource exhaustion.
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Log that we're starting up
	log.Printf("Starting server on port %s", port)
	log.Printf("Access the application at http://localhost:%s", port)
	
	// Start the server. ListenAndServe blocks until the server shuts down.
	// If there's an error starting the server (for example, if the port is
	// already in use), ListenAndServe returns the error and we log it and exit.
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

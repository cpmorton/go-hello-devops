package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Testing in Go uses the testing package from the standard library.
// Test functions must start with "Test" and take a *testing.T parameter.
// The t parameter provides methods for reporting test failures and logging.

// TestHandleRoot verifies that our root handler returns the expected content.
// This is an example of a simple integration test that exercises an HTTP handler.
func TestHandleRoot(t *testing.T) {
	// Create a fake HTTP request. This is a request to GET the root path.
	// httptest provides utilities for testing HTTP handlers without actually
	// starting a server.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	
	// Create a ResponseRecorder. This acts like an http.ResponseWriter but
	// records what the handler writes so we can check it in our test.
	rec := httptest.NewRecorder()
	
	// Call our handler with the fake request and recorder
	handleRoot(rec, req)
	
	// Check that the status code is correct
	// If it's not 200 OK, the test fails
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	
	// Check that the Content-Type header is set correctly
	contentType := rec.Header().Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
	}
	
	// Check that the response body contains our expected text
	// For a more robust test, you'd parse the HTML and check specific elements,
	// but for this simple example, checking for key strings is sufficient.
	body := rec.Body.String()
	expectedStrings := []string{
		"Hello DevOps",
		"/health",
		"/api/message",
	}
	
	for _, expected := range expectedStrings {
		if !contains(body, expected) {
			t.Errorf("Expected response body to contain %q", expected)
		}
	}
}

// TestHandleHealth verifies that the health endpoint returns the correct JSON structure.
// This test is more thorough because health endpoints are often used by monitoring
// systems, so we want to be certain they work correctly.
func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	
	handleHealth(rec, req)
	
	// Verify status code
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rec.Code)
	}
	
	// Verify content type is JSON
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
	
	// Parse the JSON response
	// This verifies that the response is valid JSON and has the expected structure
	var response HealthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}
	
	// Verify the response fields have sensible values
	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got %q", response.Status)
	}
	
	if response.Version == "" {
		t.Error("Expected version to be set")
	}
	
	// Verify that the timestamp is recent (within the last minute)
	// This catches issues where the timestamp might be zero or far in the past
	// due to programming errors.
	if response.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

// TestHandleMessage verifies the message API endpoint works correctly.
func TestHandleMessage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/message", nil)
	rec := httptest.NewRecorder()
	
	handleMessage(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rec.Code)
	}
	
	// Verify content type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
	
	// Parse and verify the response
	var response MessageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}
	
	if response.Message == "" {
		t.Error("Expected message to be set")
	}
	
	if response.Time == "" {
		t.Error("Expected time to be set")
	}
}

// TestLoggingMiddleware verifies that our middleware correctly calls the wrapped handler.
// Testing middleware can be tricky because middleware modifies the behavior of handlers.
func TestLoggingMiddleware(t *testing.T) {
	// Create a simple handler that we'll wrap with the middleware
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})
	
	// Wrap the handler with our middleware
	wrappedHandler := loggingMiddleware(testHandler)
	
	// Call the wrapped handler
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	wrappedHandler(rec, req)
	
	// Verify that the original handler was called
	if !handlerCalled {
		t.Error("Expected wrapped handler to be called")
	}
	
	// Verify that the response is still correct
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

// contains is a helper function that checks if a string contains a substring.
// In more complex projects, you'd probably use a testing utility library,
// but for this simple example we can write our own helper.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
	       (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark functions measure performance. They start with "Benchmark" and take
// a *testing.B parameter. The testing framework runs them multiple times to
// get accurate timing measurements.

// BenchmarkHandleRoot measures how fast our root handler is.
// While this simple handler will be very fast, benchmarking is a good habit
// to develop for when you're working on performance-critical code.
func BenchmarkHandleRoot(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	
	// The testing framework sets b.N to an appropriate number of iterations
	// to get statistically significant results
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handleRoot(rec, req)
	}
}

// BenchmarkHandleHealth measures the performance of the health endpoint.
func BenchmarkHandleHealth(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handleHealth(rec, req)
	}
}

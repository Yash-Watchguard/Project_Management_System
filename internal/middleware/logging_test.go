package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"


)


func TestLoggingMiddleWare(t *testing.T) {
	called := false

	// Mock handler to test if next handler is called
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Wrap with logging middleware
	handler := LoggingMiddleWare(mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/testpath", nil)
	rr := httptest.NewRecorder()

	start := time.Now()
	handler.ServeHTTP(rr, req)
	duration := time.Since(start)

	// Check if next handler was called
	if !called {
		t.Error("expected next handler to be called, but it was not")
	}

	// Check response status
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check response body
	if rr.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %q", rr.Body.String())
	}

	// Optional: check that duration is reasonable (<1s)
	if duration > time.Second {
		t.Errorf("handler took too long: %v", duration)
	}
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	
)


func TestCorsMiddleWare(t *testing.T) {
	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusTeapot) // 418 to check if called
	})

	middleware := CorsMiddleWare(nextHandler)

	// Test normal GET request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("Expected status %d, got %d", http.StatusTeapot, resp.StatusCode)
	}
	if !called {
		t.Errorf("Next handler was not called for GET request")
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("CORS header missing or incorrect")
	}
	if resp.Header.Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("CORS headers incorrect")
	}
	if resp.Header.Get("Access-Control-Allow-Methods") != "GET, POST, PUT, PATCH, OPTIONS, DELETE" {
		t.Errorf("CORS methods incorrect")
	}

	// Test OPTIONS preflight request
	called = false
	req = httptest.NewRequest(http.MethodOptions, "/", nil)
	w = httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d for OPTIONS, got %d", http.StatusOK, resp.StatusCode)
	}
	if called {
		t.Errorf("Next handler should not be called for OPTIONS request")
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("CORS header missing or incorrect for OPTIONS")
	}
	if resp.Header.Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("CORS headers incorrect for OPTIONS")
	}
	if resp.Header.Get("Access-Control-Allow-Methods") != "GET, POST, PUT, PATCH, OPTIONS, DELETE" {
		t.Errorf("CORS methods incorrect for OPTIONS")
	}
}

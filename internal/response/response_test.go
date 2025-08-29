package response

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestSuccessResponse(t *testing.T) {
    t.Run("with data", func(t *testing.T) {
        w := httptest.NewRecorder()
        data := map[string]string{"foo": "bar"}

        SuccessResponse(w, data, "ok", http.StatusCreated)

        if w.Code != http.StatusCreated {
            t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
        }

        if ct := w.Header().Get("Content-Type"); ct != "application/json" {
            t.Fatalf("expected Content-Type application/json, got %q", ct)
        }

        var got map[string]interface{}
        if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
            t.Fatalf("failed to decode response body: %v", err)
        }

        if got["status"] != "Success" {
            t.Fatalf("expected status field Success, got %v", got["status"])
        }
        if got["message"] != "ok" {
            t.Fatalf("expected message 'ok', got %v", got["message"])
        }

       
    })

}

func TestErrorResponse(t *testing.T) {
	// Prepare recorder (simulates http.ResponseWriter)
	rr := httptest.NewRecorder()

	// Call the function
	ErrorResponse(rr, http.StatusBadRequest, "invalid input", 4001)

	// Check status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	// Check content type
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", got)
	}

	// Decode response body
	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Validate fields
	if resp.Status != "fail" {
		t.Errorf("expected status 'fail', got %s", resp.Status)
	}
	if resp.Message != "invalid input" {
		t.Errorf("expected message 'invalid input', got %s", resp.Message)
	}
}

package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
)

// --- Mock services ---
type mockAuthService struct{ service1.AuthServiceInterface }
type mockUserService struct{ service1.UserServiceInterface }
type mockProjectService struct{ service1.ProjectServiceInterface }
type mockTaskService struct{ service1.TaskServiceInterface }
type mockCommentService struct{ service1.CommentServiceInterface }

// --- Test SetupRouter ---
func TestSetupRouter(t *testing.T) {
	// create router with mocks
	router := SetupRouter(
		&mockAuthService{},
		&mockUserService{},
		&mockProjectService{},
		&mockTaskService{},
		&mockCommentService{},
	)

	// Test: public route (signup)
	req := httptest.NewRequest(http.MethodPost, "/v1/signup", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Errorf("expected signup route to exist, got 404")
	}

	// Test: public route (login)
	req = httptest.NewRequest(http.MethodPost, "/v1/login", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Errorf("expected login route to exist, got 404")
	}

	// Test: protected route (requires AuthMiddleware)
	req = httptest.NewRequest(http.MethodGet, "/v1/users/", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// since middleware may block without token, expect 401 Unauthorized
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for protected route, got %d", rr.Code)
	}
}
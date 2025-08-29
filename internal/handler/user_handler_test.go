package handler_test

import (
	// "bytes"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/Yash-Watchguard/Tasknest/internal/handler"
	"github.com/Yash-Watchguard/Tasknest/internal/mocks"
	"github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"go.uber.org/mock/gomock"
)

func newUserHandlerWithMock(t *testing.T) (*handler.UserHandler, *mocks.MockUserServiceInterface, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockUserServiceInterface(ctrl)
	h := handler.NewUserHandler(svc)
	return h, svc, ctrl
}
func assertAnError() error {
	return assertErr("mock error")
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
func TestGetuser(t *testing.T) {
	h, mockSvc, ctrl := newUserHandlerWithMock(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		path       string
		ctx        context.Context
		mock       func()
		wantStatus int
	}{
		{
			name:       "missing userId in context -> 401",
			path:       "/v1/users/u2",
			ctx:        context.Background(),
			mock:       func() {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing role in context -> 401",
			path: "/v1/users/u2",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), ContextKey.UserId, "u1")
			}(),
			mock:       func() {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "non-admin trying to view another user -> 403",
			path: "/v1/users/u2",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), ContextKey.UserId, "u1")
				return context.WithValue(ctx, ContextKey.UserRole, roles.Employee)
			}(),
			mock:       func() {},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "single user not found -> 404",
			path: "/v1/users/u1",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), ContextKey.UserId, "u1")
				return context.WithValue(ctx, ContextKey.UserRole, roles.Admin)
			}(),
			mock: func() {
				// MUST return []user.User, not map or other type
				mockSvc.EXPECT().ViewProfile("u1").Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "single user success -> 200",
			path: "/v1/users/u1",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), ContextKey.UserId, "u1")
				return context.WithValue(ctx, ContextKey.UserRole, roles.Admin)
			}(),
			mock: func() {
				// Return correct concrete type: []user.User
				mockSvc.EXPECT().ViewProfile("u1").Return([]user.User{{}}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "list users forbidden for non-admin -> 403",
			path: "/v1/users/", // id empty -> list branch
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), ContextKey.UserId, "u1")
				return context.WithValue(ctx, ContextKey.UserRole, roles.Employee)
			}(),
			mock:       func() {},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "admin list users error -> 404",
			path: "/v1/users/",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), ContextKey.UserId, "u1")
				return context.WithValue(ctx, ContextKey.UserRole, roles.Admin)
			}(),
			mock: func() {
				// MUST return []user.User, error
				mockSvc.EXPECT().ViewAllUsers().Return(nil, errors.New("no users"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "admin list users success -> 200",
			path: "/v1/users/",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), ContextKey.UserId, "u1")
				return context.WithValue(ctx, ContextKey.UserRole, roles.Admin)
			}(),
			mock: func() {
				mockSvc.EXPECT().ViewAllUsers().Return([]user.User{{}, {}}, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil).WithContext(tc.ctx)
			rec := httptest.NewRecorder()

			h.Getuser(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name       string
		pathID     string
		ctxUserID  any
		ctxRole    any
		mockFunc   func(m *mocks.MockUserServiceInterface)
		wantStatus int
	}{
		{
			name:       "missing id in path",
			pathID:     "",
			ctxUserID:  "123",
			ctxRole:    roles.Employee,
			mockFunc:   func(m *mocks.MockUserServiceInterface) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing user id in context",
			pathID:     "123",
			ctxUserID:  nil,
			ctxRole:    roles.Employee,
			mockFunc:   func(m *mocks.MockUserServiceInterface) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "missing role in context",
			pathID:     "123",
			ctxUserID:  "123",
			ctxRole:    nil,
			mockFunc:   func(m *mocks.MockUserServiceInterface) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:      "unauthorized delete attempt by employee",
			pathID:    "456",
			ctxUserID: "123",
			ctxRole:   roles.Employee,
			mockFunc:  func(m *mocks.MockUserServiceInterface) {},
			wantStatus: http.StatusForbidden,
		},
		{
			name:      "service error while deleting user",
			pathID:    "123",
			ctxUserID: "123",
			ctxRole:   roles.Employee,
			mockFunc: func(m *mocks.MockUserServiceInterface) {
				m.EXPECT().DeleteUser("123").Return(assertAnError())
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:      "successful self delete",
			pathID:    "123",
			ctxUserID: "123",
			ctxRole:   roles.Employee,
			mockFunc: func(m *mocks.MockUserServiceInterface) {
				m.EXPECT().DeleteUser("123").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:      "successful admin delete",
			pathID:    "456",
			ctxUserID: "123",
			ctxRole:   roles.Admin,
			mockFunc: func(m *mocks.MockUserServiceInterface) {
				m.EXPECT().DeleteUser("456").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockUserServiceInterface(ctrl)
			tt.mockFunc(mockService)

			uh := handler.NewUserHandler(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/v1/users/"+tt.pathID, nil)
			// Inject context values
			ctx := req.Context()
			if tt.ctxUserID != nil {
				ctx = context.WithValue(ctx, ContextKey.UserId, tt.ctxUserID)
			}
			if tt.ctxRole != nil {
				ctx = context.WithValue(ctx, ContextKey.UserRole, tt.ctxRole)
			}
			req = req.WithContext(ctx)

			// Force set PathValue for compatibility
			req.SetPathValue("id", tt.pathID)

			rr := httptest.NewRecorder()
			uh.DeleteUser(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_PromoteEmployee(t *testing.T) {
	tests := []struct {
		name       string
		role       roles.Role
		pathID     string
		mockFunc   func(m *mocks.MockUserServiceInterface)
		wantStatus int
	}{
		{
			name:       "non-admin cannot promote",
			role:       roles.Employee,
			pathID:     "123",
			mockFunc:   func(m *mocks.MockUserServiceInterface) {}, // no call expected
			wantStatus: http.StatusForbidden,
		},
		{
			name:   "admin promotion fails due to service error",
			role:   roles.Admin,
			pathID: "123",
			mockFunc: func(m *mocks.MockUserServiceInterface) {
				m.EXPECT().PromoteEmployee("123").Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:   "admin promotes successfully",
			role:   roles.Admin,
			pathID: "123",
			mockFunc: func(m *mocks.MockUserServiceInterface) {
				m.EXPECT().PromoteEmployee("123").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, svc, ctrl := newUserHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockFunc(svc)

			req := httptest.NewRequest(http.MethodPut, "/v1/users/"+tt.pathID+"/promote", nil)
			// Inject context with role
			ctx := context.WithValue(req.Context(), ContextKey.UserRole, tt.role)
			req = req.WithContext(ctx)
			// Simulate chi/go1.22 style path param
			req.SetPathValue("id", tt.pathID)

			rr := httptest.NewRecorder()
			h.PromoteEmployee(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name       string
		ctxUserID  string
		ctxRole    roles.Role
		pathID     string
		body       interface{}
		mockFunc   func(m *mocks.MockUserServiceInterface)
		wantStatus int
	}{
		{
			name:       "unauthorized update attempt by employee",
			ctxUserID:  "123",
			ctxRole:    roles.Employee,
			pathID:     "456",
			body:       map[string]interface{}{"name": "new name"},
			mockFunc:   func(m *mocks.MockUserServiceInterface) {},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "invalid JSON body",
			ctxUserID:  "123",
			ctxRole:    roles.Employee,
			pathID:     "123",
			body:       "not-a-json", // will break JSON decoding
			mockFunc:   func(m *mocks.MockUserServiceInterface) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:      "service error during update",
			ctxUserID: "123",
			ctxRole:   roles.Employee,
			pathID:    "123",
			body:      map[string]interface{}{"name": "bad update"},
			mockFunc: func(m *mocks.MockUserServiceInterface) {
				m.EXPECT().UpdateUser("123", gomockAnyMap()).Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:      "successful self update",
			ctxUserID: "123",
			ctxRole:   roles.Employee,
			pathID:    "123",
			body:      map[string]interface{}{"name": "updated name"},
			mockFunc: func(m *mocks.MockUserServiceInterface) {
				m.EXPECT().UpdateUser("123", gomockAnyMap()).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		// {
		// 	name:      "successful admin update",
		// 	ctxUserID: "admin-1",
		// 	ctxRole:   roles.Admin,
		// 	pathID:    "456",
		// 	body:      map[string]interface{}{"email": "new@example.com"},
		// 	mockFunc: func(m *mocks.MockUserServiceInterface) {
		// 		m.EXPECT().UpdateUser("admin-1", gomockAnyMap()).Return(nil)
		// 	},
		// 	wantStatus: http.StatusOK,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, svc, ctrl := newUserHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockFunc(svc)

			// prepare request body
			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				b, _ := json.Marshal(v)
				bodyBytes = b
			}

			req := httptest.NewRequest(http.MethodPut, "/v1/users/"+tt.pathID, bytes.NewReader(bodyBytes))
			// Inject context values
			ctx := context.WithValue(req.Context(), ContextKey.UserId, tt.ctxUserID)
			ctx = context.WithValue(ctx, ContextKey.UserRole, tt.ctxRole)
			req = req.WithContext(ctx)
			// Simulate chi/go1.22 style path param
			req.SetPathValue("id", tt.pathID)

			rr := httptest.NewRecorder()
			h.UpdateUser(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}

func gomockAnyMap() gomock.Matcher {
	return gomock.AssignableToTypeOf(map[string]interface{}{})
}
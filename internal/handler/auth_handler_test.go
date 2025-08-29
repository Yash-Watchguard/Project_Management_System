package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Yash-Watchguard/Tasknest/internal/mocks"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"

	"go.uber.org/mock/gomock"
)
func Setup(ctrl *gomock.Controller) (*authHandler, *mocks.MockUserServiceInterface, *mocks.MockAuthServiceInterface) {
    mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
    mockUserService := mocks.NewMockUserServiceInterface(ctrl)

    authHandler := NewAuthHandler(mockAuthService, mockUserService)
    return authHandler, mockUserService, mockAuthService
}

func makeBody(data interface{}) *bytes.Reader {
	b, _ := json.Marshal(data)
	return bytes.NewReader(b)
}

func TestSignup(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    authHandler, mockUserService, mockAuthService := Setup(ctrl)

    tests := []struct {
        name   string
        method string
        body   interface{}
        mock   func()
        want   int
    }{
        {
            name:   "invalid method",
            method: http.MethodGet,
            body:   nil,
            mock:   func() {},
            want:   http.StatusMethodNotAllowed,
        },
        {
            name:   "invalid json",
            method: http.MethodPost,
            body:   "not-a-json",
            mock:   func() {},
            want:   http.StatusBadRequest,
        },
        {
            name:   "missing required fields",
            method: http.MethodPost,
            body: map[string]string{
                "name":  "yash",
                "email": "yash@gmail.com",
            },
            mock: func() {
                mockUserService.EXPECT().CheckUserExist("yash@gmail.com").Return(false)
            },
            want: http.StatusBadRequest,
        },
        {
            name:   "user already exists",
            method: http.MethodPost,
            body: map[string]string{
                "name":     "yash",
                "email":    "yash@gmail.com",
                "password": "ValidPassword1!",
            },
            mock: func() {
                mockUserService.EXPECT().CheckUserExist("yash@gmail.com").Return(true)
            },
            want: http.StatusConflict,
        },
        {
            name:   "invalid password",
            method: http.MethodPost,
            body: map[string]string{
                "name":     "yash",
                "email":    "yash@gmail.com",
                "password": "short",
            },
            mock: func() {
                mockUserService.EXPECT().CheckUserExist("yash@gmail.com").Return(false)
            },
            want: http.StatusBadRequest,
        },
        {
            name:   "auth service error",
            method: http.MethodPost,
            body: map[string]string{
                "name":     "yash",
                "email":    "yash@gmail.com",
                "password": "ValidPassword1!",
            },
            mock: func() {
                mockUserService.EXPECT().CheckUserExist("yash@gmail.com").Return(false)
                mockAuthService.EXPECT().Signup(gomock.Any()).Return(errors.New("internal error"))
            },
            want: http.StatusInternalServerError,
        },
        {
            name:   "success",
            method: http.MethodPost,
            body: map[string]string{
                "name":     "yash",
                "email":    "yash@gmail.com",
                "password": "ValidPassword1!",
            },
            mock: func() {
                mockUserService.EXPECT().CheckUserExist("yash@gmail.com").Return(false)
                mockAuthService.EXPECT().Signup(gomock.Any()).Return(nil)
            },
            want: http.StatusCreated,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            var req *http.Request
            switch v := tc.body.(type) {
            case string:
                req = httptest.NewRequest(tc.method, "/signup", bytes.NewBufferString(v))
            case nil:
                req = httptest.NewRequest(tc.method, "/signup", nil)
            default:
                req = httptest.NewRequest(tc.method, "/signup", makeBody(v))
            }

            tc.mock()
            w := httptest.NewRecorder()
            authHandler.Signup(w, req)

            if w.Code != tc.want {
                t.Fatalf("%s: expected status %d, got %d; body=%s", tc.name, tc.want, w.Code, w.Body.String())
            }
        })
    }
}


func TestLogin(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    authHandler, mockUserService,_ := Setup(ctrl)

    
    origGen := GenerateJwt
    defer func() { GenerateJwt = origGen }()

    tests := []struct {
        name   string
        method string
        body   interface{}
        mock   func()
        want   int
    }{
        {
            name:   "invalid method",
            method: http.MethodPost,
            body:   nil,
            mock:   func() {},
            want:   http.StatusMethodNotAllowed,
        },
        {
            name:   "invalid json",
            method: http.MethodGet,
            body:   "not-a-json",
            mock:   func() {},
            want:   http.StatusBadRequest,
        },
        {
            name:   "validation error",
            method: http.MethodGet,
            body: map[string]string{
                "name": "yash",
                "email": "bad-email",
                "password": "short",
            },
            mock: func() {},
            want: http.StatusBadRequest,
        },
        {
            name:   "user not present",
            method: http.MethodGet,
            body: map[string]string{
                "name": "yash",
                "email": "yash@gmail.com",
                "password": "ValidPassword1!",
            },
            mock: func() {
                mockUserService.EXPECT().IsUserPresent("yash", "yash@gmail.com", "ValidPassword1!").Return(nil, errors.New("not found"))
            },
            want: http.StatusUnauthorized,
        },
        {
            name:   "jwt generation error",
            method: http.MethodGet,
            body: map[string]string{
                "name": "yash",
                "email": "yash@gmail.com",
                "password": "ValidPassword1!",
            },
            mock: func() {
                mockUserService.EXPECT().IsUserPresent("yash", "yash@gmail.com", "ValidPassword1!").Return(&user.User{Id: "uid123", Role: roles.Employee}, nil)
                // override GenerateJwt to return error
                GenerateJwt = func(userId string, role roles.Role) (string, error) { return "", errors.New("jwt error") }
            },
            want: http.StatusInternalServerError,
        },
        {
            name:   "success",
            method: http.MethodGet,
            body: map[string]string{
                "name": "yash",
                "email": "yash@gmail.com",
                "password": "ValidPassword1!",
            },
            mock: func() {
                mockUserService.EXPECT().IsUserPresent("yash", "yash@gmail.com", "ValidPassword1!").Return(&user.User{Id: "uid123", Role: roles.Employee}, nil)
                // valid token
                GenerateJwt = func(userId string, role roles.Role) (string, error) { return "token123", nil }
            },
            want: http.StatusCreated,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            var req *http.Request
            switch v := tc.body.(type) {
            case string:
                req = httptest.NewRequest(tc.method, "/login", bytes.NewBufferString(v))
            case nil:
                req = httptest.NewRequest(tc.method, "/login", nil)
            default:
                req = httptest.NewRequest(tc.method, "/login", makeBody(v))
            }

            tc.mock()
            w := httptest.NewRecorder()
            authHandler.Login(w, req)

            if w.Code != tc.want {
                t.Fatalf("%s: expected status %d, got %d; body=%s", tc.name, tc.want, w.Code, w.Body.String())
            }
        })
    }
}





    






    




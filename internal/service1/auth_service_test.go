package service1

import (
	"testing"

	"github.com/Yash-Watchguard/Tasknest/internal/mocks"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"go.uber.org/mock/gomock"
)

func TestSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authService := NewAuthService(mockRepo)

	u := &user.User{Id: "1", Name: "Yash", Email: "yash@test.com",Role: roles.Employee}

	mockRepo.EXPECT().SaveUser(u).Return(nil)

	if err := authService.Signup(u); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authService := NewAuthService(mockRepo)

	// expectation
	mockRepo.EXPECT().
		IsUserPresent("Yash", "yash@test.com", "1234").
		Return(&user.User{Id: "1", Name: "Yash", Email: "yash@test.com"}, nil)

	u, err := authService.Login("Yash", "yash@test.com", "1234")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u == nil || u.Name != "Yash" {
		t.Errorf("expected user Yash, got %v", u)
	}
}
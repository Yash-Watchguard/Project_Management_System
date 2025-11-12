package service1

import (
	"context"
	"errors"
	"testing"

	"github.com/Yash-Watchguard/Tasknest/internal/mocks"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"go.uber.org/mock/gomock"
)



func TestViewProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	expectedUsers := []user.User{
		{Id: "1", Name: "Yash", Email: "yash@test.com", Role: roles.Employee},
	}

	// Expect repo call
	mockRepo.EXPECT().ViewProfile("25185bed-4bc2-4dd9-9b40-6de2ef889c18").Return(expectedUsers, nil)

	// Call service
	users, err := userService.ViewProfile("25185bed-4bc2-4dd9-9b40-6de2ef889c18")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 1 || users[0].Name != "Yash" {
		t.Errorf("expected user Yash, got %v", users)
	}
}

func TestUserService_ViewAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	expectedUsers := []user.User{
		{Id: "1", Name: "Yash", Email: "yash@test.com", Role: roles.Manager},
		{Id: "2", Name: "Raj", Email: "raj@test.com", Role: roles.Employee},
	}

	// Expect repo call
	mockRepo.EXPECT().GetAllUsers().Return(expectedUsers, nil)

	// Call service
	users, err := userService.ViewAllUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 2 || users[1].Name != "Raj" {
		t.Errorf("expected 2 users, got %v", users)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)
	mp:=make(map[string]interface{})
	mp["status"]="InActive"

	userID := "123"

	// Expectation: repo should be called with the given userID
	mockRepo.EXPECT().UpdateProfile(userID,mp).Return(nil)

	// Call service
	err := userService.DeleteUser(userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUserService_GetAllManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	tests := []struct {
		name      string
		role      roles.Role
		setupMock func()
		wantErr   bool
	}{
		{
			name: "Authorized - returns managers",
			role: roles.Admin, // Admin role is allowed (0)
			setupMock: func() {
				mockRepo.EXPECT().
					GetAllManager().
					Return([]user.User{{Id: "1", Name: "Manager1"}}, nil)
			},
			wantErr: false,
		},
		{
			name:      "Unauthorized - returns error",
			role:      roles.Employee, // Non-admin
			setupMock: func() {},     // no repo call expected
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), ContextKey.UserRole, tt.role)

			tt.setupMock()

			users, err := userService.GetAllManager(ctx)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(users) == 0 || users[0].Name != "Manager1" {
					t.Errorf("expected Manager1, got %v", users)
				}
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	existingUser := &user.User{
		Id:    "1",
		Name:  "OldName",
		Email: "old@test.com",
	}

	tests := []struct {
		name      string
		id        string
		updates   map[string]interface{}
		setupMock func()
		wantErr   bool
	}{
		{
			name:    "Update name successfully",
			id:      "1",
			updates: map[string]interface{}{"name": "NewName"},
			setupMock: func() {
				mockRepo.EXPECT().
					ViewProfile("1").
					Return([]user.User{*existingUser}, nil)

				mockRepo.EXPECT().
					UpdateProfile("1", gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "Update email - invalid format",
			id:      "1",
			updates: map[string]interface{}{"email": "invalidEmail"},
			setupMock: func() {
				mockRepo.EXPECT().
					ViewProfile("1").
					Return([]user.User{*existingUser}, nil)
			},
			wantErr: true,
		},
		{
			name:    "Update email - already exists",
			id:      "1",
			updates: map[string]interface{}{"email": "exists@test.com"},
			setupMock: func() {
				mockRepo.EXPECT().
					ViewProfile("1").
					Return([]user.User{*existingUser}, nil)

				mockRepo.EXPECT().
					GetUserByEmail("exists@test.com").
					Return(&user.User{Id: "2", Email: "exists@test.com"}, nil)
			},
			wantErr: true,
		},
		{
			name:    "Update password - too weak",
			id:      "1",
			updates: map[string]interface{}{"password": "123"},
			setupMock: func() {
				mockRepo.EXPECT().
					ViewProfile("1").
					Return([]user.User{*existingUser}, nil)
			},
			wantErr: true,
		},
		{
			name:    "Update phone number successfully",
			id:      "1",
			updates: map[string]interface{}{"phoneNumber": "1234567890"},
			setupMock: func() {
				mockRepo.EXPECT().
					ViewProfile("1").
					Return([]user.User{*existingUser}, nil)

				mockRepo.EXPECT().
					UpdateProfile("1", gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "No valid fields to update",
			id:      "1",
			updates: map[string]interface{}{},
			setupMock: func() {
				mockRepo.EXPECT().
					ViewProfile("1").
					Return([]user.User{*existingUser}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := service.UpdateUser(tt.id, tt.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_PromoteEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	tests := []struct {
		name      string
		employeeId string
		setupMock func()
		wantErr   bool
	}{
		{
			name:      "Promote employee successfully",
			employeeId: "emp1",
			setupMock: func() {
				mockRepo.EXPECT().
					PromoteEmployee("emp1").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Promote employee fails",
			employeeId: "emp2",
			setupMock: func() {
				mockRepo.EXPECT().
					PromoteEmployee("emp2").
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := service.PromoteEmployee(tt.employeeId)
			if (err != nil) != tt.wantErr {
				t.Errorf("PromoteEmployee() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_ViewAllEmplpyee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	tests := []struct {
		name         string
		role         roles.Role
		mockSetup    func()
		expectedErr  string
		expectedUserCount int
	}{
		{
			name: "unauthorized role",
			role: 2,
			mockSetup: func() {},
			expectedErr: "unauthorized access",
			expectedUserCount: 0,
		},
		{
			name: "success case",
			role: 0,
			mockSetup: func() {
				mockRepo.EXPECT().
					ViewAllEmployee().
					Return([]user.User{
						{Id: "1", Name: "Yash", Email: "yash@test.com", Role: roles.Employee},
					}, nil)
			},
			expectedErr: "",
			expectedUserCount: 1,
		},
		{
			name: "repo failure",
			role: 0,
			mockSetup: func() {
				mockRepo.EXPECT().
					ViewAllEmployee().
					Return(nil, errors.New("db error"))
			},
			expectedErr: "db error",
			expectedUserCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), ContextKey.UserRole, tt.role)

			tt.mockSetup()
			users, err := service.ViewAllEmplpyee(ctx)

			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(users) != tt.expectedUserCount {
				t.Errorf("expected %d users, got %d", tt.expectedUserCount, len(users))
			}
		})
	}
}

func TestUserService_CheckUserExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	tests := []struct {
		name        string
		mockSetup   func()
		expectedRes bool
	}{
		{
			name: "user exists",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByEmail("yash@test.com").
					Return(&user.User{Id: "1", Name: "Yash"}, nil)
			},
			expectedRes: true,
		},
		{
			name: "user does not exist (repo returns nil user)",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByEmail("yash@test.com").
					Return(nil, nil)
			},
			expectedRes: false,
		},
		{
			name: "repo returns error",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByEmail("yash@test.com").
					Return(nil, errors.New("db error"))
			},
			expectedRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result := service.CheckUserExist("yash@test.com")

			if result != tt.expectedRes {
				t.Errorf("expected %v, got %v", tt.expectedRes, result)
			}
		})
	}
}

func TestUserService_IsUserPresent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	tests := []struct {
		name        string
		mockSetup   func()
		expectedErr bool
		expectedNil bool
	}{
		{
			name: "user found",
			mockSetup: func() {
				mockRepo.EXPECT().
					IsUserPresent("Yash", "yash@test.com", "1234").
					Return(&user.User{Id: "1", Name: "Yash"}, nil)
			},
			expectedErr: false,
			expectedNil: false,
		},
		{
			name: "user not found (repo error)",
			mockSetup: func() {
				mockRepo.EXPECT().
					IsUserPresent("Yash", "yash@test.com", "1234").
					Return(nil, errors.New("not found"))
			},
			expectedErr: true,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			u, err := service.IsUserPresent("yash","yash@test.com", "1234")

			if tt.expectedErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectedNil && u != nil {
				t.Errorf("expected nil user, got %+v", u)
			}
		})
	}
}







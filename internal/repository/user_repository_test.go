package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"

	// "github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	// "github.com/Yash-Watchguard/Tasknest/internal/repository"
	// "github.com/Yash-Watchguard/Tasknest/internal/repository"
)

// sqlDB is a small wrapper to satisfy UserRepoInterface using *sql.DB

func TestSaveUser(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock, u *user.User)
		user        *user.User
		expectedErr string
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock, u *user.User) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(u.Id, u.Role, u.Name, u.Password, u.PhoneNumber, u.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			user: &user.User{
				Id:          "u1",
				Role:        roles.Employee,
				Name:        "Test User",
				Password:    "pass123",
				PhoneNumber: "1234567890",
				Email:       "test@example.com",
			},
			expectedErr: "",
		},
		{
			name: "duplicate phone number",
			setupMock: func(mock sqlmock.Sqlmock, u *user.User) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(u.Id, u.Role, u.Name, u.Password, u.PhoneNumber, u.Email).
					WillReturnError(&mysql.MySQLError{
						Number:  1062,
						Message: "Duplicate entry '1234567890' for key 'phone_number'",
					})
			},
			user: &user.User{
				Id:          "u2",
				Role:        roles.Employee,
				Name:        "User Two",
				Password:    "pass456",
				PhoneNumber: "1234567890",
				Email:       "u2@example.com",
			},
			expectedErr: "phone number already exists",
		},
		{
			name: "duplicate email",
			setupMock: func(mock sqlmock.Sqlmock, u *user.User) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(u.Id, u.Role, u.Name, u.Password, u.PhoneNumber, u.Email).
					WillReturnError(&mysql.MySQLError{
						Number:  1062,
						Message: "Duplicate entry 'u3@example.com' for key 'email'",
					})
			},
			user: &user.User{
				Id:          "u3",
				Role:        roles.Employee,
				Name:        "User Three",
				Password:    "pass789",
				PhoneNumber: "9876543210",
				Email:       "u3@example.com",
			},
			expectedErr: "email already exists",
		},
		{
			name: "duplicate other field",
			setupMock: func(mock sqlmock.Sqlmock, u *user.User) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(u.Id, u.Role, u.Name, u.Password, u.PhoneNumber, u.Email).
					WillReturnError(&mysql.MySQLError{
						Number:  1062,
						Message: "Duplicate entry 'someid' for key 'PRIMARY'",
					})
			},
			user: &user.User{
				Id:          "u4",
				Role:        roles.Employee,
				Name:        "User Four",
				Password:    "pass999",
				PhoneNumber: "5555555555",
				Email:       "u4@example.com",
			},
			expectedErr: "duplicate entry",
		},
		{
			name: "unexpected db error",
			setupMock: func(mock sqlmock.Sqlmock, u *user.User) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(u.Id, u.Role, u.Name, u.Password, u.PhoneNumber, u.Email).
					WillReturnError(errors.New("db connection lost"))
			},
			user: &user.User{
				Id:          "u5",
				Role:        roles.Employee,
				Name:        "User Five",
				Password:    "pass111",
				PhoneNumber: "6666666666",
				Email:       "u5@example.com",
			},
			expectedErr: "db connection lost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewUserRepo(db)

			tt.setupMock(mock, tt.user)

			err = repo.SaveUser(tt.user)

			if tt.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			} else {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}
			
func TestIsUserPresent(t *testing.T) {
	query := regexp.QuoteMeta("SELECT id, name, email, password, role, phone_number FROM users WHERE name = ? AND email = ?")

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		password    string
		expectedErr string
	}{
		{
			name: "user not found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs("yash", "yash@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			password:    "pass",
			expectedErr: "no user found",
		},
	
		{
			name: "invalid password",
			setupMock: func(mock sqlmock.Sqlmock) {
				hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "phone_number"}).
					AddRow("u1", "yash", "yash@example.com", string(hashed), roles.Employee, "12345")
				mock.ExpectQuery(query).
					WithArgs("yash", "yash@example.com").
					WillReturnRows(rows)
			},
			password:    "wrongpass",
			expectedErr: "invalid password",
		},
		{
			name: "successfully found user",
			setupMock: func(mock sqlmock.Sqlmock) {
				hashed, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "phone_number"}).
					AddRow("u2", "yash", "yash@example.com", string(hashed), roles.Manager, "67890")
				mock.ExpectQuery(query).
					WithArgs("yash", "yash@example.com").
					WillReturnRows(rows)
			},
			password:    "mypassword",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewUserRepo(db)

			tt.setupMock(mock)

			u, err := repo.IsUserPresent("yash", "yash@example.com", tt.password)

			if tt.expectedErr == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if u == nil {
					t.Errorf("expected user, got nil")
				}
			} else {
				if err == nil || err.Error() != tt.expectedErr {
					t.Errorf("expected error %q, got %v", tt.expectedErr, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestViewProfile(t *testing.T) {
	query := regexp.QuoteMeta("SELECT id, name, email, role, phone_number FROM users WHERE id = ?")
	tests := []struct {
		name          string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedUser  []user.User
		expectedError error
	}{
		{
			name: "user not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs("u1").
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
		{
			name: "user found successfully",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "phone_number"}).
					AddRow("u1", "yash", "yash@example.com", roles.Employee, "1234567890")

				mock.ExpectQuery(query).
					WithArgs("u1").
					WillReturnRows(rows)
			},
			expectedUser: []user.User{
				{
					Id:          "u1",
					Name:        "yash",
					Email:       "yash@example.com",
					Role:        roles.Employee,
					PhoneNumber: "1234567890",
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock db: %v", err)
			}
			defer db.Close()

			repo := NewUserRepo(db)

			tt.mockSetup(mock)

			result, err := repo.ViewProfile("u1")

			// check error
			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}

			// check result
			if len(result) != len(tt.expectedUser) {
				t.Errorf("expected %d users, got %d", len(tt.expectedUser), len(result))
			} else {
				for i := range result {
					if result[i] != tt.expectedUser[i] {
						t.Errorf("expected user %+v, got %+v", tt.expectedUser[i], result[i])
					}
				}
			}

			// check expectations
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	query := regexp.QuoteMeta("SELECT id, name, email, role, phone_number FROM users")

	tests := []struct {
		name         string
		setupMock    func(mock sqlmock.Sqlmock)
		expected     []user.User
		expectedErr  string // empty => no error
	}{
		{
			name: "query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).WillReturnError(errors.New("db down"))
			},
			expected:    nil,
			expectedErr: "db down",
		},
		{
			name: "no users present",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "phone_number"})
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: "no users present",
		},
		{
			name: "success multiple users",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "phone_number"}).
					AddRow("u1", "Alice", "alice@example.com", roles.Employee, "11111").
					AddRow("u2", "Bob", "bob@example.com", roles.Manager, "22222")
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expected: []user.User{
				{Id: "u1", Name: "Alice", Email: "alice@example.com", Role: roles.Employee, PhoneNumber: "11111"},
				{Id: "u2", Name: "Bob", Email: "bob@example.com", Role: roles.Manager, PhoneNumber: "22222"},
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock db: %v", err)
			}
			defer db.Close()

			repo := NewUserRepo(db)

			tt.setupMock(mock)

			got, err := repo.GetAllUsers()

			if tt.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(got) != len(tt.expected) {
					t.Fatalf("expected %d users, got %d", len(tt.expected), len(got))
				}
				for i := range got {
					if got[i] != tt.expected[i] {
						t.Errorf("expected user %+v, got %+v", tt.expected[i], got[i])
					}
				}
			} else {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestDeleteUserById(t *testing.T) {
	query := regexp.QuoteMeta("DELETE FROM users WHERE id = ?")

	tests := []struct {
		name        string
		userId      string
		setupMock   func(mock sqlmock.Sqlmock)
		expectedErr string // empty => success
	}{
		{
			name:   "successful delete",
			userId: "u1",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs("u1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: "",
		},
		{
			name:   "exec error",
			userId: "bad",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs("bad").
					WillReturnError(errors.New("exec failed"))
			},
			expectedErr: "please enter valid user id",
		},
		{
			name:   "no rows affected",
			userId: "u404",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs("u404").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: "please enter valid user id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock db: %v", err)
			}
			defer db.Close()

			repo := NewUserRepo(db)

			tt.setupMock(mock)

			err = repo.DeleteUserById(tt.userId)

			if tt.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			} else {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPromoteEmployee(t *testing.T) {
	query := regexp.QuoteMeta("UPDATE users SET role = ? WHERE id = ?")

	tests := []struct {
		name        string
		employeeId  string
		setupMock   func(mock sqlmock.Sqlmock)
		expectedErr string // empty => success
	}{
		{
			name:       "successful promotion",
			employeeId: "emp1",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(roles.Manager, "emp1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: "",
		},
		{
			name:       "exec error",
			employeeId: "bad",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(roles.Manager, "bad").
					WillReturnError(errors.New("exec failed"))
			},
			expectedErr: "exec failed",
		},
		{
			name:       "no rows affected",
			employeeId: "emp404",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(roles.Manager, "emp404").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: "user not found",
		},
		{
			name:       "rows affected error",
			employeeId: "emp2",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(roles.Manager, "emp2").
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected failed")))
			},
			expectedErr: "rows affected failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock db: %v", err)
			}
			defer db.Close()

			repo := NewUserRepo(db)

			tt.setupMock(mock)

			err = repo.PromoteEmployee(tt.employeeId)

			if tt.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			} else {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}


func TestViewAllEmployee(t *testing.T) {
	query := regexp.QuoteMeta(
		"SELECT id, name, email, role, phone_number, password FROM users WHERE role = ?",
	)

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		expectedErr string
		expectUsers bool
	}{
		{
			name: "query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs(roles.Employee).
					WillReturnError(errors.New("query failed"))
			},
			expectedErr: "query failed",
			expectUsers: false,
		},
		{
			name: "scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "phone_number", "password"}).
					AddRow("e1", "Alice", "alice@mail.com", roles.Employee, 12345, nil) // password nil → scan fails
				mock.ExpectQuery(query).
					WithArgs(roles.Employee).
					WillReturnRows(rows)
			},
			expectedErr: "sql: Scan error",
			expectUsers: false,
		},
		{
			name: "no employees found",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "phone_number", "password"})
				mock.ExpectQuery(query).
					WithArgs(roles.Employee).
					WillReturnRows(rows)
			},
			expectedErr: "no employees found",
			expectUsers: false,
		},
		{
			name: "successfully fetched employees",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "phone_number", "password"}).
					AddRow("e1", "Alice", "alice@mail.com", roles.Employee, "12345", "hashedpass").
					AddRow("e2", "Bob", "bob@mail.com", roles.Employee, "67890", "hashedpass2")
				mock.ExpectQuery(query).
					WithArgs(roles.Employee).
					WillReturnRows(rows)
			},
			expectedErr: "",
			expectUsers: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock db: %v", err)
			}
			defer db.Close()

			repo := NewUserRepo(db)
			tt.setupMock(mock)

			users, err := repo.ViewAllEmployee()

			if tt.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if !tt.expectUsers || len(users) == 0 {
					t.Fatalf("expected users, got none")
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	query := regexp.QuoteMeta(
		"SELECT id, name, email, role, phone_number FROM users WHERE email = ?",
	)

	tests := []struct {
		name        string
		email       string
		setupMock   func(mock sqlmock.Sqlmock)
		expectedErr string
		expectUser  bool
	}{
		{
			name:  "user not found",
			email: "nouser@mail.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs("nouser@mail.com").
					WillReturnError(sql.ErrNoRows)
			},
			expectedErr: "user not found",
			expectUser:  false,
		},
		{
			name:  "scan error",
			email: "bad@mail.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				// mismatched column type → scan fails
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "phone_number"}).
					AddRow(1, "Alice", "bad@mail.com", roles.Employee, nil)
				mock.ExpectQuery(query).
					WithArgs("bad@mail.com").
					WillReturnRows(rows)
			},
			expectedErr: "sql: Scan error",
			expectUser:  false,
		},
		{
			name:  "successfully fetched user",
			email: "alice@mail.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "phone_number"}).
					AddRow("u1", "Alice", "alice@mail.com", roles.Employee, "12345")
				mock.ExpectQuery(query).
					WithArgs("alice@mail.com").
					WillReturnRows(rows)
			},
			expectedErr: "",
			expectUser:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock db: %v", err)
			}
			defer db.Close()

			repo := NewUserRepo(db)
			tt.setupMock(mock)

			user, err := repo.GetUserByEmail(tt.email)

			if tt.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if !tt.expectUser || user == nil {
					t.Fatalf("expected user, got nil")
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}


func TestUpdateProfile(t *testing.T) {
	tests := []struct {
		name        string
		userId      string
		updates     map[string]interface{}
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name:   "success - update name and email",
			userId: "u1",
			updates: map[string]interface{}{
				"name":  "New Name",
				"email": "new@example.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET name = \\?, email = \\? WHERE id = \\?").
					WithArgs("New Name", "new@example.com", "u1").
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
			},
			expectError: false,
		},
		{
			name:   "failure - invalid field",
			userId: "u2",
			updates: map[string]interface{}{
				"invalid_field": "oops",
			},
			setupMock:   func(mock sqlmock.Sqlmock) {}, // no DB call expected
			expectError: true,
			errorMsg:    "invalid field update: invalid_field",
		},
		{
			name:        "failure - no updates",
			userId:      "u3",
			updates:     map[string]interface{}{},
			setupMock:   func(mock sqlmock.Sqlmock) {}, // no DB call expected
			expectError: true,
			errorMsg:    "no valid fields to update",
		},
		{
			name:   "failure - duplicate email",
			userId: "u4",
			updates: map[string]interface{}{
				"email": "dup@example.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET email = \\? WHERE id = \\?").
					WithArgs("dup@example.com", "u4").
					WillReturnError(errors.New("Duplicate entry 'dup@example.com' for key 'email'"))
			},
			expectError: true,
			errorMsg:    "email already exists",
		},
		{
			name:   "failure - duplicate phone_number",
			userId: "u5",
			updates: map[string]interface{}{
				"phone_number": "1234567890",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET phone_number = \\? WHERE id = \\?").
					WithArgs("1234567890", "u5").
					WillReturnError(errors.New("Duplicate entry '1234567890' for key 'phone_number'"))
			},
			expectError: true,
			errorMsg:    "phone number already exists",
		},
		{
			name:   "failure - user not found",
			userId: "u6",
			updates: map[string]interface{}{
				"name": "NotFound",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET name = \\? WHERE id = \\?").
					WithArgs("NotFound", "u6").
					WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
			},
			expectError: true,
			errorMsg:    "user not found",
		},
		{
			name:   "failure - db error",
			userId: "u7",
			updates: map[string]interface{}{
				"name": "Broken",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET name = \\? WHERE id = \\?").
					WithArgs("Broken", "u7").
					WillReturnError(errors.New("db failure"))
			},
			expectError: true,
			errorMsg:    "db failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewUserRepo(db)

			tt.setupMock(mock)

			err = repo.UpdateProfile(tt.userId, tt.updates)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetAllManager(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		expected    []user.User
		expectError bool
	}{
		{
			name: "success - multiple managers",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("m1", "Manager One").
					AddRow("m2", "Manager Two")

				mock.ExpectQuery("SELECT id, name FROM users WHERE role = ?").
					WithArgs(roles.Manager).
					WillReturnRows(rows)
			},
			expected: []user.User{
				{Id: "m1", Name: "Manager One"},
				{Id: "m2", Name: "Manager Two"},
			},
			expectError: false,
		},
		{
			name: "success - no managers",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery("SELECT id, name FROM users WHERE role = ?").
					WithArgs(roles.Manager).
					WillReturnRows(rows)
			},
			expected:    []user.User{}, // should just return empty slice
			expectError: false,
		},
		{
			name: "failure - db query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name FROM users WHERE role = ?").
					WithArgs(roles.Manager).
					WillReturnError(errors.New("db query failed"))
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "failure - row scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				// return wrong column types (e.g. int instead of string)
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(123, nil)

				mock.ExpectQuery("SELECT id, name FROM users WHERE role = ?").
					WithArgs(roles.Manager).
					WillReturnRows(rows)
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewUserRepo(db)

			tt.setupMock(mock)

			managers, err := repo.GetAllManager()

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(managers) != len(tt.expected) {
					t.Fatalf("expected %d managers, got %d", len(tt.expected), len(managers))
				}
				for i := range managers {
					if managers[i].Id != tt.expected[i].Id || managers[i].Name != tt.expected[i].Name {
						t.Fatalf("expected %+v, got %+v", tt.expected[i], managers[i])
					}
				}
			}

			// check expectations
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}




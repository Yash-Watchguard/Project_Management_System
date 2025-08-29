package repository

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
)

func TestAddProject(t *testing.T) {
	tests := []struct {
		name        string
		project     project.Project
		mockSetup   func(mock sqlmock.Sqlmock, proj project.Project)
		expectError string
	}{
		{
			name: "success - project inserted",
			project: project.Project{
				ProjectId:          "p1",
				ProjectName:        "Test Project",
				ProjectDescription: "Test Description",
				Deadline:           time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
				CreatedBy:          "admin1",
				AssignedManager:    "mgr1",
			},
			mockSetup: func(mock sqlmock.Sqlmock, proj project.Project) {
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO projects (project_id, project_name, project_description, deadline, created_by, assigned_manager_id) VALUES (?, ?, ?, ?, ?, ?)`)).
					WithArgs(proj.ProjectId, proj.ProjectName, proj.ProjectDescription, proj.Deadline, proj.CreatedBy, proj.AssignedManager).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: "",
		},
		{
			name: "failure - db error",
			project: project.Project{
				ProjectId:          "p2",
				ProjectName:        "Fail Project",
				ProjectDescription: "Desc",
				Deadline:           time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
				CreatedBy:          "admin2",
				AssignedManager:    "mgr2",
			},
			mockSetup: func(mock sqlmock.Sqlmock, proj project.Project) {
				mock.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO projects (project_id, project_name, project_description, deadline, created_by, assigned_manager_id) VALUES (?, ?, ?, ?, ?, ?)`)).
					WithArgs(proj.ProjectId, proj.ProjectName, proj.ProjectDescription, proj.Deadline, proj.CreatedBy, proj.AssignedManager).
					WillReturnError(sqlmock.ErrCancelled)
			},
			expectError: "canceling query due to user request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo := NewProjectRepo(db)

			tt.mockSetup(mock, tt.project)

			err := repo.AddProject(tt.project)

			if tt.expectError == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.expectError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectError) {
					t.Fatalf("expected error containing %q, got %v", tt.expectError, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}


func TestViewAllProjects(t *testing.T) {
	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockError   error
		expectError bool
		expectCount int
	}{
		{
			name: "success - return projects",
			mockRows: sqlmock.NewRows([]string{
				"project_id", "project_name", "project_description", "deadline", "created_by", "assigned_manager_id",
			}).
				AddRow("p1", "Project One", "First project", "2025-12-31", "admin1", "mgr1").
				AddRow("p2", "Project Two", "Second project", "2026-01-15", "admin2", "mgr2"),
			expectError: false,
			expectCount: 2,
		},
		{
			name:        "failure - db error",
			mockRows:    nil,
			mockError:   sqlmock.ErrCancelled,
			expectError: true,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewProjectRepo(db)

			if tt.mockError != nil {
				mock.ExpectQuery("SELECT project_id, project_name, project_description, deadline, created_by, assigned_manager_id FROM projects").
					WillReturnError(tt.mockError)
			} else {
				mock.ExpectQuery("SELECT project_id, project_name, project_description, deadline, created_by, assigned_manager_id FROM projects").
					WillReturnRows(tt.mockRows)
			}

			projects, err := repo.ViewAllProjects()

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(projects) != tt.expectCount {
				t.Fatalf("expected %d projects, got %d", tt.expectCount, len(projects))
			}

			// Verify one sample project parsing
			if tt.expectCount > 0 {
				expectedDeadline, _ := time.Parse("2006-01-02", "2025-12-31")
				if projects[0].Deadline != expectedDeadline {
					t.Errorf("expected deadline %v, got %v", expectedDeadline, projects[0].Deadline)
				}
			}
		})
	}
}


func TestDeleteProject(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "success - project deleted",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Check project existence -> exists
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM projects WHERE project_id = \\?\\)").
					WithArgs("p1").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))

				// Delete tasks
				mock.ExpectExec("DELETE FROM tasks WHERE projectid = \\?").
					WithArgs("p1").
					WillReturnResult(sqlmock.NewResult(0, 2)) // 2 tasks deleted

				// Delete project
				mock.ExpectExec("DELETE FROM projects WHERE project_id = \\?").
					WithArgs("p1").
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row deleted
			},
			expectError: false,
		},
		{
			name: "failure - project not found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM projects WHERE project_id = \\?\\)").
					WithArgs("p1").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(0))
			},
			expectError: true,
			errorMsg:    "project not found",
		},
		{
			name: "failure - error checking existence",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM projects WHERE project_id = \\?\\)").
					WithArgs("p1").
					WillReturnError(errors.New("db error"))
			},
			expectError: true,
			errorMsg:    "failed to check project existence",
		},
		{
			name: "failure - error deleting tasks",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM projects WHERE project_id = \\?\\)").
					WithArgs("p1").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))

				mock.ExpectExec("DELETE FROM tasks WHERE projectid = \\?").
					WithArgs("p1").
					WillReturnError(errors.New("db error"))
			},
			expectError: true,
			errorMsg:    "failed to delete tasks for the project",
		},
		{
			name: "failure - error deleting project",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM projects WHERE project_id = \\?\\)").
					WithArgs("p1").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))

				mock.ExpectExec("DELETE FROM tasks WHERE projectid = \\?").
					WithArgs("p1").
					WillReturnResult(sqlmock.NewResult(0, 2))

				mock.ExpectExec("DELETE FROM projects WHERE project_id = \\?").
					WithArgs("p1").
					WillReturnError(errors.New("db error"))
			},
			expectError: true,
			errorMsg:    "failed to delete project",
		},
		{
			name: "failure - no project deleted",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM projects WHERE project_id = \\?\\)").
					WithArgs("p1").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))

				mock.ExpectExec("DELETE FROM tasks WHERE projectid = \\?").
					WithArgs("p1").
					WillReturnResult(sqlmock.NewResult(0, 2))

				mock.ExpectExec("DELETE FROM projects WHERE project_id = \\?").
					WithArgs("p1").
					WillReturnResult(sqlmock.NewResult(0, 0)) // no rows affected
			},
			expectError: true,
			errorMsg:    "no project deleted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewProjectRepo(db)

			tt.setupMock(mock)

			err = repo.DeleteProject("p1")
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Fatalf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestViewAssignedProject(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		expected    []project.Project
		expectError bool
	}{
		{
			name: "success - assigned projects found",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"project_id", "project_name", "project_description", "deadline", "created_by", "assigned_manager_id",
				}).
					AddRow("p1", "Project One", "Description One", "2025-08-30", "admin1", "mgr1").
					AddRow("p2", "Project Two", "Description Two", "2025-09-15", "admin2", "mgr1")

				mock.ExpectQuery("SELECT project_id, project_name, project_description, deadline, created_by, assigned_manager_id FROM projects WHERE assigned_manager_id = \\?").
					WithArgs("mgr1").
					WillReturnRows(rows)
			},
			expected: []project.Project{
				{
					ProjectId:          "p1",
					ProjectName:        "Project One",
					ProjectDescription: "Description One",
					Deadline:           time.Date(2025, 8, 30, 0, 0, 0, 0, time.UTC),
					CreatedBy:          "admin1",
					AssignedManager:    "mgr1",
				},
				{
					ProjectId:          "p2",
					ProjectName:        "Project Two",
					ProjectDescription: "Description Two",
					Deadline:           time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
					CreatedBy:          "admin2",
					AssignedManager:    "mgr1",
				},
			},
			expectError: false,
		},
		{
			name: "failure - db query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT project_id, project_name, project_description, deadline, created_by, assigned_manager_id FROM projects WHERE assigned_manager_id = \\?").
					WithArgs("mgr1").
					WillReturnError(sqlmock.ErrCancelled)
			},
			expectError: true,
		},
		{
			name: "failure - invalid deadline format",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"project_id", "project_name", "project_description", "deadline", "created_by", "assigned_manager_id",
				}).
					AddRow("p1", "Project One", "Description One", "30-08-2025", "admin1", "mgr1") // invalid format

				mock.ExpectQuery("SELECT project_id, project_name, project_description, deadline, created_by, assigned_manager_id FROM projects WHERE assigned_manager_id = \\?").
					WithArgs("mgr1").
					WillReturnRows(rows)
			},
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

			repo := NewProjectRepo(db)

			tt.setupMock(mock)

			projects, err := repo.ViewAssignedProject("mgr1")

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(projects) != len(tt.expected) {
					t.Fatalf("expected %d projects, got %d", len(tt.expected), len(projects))
				}
				for i := range tt.expected {
					if projects[i].ProjectId != tt.expected[i].ProjectId ||
						projects[i].ProjectName != tt.expected[i].ProjectName ||
						projects[i].ProjectDescription != tt.expected[i].ProjectDescription ||
						!projects[i].Deadline.Equal(tt.expected[i].Deadline) ||
						projects[i].CreatedBy != tt.expected[i].CreatedBy ||
						projects[i].AssignedManager != tt.expected[i].AssignedManager {
						t.Errorf("expected %+v, got %+v", tt.expected[i], projects[i])
					}
				}
			}
		})
	}
}

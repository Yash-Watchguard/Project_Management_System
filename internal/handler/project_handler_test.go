package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// "github.com/Yash-Watchguard/Tasknest/internal/handler"
	"github.com/Yash-Watchguard/Tasknest/internal/mocks"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"

	// "github.com/Yash-Watchguard/Tasknest/internal/model/task"
	// status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"go.uber.org/mock/gomock"
)

func TestProjectStatus(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		role         roles.Role
		mockSetup    func(svc *mocks.MockTaskServiceInterface)
		expectedCode int
		expectBody   string // substring check
	}{
		{
			name:         "Invalid project ID length",
			projectID:    "p1",
			role:         roles.Admin,
			mockSetup:    func(svc *mocks.MockTaskServiceInterface) {},
			expectedCode: http.StatusBadRequest,
			expectBody:   "Invalid Project Id",
		},
		{
			name:         "Employee forbidden",
			projectID:    "12345678-1234-1234-1234-123456789012", // valid UUID length (36)
			role:         roles.Employee,
			mockSetup:    func(svc *mocks.MockTaskServiceInterface) {},
			expectedCode: http.StatusForbidden,
			expectBody:   "Access denied",
		},
		{
			name:      "Service error from ViewAllTask",
			projectID: "12345678-1234-1234-1234-123456789012",
			role:      roles.Admin,
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().ViewAllTask("12345678-1234-1234-1234-123456789012").Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectBody:   "Error fatching the tasks",
		},
		{
			name:      "Some tasks done",
			projectID: "12345678-1234-1234-1234-123456789012",
			role:      roles.Admin,
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().ViewAllTask("12345678-1234-1234-1234-123456789012").Return([]task.Task{
					{TaskStatus: status.Done},
					{TaskStatus: status.Pending},
					{TaskStatus: status.Done},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectBody:   `"completedTasks":2`,
		},
		{
			name:      "No tasks found (edge case)",
			projectID: "12345678-1234-1234-1234-123456789012",
			role:      roles.Admin,
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().ViewAllTask("12345678-1234-1234-1234-123456789012").Return([]task.Task{}, nil)
			},
			expectedCode: http.StatusOK, // handler doesn't return 404 (your branch is commented out)
			expectBody:   `"totalTasks":0`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc:= mocks.NewMockTaskServiceInterface(ctrl)
			tt.mockSetup(svc)

			h := &ProjectHandler{nil,nil,svc}

			req := httptest.NewRequest(http.MethodGet, "/v1/projects/"+tt.projectID+"/status", nil)
			ctx := context.WithValue(req.Context(), ContextKey.UserRole, tt.role)
			req = req.WithContext(ctx)

			// set path value (Go 1.21+)
			req.SetPathValue("project_id", tt.projectID)

			w := httptest.NewRecorder()
			h.ProjectStatus(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
			if tt.expectBody != "" && !strings.Contains(w.Body.String(), tt.expectBody) {
				t.Errorf("%s: expected body to contain %q, got %q", tt.name, tt.expectBody, w.Body.String())
			}
		})
	}

}

func TestGetProjects(t *testing.T) {
    tests := []struct {
        name         string
        urlPath      string
        role         interface{}
        mockSetup    func(svc *mocks.MockProjectServiceInterface)
        expectedCode int
        expectBody   string
    }{
        {
            name:      "User role not found",
            urlPath:   "/v1/projects",
            role:      nil,
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {},
            expectedCode: http.StatusUnauthorized,
            expectBody:   "User not authenticated",
        },
        {
            name:      "Invalid path",
            urlPath:   "/invalid/path",
            role:      roles.Admin,
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {},
            expectedCode: http.StatusBadRequest,
            expectBody:   "Invalid path",
        },
        {
            name:      "Non-admin forbidden to view all projects",
            urlPath:   "/v1/projects",
            role:      roles.Manager,
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {},
            expectedCode: http.StatusForbidden,
            expectBody:   "Access denied",
        },
        {
            name:      "Admin views all projects, service error",
            urlPath:   "/v1/projects",
            role:      roles.Admin,
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {
                svc.EXPECT().ViewAllProjects().Return(nil, errors.New("db error"))
            },
            expectedCode: http.StatusInternalServerError,
            expectBody:   "Failed to fetch projects",
        },
        // {
        //     name:      "Admin views all projects, no projects found",
        //     urlPath:   "/v1/projects",
        //     role:      roles.Admin,
        //     mockSetup: func(svc *mocks.MockProjectServiceInterface) {
        //         svc.EXPECT().ViewAllProjects().Return([]project.Project{}, nil)
        //     },
        //     expectedCode: http.StatusNotFound,
        //     expectBody:   "No projects assigned",
        // },
        {
            name:      "Admin views all projects, success",
            urlPath:   "/v1/projects",
            role:      roles.Admin,
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {
                svc.EXPECT().ViewAllProjects().Return([]project.Project{{ProjectId: "p1"}}, nil)
            },
            expectedCode: http.StatusOK,
            expectBody:   `"Projects Retrived Successfully"`,
        },
        {
            name:      "Assigned user, service error",
            urlPath:   "/v1/projects/emp1",
            role:      roles.Employee,
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {
                svc.EXPECT().ViewAssignedProject("emp1").Return(nil, errors.New("db error"))
            },
            expectedCode: http.StatusInternalServerError,
            expectBody:   "Failed to fetch Assigned projects",
        },
        {
            name:      "Assigned user, no projects found",
            urlPath:   "/v1/projects/emp1",
            role:      roles.Employee,
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {
                svc.EXPECT().ViewAssignedProject("emp1").Return([]project.Project{}, nil)
            },
            expectedCode: http.StatusNotFound,
            expectBody:   "No projects assigned",
        },
        {
            name:      "Assigned user, success",
            urlPath:   "/v1/projects/emp1",
            role:      roles.Employee,
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {
                svc.EXPECT().ViewAssignedProject("emp1").Return([]project.Project{{ProjectId: "p1"}}, nil)
            },
            expectedCode: http.StatusOK,
            expectBody:   `"Projects retrived successfully"`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            svc := mocks.NewMockProjectServiceInterface(ctrl)
            h := &ProjectHandler{projectService: svc}

            tt.mockSetup(svc)

            req := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)
            ctx := req.Context()
            if tt.role != nil {
                ctx = context.WithValue(ctx, ContextKey.UserRole, tt.role)
            }
            req = req.WithContext(ctx)

            w := httptest.NewRecorder()
            h.GetProjects(w, req)

            if w.Code != tt.expectedCode {
                t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
            }
            if tt.expectBody != "" && !strings.Contains(w.Body.String(), tt.expectBody) {
                t.Errorf("%s: expected body to contain %q, got %q", tt.name, tt.expectBody, w.Body.String())
            }
        })
    }
}


func TestCreateProject(t *testing.T) {
	tests := []struct {
		name         string
		role         interface{}
		userId       string
		body         string
		mockSetup    func(projSvc *mocks.MockProjectServiceInterface, userSvc *mocks.MockUserServiceInterface)
		expectedCode int
		expectBody   string
	}{
		{
			name:      "User role not found",
			role:      nil,
			userId:    "admin1",
			body:      `{}`,
			mockSetup: func(projSvc *mocks.MockProjectServiceInterface, userSvc *mocks.MockUserServiceInterface) {},
			expectedCode: http.StatusUnauthorized,
			expectBody:   "User not authenticated",
		},
		{
			name:      "Non-admin forbidden to create",
			role:      roles.Manager,
			userId:    "mgr1",
			body:      `{}`,
			mockSetup: func(projSvc *mocks.MockProjectServiceInterface, userSvc *mocks.MockUserServiceInterface) {},
			expectedCode: http.StatusForbidden,
			expectBody:   "Only admin can create projects",
		},
		{
			name:      "Invalid request body",
			role:      roles.Admin,
			userId:    "admin1",
			body:      `{invalid json}`,
			mockSetup: func(projSvc *mocks.MockProjectServiceInterface, userSvc *mocks.MockUserServiceInterface) {},
			expectedCode: http.StatusBadRequest,
			expectBody:   "Invalid input",
		},
		
		{
			name:   "Service error from AddProject",
			role:   roles.Admin,
			userId: "admin1",
			body:   `{"projectName": "Test", "projectDescription": "Desc", "deadline": "2025-09-01", "assignedManagerId": "mgr1"}`,
			mockSetup: func(projSvc *mocks.MockProjectServiceInterface, userSvc *mocks.MockUserServiceInterface) {
				userSvc.EXPECT().ViewProfile("mgr1").Return([]user.User{
					{Id: "mgr1", Status: user.Active},
				}, nil)
				projSvc.EXPECT().AddProject(gomock.Any()).Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectBody:   "Error creating project",
		},
		{
			name:   "Success",
			role:   roles.Admin,
			userId: "admin1",
			body:   `{"projectName": "Test", "projectDescription": "Desc", "deadline": "2025-09-01", "assignedManagerId": "mgr1"}`,
			mockSetup: func(projSvc *mocks.MockProjectServiceInterface, userSvc *mocks.MockUserServiceInterface) {
				userSvc.EXPECT().ViewProfile("mgr1").Return([]user.User{
					{Id: "mgr1", Status: user.Active},
				}, nil)
				projSvc.EXPECT().AddProject(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusCreated,
			expectBody:   "Project created successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			projSvc := mocks.NewMockProjectServiceInterface(ctrl)
			userSvc := mocks.NewMockUserServiceInterface(ctrl)
			h := &ProjectHandler{projectService: projSvc, userService: userSvc}

			tt.mockSetup(projSvc, userSvc)

			req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(tt.body))
			ctx := req.Context()
			if tt.role != nil {
				ctx = context.WithValue(ctx, ContextKey.UserRole, tt.role)
			}
			ctx = context.WithValue(ctx, ContextKey.UserId, tt.userId)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.CreateProject(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
			if tt.expectBody != "" && !strings.Contains(w.Body.String(), tt.expectBody) {
				t.Errorf("%s: expected body to contain %q, got %q", tt.name, tt.expectBody, w.Body.String())
			}
		})
	}
}

func TestDeleteProject(t *testing.T) {
    tests := []struct {
        name         string
        role         interface{}
        projectId    string
        mockSetup    func(svc *mocks.MockProjectServiceInterface)
        expectedCode int
        expectBody   string
    }{
        {
            name:      "User role not found",
            role:      nil,
            projectId: "p1",
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {},
            expectedCode: http.StatusUnauthorized,
            expectBody:   "User not authenticated",
        },
        {
            name:      "Non-admin forbidden to delete",
            role:      roles.Manager,
            projectId: "p1",
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {},
            expectedCode: http.StatusForbidden,
            expectBody:   "Only admin can delete projects",
        },
        {
            name:      "Missing project_id",
            role:      roles.Admin,
            projectId: "",
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {},
            expectedCode: http.StatusBadRequest,
            expectBody:   "Project ID is required",
        },
        {
            name:      "Service error on delete",
            role:      roles.Admin,
            projectId: "p1",
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {
                svc.EXPECT().DeleteProject("p1").Return(errors.New("db error"))
            },
            expectedCode: http.StatusInternalServerError,
            expectBody:   "Failed to delete project",
        },
        {
            name:      "Successful delete",
            role:      roles.Admin,
            projectId: "p1",
            mockSetup: func(svc *mocks.MockProjectServiceInterface) {
                svc.EXPECT().DeleteProject("p1").Return(nil)
            },
            expectedCode: http.StatusOK,
            expectBody:   "Project deleted successfully",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            svc := mocks.NewMockProjectServiceInterface(ctrl)
            h := &ProjectHandler{projectService: svc}
            defer ctrl.Finish()

            tt.mockSetup(svc)

            req := httptest.NewRequest(http.MethodDelete, "/v1/projects/"+tt.projectId, nil)
            if tt.projectId != "" {
                req.SetPathValue("project_id", tt.projectId)
            }
            ctx := req.Context()
            if tt.role != nil {
                ctx = context.WithValue(ctx, ContextKey.UserRole, tt.role)
            }
            req = req.WithContext(ctx)

            w := httptest.NewRecorder()
            h.DeleteProject(w, req)

            if w.Code != tt.expectedCode {
                t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
            }
            if tt.expectBody != "" && !strings.Contains(w.Body.String(), tt.expectBody) {
                t.Errorf("%s: expected body to contain %q, got %q", tt.name, tt.expectBody, w.Body.String())
            }
        })
    }
}
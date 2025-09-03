package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	"strings"
	"testing"

	"github.com/Yash-Watchguard/Tasknest/internal/handler"
	"github.com/Yash-Watchguard/Tasknest/internal/mocks"
	"go.uber.org/mock/gomock"

	// "github.com/Yash-Watchguard/Tasknest/internal/model/ContextKey"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	// "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	// "github.com/Yash-Watchguard/Tasknest/mocks"
)

func newTaskHandlerWithMock(t *testing.T) (*handler.TaskHandler, *mocks.MockTaskServiceInterface,*mocks.MockUserServiceInterface, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockTaskServiceInterface(ctrl)
	svc1:=mocks.NewMockUserServiceInterface(ctrl)
	h := handler.NewTaskHandler(svc,svc1)
	return h, svc,svc1, ctrl
}

func TestGetTask(t *testing.T) {
	validProjectID := "123e4567-e89b-12d3-a456-426614174000" // 36-char UUID

	tests := []struct {
		name         string
		userIdQuery  string
		role         roles.Role
		employeeId   string
		mockSetup    func(svc *mocks.MockTaskServiceInterface)
		expectedCode int
	}{
		{
			name:         "Employee unauthorized when no userId",
			userIdQuery:  "",
			role:         roles.Employee,
			employeeId:   "emp1",
			mockSetup:    func(svc *mocks.MockTaskServiceInterface) {},
			expectedCode: http.StatusForbidden,
		},
		{
			name:        "Manager gets error from ViewAllTask",
			userIdQuery: "",
			role:        roles.Manager,
			employeeId:  "mgr1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().ViewAllTask(validProjectID).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "Manager gets empty tasks",
			userIdQuery: "",
			role:        roles.Manager,
			employeeId:  "mgr1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().ViewAllTask(validProjectID).Return([]task.Task{}, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "Manager gets tasks successfully",
			userIdQuery: "",
			role:        roles.Manager,
			employeeId:  "mgr1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().ViewAllTask(validProjectID).Return([]task.Task{{TaskId: "t1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Employee requesting others tasks forbidden",
			userIdQuery:  "emp2",
			role:         roles.Employee,
			employeeId:   "emp1",
			mockSetup:    func(svc *mocks.MockTaskServiceInterface) {},
			expectedCode: http.StatusForbidden,
		},
		{
			name:        "Employee gets error from ViewAllAssignedTasksInProject",
			userIdQuery: "emp1",
			role:        roles.Employee,
			employeeId:  "emp1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().ViewAllAssignedTasksInProject(validProjectID, "emp1").
					Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "Employee gets empty tasks",
			userIdQuery: "emp1",
			role:        roles.Employee,
			employeeId:  "emp1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().ViewAllAssignedTasksInProject(validProjectID, "emp1").
					Return([]task.Task{}, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "Employee gets tasks successfully",
			userIdQuery: "emp1",
			role:        roles.Employee,
			employeeId:  "emp1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().ViewAllAssignedTasksInProject(validProjectID, "emp1").
					Return([]task.Task{{TaskId: "t1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, svc,_, ctrl := newTaskHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockSetup(svc)

			req := httptest.NewRequest(http.MethodGet, "/tasks?assigned_id="+tt.userIdQuery, nil)
			req.SetPathValue("project_id", validProjectID)

			ctx := context.WithValue(req.Context(), ContextKey.UserRole, tt.role)
			ctx = context.WithValue(ctx, ContextKey.UserId, tt.employeeId)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.GetTask(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d. Body: %s", 
					tt.name, tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}


func TestGetAssignedTask(t *testing.T) {
	tests := []struct {
		name         string
		role         roles.Role
		userId       string
		empId        string
		mockSetup    func(svc *mocks.MockTaskServiceInterface)
		expectedCode int
	}{
		{
			name:   "Employee gets assigned tasks successfully",
			role:   roles.Employee,
			userId: "emp1",
			empId:  "emp1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().
					GetAssigenedTask("emp1").
					Return([]task.Task{{TaskId: "t1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "Employee no assigned tasks found",
			role:   roles.Employee,
			userId: "emp1",
			empId:  "emp1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().
					GetAssigenedTask("emp1").
					Return([]task.Task{}, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "Employee service error on assigned tasks",
			role:   roles.Employee,
			userId: "emp1",
			empId:  "emp1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().
					GetAssigenedTask("emp1").
					Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "Employee forbidden to access other's tasks",
			role:   roles.Employee,
			userId: "emp2",
			empId:  "emp1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				// No call expected
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name:   "Manager forbidden to access assigned tasks",
			role:   roles.Manager,
			userId: "mgr1",
			empId:  "mgr1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				// No call expected
			},
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, svc,_, ctrl := newTaskHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockSetup(svc)

			req := httptest.NewRequest(http.MethodGet, "/v1/projects/tasks/"+tt.empId, nil)
			req.SetPathValue("employee_id", tt.empId) // <-- This is crucial!

			ctx := context.WithValue(req.Context(), ContextKey.UserRole, tt.role)
			ctx = context.WithValue(ctx, ContextKey.UserId, tt.userId)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.AssignedTasks(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
		})
	}
}


func TestCreateTask(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		setupContext func(req *http.Request) *http.Request
		mockSetup    func(taskSvc *mocks.MockTaskServiceInterface, userSvc *mocks.MockUserServiceInterface)
		expectedCode int
	}{
		{
			name: "Valid request, task created",
			body: `{
				"title": "Test Task",
				"description": "Test Description",
				"acceptance_criteria": "Criteria",
				"deadline": "2025-09-01",
				"priority": "High",
				"assigned_to": "emp1"
			}`,
			setupContext: func(req *http.Request) *http.Request {
				req.SetPathValue("project_id", "p1")
				ctx := context.WithValue(req.Context(), ContextKey.UserId, "mgr1")
				return req.WithContext(ctx)
			},
			mockSetup: func(taskSvc *mocks.MockTaskServiceInterface, userSvc *mocks.MockUserServiceInterface) {
				userSvc.EXPECT().ViewProfile("emp1").
					Return([]user.User{{Status: user.Active}}, nil)

				taskSvc.EXPECT().CreateTask(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Invalid request body",
			body: `{invalid json}`,
			setupContext: func(req *http.Request) *http.Request {
				req.SetPathValue("project_id", "p1")
				ctx := context.WithValue(req.Context(), ContextKey.UserId, "mgr1")
				return req.WithContext(ctx)
			},
			mockSetup:    func(taskSvc *mocks.MockTaskServiceInterface, userSvc *mocks.MockUserServiceInterface) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid deadline format",
			body: `{
				"title": "Test Task",
				"description": "Test Description",
				"acceptance_criteria": "Criteria",
				"deadline": "01-09-2025",
				"priority": "High",
				"assigned_to": "emp1"
			}`,
			setupContext: func(req *http.Request) *http.Request {
				req.SetPathValue("project_id", "p1")
				ctx := context.WithValue(req.Context(), ContextKey.UserId, "mgr1")
				return req.WithContext(ctx)
			},
			mockSetup: func(taskSvc *mocks.MockTaskServiceInterface, userSvc *mocks.MockUserServiceInterface) {
				userSvc.EXPECT().ViewProfile("emp1").
					Return([]user.User{{Status: user.Active}}, nil)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid priority value",
			body: `{
				"title": "Test Task",
				"description": "Test Description",
				"acceptance_criteria": "Criteria",
				"deadline": "2025-09-01",
				"priority": "Urgent",
				"assigned_to": "emp1"
			}`,
			setupContext: func(req *http.Request) *http.Request {
				req.SetPathValue("project_id", "p1")
				ctx := context.WithValue(req.Context(), ContextKey.UserId, "mgr1")
				return req.WithContext(ctx)
			},
			mockSetup: func(taskSvc *mocks.MockTaskServiceInterface, userSvc *mocks.MockUserServiceInterface) {
				userSvc.EXPECT().ViewProfile("emp1").
					Return([]user.User{{Status: user.Active}}, nil)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Service error on create",
			body: `{
				"title": "Test Task",
				"description": "Test Description",
				"acceptance_criteria": "Criteria",
				"deadline": "2025-09-01",
				"priority": "High",
				"assigned_to": "emp1"
			}`,
			setupContext: func(req *http.Request) *http.Request {
				req.SetPathValue("project_id", "p1")
				ctx := context.WithValue(req.Context(), ContextKey.UserId, "mgr1")
				return req.WithContext(ctx)
			},
			mockSetup: func(taskSvc *mocks.MockTaskServiceInterface, userSvc *mocks.MockUserServiceInterface) {
				userSvc.EXPECT().ViewProfile("emp1").
					Return([]user.User{{Status: user.Active}}, nil)

				taskSvc.EXPECT().CreateTask(gomock.Any()).Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mockTaskSvc, mockUserSvc, ctrl := newTaskHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockSetup(mockTaskSvc, mockUserSvc)

			req := httptest.NewRequest(http.MethodPost, "/v1/projects/p1/tasks", strings.NewReader(tt.body))
			req = tt.setupContext(req)

			w := httptest.NewRecorder()
			h.CreateTask(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
		})
	}
}



func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name         string
		role         roles.Role
		userId       string
		projectId    string
		taskId       string
		mockSetup    func(svc *mocks.MockTaskServiceInterface)
		expectedCode int
	}{
		{
			name:         "Employee forbidden to delete",
			role:         roles.Employee,
			userId:       "emp1",
			projectId:    "p1",
			taskId:       "t1",
			mockSetup:    func(svc *mocks.MockTaskServiceInterface) {},
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "Missing project_id",
			role:         roles.Manager,
			userId:       "mgr1",
			projectId:    "",
			taskId:       "t1",
			mockSetup:    func(svc *mocks.MockTaskServiceInterface) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Missing task_id",
			role:         roles.Manager,
			userId:       "mgr1",
			projectId:    "p1",
			taskId:       "",
			mockSetup:    func(svc *mocks.MockTaskServiceInterface) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "Service error on delete",
			role:      roles.Manager,
			userId:    "mgr1",
			projectId: "p1",
			taskId:    "t1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().DeleteTask("mgr1", "t1").Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:      "Successful delete",
			role:      roles.Manager,
			userId:    "mgr1",
			projectId: "p1",
			taskId:    "t1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().DeleteTask("mgr1", "t1").Return(nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, svc, _,ctrl := newTaskHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockSetup(svc)

			req := httptest.NewRequest(http.MethodDelete, "/v1/projects/"+tt.projectId+"/tasks/"+tt.taskId, nil)
			req.SetPathValue("project_id", tt.projectId)
			req.SetPathValue("task_id", tt.taskId)
			ctx := context.WithValue(req.Context(), ContextKey.UserRole, tt.role)
			ctx = context.WithValue(ctx, ContextKey.UserId, tt.userId)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.DeleteTask(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
		})
	}
}

func TestUpdateStatus(t *testing.T) {

	tests := []struct {
		name         string
		body         string
		userId       string
		taskId       string
		mockSetup    func(svc *mocks.MockTaskServiceInterface)
		expectedCode int
	}{
		{
			name:   "Valid request, status updated",
			body:   `{"status": "pending"}`, // Use the exact string your parser expects
			userId: "emp1",
			taskId: "t1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().UpdateTaskStatus("emp1", "t1", gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid request body",
			body:         "{invalid json}",
			userId:       "emp1",
			taskId:       "t1",
			mockSetup:    func(svc *mocks.MockTaskServiceInterface) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid status value",
			body:         "{\"status\": \"NotAStatus\"}",
			userId:       "emp1",
			taskId:       "t1",
			mockSetup:    func(svc *mocks.MockTaskServiceInterface) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "Service error on update",
			body:   `{"status": "pending"}`, // Use the exact string your parser expects
			userId: "emp1",
			taskId: "t1",
			mockSetup: func(svc *mocks.MockTaskServiceInterface) {
				svc.EXPECT().UpdateTaskStatus("emp1", "t1", gomock.Any()).Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, svc, _,ctrl := newTaskHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockSetup(svc)

			req := httptest.NewRequest(http.MethodPut, "/v1/projects/tasks/"+tt.taskId+"/status", strings.NewReader(tt.body))
			req.SetPathValue("task_id", tt.taskId)
			ctx := context.WithValue(req.Context(), ContextKey.UserId, tt.userId)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.UpdateStatus(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
		})
	}
}

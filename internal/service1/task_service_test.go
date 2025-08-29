package service1

import (
	"testing"

	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/mocks"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"go.uber.org/mock/gomock"
)

func TestTaskService_ViewAllTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepo(ctrl)
	taskService :=NewTaskService(mockRepo)

	projectID := "p1"
	expectedTasks := []task.Task{
		{TaskId: "t1", Title: "Task 1"},
	}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func() {
				mockRepo.EXPECT().ViewAllTask(projectID).Return(expectedTasks, nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			mockSetup: func() {
				mockRepo.EXPECT().ViewAllTask(projectID).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			tasks, err := taskService.ViewAllTask(projectID)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr && len(tasks) != 1 {
				t.Errorf("expected 1 task, got %d", len(tasks))
			}
		})
	}
}

func TestTaskService_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepo(ctrl)
	taskService := NewTaskService(mockRepo)

	newTask := task.Task{TaskId: "t1", Title: "New Task"}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func() {
				mockRepo.EXPECT().SaveTask(newTask).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			mockSetup: func() {
				mockRepo.EXPECT().SaveTask(newTask).Return(errors.New("save failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := taskService.CreateTask(newTask)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepo(ctrl)
	taskService := NewTaskService(mockRepo)

	taskID := "task-123"

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func() {
				mockRepo.EXPECT().DeleteTask(taskID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			mockSetup: func() {
				mockRepo.EXPECT().DeleteTask(taskID).Return(errors.New("delete failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := taskService.DeleteTask("manager-1", taskID)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestTaskService_GetAssigenedTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepo(ctrl)
	taskService := NewTaskService(mockRepo)

	empID := "emp-1"
	expectedTasks := []task.Task{
		{TaskId: "t1", Title: "Task 1"},
		{TaskId: "t2", Title: "Task 2"},
	}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func() {
				mockRepo.EXPECT().ViewAssignedTask(empID).Return(expectedTasks, nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			mockSetup: func() {
				mockRepo.EXPECT().ViewAssignedTask(empID).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			tasks, err := taskService.GetAssigenedTask(empID)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}

			if !tt.wantErr && len(tasks) != len(expectedTasks) {
				t.Errorf("expected %d tasks, got %d", len(expectedTasks), len(tasks))
			}
		})
	}
}


func TestTaskService_UpdateTaskStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepo(ctrl)
	taskService := NewTaskService(mockRepo)

	userID := "u1"
	taskID := "t1"
	newStatus := status.InProgress

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func() {
				mockRepo.EXPECT().UpdateTaskStatus(userID, taskID, newStatus).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			mockSetup: func() {
				mockRepo.EXPECT().UpdateTaskStatus(userID, taskID, newStatus).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := taskService.UpdateTaskStatus(userID, taskID, newStatus)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestTaskService_ViewAllAssignedTasksInProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepo(ctrl)
	taskService := NewTaskService(mockRepo)

	projectID := "p1"
	empID := "e1"
	expectedTasks := []task.Task{
		{TaskId: "t1", Title: "Task 1"},
	}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
		wantLen   int
	}{
		{
			name: "success",
			mockSetup: func() {
				mockRepo.EXPECT().
					ViewAllAssignedTasksInProject(projectID, empID).
					Return(expectedTasks, nil)
			},
			wantErr: false,
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			tasks, err := taskService.ViewAllAssignedTasksInProject(projectID, empID)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
			if len(tasks) != tt.wantLen {
				t.Errorf("expected %d tasks, got %d", tt.wantLen, len(tasks))
			}
		})
	}
}





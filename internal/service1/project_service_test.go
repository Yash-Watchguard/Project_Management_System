package service1

import (

	"testing"

	"github.com/Yash-Watchguard/Tasknest/internal/mocks"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"go.uber.org/mock/gomock"
)

func TestAddProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjectRepository(ctrl)
	service := NewProjectService(mockRepo)

	p := project.Project{ProjectId: "p1", ProjectName: "Project 1"}

	mockRepo.EXPECT().AddProject(p).Return(nil)

	if err := service.AddProject(p); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestViewAllProjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjectRepository(ctrl)
	service := NewProjectService(mockRepo)

	expectedProjects := []project.Project{
		{ProjectId: "p1", ProjectName: "Project 1"},
	}

	mockRepo.EXPECT().ViewAllProjects().Return(expectedProjects, nil)

	projects, err := service.ViewAllProjects()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(projects) != len(expectedProjects) {
		t.Errorf("expected %d projects, got %d", len(expectedProjects), len(projects))
	}
}

func TestDeleteProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjectRepository(ctrl)
	service := NewProjectService(mockRepo)

	projectID := "p1"

	
	mockRepo.EXPECT().DeleteProject(projectID).Return(nil)

	if err := service.DeleteProject(projectID); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestViewAssignedProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjectRepository(ctrl)
	service :=NewProjectService(mockRepo)

	userID := "u1"
	expectedProjects := []project.Project{
		{ProjectId: "p1", ProjectName: "Project 1"},
	}

	mockRepo.EXPECT().ViewAssignedProject(userID).Return(expectedProjects, nil)

	projects, err := service.ViewAssignedProject(userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(projects) != len(expectedProjects) {
		t.Errorf("expected %d projects, got %d", len(expectedProjects), len(projects))
	}
}









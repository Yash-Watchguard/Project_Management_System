package service1

import (
	"context"
	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
)

type ProjectService struct{
	projectRepo interfaces.ProjectRepository
}

func NewProjectService(projectRepo interfaces.ProjectRepository)*ProjectService{
	return &ProjectService{projectRepo: projectRepo}
}

func (ps *ProjectService) AddProject(project project.Project) error {
	return ps.projectRepo.AddProject(project)
}

func (ps *ProjectService) ViewAllProjects(ctx context.Context) ([]project.Project, error) {
	var projects []project.Project
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userRole != 0 {
		return projects, errors.New("unauthorized access")
	}
	return ps.projectRepo.ViewAllProjects()
}

func (ps *ProjectService) DeleteProject(ctx context.Context, projectID string) error {
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userRole != 0 {
		return errors.New("unauthorized access")
	}
	return ps.projectRepo.DeleteProject(projectID)
}
func (ps *ProjectService) ViewAssignedProject(ctx context.Context)([]project.Project,error){
	userRole:=ctx.Value(ContextKey.UserRole).(roles.Role)
	
	var projects []project.Project
	if userRole!=1 {
		return projects,errors.New("unauthorized access")
	}
	projects,err:=ps.projectRepo.ViewAllProjects()

	return projects,err
}

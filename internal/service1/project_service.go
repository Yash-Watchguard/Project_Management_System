package service1

import (
	

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	
	
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	
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

func (ps *ProjectService) ViewAllProjects() ([]project.Project, error) {
	return ps.projectRepo.ViewAllProjects()
}

func (ps *ProjectService) DeleteProject( projectID string) error {
	
	return ps.projectRepo.DeleteProject(projectID)
	
}
func (ps *ProjectService) ViewAssignedProject(userId string)([]project.Project,error){
	
	projects,err:=ps.projectRepo.ViewAssignedProject(userId)

	return projects,err
}

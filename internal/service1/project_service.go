package service1

import (
	"errors"
	"time"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	"github.com/Yash-Watchguard/Tasknest/internal/util"

	
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
)

//go:generate mockgen -source=project_service.go -destination=../mocks/mock_projectservice.go -package=mocks
type ProjectServiceInterface interface{
   AddProject(project project.Project) error
   ViewAllProjects() ([]project.Project, error)
   DeleteProject( creatorId,managerId,projectID string) error
   ViewAssignedProject(userId string)([]project.Project,error)
   UpdateProject(creatorId,projectId,managerId string,updates map[string]any)(error)
}

type ProjectService struct{
	projectRepo interfaces.ProjectRepository
}

func NewProjectService(projectRepo interfaces.ProjectRepository)ProjectServiceInterface{
	return &ProjectService{projectRepo: projectRepo}
}

func(ps *ProjectService)UpdateProject(creatorId,projectId,managerId string,updates map[string]any)(error){
     finalUpdates := make(map[string]interface{})

	 if title,ok:=updates["Project_name"].(string);ok{
		if len(title)!=0 {
			finalUpdates["Project_name"]=title
		}
	 }
	 if status,ok:=updates["Project_description"].(string);ok{
		if len(status)!=0 {
			
			finalUpdates["Project_description"]=status
		}
	 }

	 if deadline,ok:=updates["Project_deadline"].(string);ok{
		if len(deadline)!=0 {
			deadlineParsed,err:=util.ParseDate(deadline)
			if err!=nil {
				return errors.New("invalid date")
			}
			finalUpdates["Project_deadline"]=deadlineParsed.Format(time.RFC3339)
		}
	 }


	 if empId,ok:=updates["Assigned_manager"].(string);ok{
		if(len(empId)!=0){
			finalUpdates["Assigned_manager"]=empId
		}
	 
}
return ps.projectRepo.UpdateProject(projectId,creatorId,managerId,finalUpdates)
}


func (ps *ProjectService) AddProject(project project.Project) error {
	return ps.projectRepo.AddProject(project)
}

func (ps *ProjectService) ViewAllProjects() ([]project.Project, error) {
	return ps.projectRepo.ViewAllProjects()
}

func (ps *ProjectService) DeleteProject( creatorId,managerId,projectID string) error {
	
	return ps.projectRepo.DeleteProject(creatorId,managerId,projectID)
	
}
func (ps *ProjectService) ViewAssignedProject(userId string)([]project.Project,error){
	
	projects,err:=ps.projectRepo.ViewAssignedProject(userId)

	return projects,err
}

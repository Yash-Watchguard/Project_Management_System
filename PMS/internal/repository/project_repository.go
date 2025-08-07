// repository/project_repo.go
package repository

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
)

type ProjectRepo struct {
	filePath string
}

func NewProjectRepo() *ProjectRepo {
	return &ProjectRepo{filePath:  "C:/Users/ygoyal/Desktop/PMS_Project/Pms/internal/data/project.json"}
}

func (pr *ProjectRepo) AddProject(newProject project.Project) error {
	var projects []project.Project

	data, err := os.ReadFile(pr.filePath)

	if err == nil {
		json.Unmarshal(data, &projects)
	}

	projects = append(projects, newProject)

	newData, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(pr.filePath, newData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ProjectRepo) ViewAllProjects() ([]project.Project, error) {
	var projects []project.Project

	data, err := os.ReadFile(pr.filePath)
	if err != nil {
		return nil, err
	}
     
	err = json.Unmarshal(data, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
func (pr *ProjectRepo) DeleteProject(projectID string) error {
	data, err := os.ReadFile(pr.filePath)
	if err != nil {
		return err
	}

	var projects []project.Project

	var newProjects []project.Project
	found := false
    _=json.Unmarshal(data,&projects)
	for _, project := range projects {
		if project.ProjectId == projectID {
			found = true
			continue
		}
		newProjects = append(newProjects, project)
	}
    if !found {
		return errors.New("project not found")
	}
	updatedData, err := json.MarshalIndent(newProjects, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pr.filePath,updatedData,0644)
}

func(pr *ProjectRepo)ViewAssignedProject(userId string)([]project.Project,error){
   var projects []project.Project

	data, err := os.ReadFile(pr.filePath)
	if err != nil {
		return nil, err
	}
    
	var assignedProjects []project.Project
	err = json.Unmarshal(data, &projects)
	if err != nil {
		return nil, errors.New("error in getting projects")
	}
	for _,pro:=range projects{
		if pro.AssignedManager==userId{
			assignedProjects=append(assignedProjects, pro)
		}
	}
	
	return assignedProjects,nil
}





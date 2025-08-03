// repository/project_repo.go
package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"errors"

	"github.com/Yash-Watchguard/Tasknest/interfaces"
	"github.com/Yash-Watchguard/Tasknest/model"
)

type ProjectRepo struct {
	filePath string
}

func NewProjectRepo() interfaces.ProjectRepository {
	return &ProjectRepo{filePath:  "C:/Users/ygoyal/Desktop/PMS_Project/Pms/data/project.json"}
}

func (pr *ProjectRepo) AddProject(project model.Project) error {
	var projects []model.Project

	// Read existing data
	data, err := ioutil.ReadFile(pr.filePath)
	if err == nil {
		json.Unmarshal(data, &projects)
	}

	projects = append(projects, project)

	// Save to file
	file, err := os.Create(pr.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	newData, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(newData)
	return err
}
func (pr *ProjectRepo) ViewAllProjects() ([]model.Project, error) {
	var projects []model.Project

	data, err := ioutil.ReadFile(pr.filePath)
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
	data, err := ioutil.ReadFile(pr.filePath)
	if err != nil {
		return err
	}

	var projects []model.Project

	var newProjects []model.Project
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
		return errors.New("Project not found")
	}
	updatedData, err := json.MarshalIndent(newProjects, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(pr.filePath,updatedData,0644)
}





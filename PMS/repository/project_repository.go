// repository/project_repo.go
package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Yash-Watchguard/Tasknest/interfaces"
	"github.com/Yash-Watchguard/Tasknest/model"
)

type ProjectRepo struct {
	filePath string
}

func NewProjectRepo() interfaces.ProjectRepository {
	return &ProjectRepo{filePath:  "C:/Users/ygoyal/Desktop/PMS/data/user.json"}
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
func ViewProject(pr *ProjectRepo)error{
	
return nil
}



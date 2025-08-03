package interfaces

import "github.com/Yash-Watchguard/Tasknest/model"

type ProjectRepository interface {
	AddProject(project model.Project)error
	ViewAllProjects()([]model.Project,error)
	DeleteProject(projectID string) error
}

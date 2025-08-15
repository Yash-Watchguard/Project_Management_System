package interfaces

import (
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
)

type ProjectRepository interface {
	AddProject(newProject project.Project)error
	ViewAllProjects()([]project.Project,error)
	ViewAssignedProject(userId string)([]project.Project,error)
	DeleteProject(projectID string) error
}

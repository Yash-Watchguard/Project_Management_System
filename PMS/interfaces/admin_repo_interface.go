package interfaces

import "github.com/Yash-Watchguard/Tasknest/model"

type Admin interface {
	ViewAllUsers()
	deleteUser(userId string) error
	GetAllManager()
	AddProject(project *model.Project)error
	ViewProject()error
}
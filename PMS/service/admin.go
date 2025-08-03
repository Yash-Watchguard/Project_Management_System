package service

import (
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/interfaces"
	"github.com/Yash-Watchguard/Tasknest/model"
)

type AdminService struct {
	userRepo    interfaces.UserRepository
	projectRepo interfaces.ProjectRepository
}

func NewAdminServices(
	userRepo interfaces.UserRepository,
	projectRepo interfaces.ProjectRepository,
) *AdminService {
	return &AdminService{
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}
func(a * AdminService)ViewProfile(user *model.User){
	a.userRepo.ViewProfile(user)
}
func (a *AdminService) ViewAllUsers() {
	users := a.userRepo.GetAllUsers()

	fmt.Println("--------- All Users ---------")
	for _, user := range users {
		fmt.Printf("ID: %s, Name: %s, Email: %s, Role: %s\n",
			user.Id, user.Name, user.Email, user.Role)
	}
	fmt.Println("Press ENTER to return to dashboard...")
    fmt.Scanln()
}
func(a *AdminService)DeleteUser(userId string)error{
	return a.userRepo.DeleteUserById(userId)
}
func(a *AdminService)GetAllManager()error{
	return a.userRepo.GetAllManager()
}

func(as *AdminService) AddProject(project model.Project) error {
	return as.projectRepo.AddProject(project)
}
func (as *AdminService) ViewAllProjects() ([]model.Project, error) {
	return as.projectRepo.ViewAllProjects()
}
func (as *AdminService) DeleteProject(projectID string) error {
	return as.projectRepo.DeleteProject(projectID)
}



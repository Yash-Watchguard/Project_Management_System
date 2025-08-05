package service

import (
	"context"
	"errors"
	"fmt"
	
    "github.com/fatih/color"
	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
)

type AdminService struct {
	userRepo    interfaces.UserRepository
	projectRepo interfaces.ProjectRepository
}

func NewAdminServices(userRepo interfaces.UserRepository,projectRepo interfaces.ProjectRepository) *AdminService {
	return &AdminService{
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}



func (a *AdminService) ViewAllUsers(ctx context.Context)error {

    userId:=ctx.Value(ContextKey.UserRole).(int)
	if userId!=0{
       return errors.New("unautherized access")
	}
	users := a.userRepo.GetAllUsers()

	fmt.Println("--------- All Users ---------")
	counter:=1
	for _, user := range users {
		fmt.Printf("%d . ID: %s, Name: %s, Email: %s, Role: %d\n",
			user.Id, user.Name, user.Email, user.Role,counter)
			counter++
	}
	fmt.Println("Press ENTER to return to dashboard...")
    fmt.Scanln()
	return nil
}


func(a *AdminService)DeleteUser(userId string)error{
	return a.userRepo.DeleteUserById(userId)
}


func(a *AdminService)GetAllManager(ctx context.Context)error{
	userId:=ctx.Value(ContextKey.UserRole).(int)
	if userId!=0{
       return errors.New("unautherized access")
	}
	return a.userRepo.GetAllManager()
}

func(as *AdminService)AddProject(project project.Project) error {
	return as.projectRepo.AddProject(project)
}
func (as *AdminService) ViewAllProjects() ([]project.Project, error) {
	return as.projectRepo.ViewAllProjects()
}
func (as *AdminService) DeleteProject(projectID string) error {
	return as.projectRepo.DeleteProject(projectID)
}



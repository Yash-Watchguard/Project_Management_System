package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/fatih/color"
)

type AdminService struct {
	userRepo    interfaces.UserRepository
	projectRepo interfaces.ProjectRepository
}

func NewAdminServices(userRepo interfaces.UserRepository, projectRepo interfaces.ProjectRepository) *AdminService {
	return &AdminService{
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}
func (a *AdminService) ViewProfile(ctx context.Context, userId string) ([]user.User, error) {
	userID := ctx.Value(ContextKey.UserId).(string)
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)


	if userID == userId || userRole == 0 {
		return a.userRepo.ViewProfile(userId)
	}
	return nil, errors.New("unauthorized access")
}

func (a *AdminService) ViewAllUsers(ctx context.Context) error {

	userID := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userID != 0 {
		return errors.New("unautherized access")
	}
	users := a.userRepo.GetAllUsers()

	fmt.Println("--------- All Users ---------")
	counter := 1
	for _, user := range users {
		fmt.Printf("%d. ID: %s, Name: %s, Email: %s, Role: %d\n",counter,user.Id, user.Name, user.Email, user.Role)
		counter++
	}
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()
	return nil
}

func (a *AdminService) DeleteUser(ctx context.Context,userId string) error {
	userID := ctx.Value(ContextKey.UserId).(string)
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

	if userId == userID || userRole == 0 {
		err:=a.userRepo.DeleteUserById(userId)
		if err != nil {
			color.Red("%v", err)
		}
	}else{
		return errors.New("unauthorized access")
	}
	return nil
}

func (a *AdminService) GetAllManager(ctx context.Context) error {
	userId := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userId != 0 {
		return errors.New("unautherized access")
	}
	return a.userRepo.GetAllManager()
}

func (as *AdminService) AddProject(project project.Project) error {
	return as.projectRepo.AddProject(project)
}


func (as *AdminService) ViewAllProjects(ctx context.Context) ([]project.Project, error) {
	var projects []project.Project
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userRole!=0{
		return projects,errors.New("unauthorized access")
	}
	return as.projectRepo.ViewAllProjects()
}


func (as *AdminService) DeleteProject(ctx context.Context,projectID string) error {
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userRole!=0{
		return errors.New("unauthorized access")
	} 
	return as.projectRepo.DeleteProject(projectID)
}

package service

import (
	"context"
	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
)

type ManagerService struct {
	userRepo    interfaces.UserRepository
	projectRepo interfaces.ProjectRepository
}

func NewManagerService(userRepo interfaces.UserRepository, projectRepo interfaces.ProjectRepository) *ManagerService {
	return &ManagerService{userRepo: userRepo, projectRepo: projectRepo}
}

func(manager *ManagerService)ViewProfile(ctx context.Context,userId string)([]user.User,error){
	userID := ctx.Value(ContextKey.UserId).(string)
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

	if userID == userId || userRole == 0 {
		return manager.userRepo.ViewProfile(userId)
	}
	return nil, errors.New("unauthorized access")
}

func(manager *ManagerService)DeleteUser(ctx context.Context,userId string)error{
	userID:=ctx.Value(ContextKey.UserId).(string)
	
	if userID!= userId{
		return errors.New("unauthorized access")
	}

	return manager.userRepo.DeleteUserById(userId)
}

func (manger *ManagerService) ViewAssignedProject(ctx context.Context)([]project.Project,error){
	userRole:=ctx.Value(ContextKey.UserRole).(roles.Role)
	var projects []project.Project
	if userRole!=1{
		return projects,errors.New("unauthorized access")
	}
	projects,err:=manger.projectRepo.ViewAllProjects()

	return projects,err
}
func(manager *ManagerService)UpdateProfile(userId string,ctx context.Context,name string, email string,password string,number string)error{
    userID:=ctx.Value(ContextKey.UserId).(string)

	if userID!=userId {
        return errors.New("unauthorized access")
	}
	return manager.userRepo.UpdateProfile(userId,name,email,password,number)
}


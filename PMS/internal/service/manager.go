package service

import (
	"context"
	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
)

type ManagerService struct {
	userRepo    interfaces.UserRepository
	projectRepo interfaces.ProjectRepository
	taskRepo    interfaces.TaskRepo
	managerRepo interfaces.ManagerRepository
}

func NewManagerService(userRepo interfaces.UserRepository, projectRepo interfaces.ProjectRepository,taskRepo interfaces.TaskRepo,managerRepo interfaces.ManagerRepository) *ManagerService {
	return &ManagerService{userRepo: userRepo, projectRepo: projectRepo,taskRepo: taskRepo,managerRepo: managerRepo}
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

func(manager *ManagerService)ViewAllTask(ctx context.Context,projectId string)([]task.Task,error){
	userRole:=ctx.Value(ContextKey.UserRole).(roles.Role)
	var tasks []task.Task
	if userRole!=1{
       return tasks,errors.New("unauthorized access")
	}
	return manager.taskRepo.ViewAllTask(projectId)
}
func(manager *ManagerService)CreateTask(ctx context.Context,managerid string,task task.Task)error{
	userID:=ctx.Value(ContextKey.UserId).(string)

	if userID!=managerid{
		return errors.New("unauthoeized access")
	}
	return manager.taskRepo.SaveTask(task)
}

func(manager *ManagerService)DeleteTask(ctx context.Context,managerId string,taskId string)error{
	userId:=ctx.Value(ContextKey.UserId).(string)

	if userId!=managerId{
		return errors.New("unauthorizrd access")
	}
	return manager.taskRepo.DeleteTask(taskId)
}
func(manager *ManagerService)ViewAllEmplpyee(ctx context.Context)([]user.User,error){
	userRole:=ctx.Value(ContextKey.UserId).(roles.Role)

    if userRole!=1{
		return []user.User{},errors.New("unauthorized access")
	}
	return manager.managerRepo.ViewAllEmployee()
}

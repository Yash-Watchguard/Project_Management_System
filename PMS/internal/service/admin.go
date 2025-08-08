package service

import (
	"context"
	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/fatih/color"
)

type AdminService struct {
	userRepo    interfaces.UserRepository
	projectRepo interfaces.ProjectRepository
	adminRepo   interfaces.AdminRepository
	taskRepo    interfaces.TaskRepo
	commentRepo interfaces.CommentRepo
}

func NewAdminServices(userRepo interfaces.UserRepository, projectRepo interfaces.ProjectRepository, adminRepo interfaces.AdminRepository, taskRepo interfaces.TaskRepo, commentRepo interfaces.CommentRepo) *AdminService {
	return &AdminService{
		userRepo:    userRepo,
		projectRepo: projectRepo,
		adminRepo:   adminRepo,
		taskRepo:    taskRepo,
		commentRepo: commentRepo,
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

func (a *AdminService) ViewAllUsers(ctx context.Context) ([]user.User, error) {

	userID := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userID != 0 {
		return []user.User{}, errors.New("unautherized access")
	}
	return a.userRepo.GetAllUsers()

}

func (a *AdminService) DeleteUser(ctx context.Context, userId string) error {
	userID := ctx.Value(ContextKey.UserId).(string)
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

	if userId == userID || userRole == 0 {
		err := a.userRepo.DeleteUserById(userId)
		if err != nil {
			color.Red("%v", err)
		}
	} else {
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
	if userRole != 0 {
		return projects, errors.New("unauthorized access")
	}
	return as.projectRepo.ViewAllProjects()
}

func (as *AdminService) DeleteProject(ctx context.Context, projectID string) error {
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userRole != 0 {
		return errors.New("unauthorized access")
	}
	return as.projectRepo.DeleteProject(projectID)
}

func (as *AdminService) PromoteEmployee(ctx context.Context, employeeId string) error {
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

	if userRole != 0 {
		return errors.New("unauthorized person")
	}
	return as.adminRepo.PromoteEmployee(employeeId)
}

func (as *AdminService) UpdateProfile(userId string, ctx context.Context, name string, email string, password string, number string) error {
	userID := ctx.Value(ContextKey.UserId).(string)

	if userID != userId {
		return errors.New("unauthorized access")
	}
	return as.userRepo.UpdateProfile(userId, name, email, password, number)
}

func (as *AdminService) ViewAllTask(ctx context.Context, projectId string) ([]task.Task, error) {
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

	if userRole != 0 {
		return []task.Task{}, errors.New("unauthorized access")
	}
	return as.taskRepo.ViewAllTask(projectId)
}

func (as *AdminService) ViewAllComment(taskId string) ([]comment.Comment, error) {

	return as.commentRepo.ViewAllComments(taskId)
}

func (as *AdminService) UpdateComment(ctx context.Context, commentId string, updatedComment string) error {

	userId := ctx.Value(ContextKey.UserId).(string)

	return as.commentRepo.UpdateComment(userId, commentId, updatedComment)
}

func (as *AdminService) AddComment(newComment comment.Comment) error {
	return as.commentRepo.AddComment(newComment)
}

func (as *AdminService) DeleteComment(ctx context.Context, commentId string) error {

	userId := ctx.Value(ContextKey.UserId).(string)
	return as.commentRepo.DeleteComment(userId, commentId)
}

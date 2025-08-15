// package service

// import (
// 	"context"
// 	"errors"

// 	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
// 	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
// 	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
// 	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
// 	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
// 	// "github.com/Yash-Watchguard/Tasknest/internal/model/project"
// 	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
// 	// "github.com/Yash-Watchguard/Tasknest/internal/model/task"
// 	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
// 	"github.com/fatih/color"
// )

// type EmployeeService struct {
// 	userRepo interfaces.UserRepository
// 	commentRepo interfaces.CommentRepo
// 	taskRepo interfaces.TaskRepo
// 	EmpRepo  interfaces.EmployeeRepo
// }

// func NewEmpService(userRepo interfaces.UserRepository, taskRepo interfaces.TaskRepo, commentRepo interfaces.CommentRepo,EmpRepo interfaces.EmployeeRepo)*EmployeeService{
// 	return &EmployeeService{userRepo: userRepo,taskRepo: taskRepo,commentRepo: commentRepo,EmpRepo: EmpRepo}
// }

// func (es *EmployeeService) ViewProfile(ctx context.Context, userId string) ([]user.User, error) {
// 	userID := ctx.Value(ContextKey.UserId).(string)
// 	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

// 	if userID == userId || userRole == 0 {
// 		return es.userRepo.ViewProfile(userId)
// 	}
// 	return nil, errors.New("unauthorized access")
// }

// func (es *EmployeeService) UpdateProfile(userId string, ctx context.Context, name string, email string, password string, number string) error {
// 	userID := ctx.Value(ContextKey.UserId).(string)

// 	if userID != userId {
// 		return errors.New("unauthorized access")
// 	}
// 	return es.userRepo.UpdateProfile(userId, name, email, password, number)
// }

// func (es *EmployeeService) DeleteEmp(ctx context.Context, userId string) error {
// 	userID := ctx.Value(ContextKey.UserId).(string)
// 	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

// 	if userId == userID || userRole == 0 {
// 		err := es.userRepo.DeleteUserById(userId)
// 		if err != nil {
// 			color.Red("%v", err)
// 		}
// 	} else {
// 		return errors.New("unauthorized access")
// 	}
// 	return nil
// }

// func(es *EmployeeService)GetAssigenedTask(ctx context.Context,empId string)([]task.Task,error){
//     userId:=ctx.Value(ContextKey.UserId).(string)

// 	if userId!=empId{
// 		return []task.Task{},errors.New("unauthorized access")
// 	}

// 	return es.EmpRepo.ViewAssignedTask(empId)
// }

// func (es *EmployeeService) ViewAllComment(taskId string) ([]comment.Comment, error) {

// 	return es.commentRepo.ViewAllComments(taskId)

// }

// func (es *EmployeeService) UpdateComment(ctx context.Context, commentId string, updatedComment string) error {

// 	userId := ctx.Value(ContextKey.UserId).(string)

// 	return es.commentRepo.UpdateComment(userId, commentId, updatedComment)
// }

// func (es *EmployeeService)AddComment(newComment comment.Comment) error {
// 	return es.commentRepo.AddComment(newComment)
// }

// func (es *EmployeeService) DeleteComment(ctx context.Context, commentId string) error {

// 	userId := ctx.Value(ContextKey.UserId).(string)
// 	return es.commentRepo.DeleteComment(userId, commentId)
// }
// func(es *EmployeeService)UpdateTaskStatus(ctx context.Context,taskId string,updatedStatus status.TaskStatus)error{
// 	userId := ctx.Value(ContextKey.UserId).(string)

// 	return es.EmpRepo.UpdateTaskStatus(userId,taskId,updatedStatus)
// }
package service
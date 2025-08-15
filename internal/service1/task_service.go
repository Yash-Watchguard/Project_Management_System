package service1

import (
	"context"
	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	
)

type TaskService struct{
	taskRepo    interfaces.TaskRepo
	empRepo  interfaces.EmployeeRepo
}

func NewTaskService(taskRepo interfaces.TaskRepo,empRepo interfaces.EmployeeRepo)*TaskService{
	return &TaskService{taskRepo: taskRepo,empRepo: empRepo}
}

func (ts *TaskService) ViewAllTask(ctx context.Context, projectId string) ([]task.Task, error) {
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

	if userRole != 0 {
		return []task.Task{}, errors.New("unauthorized access")
	}
	return ts.taskRepo.ViewAllTask(projectId)
}
func(ts *TaskService)CreateTask(ctx context.Context,managerid string,task task.Task)error{
	userID:=ctx.Value(ContextKey.UserId).(string)

	if userID!=managerid{
		return errors.New("unauthoeized access")
	}
	return ts.taskRepo.SaveTask(task)
}

func(ts *TaskService)DeleteTask(ctx context.Context,managerId string,taskId string)error{
	userId:=ctx.Value(ContextKey.UserId).(string)

	if userId!=managerId{
		return errors.New("unauthorizrd access")
	}
	return ts.taskRepo.DeleteTask(taskId)
}

func(ts *TaskService)GetAssigenedTask(ctx context.Context,empId string)([]task.Task,error){
    userId:=ctx.Value(ContextKey.UserId).(string)

	if userId!=empId{
		return []task.Task{},errors.New("unauthorized access")
	}

	return ts.taskRepo.ViewAssignedTask(empId)
}

func(ts *TaskService)UpdateTaskStatus(ctx context.Context,taskId string,updatedStatus status.TaskStatus)error{
	userId := ctx.Value(ContextKey.UserId).(string)

	return ts.taskRepo.UpdateTaskStatus(userId,taskId,updatedStatus)
}

package service1

import (
	

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	
	
	
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	
)

type TaskServiceInterface interface{
	ViewAllTask( projectId string) ([]task.Task, error)
	CreateTask(task task.Task)error
	DeleteTask(managerId string,taskId string)error
	GetAssigenedTask(empId string)([]task.Task,error)
	UpdateTaskStatus(userId string,taskId string,updatedStatus status.TaskStatus)error
	ViewAllAssignedTasksInProject(projectId string,emp string)([]task.Task,error)
}
type TaskService struct{
	taskRepo    interfaces.TaskRepo
	
}

func NewTaskService(taskRepo interfaces.TaskRepo)TaskServiceInterface{
	return &TaskService{taskRepo: taskRepo}
}

func (ts *TaskService) ViewAllTask( projectId string) ([]task.Task, error) {
	return ts.taskRepo.ViewAllTask(projectId)
}
func(ts *TaskService)CreateTask(task task.Task)error{
	
	return ts.taskRepo.SaveTask(task)
}

func(ts *TaskService)DeleteTask(managerId string,taskId string)error{
	
	return ts.taskRepo.DeleteTask(taskId)
}

func(ts *TaskService)GetAssigenedTask(empId string)([]task.Task,error){
	return ts.taskRepo.ViewAssignedTask(empId)
}

func(ts *TaskService)UpdateTaskStatus(userId string,taskId string,updatedStatus status.TaskStatus)error{
	

	return ts.taskRepo.UpdateTaskStatus(userId,taskId,updatedStatus)
}

func(ts *TaskService)ViewAllAssignedTasksInProject(projectId string,emp string)([]task.Task,error){
	return ts.taskRepo.ViewAllAssignedTasksInProject(projectId,emp)
}

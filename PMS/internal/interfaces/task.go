package interfaces

import "github.com/Yash-Watchguard/Tasknest/internal/model/task"

type TaskRepo interface {
	ViewAllTask(projectId string) ([]task.Task,error)
	SaveTask(task task.Task)error
	DeleteTask(taskId string)error
}
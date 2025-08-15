package interfaces

import( "github.com/Yash-Watchguard/Tasknest/internal/model/task"
status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
)

type TaskRepo interface {
	ViewAllTask(projectId string) ([]task.Task,error)
	SaveTask(task task.Task)error
	DeleteTask(taskId string)error
	ViewAssignedTask(empId string)([]task.Task,error)
	UpdateTaskStatus(userId string,empId string,updatedStatus status.TaskStatus)(error)
}
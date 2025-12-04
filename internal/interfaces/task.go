package interfaces

import( "github.com/Yash-Watchguard/Tasknest/internal/model/task"
status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
)
//go:generate mockgen -source=task.go -destination=../mocks/mock_taskrepository.go -package=mocks
type TaskRepo interface {
	ViewAllTask(projectId string) ([]task.Task,error)
	SaveTask(task task.Task)error
	DeleteTask(projectId,taskId,managerId,empId string)error
	ViewAssignedTask(empId string)([]task.Task,error)
	UpdateTaskStatus(userId string,empId string,updatedStatus status.TaskStatus)(error)
	ViewAllAssignedTasksInProject(projectId string,empId string)([]task.Task,error)
	ViewAllManagerTask(managerId string)([]task.Task,error)
	UpdateTask(projecId,taskId string,managerId string,updatedfields map[string]interface{})(error)
	UpdateTaskEmailId(taskId string,managerId string, updates map[string]interface{}) error
}
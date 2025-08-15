package interfaces

import (
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
)

type EmployeeRepo interface {
	ViewAssignedTask(empId string)([]task.Task,error)
	UpdateTaskStatus(userId string,empId string,updatedStatus status.TaskStatus)(error)
}
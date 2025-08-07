package task

import (
	"time"

	Priority "github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
)

type Task struct {
	TaskId      string `json:"task_id"`
	Tile        string `json:"titel"`
	Description string `json:"description"`
	Deadline    time.Time `json:"deadline"`
	TaskPriority Priority.Priority `json:"taskpriority"`
	TaskStatus   status.TaskStatus  `json:"taskstatus"`
	AssignedTo   string             `json:"assignesto"`
	ProjectId    string              `json:"projectid"`
	CreatesBy    string              `json:"ceatedby"`
}
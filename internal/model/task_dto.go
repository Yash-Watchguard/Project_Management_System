package model

import "time"

type TaskDto struct {
	TaskId             string    `json:"task_id" db:"task_id"`
	Title              string    `json:"title" db:"title"`
	Description        string    `json:"description" db:"description"`
	AcceptanceCriteria string    `json:"acceptance_criteria" db:"acceptance_criteria"`
	Deadline           time.Time `json:"deadline" db:"deadline"`
	TaskPriority       string    `json:"taskpriority" db:"task_priority"`
	TaskStatus         string    `json:"taskstatus" db:"task_status"`
	AssignedTo         string    `json:"assigned_to" db:"assigned_to"`
	ProjectId          string    `json:"project_id" db:"project_id"`
	CreatedBy          string    `json:"created_by" db:"created_by"`
}
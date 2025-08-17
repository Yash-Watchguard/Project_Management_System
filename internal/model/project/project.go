package project

import "time"

type Project struct {
	ProjectId          string    `json:"project_id" db:"project_id"`              
	ProjectName        string    `json:"project_name" db:"project_name"`           
	ProjectDescription string    `json:"project_description" db:"project_description"` 
	Deadline           time.Time `json:"deadline" db:"deadline"`                  
	CreatedBy          string    `json:"created_by" db:"created_by"`               
	AssignedManager    string    `json:"assigned_manager_id" db:"assigned_manager"`
}
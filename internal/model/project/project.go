package project

import "time"

type Project struct {
	ProjectId          string `json:"projectid"`
	ProjectName        string `json:"project_name"`
	ProjectDescription string `json:"project_description"`
	Deadline           time.Time `json:"deadline"`
	CreatedBy          string  `json:"created_by"`
	AssignedManager    string   `json:"assigned_manager_id"`
}
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

type DynamoProject struct {
	PK                  string            `json:"PK" dynamodbav:"PK"` 
	SK                  string            `json:"SK" dynamodbav:"SK"` 
	Assigned_manager    string            `json:"task_id" dynamodbav:"Assigned_manager"`
	Created_by          string            `json:"title" dynamodbav:"Created_by"`
	Project_deadline    string            `json:"description" dynamodbav:"Project_deadline"`
	Project_description string            `json:"acceptance_criteria" dynamodbav:"Project_description"`
	Project_id          string         `json:"deadline" dynamodbav:"Project_id"`
	Project_name        string `json:"taskpriority" dynamodbav:"Project_name"`
}
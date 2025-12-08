package service1

import (
	"errors"
	"time"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"

	Priority "github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	
	"github.com/Yash-Watchguard/Tasknest/internal/util"

	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
)

//go:generate mockgen -source=task_service.go -destination=../mocks/mock_taskservice.go -package=mocks
type TaskServiceInterface interface{
	ViewAllTask( projectId string) ([]task.Task, error)
	CreateTask(task task.Task)error
	DeleteTask(projectId,taskId,managerId string,empId string)error
	GetAssigenedTask(empId string)([]task.Task,error)
	UpdateTaskStatus(userId string,taskId string,updatedStatus status.TaskStatus)error
	ViewAllTasksInProject(projectId string,emp string)([]task.Task,error)
	GetAllManagerTask(managerId string)([]task.Task,error)
	UpdateTask(projectId,taskId string,managerId string,updates map[string]interface{})error
	GetSingleTask(creatorId,projectId,taskId string)([]task.Task,error)
	
}
type TaskService struct{
	taskRepo    interfaces.TaskRepo
}


func NewTaskService(taskRepo interfaces.TaskRepo)TaskServiceInterface{
	return &TaskService{taskRepo: taskRepo}
}
func(ts *TaskService)GetAllManagerTask(managerId string)([]task.Task,error){
	return ts.taskRepo.ViewAllManagerTask(managerId)
}
func (ts *TaskService) ViewAllTask( projectId string) ([]task.Task, error) {
	return ts.taskRepo.ViewAllTask(projectId)
}
func(ts *TaskService)CreateTask(task task.Task)error{
	
	return ts.taskRepo.SaveTask(task)
}

func(ts *TaskService)DeleteTask(projectId,taskId,managerId string,empId string)error{
	
	return ts.taskRepo.DeleteTask(projectId,taskId,managerId,empId)
}

func(ts *TaskService)GetAssigenedTask(empId string)([]task.Task,error){
	return ts.taskRepo.ViewAssignedTask(empId)
}

func(ts *TaskService)UpdateTaskStatus(userId string,taskId string,updatedStatus status.TaskStatus)error{
	

	return ts.taskRepo.UpdateTaskStatus(userId,taskId,updatedStatus)
}

func(ts *TaskService)ViewAllTasksInProject(projectId string,creator_id string)([]task.Task,error){
	return ts.taskRepo.ViewAllTasksInProject(projectId,creator_id)
}

func(ts *TaskService)UpdateTask(projectId,taskId string,managerId string,updates map[string]interface{})error{
     finalUpdates := make(map[string]interface{})

	 if title,ok:=updates["titel"].(string);ok{
		if len(title)!=0 {
			finalUpdates["Title"]=title
		}
	 }
	 if status,ok:=updates["status"].(string);ok{
		if len(status)!=0 {
			
			finalUpdates["TaskStatus"]=status
		}
	 }
	 if description,ok:=updates["description"].(string);ok{
		if len(description)!=0 {
			finalUpdates["Description"]=description
		}
	 }
	  if AcceptanceCriteria,ok:=updates["acceptanceCriteria"].(string);ok{
		if len(AcceptanceCriteria)!=0 {
			finalUpdates["AcceptanceCriteria"]=AcceptanceCriteria
		}
	 }

	 if deadline,ok:=updates["deadline"].(string);ok{
		if len(deadline)!=0 {
			deadlineParsed,err:=util.ParseDate(deadline)
			if err!=nil {
				return errors.New("invalid date")
			}
			finalUpdates["Deadline"]=deadlineParsed.Format(time.RFC3339)
		}
	 }

	 if taskpriorityStr,ok:=updates["task_priority"].(string);ok{
		if len(taskpriorityStr)!=0 {
			priority, err := Priority.PriorityParser(taskpriorityStr)
			if err != nil {
				return errors.New("invalid priority")
			}
			finalUpdates["TaskPriority"]=Priority.GetPriority(priority)
		}
	 }

	 if empId,ok:=updates["empId"].(string);ok{
		if(len(empId)!=0){
			finalUpdates["AssignedTo"]=empId
		}
	 
}
return ts.taskRepo.UpdateTask(projectId,taskId,managerId,finalUpdates)
}

func(ts *TaskService)GetSingleTask(creatorId , projectId , taskId string)([]task.Task,error){
	return ts.taskRepo.GetSingleTask(creatorId,projectId,taskId)
}


/*


TaskId             string            `json:"task_id" db:"task_id"`                    
	Title              string            `json:"title" db:"title"`                        
	Description        string            `json:"description" db:"description"`             
	AcceptanceCriteria string            `json:"acceptance_criteria" db:"acceptance_criteria"`
	Deadline           time.Time         `json:"deadline" db:"deadline"`                  
	TaskPriority       Priority.Priority `json:"taskpriority" db:"task_priority"`         
	TaskStatus         status.TaskStatus `json:"taskstatus" db:"task_status"`              
	AssignedTo         string            `json:"assigned_to" db:"assigned_to"`            
	ProjectId          string            `json:"project_id" db:"project_id"`              
	CreatedBy          string            `json:"created_by" db:"created_by"`  



*/
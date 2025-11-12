package service1

import (
	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"

	Priority "github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/Yash-Watchguard/Tasknest/internal/util"

	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
)

//go:generate mockgen -source=task_service.go -destination=../mocks/mock_taskservice.go -package=mocks
type TaskServiceInterface interface{
	ViewAllTask( projectId string) ([]task.Task, error)
	CreateTask(task task.Task)error
	DeleteTask(managerId string,taskId string)error
	GetAssigenedTask(empId string)([]task.Task,error)
	UpdateTaskStatus(userId string,taskId string,updatedStatus status.TaskStatus)error
	ViewAllAssignedTasksInProject(projectId string,emp string)([]task.Task,error)
	GetAllManagerTask(managerId string)([]task.Task,error)
	UpdateTask(taskId string,updates map[string]interface{})error
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

func(ts *TaskService)DeleteTask(managerId string,taskId string)error{
	
	return ts.taskRepo.DeleteTask(taskId)
}

func(ts *TaskService)GetAssigenedTask(empId string)([]task.Task,error){
	return ts.taskRepo.ViewAssignedTask(empId)
}

func(ts *TaskService)UpdateTaskStatus(userId string,taskId string,updatedStatus status.TaskStatus)error{
	

	return ts.taskRepo.UpdateTaskStatus(userId,taskId,updatedStatus)
}

func(ts *TaskService)ViewAllAssignedTasksInProject(projectId string,emp string)([]task.Task,error){
	return ts.taskRepo.ViewAllAssignedTasksInProject(projectId,emp)
}

func(ts *TaskService)UpdateTask(taskId string,updates map[string]interface{})error{
     finalUpdates := make(map[string]interface{})

	 if titel,ok:=updates["titel"].(string);ok{
		if len(titel)!=0 {
			finalUpdates["title"]=titel
		}
	 }
	 if description,ok:=updates["description"].(string);ok{
		if len(description)!=0 {
			finalUpdates["description"]=description
		}
	 }
	  if AcceptanceCriteria,ok:=updates["acceptanceCriteria"].(string);ok{
		if len(AcceptanceCriteria)!=0 {
			finalUpdates["acceptance_criteria"]=AcceptanceCriteria
		}
	 }

	 if deadline,ok:=updates["deadline"].(string);ok{
		if len(deadline)!=0 {
			deadline,err:=util.ParseDate(deadline)
			if err!=nil {
				return errors.New("invalid date")
			}
			finalUpdates["deadline"]=deadline
		}
	 }

	 if taskpriorityStr,ok:=updates["task_priority"].(string);ok{
		if len(taskpriorityStr)!=0 {
			priority, err := Priority.PriorityParser(taskpriorityStr)
			if err != nil {
				return errors.New("invalid priority")
			}
			finalUpdates["taskpriority"]=priority
		}
	 }

	 if empId,ok:=updates["empId"].(string);ok{
		if(len(empId)!=0){
           finalUpdates["assignesto"]=empId
		}
	 }

	 if(len(finalUpdates)==0){
		return errors.New("no valid fields to update")
	 }
     
	 return ts.taskRepo.UpdateTask(taskId,finalUpdates)
	 
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
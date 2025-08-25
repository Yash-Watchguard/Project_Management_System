package handler

import (
	// "context"
	// "errors"
	// "fmt"
	"net/http"
	// "strings"
	"encoding/json"
	"time"

	"github.com/Yash-Watchguard/Tasknest/internal/logger"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	Priority "github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/Yash-Watchguard/Tasknest/internal/util"

	// "github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	// "github.com/Yash-Watchguard/Tasknest/internal/model/task"
	// status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	// "github.com/fatih/color"
)

type TaskHandler struct {
	taskService *service1.TaskService
}

func NewTaskHandler(taskService *service1.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func(th * TaskHandler)GetTask(w http.ResponseWriter,r *http.Request){
	userId:=r.URL.Query().Get("assigned_id")

	projectId:=r.PathValue("project_id")
	role:=r.Context().Value(ContextKey.UserRole).(roles.Role)
	employeeId:=r.Context().Value(ContextKey.UserId).(string)

	if userId==""{
        if role==roles.Employee{
            logger.Error("unauthorized to get task")
			response.ErrorResponse(w,http.StatusForbidden,"Unauthorized to get tasks",403)
			return
		}

		// get all the tasks of the project
		tasks,err:=th.taskService.ViewAllTask(projectId)
		if err!=nil{
			logger.Error("error getting the tasks")
			response.ErrorResponse(w,http.StatusInternalServerError,"Error in fetching the tasks",500)
			return
		}

		if len(tasks)==0{
		logger.Error("No task Created")
		response.ErrorResponse(w,http.StatusNotFound,"No task Created",404)
		return
	    }

        logger.Info("Tasks retrived Successfully")
		response.SuccessResponse(w,tasks,"Tasks retrived Successfully",http.StatusOK)
		return

	}

        if userId!=employeeId{
			logger.Error("unauthorized to get tasks")
			response.ErrorResponse(w,http.StatusForbidden,"Unauthorized to get tasks",403)
			return
		}
		tasks,err:=th.taskService.ViewAllAssignedTasksInProject(projectId,userId)
        if err!=nil{
			logger.Error("error getting the tasks")
			response.ErrorResponse(w,http.StatusInternalServerError,"Error in fetching the tasks",500)
			return
		}
		if len(tasks)==0{
		logger.Error("No task assigned")
		response.ErrorResponse(w,http.StatusNotFound,"No task Assigned",404)
		return
	    }
        logger.Info("Tasks retrived Successfully")
		response.SuccessResponse(w, tasks, "Tasks retrived Successfully",http.StatusOK)
	
}
func(th *TaskHandler)AssignedTasks(w http.ResponseWriter,r *http.Request){
	empId:=r.PathValue("employee_id")

	userId:=r.Context().Value(ContextKey.UserId).(string)

	if empId!=userId{
		logger.Error("unauthorized to get task")
		response.ErrorResponse(w,http.StatusForbidden,"Unauthorized to get tasks",403)
		return
	}

	tasks,err:=th.taskService.GetAssigenedTask(empId)
	if err!=nil{
		logger.Error("error getting the tasks")
		response.ErrorResponse(w,http.StatusInternalServerError,"Error in fetching the tasks",500)
		return
	}
    if len(tasks)==0{
		logger.Error("No task assigned")
		response.ErrorResponse(w,http.StatusNotFound,"No task Assigned",404)
		return
	}
	logger.Info("Tasks retrived Successfully")
	response.SuccessResponse(w, tasks, "Tasks retrived Successfully",http.StatusOK)

}
func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {


	// apply rbac
    
    projectId := r.PathValue("project_id")
    managerId := r.Context().Value(ContextKey.UserId).(string)

    // Define expected request body
    var req struct {
        Title              string `json:"title"`
        Description        string `json:"description"`
        AcceptanceCriteria string `json:"acceptance_criteria"`
        Deadline           string `json:"deadline"`   
        Priority           string `json:"priority"`   
        AssignedTo         string `json:"assigned_to"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logger.Error("invalid request body")
        response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 400)
        return
    }

    // deadline, err := time.Parse("2006-01-02", req.Deadline)
	var deadline time.Time
	deadline,err:=util.ParseDate(req.Deadline)
    if err != nil {
        logger.Error("invalid deadline format")
        response.ErrorResponse(w, http.StatusBadRequest, "Invalid deadline format (use YYYY-MM-DD) or may be deadline in past", 400)
        return
    }

    // Parse priority
	
    priority, err := Priority.PriorityParser(req.Priority)
    if err != nil {
        logger.Error("invalid priority")
        response.ErrorResponse(w, http.StatusBadRequest, "Invalid priority. Use Low, Medium, or High.", 400)
        return
    }

    // Generate task ID
    taskId := GenerateUUID()

    // Create task object
    newTask := task.Task{
        TaskId:             taskId,
        Title:              req.Title,
        Description:        req.Description,
        AcceptanceCriteria: req.AcceptanceCriteria,
        Deadline:           deadline,
        TaskPriority:       priority,
        TaskStatus:         status.Pending,
        AssignedTo:         req.AssignedTo,
        ProjectId:          projectId,
        CreatedBy:          managerId,
    }

    if err := th.taskService.CreateTask(newTask); err != nil {
        logger.Error("failed to create task")
        response.ErrorResponse(w, http.StatusInternalServerError, "Failed to create task", 500)
        return
    }

    logger.Info("Task created successfully")
    response.SuccessResponse(w, newTask, "Task created successfully", http.StatusCreated)
}

func (th *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    projectId := r.PathValue("project_id")
    taskId := r.PathValue("task_id")
    managerId := r.Context().Value(ContextKey.UserId).(string)
    role := r.Context().Value(ContextKey.UserRole).(roles.Role)

    if role == roles.Employee {
        logger.Error("unauthorized delete attempt by employee")
        response.ErrorResponse(w, http.StatusForbidden, "Unauthorized to delete tasks", 403)
        return
    }

    if projectId == "" || taskId == "" {
        logger.Error("missing project_id or task_id in request")
        response.ErrorResponse(w, http.StatusBadRequest, "Missing project_id or task_id", 400)
        return
    }
    
    if err := th.taskService.DeleteTask(managerId, taskId); err != nil {
        logger.Error("failed to delete task")
        response.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete task", 500)
        return
    }

    logger.Info("Task deleted successfully")
    response.SuccessResponse(w, nil, "Task deleted successfully", http.StatusOK)
}

func (th *TaskHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {

	// apply rbac
    // first get the task by task_id so write the function that get the only single task from the tasks
    userId:=r.Context().Value(ContextKey.UserId).(string)
    taskId := r.PathValue("task_id")
   
    var req struct {
        Status string `json:"status"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logger.Error("invalid request body")
        response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 400)
        return
    }

   
	
    newStatus,err := status.GetStatusFromString(req.Status)
	
    if err!=nil {
		logger.Error("Invalid status value")
        response.ErrorResponse(w, http.StatusBadRequest, "Invalid status value", 400)
        return
    }
    if err := th.taskService.UpdateTaskStatus( userId,taskId, newStatus); err != nil {
        logger.Error("failed to update task status")
        response.ErrorResponse(w, http.StatusInternalServerError, "Failed to update task status", 500)
        return
    }
    logger.Info("Task status updated successfully")
    response.SuccessResponse(w, nil, "Task status updated successfully", http.StatusOK)
}
















// func (th *TaskHandler) ViewAllTask(ctx context.Context, projectId string) error {
// 	if projectId == "" {
// 		return errors.New("projectId cannot be empty")
// 	}

// 	// Fetch tasks for the project
// 	tasks, err := th.taskService.ViewAllTask(ctx, projectId)
// 	if err != nil {
// 		color.Red("Failed to fetch tasks: %v", err)
// 		return err
// 	}

// 	// If no tasks exist
// 	if len(tasks) == 0 {

// 		return errors.New("no task found for this project")
// 	}

// 	// Display all tasks
// 	for i, task := range tasks {
// 		color.Cyan("------------ Task %d ------------", i+1)
// 		fmt.Printf("%-20s: %v\n", "Id", task.TaskId)
// 		fmt.Printf("%-20s: %v\n", "Title", task.Title)
// 		fmt.Printf("%-20s: %v\n", "Description", task.Description)
// 		fmt.Printf("%-20s: %v\n", "Acceptance_Criteria", task.AcceptanceCriteria)
// 		fmt.Printf("%-20s: %v\n", "Priority", Priority.GetPriority(task.TaskPriority))
// 		fmt.Printf("%-20s: %v\n", "Assigned To", task.AssignedTo)
// 		fmt.Printf("%-20s: %v\n", "Status", status.GetStatusString(task.TaskStatus))
// 		fmt.Printf("%-20s: %v\n", "Deadline", task.Deadline)
// 		fmt.Println()
// 	}

// 	return nil
// }
// func (th *TaskHandler) CreateTask(ctx context.Context, projectId string) error {
// 	managerId := ctx.Value(ContextKey.UserId).(string)
// 	taskId := GenerateUUID()

// 	title, err := GetInput("Enter Task Title : ")
// 	if err != nil {
// 		return err
// 	}

// 	description, err := GetInput("Enter Task Description : ")
// 	if err != nil {
// 		return err
// 	}

// 	acceptanceCriteria, err := GetInput("Enter Task Acceptance Criteria : ")
// 	if err != nil {
// 		return err
// 	}

// 	var deadline time.Time
// 	for {
// 		deadlineStr, err := GetInput("Enter Deadline in YYYY-MM-DD : ")
// 		if err != nil {
// 			return err
// 		}

// 		deadline, err = TimeParser(deadlineStr)
// 		if err != nil {
// 			color.Red("%v", err)
// 		} else {
// 			break
// 		}
// 	}

// 	var priority Priority.Priority
// 	for {
// 		priorityStr, err := GetInput("Enter Priority => Low/Medium/High : ")
// 		if err != nil {
// 			return err
// 		}

// 		priority, err = Priority.PriorityParser(priorityStr)
// 		if err != nil {
// 			color.Red("Invalid priority. Choose Low, Medium, or High.")
// 		} else {
// 			break
// 		}
// 	}

// 	assignedTo, err := GetInput("Enter Employee ID to assign this task to : ")
// 	if err != nil {
// 		return err
// 	}

// 	newTask := task.Task{
// 		TaskId:             taskId,
// 		Title:              title,
// 		Description:        description,
// 		AcceptanceCriteria: acceptanceCriteria,
// 		Deadline:           deadline,
// 		TaskPriority:       priority,
// 		TaskStatus:         status.Pending,
// 		AssignedTo:         assignedTo,
// 		ProjectId:          projectId,
// 		CreatedBy:          managerId,
// 	}

// 	// Call service method
// 	if err := th.taskService.CreateTask(ctx, managerId, newTask); err != nil {
// 		return err
// 	}

// 	color.Green("Task created successfully!")

// 	Pause()
// 	return nil
// }
// func (th *TaskHandler) DeleteTask(ctx context.Context, projectId string) error {
// 	managerId := ctx.Value(ContextKey.UserId).(string)

// 	// Fetch all tasks for the given project
// 	projectTasks, err := th.taskService.ViewAllTask(ctx, projectId)
// 	if err != nil {
// 		return err
// 	}

// 	if len(projectTasks) == 0 {
// 		return errors.New("no task created for this project")
// 	}

// 	// Display tasks
// 	for i, task := range projectTasks {
// 		color.Yellow("%d. Name: %s  ID: %s", i+1, task.Title, task.TaskId)
// 	}

// 	// Get task ID to delete
// 	taskId, err := GetInput("Enter Task Id to delete: ")
// 	if err != nil {
// 		return err
// 	}

// 	// Delete task using service
// 	if err := th.taskService.DeleteTask(ctx, managerId, taskId); err != nil {
// 		return err
// 	}

// 	color.Green("Task deleted successfully!")

// 	Pause()
// 	return nil
// }
// func (th *TaskHandler) GetAssignedTask(ctx context.Context, userId string) error {
// 	tasks, err := th.taskService.GetAssigenedTask(ctx, userId)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch assigned tasks: %v", err)
// 	}

// 	if len(tasks) == 0 {
// 		color.Yellow("No tasks assigned.")
// 		return nil
// 	}

// 	for i, task := range tasks {
// 		color.Cyan("------------ Task %d ------------", i+1)
// 		fmt.Printf("%-20s : %v\n", "Task ID", task.TaskId)
// 		fmt.Printf("%-20s : %v\n", "Title", task.Title)
// 		fmt.Printf("%-20s : %v\n", "Description", task.Description)
// 		fmt.Printf("%-20s : %v\n", "Priority", Priority.GetPriority(task.TaskPriority))
// 		fmt.Printf("%-20s : %v\n", "Acceptance Criteria", task.AcceptanceCriteria)
// 		fmt.Printf("%-20s : %v\n", "Status", status.GetStatusString(task.TaskStatus))
// 		fmt.Printf("%-20s : %v\n", "Created By", task.CreatedBy)
// 		fmt.Printf("%-20s : %v\n", "Deadline", task.Deadline.Format("2006-01-02 15:04:05"))
// 		fmt.Println()
// 	}

// 	return nil
// }

// func (th *TaskHandler) UpdateTaskStatus(ctx context.Context, taskId string) error {

// 	updatedStatus, err := GetInput("Enter Updated status (pending/in progress/done) : ")
// 	if err != nil {
// 		color.Red("error in getting input")
// 	}
// 	newStatus := status.GetStatusFromString(updatedStatus)
// 	err = th.taskService.UpdateTaskStatus(ctx, taskId, newStatus)
// 	if err != nil {
// 		return err
// 	}
// 	color.Green("Task Status Updated Successfully")
// 	Pause()
// 	return nil
// }

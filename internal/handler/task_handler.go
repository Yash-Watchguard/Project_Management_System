package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Yash-Watchguard/Tasknest/internal/logger"
	"github.com/Yash-Watchguard/Tasknest/internal/model"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	Priority "github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/util"

	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"

	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
)

type TaskHandler struct {
	taskService service1.TaskServiceInterface
	userService service1.UserServiceInterface
}

func NewTaskHandler(taskService service1.TaskServiceInterface, userService service1.UserServiceInterface) *TaskHandler {
	return &TaskHandler{taskService: taskService, userService: userService}
}
func(th *TaskHandler)GetAllMangerCreatedTask(w http.ResponseWriter, r *http.Request){
	user_Id :=r.PathValue("manager_id")

	newTasks, err := th.taskService.GetAllManagerTask(user_Id)
		if err != nil {
			logger.Error("error getting the tasks")
			response.ErrorResponse(w, http.StatusInternalServerError, "Error in fetching the tasks", 500)
			return
		}

		if len(newTasks) == 0 {
			logger.Error("No task Created")
			response.ErrorResponse(w, http.StatusNotFound, "No task Created", 1000)
			return
		}
		var tasks []model.TaskDto
		for _, newTask := range newTasks {
			taskDto := model.TaskDto{
				TaskId:             newTask.TaskId,
				Title:              newTask.Title,
				Description:        newTask.Description,
				AcceptanceCriteria: newTask.AcceptanceCriteria,
				Deadline:           newTask.Deadline,
				TaskPriority:       Priority.GetPriority(newTask.TaskPriority),
				TaskStatus:         status.GetStatusString(newTask.TaskStatus),
				AssignedTo:         newTask.AssignedTo,
				ProjectId:          newTask.ProjectId,
				CreatedBy:          newTask.CreatedBy,
			}
			tasks = append(tasks, taskDto)
		}

		logger.Info("Tasks retrived Successfully")
		response.SuccessResponse(w, tasks, "Tasks retrived Successfully", http.StatusOK)
}
func (th *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("assigned_id")

	projectId := r.PathValue("project_id")

	if len(projectId) != 36 {
		logger.Error("in url projectId is Invalid")
		response.ErrorResponse(w, http.StatusBadRequest, "ProjectId is Invalid", 1000)
		return
	}

	role := r.Context().Value(ContextKey.UserRole).(roles.Role)
	employeeId := r.Context().Value(ContextKey.UserId).(string)

	if userId == "" {
		if role == roles.Employee {
			logger.Error("unauthorized to get task")
			response.ErrorResponse(w, http.StatusForbidden, "Unauthorized to get tasks", 403)
			return
		}

		// get all the tasks of the project
		newTasks, err := th.taskService.ViewAllTask(projectId)
		if err != nil {
			logger.Error("error getting the tasks")
			response.ErrorResponse(w, http.StatusInternalServerError, "Error in fetching the tasks", 500)
			return
		}

		if len(newTasks) == 0 {
			logger.Error("No task Created")
			response.ErrorResponse(w, http.StatusNotFound, "No task Created", 1000)
			return
		}
		var tasks []model.TaskDto
		for _, newTask := range newTasks {
			taskDto := model.TaskDto{
				TaskId:             newTask.TaskId,
				Title:              newTask.Title,
				Description:        newTask.Description,
				AcceptanceCriteria: newTask.AcceptanceCriteria,
				Deadline:           newTask.Deadline,
				TaskPriority:       Priority.GetPriority(newTask.TaskPriority),
				TaskStatus:         status.GetStatusString(newTask.TaskStatus),
				AssignedTo:         newTask.AssignedTo,
				ProjectId:          newTask.ProjectId,
				CreatedBy:          newTask.CreatedBy,
			}
			tasks = append(tasks, taskDto)
		}

		logger.Info("Tasks retrived Successfully")
		response.SuccessResponse(w, tasks, "Tasks retrived Successfully", http.StatusOK)
		return

	}

	if userId != employeeId {
		if role != roles.Admin {
			logger.Error("unauthorized to get tasks")
			response.ErrorResponse(w, http.StatusForbidden, "Unauthorized to get tasks", 403)
			return
		}
	}
	newTasks, err := th.taskService.ViewAllTasksInProject(projectId, userId)
	if err != nil {
		logger.Error("error getting the tasks")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error in fetching the tasks", 500)
		return
	}
	if len(newTasks) == 0 {
		logger.Error("No task assigned")
		response.ErrorResponse(w, http.StatusNotFound, "No task Assigned", 404)
		return
	}

	var tasks []model.TaskDto
	for _, newTask := range newTasks {
		taskDto := model.TaskDto{
			TaskId:             newTask.TaskId,
			Title:              newTask.Title,
			Description:        newTask.Description,
			AcceptanceCriteria: newTask.AcceptanceCriteria,
			Deadline:           newTask.Deadline,
			TaskPriority:       Priority.GetPriority(newTask.TaskPriority),
			TaskStatus:         status.GetStatusString(newTask.TaskStatus),
			AssignedTo:         newTask.AssignedTo,
			ProjectId:          newTask.ProjectId,
			CreatedBy:          newTask.CreatedBy,
		}
		tasks = append(tasks, taskDto)
	}

	logger.Info("Tasks retrived Successfully")
	response.SuccessResponse(w, tasks, "Tasks retrived Successfully", http.StatusOK)

}

func (th *TaskHandler) AssignedTasks(w http.ResponseWriter, r *http.Request) {
	empId := r.PathValue("employee_id")
	userId := r.Context().Value(ContextKey.UserId).(string)
	role := r.Context().Value(ContextKey.UserRole).(roles.Role)

	if empId != userId || role != roles.Employee {
		if role != roles.Admin {
			logger.Error("unauthorized to get task")
			response.ErrorResponse(w, http.StatusForbidden, "Unauthorized to get tasks", 403)
			return
		}
	}

	newTasks, err := th.taskService.GetAssigenedTask(empId)
	if err != nil {
		logger.Error("error getting the tasks")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error in fetching the tasks", 500)
		return
	}
	if len(newTasks) == 0 {
		logger.Error("No task assigned")
		response.ErrorResponse(w, http.StatusNotFound, "No task Assigned", 404)
		return
	}
	var tasks []model.TaskDto
	for _, newTask := range newTasks {
		taskDto := model.TaskDto{
			TaskId:             newTask.TaskId,
			Title:              newTask.Title,
			Description:        newTask.Description,
			AcceptanceCriteria: newTask.AcceptanceCriteria,
			Deadline:           newTask.Deadline,
			TaskPriority:       Priority.GetPriority(newTask.TaskPriority),
			TaskStatus:         status.GetStatusString(newTask.TaskStatus),
			AssignedTo:         newTask.AssignedTo,
			ProjectId:          newTask.ProjectId,
			CreatedBy:          newTask.CreatedBy,
		}
		tasks = append(tasks, taskDto)
	}

	logger.Info("Tasks retrived Successfully")
	response.SuccessResponse(w, tasks, "Tasks retrived Successfully", http.StatusOK)

}
func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("project_id")

	if projectId == "" {
		logger.Error("missing project_id in request path")
		response.ErrorResponse(w, http.StatusBadRequest, "Missing project_id in request path", 400)
		return
	}

	managerId := r.Context().Value(ContextKey.UserId).(string)

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

	AssignedEmpId := req.AssignedTo
	users, err := th.userService.ViewProfile(AssignedEmpId)
	if err != nil {
		logger.Error("invalid employee id")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid Employee Id", 1000)
		return
	}
	if users[0].Status == user.InActive {
		logger.Error("user is not available")
		response.ErrorResponse(w, http.StatusForbidden, "Employee Is Not Active", 1000)
		return
	}

	deadline, err := util.ParseDate(req.Deadline)
	if err != nil {
		logger.Error("invalid deadline format")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid deadline format (use YYYY-MM-DD)", 1000)
		return
	}

	priority, err := Priority.PriorityParser(req.Priority)
	if err != nil {
		logger.Error("invalid priority")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid priority. Use Low, Medium, or High.", 1000)
		return
	}

	req.AssignedTo = strings.TrimSpace(req.AssignedTo)
	taskId := GenerateUUID()

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
	taskDto := model.TaskDto{
		TaskId:             newTask.TaskId,
		Title:              newTask.Title,
		Description:        newTask.Description,
		AcceptanceCriteria: newTask.AcceptanceCriteria,
		Deadline:           newTask.Deadline,
		TaskPriority:       Priority.GetPriority(newTask.TaskPriority),
		TaskStatus:         status.GetStatusString(newTask.TaskStatus),
		AssignedTo:         newTask.AssignedTo,
		ProjectId:          newTask.ProjectId,
		CreatedBy:          newTask.CreatedBy,
	}
	logger.Info("Task created successfully")
	response.SuccessResponse(w, taskDto, "Task created successfully", http.StatusCreated)
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

	if err := th.taskService.DeleteTask(managerId, taskId,"","",); err != nil {
		logger.Error("failed to delete task")
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete task", 500)
		return
	}

	logger.Info("Task deleted successfully")
	response.SuccessResponse(w, nil, "Task deleted successfully", http.StatusOK)
}

func (th *TaskHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	
	userId := r.Context().Value(ContextKey.UserId).(string)
	taskId := r.PathValue("task_id")

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("invalid request body")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 400)
		return
	}

	newStatus, err := status.GetStatusFromString(req.Status)

	if err != nil {
		logger.Error("Invalid status value")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid status value", 400)
		return
	}
	if err := th.taskService.UpdateTaskStatus(userId, taskId, newStatus); err != nil {
		logger.Error("failed to update task status")
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to update task status", 500)
		return
	}
	logger.Info("Task status updated successfully")
	response.SuccessResponse(w, nil, "Task status updated successfully", http.StatusOK)
}

func(th *TaskHandler) UpdateTask(w http.ResponseWriter,r* http.Request){

	 projectId := r.PathValue("project_id")
	 taskId:=r.PathValue("task_id");
	 managerId := r.Context().Value(ContextKey.UserId).(string)

	 var updates map[string]interface{}

	 if err:=json.NewDecoder(r.Body).Decode(&updates);err!=nil{
		response.ErrorResponse(w,http.StatusBadRequest,"Invalid request body",1000)
		return
	 }

	 err:=th.taskService.UpdateTask(projectId,taskId,managerId,updates)

	 if err!=nil{
        response.ErrorResponse(w,http.StatusInternalServerError,err.Error(),1000)
		return
	 }
     response.SuccessResponse(w,nil,"Task Updated SuccessFully",http.StatusOK)

}

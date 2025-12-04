package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/middleware"
	"github.com/Yash-Watchguard/Tasknest/internal/model"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	Priority "github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/util"

	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var taskService service1.TaskServiceInterface
var userService service1.UserServiceInterface


func init() {
	dynamDbClient := config.GetDyanoDbCliebt()

	taskRepo := repository.NewTaskRepo(dynamDbClient, "TaskNest")
	taskService= service1.NewTaskService(taskRepo)

	userRepo:= repository.NewUserRepo(dynamDbClient,"TaskNest")
	userService=service1.NewUserService(userRepo)
}

func main() {
	lambda.Start(middleware.LambdaAuthMiddleWare(handler))
}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
       projectId := request.PathParameters["project_id"]
	if projectId == "" {
		return response.LambdaErrorResponse(nil, "Missing project_id in path", 400, 400), nil
	}

	managerId, ok := ctx.Value(ContextKey.UserId).(string) 
	if !ok || managerId == "" {
		return response.LambdaErrorResponse(nil, "Unauthorized: Manager ID not found in context", 401, 401), nil
	}

	var req struct {
		Title              string `json:"title"`
		Description        string `json:"description"`
		AcceptanceCriteria string `json:"acceptance_criteria"`
		Deadline           string `json:"deadline"`
		Priority           string `json:"priority"`
		AssignedTo         string `json:"assigned_to"`
	}

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.LambdaErrorResponse(nil, "Invalid request body", 400, 400), nil
	}

	assignedEmpId := strings.TrimSpace(req.AssignedTo)
	users, err := userService.ViewProfile(assignedEmpId) // You need to initialize userService
	if err != nil || len(users) == 0 {
		return response.LambdaErrorResponse(nil, "Invalid Employee Id", 1000, 400), nil
	}
	if users[0].Status == user.InActive {
		return response.LambdaErrorResponse(nil, "Employee Is Not Active", 1000, 403), nil
	}

	deadline, err := util.ParseDate(req.Deadline)
	if err != nil {
		return response.LambdaErrorResponse(nil, "Invalid deadline format (use YYYY-MM-DD)", 1000, 400), nil
	}

	taskPriority, err := Priority.PriorityParser(req.Priority)
	if err != nil {
		return response.LambdaErrorResponse(nil, "Invalid priority. Use Low, Medium, or High.", 1000, 400), nil
	}

	taskId := util.GenerateUniqueUUID() // Make sure you have a UUID generator

	newTask := task.Task{
		TaskId:             taskId,
		Title:              req.Title,
		Description:        req.Description,
		AcceptanceCriteria: req.AcceptanceCriteria,
		Deadline:           deadline,
		TaskPriority:       taskPriority,
		TaskStatus:         status.Pending,
		AssignedTo:         assignedEmpId,
		ProjectId:          projectId,
		CreatedBy:          managerId,
	}

	if err := taskService.CreateTask(newTask); err != nil { // Initialize taskService
		return response.LambdaErrorResponse(nil, "Failed to create task", 500, 500), nil
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

	return response.LambdaSuccessResponse(taskDto, "Task created successfully", 200, 201), nil
}
package main

import (
	"context"
	"encoding/json"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/middleware"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)
var taskService service1.TaskServiceInterface

func init() {
	dynamoClinent := config.GetDyanoDbCliebt()

	taskRepo:= repository.NewTaskRepo(dynamoClinent,"TaskNest")

	taskService= service1.NewTaskService(taskRepo)
}

func main() {
     lambda.Start(middleware.WithCORS(middleware.LambdaAuthMiddleWare(handler)))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
	taskId := req.PathParameters["task_id"]
	managerId := req.PathParameters["manager_id"]
	projectId :=req.PathParameters["project_id"]
	
	if taskId == "" {
		return response.LambdaErrorResponse(nil, "Missing task_id in path", 400, 400), nil
	}


	var updates map[string]interface{}
	if err := json.Unmarshal([]byte(req.Body), &updates); err != nil {
		return response.LambdaErrorResponse(nil, "Invalid request body", 400, 400), nil
	}

	err := taskService.UpdateTask(projectId,taskId,managerId, updates)
	if err != nil {
		return response.LambdaErrorResponse(nil, err.Error(), 500, 500), nil
	}

	return response.LambdaSuccessResponse(nil, "Task updated successfully", 200, 200), nil
}

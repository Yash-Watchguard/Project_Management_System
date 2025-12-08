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

var projectService service1.ProjectServiceInterface

func init() {
dynamoClient:=config.GetDyanoDbCliebt()
taskRepo := repository.NewTaskRepo(dynamoClient,"TaskNest")

projectRepo:=repository.NewProjectRepo(dynamoClient,"TaskNest",*taskRepo)

projectService=service1.NewProjectService(projectRepo)
}

func main() {
   lambda.Start(middleware.WithCORS(middleware.LambdaAuthMiddleWare(handler)))
}

func handler(ctx context.Context,req events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,error) {
    creatorId := req.PathParameters["creator_id"]
	projectId := req.PathParameters["project_id"]
	managerId :=req.PathParameters["manager_id"]

	var updates map[string]interface{}
	if err := json.Unmarshal([]byte(req.Body), &updates); err != nil {
		return response.LambdaErrorResponse(nil, "Invalid request body", 400, 400), nil
	}

	err:= projectService.UpdateProject(creatorId,projectId,managerId,updates)
	if err != nil {
		return response.LambdaErrorResponse(nil, "error in updating the project"+err.Error(), 500, 500), nil
	}

	return response.LambdaSuccessResponse(nil, "Project updated successfully", 200, 200), nil
}
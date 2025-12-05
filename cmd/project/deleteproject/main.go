package main

import (
	"context"
	"net/http"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/middleware"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
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
   lambda.Start(middleware.LambdaAuthMiddleWare(handler))
}

func handler(ctx context.Context,req events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,error) {
  role, _ := ctx.Value(ContextKey.UserRole).(roles.Role)

    if role != roles.Admin {
     
        return response.LambdaErrorResponse(nil,"Only admin can delete projects",http.StatusForbidden,http.StatusForbidden),nil
        
    }

    
    
    projectId := req.PathParameters["project_id"]
	managerId := req.PathParameters["manager_id"]
	creatorId:= ctx.Value(ContextKey.UserId).(string)
    if projectId == "" || managerId =="" {
        return response.LambdaErrorResponse(nil,"invalid request",http.StatusBadRequest,http.StatusBadRequest),nil
    }

    
    err := projectService.DeleteProject(creatorId,managerId,projectId)
    if err != nil {
       return response.LambdaErrorResponse(nil,"internalser"+err.Error(),http.StatusInternalServerError,http.StatusInternalServerError),nil
    }

    return response.LambdaSuccessResponse(nil, "Project deleted successfully", http.StatusOK,http.StatusOK),nil
}
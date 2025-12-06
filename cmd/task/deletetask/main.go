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

var taskService service1.TaskServiceInterface

func init() {
   dynamoClient:= config.GetDyanoDbCliebt()

   taskRepo:=repository.NewTaskRepo(dynamoClient,"TaskNest")

   taskService=service1.NewTaskService(taskRepo)
}

func main(){
    lambda.Start(middleware.WithCORS(middleware.LambdaAuthMiddleWare(handler)))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    projectId := req.PathParameters["project_id"]
    taskId := req.PathParameters["task_id"]
	managerId:= req.PathParameters["manager_id"]
	empId:= req.PathParameters["emp_id"]

   
    roleStr:= ctx.Value(ContextKey.UserRole).(roles.Role)

    if roleStr == roles.Employee {
        return response.LambdaErrorResponse(nil,"Unauthorized to delete tasks",1000, http.StatusForbidden), nil
    }

    if projectId == "" || taskId == "" || managerId =="" || empId=="" {
        return response.LambdaErrorResponse(nil,"invalid request", 1000,http.StatusBadRequest), nil
    }

    if err := taskService.DeleteTask(projectId,taskId,managerId, empId); err != nil {
        return response.LambdaErrorResponse(nil,"Failed to delete task"+err.Error(),1000, http.StatusInternalServerError), nil
    }

    return response.LambdaSuccessResponse(nil, "Task deleted successfully",1000, http.StatusOK), nil
}

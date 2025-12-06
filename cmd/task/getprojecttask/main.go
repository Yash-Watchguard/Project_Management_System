package main

import (
	"context"
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/logger"
	"github.com/Yash-Watchguard/Tasknest/internal/middleware"
	"github.com/Yash-Watchguard/Tasknest/internal/model"
	
	Priority "github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	
	Status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
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

	creatorId := req.PathParameters["creator_id"]

    projectId := req.PathParameters["project_id"]
	

	newTasks, err := taskService.ViewAllTasksInProject(projectId,creatorId)
	
	if err != nil {
fmt.Print(newTasks)
		logger.Error("error getting the tasks: " + err.Error())
		return response.LambdaErrorResponse(nil, "Error in fetching the tasks", 500, 500), nil
	}

	if len(newTasks) == 0 {
		logger.Info("No task assigned")
		
		return response.LambdaErrorResponse(nil, "No task Assigned", 404, 404), nil
	}

	var tasks []model.TaskDto
	for _, t := range newTasks {
		tasks = append(tasks, model.TaskDto{
			TaskId:             t.TaskId,
			Title:              t.Title,
			Description:        t.Description,
			AcceptanceCriteria: t.AcceptanceCriteria,
			Deadline:           t.Deadline,
			TaskPriority:       Priority.GetPriority(t.TaskPriority),
			TaskStatus:         Status.GetStatusString(t.TaskStatus),
			AssignedTo:         t.AssignedTo,
			ProjectId:          t.ProjectId,
			CreatedBy:          t.CreatedBy,
		})
	}

	logger.Info("Tasks retrieved successfully")
	return response.LambdaSuccessResponse(tasks, "Tasks retrieved successfully", 200, 200), nil
}

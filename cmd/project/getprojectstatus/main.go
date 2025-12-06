package main

import (
	"context"
	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/middleware"
	Status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)


var taskService service1.TaskServiceInterface

func init() {
dynamoClient:=config.GetDyanoDbCliebt()
taskRepo := repository.NewTaskRepo(dynamoClient,"TaskNest")


taskService= service1.NewTaskService(taskRepo)

}

func main() {
   lambda.Start(middleware.WithCORS(middleware.LambdaAuthMiddleWare(handler)))
}

func handler(ctx context.Context,req events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,error) {
    creatorId := req.PathParameters["creator_id"]

    projectId := req.PathParameters["project_id"]
	

	projectTasks, err := taskService.ViewAllTasksInProject(projectId,creatorId)

	if err!=nil{
		return response.LambdaErrorResponse(nil, "Error in fetching the tasks", 500, 500), nil
	}
    //   calculate the project status on the basis of the tasks and task  status
     
    
	// if len(projectTasks) == 0 {
	// 	logger.Error("No task found")
	// 	response.ErrorResponse(w,http.StatusNotFound,"No Task",1000)
	// 	return
	// }
    
	total := len(projectTasks)
	done := 0

    if(total>0){
	for _, t := range projectTasks {
		if t.TaskStatus == Status.Done {
			done++
		}
	}
}
    percerntDone:=float64(0)
	if(total>0){
	percerntDone = (float64(done) / float64(total)) * 100
	}

	type ProjectStatusResponse struct {
        ProjectID            string  `json:"projectId"`
        CompletedTasks       int     `json:"completedTasks"`
        TotalTasks           int     `json:"totalTasks"`
        CompletionPercentage float64 `json:"completionPercentage"`
    }

    statusResponse := &ProjectStatusResponse{
        ProjectID:            projectId,
        CompletedTasks:       done,
        TotalTasks:           total,
        CompletionPercentage: percerntDone,
    }

	
	return response.LambdaSuccessResponse(statusResponse,"status fatced",200,200),nil
}
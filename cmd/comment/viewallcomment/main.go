package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/middleware"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var commentService service1.CommentServiceInterface
var userserVice service1.UserServiceInterface

func init() {
	dynamDbClient := config.GetDyanoDbCliebt()

	commentRepo := repository.NewCommentRepo(*dynamDbClient, "TaskNest")
	useerRepo := repository.NewUserRepo(dynamDbClient, "TaskNest")

	commentService = service1.NewCommentService(commentRepo)
	userserVice = service1.NewUserService(useerRepo)

}

func main() {
	lambda.Start(middleware.LambdaAuthMiddleWare(handler))
}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	taskId := request.PathParameters["task_id"]

	if taskId == "" {
		return response.LambdaErrorResponse(nil, "taskId missing in request path", 1001, http.StatusBadRequest), nil

	}

	comments, err := commentService.ViewAllComment(taskId)

	if err != nil {
		return response.LambdaErrorResponse(nil, fmt.Sprint("Internal server error", err.Error()), 1001, http.StatusInternalServerError), nil
	}
	if len(comments) == 0 {
		return response.LambdaErrorResponse(nil, "no comments found for this task", 1001, http.StatusNotFound), nil
	}
	type CommentResponse struct {
		CommentId string `json:"comment_id"`
		CreatedBy string `json:"created_by"`
		Content   string `json:"content"`
	}

	var resp []CommentResponse
	for _, comment := range comments {

		user, err := userserVice.ViewProfile(comment.CreatedBy)
		if err != nil {
			return response.LambdaErrorResponse(nil, "error in getting the user of the commenter id", 1001, http.StatusInternalServerError), nil
		}
		if len(user) == 0 {
			return response.LambdaErrorResponse(nil, fmt.Sprintf("No user found with ID %s", comment.CreatedBy), 1001, http.StatusNotFound), nil
		}

		resp = append(resp, CommentResponse{
			CommentId: comment.CommentId,
			CreatedBy: user[0].Name,
			Content:   comment.Content,
		})
	}

	return response.LambdaSuccessResponse(resp, "comments retrieved successfully", 1001, http.StatusOK), nil

}

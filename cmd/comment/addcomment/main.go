package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/logger"
	"github.com/Yash-Watchguard/Tasknest/internal/middleware"
	"github.com/Yash-Watchguard/Tasknest/internal/model"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/Yash-Watchguard/Tasknest/internal/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var commentService service1.CommentServiceInterface

func init(){
   dynmoDbClient := config.GetDyanoDbCliebt()

   commentRepo := repository.NewCommentRepo(*dynmoDbClient,"TaskNest")

   commentService = service1.NewCommentService(commentRepo)
}
func main() {
	lambda.Start(middleware.LambdaAuthMiddleWare(handler))
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    taskId := req.PathParameters["task_id"]
    if taskId == "" {
        logger.Error("missing taskId in request path")
        return response.LambdaErrorResponse(nil, "taskId cannot be empty", 400, http.StatusBadRequest), nil
    }

    createdBy, ok := ctx.Value(ContextKey.UserId).(string)
    if !ok || createdBy == "" {
        logger.Error("invalid or missing user ID in context")
        return response.LambdaErrorResponse(nil, "Unauthorized: missing user ID", 401, http.StatusUnauthorized), nil
    }

    var body struct {
        Content string `json:"content"`
    }

    if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
        logger.Error("invalid request body")
        return response.LambdaErrorResponse(nil, "Invalid request body", 400, http.StatusBadRequest), nil
    }

    if body.Content == "" {
        logger.Error("empty comment received")
        return response.LambdaErrorResponse(nil, "Comment cannot be empty", 400, http.StatusBadRequest), nil
    }

    newComment := comment.Comment{
        CommentId: util.GenerateUniqueUUID(),
        Content:   body.Content,
        CreatedBy: createdBy,
        TaskId:    taskId,
    }

    if err := commentService.AddComment(newComment); err != nil {
        logger.Error("failed to save comment")
        return response.LambdaErrorResponse(nil, "Failed to add comment", 500, http.StatusInternalServerError), nil
    }

    commentDto := model.CommentDto{
        CommentId: newComment.CommentId,
        Content:   newComment.Content,
    }

    logger.Info("comment added successfully")
    return response.LambdaSuccessResponse(
        commentDto,
        "Comment added successfully!",
        201,
        http.StatusCreated,
    ), nil
}

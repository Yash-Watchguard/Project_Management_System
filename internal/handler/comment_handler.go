package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Yash-Watchguard/Tasknest/internal/logger"
	"github.com/Yash-Watchguard/Tasknest/internal/model"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/Yash-Watchguard/Tasknest/internal/util"
)

type CommentHandler struct {
	userservice    service1.UserServiceInterface
	commentService service1.CommentServiceInterface
}

func NewCommentHandler(commentService service1.CommentServiceInterface, userService service1.UserServiceInterface) *CommentHandler {
	return &CommentHandler{commentService: commentService, userservice: userService}
}

func (ch *CommentHandler) ViewAllComment(w http.ResponseWriter, r *http.Request) {
    taskId := r.PathValue("task_id")

    if taskId == "" {
        logger.Error("taskId missing in request path")
        response.ErrorResponse(w, http.StatusBadRequest, "taskId cannot be empty", 400)
        return
    }

    comments, err := ch.commentService.ViewAllComment(taskId)
    if err != nil {
        logger.Error("failed to fetch comments")
        response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch comments", 500)
        return
    }

    if len(comments) == 0 {
        logger.Error("no comments found for task")
        response.ErrorResponse(w, http.StatusNotFound, "No comments found for this task", 404)
        return
    }

    // prepare enriched response
    type CommentResponse struct {
        CommentId string `json:"comment_id"`
        CreatedBy string `json:"created_by"`
        Content   string `json:"content"`
    }

    var resp []CommentResponse
    for _, comment := range comments {
        user, err := ch.userservice.ViewProfile(comment.CreatedBy)
        if err != nil {
            logger.Error("failed to fetch user profile")
            response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user profile", 500)
            return
        }
        if len(user) == 0 {
            logger.Error("no user found with given ID")
            response.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("No user found with ID %s", comment.CreatedBy), 404)
            return
        }

        resp = append(resp, CommentResponse{
            CommentId: comment.CommentId,
            CreatedBy: user[0].Name,
            Content:   comment.Content,
        })
    }

    logger.Info("comments retrieved successfully")
    response.SuccessResponse(w, resp, "Comments retrieved successfully", http.StatusOK)
}

func (ch *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    taskId := r.PathValue("task_id")
    if taskId == "" {
        logger.Error("missing taskId in request path")
        response.ErrorResponse(w, http.StatusBadRequest, "taskId cannot be empty", 400)
        return
    }

    createdBy, ok := ctx.Value(ContextKey.UserId).(string)
    if !ok || createdBy == "" {
        logger.Error("invalid or missing user ID in context")
        response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: missing user ID", 401)
        return
    }

    var req struct {
        Content string `json:"content"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logger.Error("invalid request body")
        response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 400)
        return
    }
    if req.Content == "" {
        logger.Error("empty comment received")
        response.ErrorResponse(w, http.StatusBadRequest, "Comment cannot be empty", 400)
        return
    }

    newComment := comment.Comment{
        CommentId: util.GenerateUniqueUUID(),
        Content:   req.Content,
        CreatedBy: createdBy,
        TaskId:    taskId,
    }

    if err := ch.commentService.AddComment(newComment); err != nil {
        logger.Error("failed to save comment")
        response.ErrorResponse(w, http.StatusInternalServerError, "Failed to add comment", 500)
        return
    }
    
    commentDto:=model.CommentDto{
        CommentId: newComment.CommentId,
        Content: newComment.Content,
    }
    logger.Info("comment added successfully")
    response.SuccessResponse(w, commentDto, " Comment added successfully!", http.StatusCreated)
}

func (ch *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    taskId := r.PathValue("task_id")
    commentId := r.PathValue("comment_id")

    if taskId == "" || commentId == "" {
        logger.Error("missing taskId or commentId in request path")
        response.ErrorResponse(w, http.StatusBadRequest, "taskId and commentId cannot be empty", 400)
        return
    }

    createdBy, ok := ctx.Value(ContextKey.UserId).(string)
    if !ok || createdBy == "" {
        logger.Error("invalid or missing user ID in context")
        response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: missing user ID", 401)
        return
    }

    // Parse request body
    var req struct {
        Content string `json:"content"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logger.Error("invalid request bod")
        response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 400)
        return
    }
    if req.Content == "" {
        logger.Error("empty content received for comment update")
        response.ErrorResponse(w, http.StatusBadRequest, "Updated comment cannot be empty", 400)
        return
    }

    // Build updated comment
    updatedComment := comment.Comment{
        CommentId: commentId,
        Content:   req.Content,
        TaskId:    taskId,
        CreatedBy: createdBy,
    }

    if err := ch.commentService.UpdateComment(ctx, updatedComment); err != nil {
    switch {
    case errors.Is(err, repository.ErrUnauthorized):
        response.ErrorResponse(w, http.StatusForbidden, err.Error(), 403)
    case errors.Is(err, repository.ErrCommentNotFound):
        response.ErrorResponse(w, http.StatusNotFound, err.Error(), 404)
    case errors.Is(err, repository.ErrWrongTask):
        response.ErrorResponse(w, http.StatusBadRequest, err.Error(), 400)
    default:
        response.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error", 500)
    }
    return
}

    logger.Info("comment updated successfully")
    response.SuccessResponse(w, nil, " Comment updated successfully!", http.StatusOK)
}


func (ch *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    commentId := r.PathValue("comment_id")

    if commentId == "" {
        logger.Error("Missing commentId in request")
        response.ErrorResponse(w, http.StatusBadRequest, "commentId is required", 400)
        return
    }

    if err := ch.commentService.DeleteComment(ctx, commentId); err != nil {
       logger.Error(" Failed to delete commentId")

       response.ErrorResponse(w,http.StatusForbidden,err.Error(),403)
       return
    }

    logger.Info("Comment deleted successfully")
    response.SuccessResponse(w, nil, " Comment deleted successfully!", http.StatusOK)
}








package handler

import (
	"context"

	"github.com/Yash-Watchguard/Tasknest/internal/service1"

	"errors"
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	
	"github.com/fatih/color"
)

type CommentHandler struct {
	userservice *service1.UserService
	commentService *service1.CommentService
}
func NewCommentHandler(commentService *service1.CommentService,userService *service1.UserService)*CommentHandler{
	return &CommentHandler{commentService: commentService,userservice: userService}
}

func (ch *CommentHandler) ViewAllComment(ctx context.Context, taskId string) error {
	if taskId == "" {
		return errors.New("taskId cannot be empty")
	}

	comments, err := ch.commentService.ViewAllComment(taskId)
	if err != nil {
		return err
	}

	if len(comments) == 0 {
		return errors.New("no comments found for this task")
	}

	for idx, comment := range comments {
		color.Blue("-------- Comment %d --------", idx+1)
		color.Yellow("Task ID: %v", comment.TaskId)
		color.Yellow("Created By: %v", comment.CreatorName)
		color.Cyan("Content: %v", comment.Content)
		color.Blue("----------------------------------------")
	}

	color.Green("Press Enter to return to the previous menu...")
	fmt.Scanln()
	return nil
}

func (ch *CommentHandler) AddComment(ctx context.Context, taskId string) error {
	if taskId == "" {
		return errors.New("taskId cannot be empty")
	}
	createdBy := ctx.Value(ContextKey.UserId).(string)
    creator,err:=ch.userservice.ViewProfile(ctx,createdBy)

	if err!=nil{
		return err
	}
	
	// Get comment content from user
	content, err := GetInput("Enter your comment:")
	if err!=nil{
		return err
	}

	
	commentId := GenerateUUID()

	newComment := comment.Comment{
		CommentId: commentId,
		Content:   content,
		CreatedBy: createdBy,
		TaskId:    taskId,
		CreatorName: creator[0].Name,
	}

	// Add comment via service
	err = ch.commentService.AddComment(newComment)
	if err != nil {
		return err
	}

	color.Green("✅ Comment added successfully!")
	color.Blue("Press Enter to go back...")
	fmt.Scanln()

	return nil
}

func (ch *CommentHandler) UpdateComment(ctx context.Context, taskId string) error {
	if taskId == "" {
		return errors.New("taskId cannot be empty")
	}

	// Ask for the comment ID to update
	commentId, err := GetInput("Enter Comment ID to update:")
	if err != nil || commentId == "" {
		return errors.New("comment ID cannot be empty")
	}

	// Ask for the new comment content
	newContent, err := GetInput("Enter updated comment content:")
	if err != nil || newContent == "" {
		return errors.New("updated comment cannot be empty")
	}

	updatedComment := comment.Comment{
		CommentId: commentId,
		Content:   newContent,
		TaskId:    taskId,
		CreatedBy: ctx.Value(ContextKey.UserId).(string),
	}

	// Call the service layer
	err = ch.commentService.UpdateComment(ctx, updatedComment)
	if err != nil {
		return err
	}

	color.Green("✅ Comment updated successfully!")
	color.Blue("Press Enter to go back...")
	fmt.Scanln()

	return nil
}

func (ch *CommentHandler) DeleteComment(ctx context.Context, taskId string) error {
	if taskId == "" {
		return errors.New("taskId cannot be empty")
	}

	// Get comment ID to delete
	commentId, err := GetInput("Enter Comment ID to delete:")
	if err != nil || commentId == "" {
		return errors.New("comment ID cannot be empty")
	}

	// Call service layer
	err = ch.commentService.DeleteComment(ctx, commentId)
	if err != nil {
		return err
	}

	color.Green("✅ Comment deleted successfully!")
	color.Blue("Press Enter to go back...")
	fmt.Scanln()

	return nil
}




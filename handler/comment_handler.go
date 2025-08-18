package handler

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"os"
	"strings"

	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"

	"github.com/fatih/color"
)

var reader = bufio.NewReader(os.Stdin)

type CommentHandler struct {
	userservice    *service1.UserService
	commentService *service1.CommentService
}

func NewCommentHandler(commentService *service1.CommentService, userService *service1.UserService) *CommentHandler {
	return &CommentHandler{commentService: commentService, userservice: userService}
}

func getInput(prompt string) (string, error) {
	fmt.Print(color.RedString(prompt))
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// Pause waits for user to press Enter
func Pause() {
	fmt.Print(color.BlueString("Press Enter to go back..."))
	reader.ReadString('\n') // ignore error intentionally
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
		user, err := ch.userservice.ViewProfile(ctx, comment.CreatedBy)
if err != nil {
    return fmt.Errorf("failed to fetch user profile: %w", err)
}
if len(user) == 0 {
    return fmt.Errorf("no user found with ID %s", comment.CreatedBy)
}


		color.Blue("-------- Comment %d --------", idx+1)
		// color.Yellow("Task ID: %v", comment.TaskId)
		color.Yellow("Created By: %v", user[0].Name)
		color.Cyan("Content: %v", comment.Content)
		color.Blue("----------------------------------------")
	}
	return nil
}

func (ch *CommentHandler) AddNewComment(ctx context.Context, taskId string) error {
	if taskId == "" {
		return errors.New("taskId cannot be empty")
	}

	createdBy, ok := ctx.Value(ContextKey.UserId).(string)
	if !ok || createdBy == "" {
		return errors.New("invalid or missing user ID in context")
	}

	content, err := getInput("Enter your comment: ")
	if err != nil {
		return err
	}
	if content == "" {
		return errors.New("comment cannot be empty")
	}

	commentId := GenerateUUID()

	newComment := comment.Comment{
		CommentId: commentId,
		Content:   content,
		CreatedBy: createdBy,
		TaskId:    taskId,
	}

	if err := ch.commentService.AddComment(newComment); err != nil {
		return err
	}

	color.Green("✅ Comment added successfully!")

	return nil
}

func (ch *CommentHandler) UpdateComment(ctx context.Context, taskId string) error {
	if taskId == "" {
		return errors.New("taskId cannot be empty")
	}

	commentId, err := getInput("Enter Comment ID to update: ")
	if err != nil || commentId == "" {
		return errors.New("comment ID cannot be empty")
	}

	newContent, err := getInput("Enter updated comment content: ")
	if err != nil || newContent == "" {
		return errors.New("updated comment cannot be empty")
	}

	updatedComment := comment.Comment{
		CommentId: commentId,
		Content:   newContent,
		TaskId:    taskId,
		CreatedBy: ctx.Value(ContextKey.UserId).(string),
	}

	if err := ch.commentService.UpdateComment(ctx, updatedComment); err != nil {
		return err
	}

	color.Green("✅ Comment updated successfully!")

	return nil
}

func (ch *CommentHandler) DeleteComment(ctx context.Context, taskId string) error {
	if taskId == "" {
		return errors.New("taskId cannot be empty")
	}

	commentId, err := getInput("Enter Comment ID to delete: ")
	if err != nil || commentId == "" {
		return errors.New("comment ID cannot be empty")
	}

	if err := ch.commentService.DeleteComment(ctx, commentId); err != nil {
		return err
	}

	color.Green("✅ Comment deleted successfully!")
	Pause()
	return nil
}

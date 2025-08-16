package handler

import (
	"context"
    "bufio"
	"os"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
"fmt"
"strings"
	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	
	"github.com/fatih/color"
)
var reader= bufio.NewReader(os.Stdin)

type CommentHandler struct {
	userservice *service1.UserService
	commentService *service1.CommentService
}
func NewCommentHandler(commentService *service1.CommentService,userService *service1.UserService)*CommentHandler{
	return &CommentHandler{commentService: commentService,userservice: userService}
}

// func (ch *CommentHandler) ViewAllComment(ctx context.Context, taskId string) error {
// 	if taskId == "" {
// 		return errors.New("taskId cannot be empty")
// 	}

// 	comments, err := ch.commentService.ViewAllComment(taskId)
// 	if err != nil {
// 		return err
// 	}

// 	if len(comments) == 0 {
// 		return errors.New("no comments found for this task")
// 	}

// 	for idx, comment := range comments {
// 		color.Blue("-------- Comment %d --------", idx+1)
// 		color.Yellow("Task ID: %v", comment.TaskId)
// 		color.Yellow("Created By: %v", comment.CreatedBy)
// 		color.Cyan("Content: %v", comment.Content)
// 		color.Blue("----------------------------------------")
// 	}

// 	color.Green("Press Enter to return to the previous menu...")
// 	reader.ReadString('\n')
// 	return nil
// }

// func (ch *CommentHandler)AddNewComment(ctx context.Context, taskId string) error {
// 	if taskId == "" {
// 		return errors.New("taskId cannot be empty")
// 	}

// 	createdBy, ok := ctx.Value(ContextKey.UserId).(string)
// 	if !ok || createdBy == "" {
// 		return errors.New("invalid or missing user ID in context")
// 	}

	
	

// 	// Get comment content from user
	
// 	fmt.Print(color.RedString("Enter your comment:"))
// 	content, _ := reader.ReadString('\n')
// 	content = strings.TrimSpace(content)

// 	if content == "" {
// 		return errors.New("comment cannot be empty")
// 	}

// 	commentId := GenerateUUID()

// 	newComment := comment.Comment{
// 		CommentId:   commentId,
// 		Content:     content,
// 		CreatedBy:   createdBy,
// 		TaskId:      taskId,
// 	}

	
// 	if err := ch.commentService.AddComment(newComment); err != nil {
// 		return err
// 	}

// 	color.Green("✅ Comment added successfully!")
// 	color.Blue("Press Enter to go back...")
// 	reader.ReadString('\n')

// 	return nil
// }


// func (ch *CommentHandler) UpdateComment(ctx context.Context, taskId string) error {
// 	if taskId == "" {
// 		return errors.New("taskId cannot be empty")
// 	}

// 	// Ask for the comment ID to update
// 	commentId, err := GetInput("Enter Comment ID to update:")
// 	if err != nil || commentId == "" {
// 		return errors.New("comment ID cannot be empty")
// 	}

// 	// Ask for the new comment content
// 	newContent, err := GetInput("Enter updated comment content:")
// 	if err != nil || newContent == "" {
// 		return errors.New("updated comment cannot be empty")
// 	}

// 	updatedComment := comment.Comment{
// 		CommentId: commentId,
// 		Content:   newContent,
// 		TaskId:    taskId,
// 		CreatedBy: ctx.Value(ContextKey.UserId).(string),
// 	}

// 	// Call the service layer
// 	err = ch.commentService.UpdateComment(ctx, updatedComment)
// 	if err != nil {
// 		return err
// 	}

// 	color.Green("✅ Comment updated successfully!")
// 	color.Blue("Press Enter to go back...")
// 	reader.ReadString('\n')

// 	return nil
// }

// func (ch *CommentHandler) DeleteComment(ctx context.Context, taskId string) error {
// 	if taskId == "" {
// 		return errors.New("taskId cannot be empty")
// 	}

// 	// Get comment ID to delete
// 	commentId, err := GetInput("Enter Comment ID to delete:")
// 	if err != nil || commentId == "" {
// 		return errors.New("comment ID cannot be empty")
// 	}

// 	// Call service layer
// 	err = ch.commentService.DeleteComment(ctx, commentId)
// 	if err != nil {
// 		return err
// 	}

// 	color.Green("✅ Comment deleted successfully!")
// 	color.Blue("Press Enter to go back...")
// 	reader.ReadString('\n')

// 	return nil
// }
func getInput(prompt string) (string, error) {
	fmt.Print(color.RedString(prompt))
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// Pause waits for user to press Enter
func pause() {
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
		color.Blue("-------- Comment %d --------", idx+1)
		color.Yellow("Task ID: %v", comment.TaskId)
		color.Yellow("Created By: %v", comment.CreatedBy)
		color.Cyan("Content: %v", comment.Content)
		color.Blue("----------------------------------------")
	}

	pause()
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
	pause()
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
	pause()
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
	pause()
	return nil
}





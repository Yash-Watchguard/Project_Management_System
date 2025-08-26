package service1

import (
	"context"
	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	
)

type CommentServiceInterface interface{
	ViewAllComment(taskId string) ([]comment.Comment, error)
	UpdateComment(ctx context.Context, updatedComment comment.Comment) error
	AddComment(newComment comment.Comment) error
	DeleteComment(ctx context.Context, commentId string) error
}

type CommentService struct{
	commentRepo interfaces.CommentRepo
}

func NewCommentService(commentRepo interfaces.CommentRepo) CommentServiceInterface{
	return &CommentService{commentRepo: commentRepo}
}

func (cs *CommentService) ViewAllComment(taskId string) ([]comment.Comment, error) {

	return cs.commentRepo.ViewAllComments(taskId)
}

func (cs *CommentService) UpdateComment(ctx context.Context, updatedComment comment.Comment) error {
	if updatedComment.CommentId == "" || updatedComment.Content == "" {
		return errors.New("comment ID and content cannot be empty")
	}

	return cs.commentRepo.UpdateComment(updatedComment)
}


func (cs *CommentService) AddComment(newComment comment.Comment) error {
	return cs.commentRepo.AddComment(newComment)
}

func (cs *CommentService) DeleteComment(ctx context.Context, commentId string) error {

	userId := ctx.Value(ContextKey.UserId).(string)
	return cs.commentRepo.DeleteComment(userId, commentId)
}
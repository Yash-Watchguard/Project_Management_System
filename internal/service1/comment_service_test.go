package service1

import (
	"context"
	"errors"
	"testing"

	"github.com/Yash-Watchguard/Tasknest/internal/mocks"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"go.uber.org/mock/gomock"
)

func TestCommentService_ViewAllComment_Table(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockCommentRepo(ctrl)
    service :=NewCommentService(mockRepo)

    taskID := "task1"
    expectedComments := []comment.Comment{
        {CommentId: "c1", Content: "First Comment"},
    }

    tests := []struct {
        name      string
        mockSetup func()
        wantErr   bool
        wantLen   int
    }{
        {
            name: "success",
            mockSetup: func() {
                mockRepo.EXPECT().ViewAllComments(taskID).Return(expectedComments, nil)
            },
            wantErr: false,
            wantLen: 1,
        },
        {
            name: "repo error",
            mockSetup: func() {
                mockRepo.EXPECT().ViewAllComments(taskID).Return(nil, errors.New("repo error"))
            },
            wantErr: true,
            wantLen: 0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            comments, err := service.ViewAllComment(taskID)
            if (err != nil) != tt.wantErr {
                t.Errorf("expected error=%v, got %v", tt.wantErr, err)
            }
            if len(comments) != tt.wantLen {
                t.Errorf("expected %d comments, got %d", tt.wantLen, len(comments))
            }
        })
    }
}

func TestCommentService_UpdateComment_Table(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCommentRepo(ctrl)
	service := NewCommentService(mockRepo)

	validComment := comment.Comment{CommentId: "c1", Content: "Updated content"}
	invalidComment := comment.Comment{CommentId: "", Content: ""}

	tests := []struct {
		name      string
		input     comment.Comment
		mockSetup func()
		wantErr   bool
	}{
		{
			name:  "success",
			input: validComment,
			mockSetup: func() {
				mockRepo.EXPECT().UpdateComment(validComment).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "validation error",
			input: invalidComment,
			mockSetup: func() {
				// no repo call expected
			},
			wantErr: true,
		},
		{
			name:  "repo error",
			input: validComment,
			mockSetup: func() {
				mockRepo.EXPECT().UpdateComment(validComment).Return(errors.New("repo error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.UpdateComment(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestCommentService_AddComment_Table(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCommentRepo(ctrl)
	service := NewCommentService(mockRepo)

	newComment := comment.Comment{CommentId: "c1", Content: "New Comment"}

	tests := []struct {
		name      string
		input     comment.Comment
		mockSetup func()
		wantErr   bool
	}{
		{
			name:  "success",
			input: newComment,
			mockSetup: func() {
				mockRepo.EXPECT().AddComment(newComment).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "repo error",
			input: newComment,
			mockSetup: func() {
				mockRepo.EXPECT().AddComment(newComment).Return(errors.New("repo error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.AddComment(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}


func TestCommentService_DeleteComment_Table(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCommentRepo(ctrl)
	service :=NewCommentService(mockRepo)

	userID := "u1"
	commentID := "c1"

	tests := []struct {
		name      string
		ctx       context.Context
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), ContextKey.UserId, userID),
			mockSetup: func() {
				mockRepo.EXPECT().DeleteComment(userID, commentID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			ctx:  context.WithValue(context.Background(), ContextKey.UserId, userID),
			mockSetup: func() {
				mockRepo.EXPECT().DeleteComment(userID, commentID).Return(errors.New("repo error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.DeleteComment(tt.ctx, commentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}



package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Yash-Watchguard/Tasknest/internal/mocks"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"go.uber.org/mock/gomock"
	repo "github.com/Yash-Watchguard/Tasknest/internal/repository"
)


func newCommentHandlerWithMock(t *testing.T) (*CommentHandler, *mocks.MockCommentServiceInterface, *mocks.MockUserServiceInterface, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	commentSvc := mocks.NewMockCommentServiceInterface(ctrl)
	userSvc := mocks.NewMockUserServiceInterface(ctrl)
	h := NewCommentHandler(commentSvc, userSvc)
	return h, commentSvc, userSvc, ctrl
}

func TestAddComment(t *testing.T) {
	tests := []struct {
		name         string
		taskId       string
		userId       string
		body         string
		mockSetup    func(cs *mocks.MockCommentServiceInterface)
		expectedCode int
		expectBody   string
	}{
		{
			name:  "Missing task_id",
			taskId: "",
			userId: "u1",
			body:  `{"content":"hi"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {},
			expectedCode: http.StatusBadRequest,
			expectBody:   "taskId cannot be empty",
		},
		{
			name:  "Missing user id",
			taskId: "t1",
			userId: "",
			body:  `{"content":"hi"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {},
			expectedCode: http.StatusUnauthorized,
			expectBody:   "Unauthorized: missing user ID",
		},
		{
			name:  "Invalid request body",
			taskId: "t1",
			userId: "u1",
			body:  `{invalid json}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {},
			expectedCode: http.StatusBadRequest,
			expectBody:   "Invalid request body",
		},
		{
			name:  "Empty content",
			taskId: "t1",
			userId: "u1",
			body:  `{"content":""}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {},
			expectedCode: http.StatusBadRequest,
			expectBody:   "Comment cannot be empty",
		},
		{
			name:  "Service error on add",
			taskId: "t1",
			userId: "u1",
			body:  `{"content":"hi"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {
				cs.EXPECT().AddComment(gomock.Any()).Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectBody:   "Failed to add comment",
		},
		{
			name:  "Success",
			taskId: "t1",
			userId: "u1",
			body:  `{"content":"hi"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {
				cs.EXPECT().AddComment(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusCreated,
			expectBody:   "Comment added successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, cs, _, ctrl := newCommentHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockSetup(cs)

			req := httptest.NewRequest(http.MethodPost, "/v1/tasks/"+tt.taskId+"/comments", strings.NewReader(tt.body))
			if tt.taskId != "" {
				req.SetPathValue("task_id", tt.taskId)
			}
			// set context user id only if provided (missing user id case leaves it absent)
			if tt.userId != "" {
				ctx := context.WithValue(req.Context(), ContextKey.UserId, tt.userId)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.AddComment(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
			if tt.expectBody != "" && !strings.Contains(w.Body.String(), tt.expectBody) {
				t.Errorf("%s: expected body to contain %q, got %q", tt.name, tt.expectBody, w.Body.String())
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	tests := []struct {
		name         string
		taskId       string
		commentId    string
		userId       string
		body         string
		mockSetup    func(cs *mocks.MockCommentServiceInterface)
		expectedCode int
		expectBody   string
	}{
		{
			name:  "Missing taskId or commentId",
			taskId: "",
			commentId: "",
			userId: "u1",
			body:  `{"content":"new"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {},
			expectedCode: http.StatusBadRequest,
			expectBody:   "taskId and commentId cannot be empty",
		},
		{
			name:  "Missing user id in context",
			taskId: "t1",
			commentId: "c1",
			userId: "",
			body:  `{"content":"new"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {},
			expectedCode: http.StatusUnauthorized,
			expectBody:   "Unauthorized: missing user ID",
		},
		{
			name:  "Invalid request body",
			taskId: "t1",
			commentId: "c1",
			userId: "u1",
			body:  `{invalid json}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {},
			expectedCode: http.StatusBadRequest,
			expectBody:   "Invalid request body",
		},
		{
			name:  "Empty content",
			taskId: "t1",
			commentId: "c1",
			userId: "u1",
			body:  `{"content":""}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {},
			expectedCode: http.StatusBadRequest,
			expectBody:   "Updated comment cannot be empty",
		},
		{
			name:  "Unauthorized service error",
			taskId: "t1",
			commentId: "c1",
			userId: "u1",
			body:  `{"content":"new"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {
				cs.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Return(repo.ErrUnauthorized)
			},
			expectedCode: http.StatusForbidden,
			expectBody:   repo.ErrUnauthorized.Error(),
		},
		{
			name:  "Comment not found service error",
			taskId: "t1",
			commentId: "c1",
			userId: "u1",
			body:  `{"content":"new"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {
				cs.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Return(repo.ErrCommentNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectBody:   repo.ErrCommentNotFound.Error(),
		},
		{
			name:  "Wrong task service error",
			taskId: "t1",
			commentId: "c1",
			userId: "u1",
			body:  `{"content":"new"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {
				cs.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Return(repo.ErrWrongTask)
			},
			expectedCode: http.StatusBadRequest,
			expectBody:   repo.ErrWrongTask.Error(),
		},
		{
			name:  "Generic service error",
			taskId: "t1",
			commentId: "c1",
			userId: "u1",
			body:  `{"content":"new"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {
				cs.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectBody:   "Internal Server Error",
		},
		{
			name:  "Success",
			taskId: "t1",
			commentId: "c1",
			userId: "u1",
			body:  `{"content":"updated content"}`,
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {
				cs.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectBody:   "Comment updated successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, cs, _, ctrl := newCommentHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockSetup(cs)

			req := httptest.NewRequest(http.MethodPut, "/v1/tasks/"+tt.taskId+"/comments/"+tt.commentId, strings.NewReader(tt.body))
			if tt.taskId != "" {
				req.SetPathValue("task_id", tt.taskId)
			}
			if tt.commentId != "" {
				req.SetPathValue("comment_id", tt.commentId)
			}
			if tt.userId != "" {
				ctx := context.WithValue(req.Context(), ContextKey.UserId, tt.userId)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.UpdateComment(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
			if tt.expectBody != "" && !strings.Contains(w.Body.String(), tt.expectBody) {
				t.Errorf("%s: expected body to contain %q, got %q", tt.name, tt.expectBody, w.Body.String())
			}
		})
	}
}

func TestViewAllComment(t *testing.T) {
	tests := []struct {
		name         string
		taskId       string
		mockSetup    func(cs *mocks.MockCommentServiceInterface, us *mocks.MockUserServiceInterface)
		expectedCode int
		expectBody   string
	}{
		{
			name:  "Missing task_id",
			taskId: "",
			mockSetup: func(cs *mocks.MockCommentServiceInterface, us *mocks.MockUserServiceInterface) {},
			expectedCode: http.StatusBadRequest,
			expectBody:   "taskId cannot be empty",
		},
		{
			name:  "Service error fetching comments",
			taskId: "t1",
			mockSetup: func(cs *mocks.MockCommentServiceInterface, us *mocks.MockUserServiceInterface) {
				cs.EXPECT().ViewAllComment("t1").Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectBody:   "Failed to fetch comments",
		},
		{
			name:  "No comments found",
			taskId: "t1",
			mockSetup: func(cs *mocks.MockCommentServiceInterface, us *mocks.MockUserServiceInterface) {
				cs.EXPECT().ViewAllComment("t1").Return([]comment.Comment{}, nil)
			},
			expectedCode: http.StatusNotFound,
			expectBody:   "No comments found for this task",
		},
		{
			name:  "User service error",
			taskId: "t1",
			mockSetup: func(cs *mocks.MockCommentServiceInterface, us *mocks.MockUserServiceInterface) {
				cs.EXPECT().ViewAllComment("t1").Return([]comment.Comment{{CommentId: "c1", CreatedBy: "u1", Content: "hello"}}, nil)
				us.EXPECT().ViewProfile("u1").Return(nil, errors.New("db"))
			},
			expectedCode: http.StatusInternalServerError,
			expectBody:   "Failed to fetch user profile",
		},
		{
			name:  "User not found",
			taskId: "t1",
			mockSetup: func(cs *mocks.MockCommentServiceInterface, us *mocks.MockUserServiceInterface) {
				cs.EXPECT().ViewAllComment("t1").Return([]comment.Comment{{CommentId: "c1", CreatedBy: "u1", Content: "hello"}}, nil)
				us.EXPECT().ViewProfile("u1").Return([]user.User{}, nil)
			},
			expectedCode: http.StatusNotFound,
			expectBody:   "No user found with ID u1",
		},
		{
			name:  "Success",
			taskId: "t1",
			mockSetup: func(cs *mocks.MockCommentServiceInterface, us *mocks.MockUserServiceInterface) {
				cs.EXPECT().ViewAllComment("t1").Return([]comment.Comment{{CommentId: "c1", CreatedBy: "u1", Content: "hello"}}, nil)
				us.EXPECT().ViewProfile("u1").Return([]user.User{{Name: "Alice"}}, nil)
			},
			expectedCode: http.StatusOK,
			expectBody:   "Comments retrieved successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, cs, us, ctrl := newCommentHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockSetup(cs, us)

			req := httptest.NewRequest(http.MethodGet, "/v1/tasks/"+tt.taskId+"/comments", nil)
			if tt.taskId != "" {
				req.SetPathValue("task_id", tt.taskId)
			}

			w := httptest.NewRecorder()
			h.ViewAllComment(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
			if tt.expectBody != "" && !strings.Contains(w.Body.String(), tt.expectBody) {
				t.Errorf("%s: expected body to contain %q, got %q", tt.name, tt.expectBody, w.Body.String())
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	tests := []struct {
		name         string
		commentId    string
		userId       string
		mockSetup    func(cs *mocks.MockCommentServiceInterface)
		expectedCode int
		expectBody   string
	}{
		{
			name:      "Missing comment_id",
			commentId: "",
			userId:    "u1",
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {},
			expectedCode: http.StatusBadRequest,
			expectBody:   "commentId is required",
		},
		{
			name:      "Service delete returns ErrDelete",
			commentId: "c1",
			userId:    "u1",
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {
				cs.EXPECT().DeleteComment(gomock.Any(), "c1").Return(repo.ErrDelete)
			},
			expectedCode: http.StatusForbidden,
			expectBody:   repo.ErrDelete.Error(),
		},
		{
			name:      "Service delete returns generic error",
			commentId: "c1",
			userId:    "u1",
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {
				cs.EXPECT().DeleteComment(gomock.Any(), "c1").Return(errors.New("db error"))
			},
			expectedCode: http.StatusForbidden,
			expectBody:   "db error",
		},
		{
			name:      "Successful delete",
			commentId: "c1",
			userId:    "u1",
			mockSetup: func(cs *mocks.MockCommentServiceInterface) {
				cs.EXPECT().DeleteComment(gomock.Any(), "c1").Return(nil)
			},
			expectedCode: http.StatusOK,
			expectBody:   "Comment deleted successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, cs, _, ctrl := newCommentHandlerWithMock(t)
			defer ctrl.Finish()

			tt.mockSetup(cs)

			req := httptest.NewRequest(http.MethodDelete, "/v1/tasks/comments/"+tt.commentId, nil)
			if tt.commentId != "" {
				req.SetPathValue("comment_id", tt.commentId)
			}
			if tt.userId != "" {
				ctx := context.WithValue(req.Context(), ContextKey.UserId, tt.userId)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.DeleteComment(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
			if tt.expectBody != "" && !strings.Contains(w.Body.String(), tt.expectBody) {
				t.Errorf("%s: expected body to contain %q, got %q", tt.name, tt.expectBody, w.Body.String())
			}
		})
	}
}


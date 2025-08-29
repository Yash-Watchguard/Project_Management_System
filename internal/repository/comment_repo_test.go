package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
)

func TestViewAllComments(t *testing.T) {
	tests := []struct {
		name        string
		taskID      string
		setupMock   func(mock sqlmock.Sqlmock)
		expected    []comment.Comment
		expectError bool
		errorMsg    string
	}{
		{
			name:   "success - comments found",
			taskID: "task1",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"comment_id", "task_id", "created_by", "comment",
				}).
					AddRow("c1", "task1", "user1", "This is the first comment").
					AddRow("c2", "task1", "user2", "Another comment")

				mock.ExpectQuery("SELECT comment_id, task_id, created_by, comment FROM comments WHERE task_id = \\?").
					WithArgs("task1").
					WillReturnRows(rows)
			},
			expected: []comment.Comment{
				{CommentId: "c1", TaskId: "task1", CreatedBy: "user1", Content: "This is the first comment"},
				{CommentId: "c2", TaskId: "task1", CreatedBy: "user2", Content: "Another comment"},
			},
			expectError: false,
		},
		{
			name:   "failure - no comments found",
			taskID: "task2",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"comment_id", "task_id", "created_by", "comment",
				})
				mock.ExpectQuery("SELECT comment_id, task_id, created_by, comment FROM comments WHERE task_id = \\?").
					WithArgs("task2").
					WillReturnRows(rows)
			},
			expectError: true,
			errorMsg:    "no comments found for this task",
		},
		{
			name:   "failure - db error",
			taskID: "task3",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT comment_id, task_id, created_by, comment FROM comments WHERE task_id = \\?").
					WithArgs("task3").
					WillReturnError(errors.New("db query failed"))
			},
			expectError: true,
			errorMsg:    "db query failed",
		},
		{
			name:   "failure - scan error",
			taskID: "task4",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Wrong number of columns -> triggers scan error
				rows := sqlmock.NewRows([]string{
					"comment_id", "task_id",
				}).AddRow("c1", "task4")

				mock.ExpectQuery("SELECT comment_id, task_id, created_by, comment FROM comments WHERE task_id = \\?").
					WithArgs("task4").
					WillReturnRows(rows)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewCommentRepo(db)

			tt.setupMock(mock)

			comments, err := repo.ViewAllComments(tt.taskID)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Fatalf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(comments) != len(tt.expected) {
					t.Fatalf("expected %d comments, got %d", len(tt.expected), len(comments))
				}
				for i := range tt.expected {
					if comments[i] != tt.expected[i] {
						t.Errorf("expected %+v, got %+v", tt.expected[i], comments[i])
					}
				}
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	tests := []struct {
		name        string
		input       comment.Comment
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "success - comment updated",
			input: comment.Comment{
				CommentId: "c1",
				TaskId:    "task1",
				CreatedBy: "user1",
				Content:   "updated content",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// first query: check existing created_by, task_id
				rows := sqlmock.NewRows([]string{"created_by", "task_id"}).
					AddRow("user1", "task1")
				mock.ExpectQuery("SELECT created_by, task_id FROM comments WHERE comment_id = \\?").
					WithArgs("c1").
					WillReturnRows(rows)

				// second query: update
				mock.ExpectExec("UPDATE comments SET comment = \\? WHERE comment_id = \\?").
					WithArgs("updated content", "c1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name: "failure - no comment found",
			input: comment.Comment{
				CommentId: "c2",
				TaskId:    "task1",
				CreatedBy: "user1",
				Content:   "whatever",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT created_by, task_id FROM comments WHERE comment_id = \\?").
					WithArgs("c2").
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
			errorMsg:    "no comments found with the given Id",
		},
		{
			name: "failure - unauthorized user",
			input: comment.Comment{
				CommentId: "c3",
				TaskId:    "task1",
				CreatedBy: "userX",
				Content:   "hack attempt",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"created_by", "task_id"}).
					AddRow("user1", "task1")
				mock.ExpectQuery("SELECT created_by, task_id FROM comments WHERE comment_id = \\?").
					WithArgs("c3").
					WillReturnRows(rows)
			},
			expectError: true,
			errorMsg:    ErrUnauthorized.Error(),
		},
		{
			name: "failure - wrong task",
			input: comment.Comment{
				CommentId: "c4",
				TaskId:    "wrongTask",
				CreatedBy: "user1",
				Content:   "wrong task update",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"created_by", "task_id"}).
					AddRow("user1", "task1")
				mock.ExpectQuery("SELECT created_by, task_id FROM comments WHERE comment_id = \\?").
					WithArgs("c4").
					WillReturnRows(rows)
			},
			expectError: true,
			errorMsg:    ErrWrongTask.Error(),
		},
		{
			name: "failure - update affects 0 rows",
			input: comment.Comment{
				CommentId: "c5",
				TaskId:    "task1",
				CreatedBy: "user1",
				Content:   "no row updated",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"created_by", "task_id"}).
					AddRow("user1", "task1")
				mock.ExpectQuery("SELECT created_by, task_id FROM comments WHERE comment_id = \\?").
					WithArgs("c5").
					WillReturnRows(rows)

				mock.ExpectExec("UPDATE comments SET comment = \\? WHERE comment_id = \\?").
					WithArgs("no row updated", "c5").
					WillReturnResult(sqlmock.NewResult(1, 0)) // 0 rows affected
			},
			expectError: true,
			errorMsg:    "no comment updated",
		},
		{
			name: "failure - db error during update",
			input: comment.Comment{
				CommentId: "c6",
				TaskId:    "task1",
				CreatedBy: "user1",
				Content:   "trigger db error",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"created_by", "task_id"}).
					AddRow("user1", "task1")
				mock.ExpectQuery("SELECT created_by, task_id FROM comments WHERE comment_id = \\?").
					WithArgs("c6").
					WillReturnRows(rows)

				mock.ExpectExec("UPDATE comments SET comment = \\? WHERE comment_id = \\?").
					WithArgs("trigger db error", "c6").
					WillReturnError(errors.New("update failed"))
			},
			expectError: true,
			errorMsg:    "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewCommentRepo(db)

			tt.setupMock(mock)

			err = repo.UpdateComment(tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Fatalf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}


func TestAddComment(t *testing.T) {
	tests := []struct {
		name        string
		input       comment.Comment
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "success - comment inserted",
			input: comment.Comment{
				CommentId: "c1",
				TaskId:    "task1",
				CreatedBy: "user1",
				Content:   "Nice work!",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO comments").
					WithArgs("c1", "task1", "user1", "Nice work!").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name: "failure - db error",
			input: comment.Comment{
				CommentId: "c2",
				TaskId:    "task2",
				CreatedBy: "user2",
				Content:   "Failing insert",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO comments").
					WithArgs("c2", "task2", "user2", "Failing insert").
					WillReturnError(errors.New("insert failed"))
			},
			expectError: true,
			errorMsg:    "insert failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewCommentRepo(db)

			tt.setupMock(mock)

			err = repo.AddComment(tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Fatalf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	tests := []struct {
		name        string
		userId      string
		commentId   string
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name:      "success - comment deleted",
			userId:    "user1",
			commentId: "c1",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM comments").
					WithArgs("c1", "user1").
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row deleted
			},
			expectError: false,
		},
		{
			name:      "failure - no rows deleted",
			userId:    "user1",
			commentId: "c2",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM comments").
					WithArgs("c2", "user1").
					WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows deleted
			},
			expectError: true,
			errorMsg:    ErrDelete.Error(),
		},
		{
			name:      "failure - db error",
			userId:    "user2",
			commentId: "c3",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM comments").
					WithArgs("c3", "user2").
					WillReturnError(errors.New("db failure"))
			},
			expectError: true,
			errorMsg:    "db failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewCommentRepo(db)

			tt.setupMock(mock)

			err = repo.DeleteComment(tt.userId, tt.commentId)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Fatalf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
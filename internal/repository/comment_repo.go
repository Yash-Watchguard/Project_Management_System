package repository

import (
	"database/sql"
	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type CommentRepo struct {
	db sql.DB
	dynamoClient dynamodb.Client
	tableName string
}

func NewCommentRepo(dynamoDbClient dynamodb.Client,tableName string) *CommentRepo {
	return &CommentRepo{dynamoClient: dynamoDbClient,tableName: tableName}
}

var (
	ErrUnauthorized    = errors.New("unauthorized to update this comment")
	ErrCommentNotFound = errors.New("no comments found with the given Id")
	ErrWrongTask       = errors.New("comment does not belong to the specified task")
    ErrDelete          = errors.New("either comment not found or you are not authorized to delete it")
)

func (cr *CommentRepo) ViewAllComments(taskId string) ([]comment.Comment, error) {
	rows, err := cr.db.Query(
		config.SelectQuery("comments", []string{"comment_id", "task_id", "created_by", "comment"}, "task_id"),
		taskId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []comment.Comment
	for rows.Next() {
		var c comment.Comment
		err = rows.Scan(&c.CommentId, &c.TaskId, &c.CreatedBy, &c.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if len(comments) == 0 {
		return nil, errors.New("no comments found for this task")
	}

	return comments, nil
}

func (cr *CommentRepo) UpdateComment(updatedComment comment.Comment) error {
	var existingCreatedBy, existingTaskId string

	query := "SELECT created_by, task_id FROM comments WHERE comment_id = ?"

	row := cr.db.QueryRow(query, updatedComment.CommentId)

	err := row.Scan(&existingCreatedBy, &existingTaskId)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no comments found with the given Id")
		}
		return err
	}

	if updatedComment.CreatedBy != existingCreatedBy {
		return ErrUnauthorized
	}
	if updatedComment.TaskId != existingTaskId {
		return ErrWrongTask
	}

	query = "UPDATE comments SET comment = ? WHERE comment_id = ?"
	result, err := cr.db.Exec(query, updatedComment.Content, updatedComment.CommentId)

	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("no comment updated")
	}
	return nil
}

func (cr *CommentRepo) AddComment(newComment comment.Comment) error {
	// Prepare the INSERT query
	query := `
        INSERT INTO comments (comment_id, task_id, created_by, comment)
        VALUES (?, ?, ?, ?)
    `

	_, err := cr.db.Exec(query, newComment.CommentId, newComment.TaskId, newComment.CreatedBy, newComment.Content)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CommentRepo) DeleteComment(userId, commentId string) error {
	// Delete query with condition on comment_id and created_by
	query := `
        DELETE FROM comments
        WHERE comment_id = ? AND created_by = ?
    `

	result, err := cr.db.Exec(query, commentId, userId)
	if err != nil {
		return err
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrDelete
	}

	return nil
}



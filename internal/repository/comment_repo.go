package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

    skPrefix := fmt.Sprintf("TASK#%sCOMMENT#", taskId)

    query := fmt.Sprintf("SELECT * FROM %s WHERE PK = ? AND begins_with(SK, ?)", cr.tableName)

    out, err := cr.dynamoClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
        Statement: aws.String(query),
        Parameters: []types.AttributeValue{
            &types.AttributeValueMemberS{Value: "COMMENTS"},  
            &types.AttributeValueMemberS{Value: skPrefix},    
        },
    })

    if err != nil {
        return nil, err
    }

    if len(out.Items) == 0 {
        return nil, errors.New("no comments found for this task")
    }

    comments := make([]comment.Comment, 0, len(out.Items))
    for _, item := range out.Items {
        var dc comment.DynamoComment
        if err := attributevalue.UnmarshalMap(item, &dc); err != nil {
            return nil, err
        }

        c := comment.Comment{
            CommentId: dc.CommentId,
            TaskId:    dc.TaskId,
            CreatedBy: dc.CreatedBy,
            Content:   dc.Content,
        }
        comments = append(comments, c)
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

    sk := fmt.Sprintf("TASK#%sCOMMENT#%s", newComment.TaskId, newComment.CommentId)

    query := `
        INSERT INTO ` + cr.tableName + ` VALUE {
            'PK': ?,
            'SK': ?,
            'CommentId': ?,
            'Content': ?,
            'CreatedBy': ?,
            'TaskId': ?
        }
    `
    _, err := cr.dynamoClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
        Statement: aws.String(query),
        Parameters: []types.AttributeValue{
            &types.AttributeValueMemberS{Value: "COMMENTS"},                 // PK
            &types.AttributeValueMemberS{Value: sk},                         // SK
            &types.AttributeValueMemberS{Value: newComment.CommentId},       // CommentId
            &types.AttributeValueMemberS{Value: newComment.Content},         // Content
            &types.AttributeValueMemberS{Value: newComment.CreatedBy},       // CreatedBy
            &types.AttributeValueMemberS{Value: newComment.TaskId},          // TaskId
        },
    })

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



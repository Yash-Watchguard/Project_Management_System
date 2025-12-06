package repository

import (

	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	Priority "github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TaskRepo struct {
	db *sql.DB
	dynamoCliet *dynamodb.Client
	tableName string
}

func NewTaskRepo(dynamoCliet *dynamodb.Client,tableName string) *TaskRepo {
	return &TaskRepo{dynamoCliet: dynamoCliet,tableName: tableName}
}

func(taskRepo *TaskRepo)ViewAllManagerTask(managerId string)([]task.Task,error){
   var tasks []task.Task

	pk := fmt.Sprintf("USER#%s", managerId)

	input := &dynamodb.QueryInput{
		TableName: aws.String(taskRepo.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: pk},
			":skPrefix": &types.AttributeValueMemberS{Value: "PROJECT#"},
		},
	}

	resp, err := taskRepo.dynamoCliet.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	for _, item := range resp.Items {

		var dynTask task.DynamoTask
		if err := attributevalue.UnmarshalMap(item, &dynTask); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}

		var t task.Task
		t.TaskId = dynTask.TaskId
		t.Title = dynTask.Title
		t.Description = dynTask.Description
		t.AcceptanceCriteria = dynTask.AcceptanceCriteria
		t.AssignedTo = dynTask.AssignedTo
		t.ProjectId = dynTask.ProjectId
		t.CreatedBy = dynTask.CreatedBy

		if dynTask.Deadline != "" {
			t.Deadline, err = time.Parse(time.RFC3339, dynTask.Deadline)
			if err != nil {
				return nil, fmt.Errorf("deadline parse error: %w", err)
			}
		}

		t.TaskPriority, _ = Priority.PriorityParser(dynTask.TaskPriority)
		t.TaskStatus, _ = status.GetStatusFromString(dynTask.TaskStatus)

		tasks = append(tasks, t)
	}
	return tasks, nil
}
func (taskRepo *TaskRepo) ViewAllTask(projectId string) ([]task.Task, error) {
	query := `SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby
              FROM tasks WHERE projectid = ?`

	rows, err := taskRepo.db.Query(query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectTasks []task.Task
	for rows.Next() {
		var t task.Task
		var deadlineBytes []byte
		err := rows.Scan(
			&t.TaskId,
			&t.Title,
			&t.Description,
			&t.AcceptanceCriteria,
			&deadlineBytes,
			&t.TaskPriority,
			&t.TaskStatus,
			&t.AssignedTo,
			&t.ProjectId,
			&t.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		if len(deadlineBytes) > 0 {
			t.Deadline, err = time.Parse("2006-01-02", string(deadlineBytes))
			if err != nil {
				return nil, err
			}
		}
		projectTasks = append(projectTasks, t)
	}
	return projectTasks, nil
}

func (taskRepo *TaskRepo) SaveTask(newTask task.Task) error {

 dynamoTask := task.DynamoTask{
		PK:                 "USER#" + newTask.CreatedBy,
		SK:                 "PROJECT#" + newTask.ProjectId + "TASK#" + newTask.TaskId,
		TaskId:             newTask.TaskId,
		Title:              newTask.Title,
		Description:        newTask.Description,
		AcceptanceCriteria: newTask.AcceptanceCriteria,
		Deadline:           newTask.Deadline.Format(time.RFC3339),
		TaskPriority:       Priority.GetPriority(newTask.TaskPriority),
		TaskStatus:         status.GetStatusString(newTask.TaskStatus),
		AssignedTo:         newTask.AssignedTo,
		ProjectId:          newTask.ProjectId,
		CreatedBy:          newTask.CreatedBy,
	}

	item, err := attributevalue.MarshalMap(dynamoTask)
	if err != nil {
		return fmt.Errorf("error marshaling task: %s", err.Error())
	}
	_, err = taskRepo.dynamoCliet.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(taskRepo.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("error saving task: %s", err.Error())
	}

	dynamoTask = task.DynamoTask{
		PK:                 "USER#" + newTask.AssignedTo,
		SK:                 "TASK#PROJECT#" + newTask.ProjectId + "TASK#" + newTask.TaskId,
		TaskId:             newTask.TaskId,
		Title:              newTask.Title,
		Description:        newTask.Description,
		AcceptanceCriteria: newTask.AcceptanceCriteria,
		Deadline:           newTask.Deadline.Format(time.RFC3339),
		TaskPriority:       Priority.GetPriority(newTask.TaskPriority),
		TaskStatus:         status.GetStatusString(newTask.TaskStatus),
		AssignedTo:         newTask.AssignedTo,
		ProjectId:          newTask.ProjectId,
		CreatedBy:          newTask.CreatedBy,
	}

	item, err = attributevalue.MarshalMap(dynamoTask)
	if err != nil {
		return fmt.Errorf("error marshaling task: %s", err.Error())
	}
	_, err = taskRepo.dynamoCliet.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(taskRepo.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("error saving task: %s", err.Error())
	}

	return nil
}


func (r *TaskRepo) DeleteTask(projectId, taskId, managerId, empId string) error {

    pkManager := "USER#" + managerId
    skManager := "PROJECT#" + projectId + "TASK#" + taskId

    pkEmployee := "USER#" + empId
    skEmployee := "TASK#PROJECT#" + projectId + "TASK#" + taskId

    deleteStmt := "DELETE FROM " + r.tableName + " WHERE PK = ? AND SK = ?"

    _, err := r.dynamoCliet.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
        Statement: aws.String(deleteStmt),
        Parameters: []types.AttributeValue{
            &types.AttributeValueMemberS{Value: pkManager},
            &types.AttributeValueMemberS{Value: skManager},
        },
    })
    if err != nil {
        return fmt.Errorf("failed to delete manager task copy: %w", err)
    }

    _, err = r.dynamoCliet.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
        Statement: aws.String(deleteStmt),
        Parameters: []types.AttributeValue{
            &types.AttributeValueMemberS{Value: pkEmployee},
            &types.AttributeValueMemberS{Value: skEmployee},
        },
    })
    if err != nil {
        return fmt.Errorf("failed to delete employee task copy: %w", err)
    }

    return nil
}


func (taskRepo *TaskRepo) ViewAssignedTask(empId string) ([]task.Task, error) {

	var tasks []task.Task

	pk := fmt.Sprintf("USER#%s", empId)

	input := &dynamodb.QueryInput{
		TableName: aws.String(taskRepo.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: pk},
			":skPrefix": &types.AttributeValueMemberS{Value: "TASK#PROJECT#"},
		},
	}

	resp, err := taskRepo.dynamoCliet.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	for _, item := range resp.Items {

		var dynTask task.DynamoTask
		if err := attributevalue.UnmarshalMap(item, &dynTask); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}
        
		status,_:= status.GetStatusFromString(dynTask.TaskStatus)
		var t task.Task
		t.TaskId = dynTask.TaskId
		t.Title = dynTask.Title
		t.Description = dynTask.Description
		t.AcceptanceCriteria = dynTask.AcceptanceCriteria
		t.AssignedTo = dynTask.AssignedTo
		t.ProjectId = dynTask.ProjectId
		t.CreatedBy = dynTask.CreatedBy
		t.TaskStatus= status

		if dynTask.Deadline != "" {
			t.Deadline, err = time.Parse(time.RFC3339, dynTask.Deadline)
			if err != nil {
				return nil, fmt.Errorf("deadline parse error: %w", err)
			}
		}

		t.TaskPriority, _ = Priority.PriorityParser(dynTask.TaskPriority)

		tasks = append(tasks, t)
	}

	return tasks, nil
}



func (taskRepo *TaskRepo) UpdateTaskStatus(empId string, taskId string, updatedStatus status.TaskStatus) error {
	query := `UPDATE tasks 
              SET taskstatus = ? 
              WHERE task_id = ?`

	res, err := taskRepo.db.Exec(query, updatedStatus, taskId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not assigned to employee or task does not exist")
	}

	return nil
}

func (taskRepo *TaskRepo)ViewAllTasksInProject(projectId string, creator string) ([]task.Task, error) {
	var tasks []task.Task

	pk := fmt.Sprintf("USER#%s", creator)

	input := &dynamodb.QueryInput{
		TableName: aws.String(taskRepo.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: pk},
			":skPrefix": &types.AttributeValueMemberS{Value: "PROJECT#"+projectId},
		},
	}

	resp, err := taskRepo.dynamoCliet.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	for _, item := range resp.Items {

		var dynTask task.DynamoTask
		if err := attributevalue.UnmarshalMap(item, &dynTask); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}

		var t task.Task
		t.TaskId = dynTask.TaskId
		t.Title = dynTask.Title
		t.Description = dynTask.Description
		t.AcceptanceCriteria = dynTask.AcceptanceCriteria
		t.AssignedTo = dynTask.AssignedTo
		t.ProjectId = dynTask.ProjectId
		t.CreatedBy = dynTask.CreatedBy

		if dynTask.Deadline != "" {
			t.Deadline, err = time.Parse(time.RFC3339, dynTask.Deadline)
			if err != nil {
				return nil, fmt.Errorf("deadline parse error: %w", err)
			}
		}

		t.TaskPriority, _ = Priority.PriorityParser(dynTask.TaskPriority)
		t.TaskStatus, _ = status.GetStatusFromString(dynTask.TaskStatus)

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepo) UpdateTask(projectId, taskId, managerId string, updates map[string]interface{}) error {

   
    pkManager := "USER#" + managerId
    skManager := "PROJECT#" + projectId + "TASK#" + taskId

    selectStmt := "SELECT * FROM " + r.tableName + " WHERE PK = ? AND SK = ?"

    out, err := r.dynamoCliet.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
        Statement: aws.String(selectStmt),
        Parameters: []types.AttributeValue{
            &types.AttributeValueMemberS{Value: pkManager},
            &types.AttributeValueMemberS{Value: skManager},
        },
    })
    if err != nil {
        return err
    }
    if len(out.Items) == 0 {
        return errors.New("task not found under manager")
    }

    var existing task.DynamoTask
    err = attributevalue.UnmarshalMap(out.Items[0], &existing)
    if err != nil {
        return err
    }

    oldAssignedTo := existing.AssignedTo

    
    updateStmt := "UPDATE " + r.tableName + " SET "
    params := []types.AttributeValue{}

    modifiableFields := []string{
        "Title",
        "Description",
        "AcceptanceCriteria",
        "Deadline",
        "TaskPriority",
        "AssignedTo",
		"TaskStatus",
    }

    for _, field := range modifiableFields {
    if val, ok := updates[field]; ok {

        // MUST wrap attribute names in double quotes
        updateStmt += fmt.Sprintf("\"%s\" = ?, ", field)

        params = append(params, &types.AttributeValueMemberS{
            Value: fmt.Sprintf("%v", val),
        })
    }
}

    if len(params) == 0 {
        return errors.New("no valid updates provided")
    }

    updateStmt = strings.TrimSuffix(updateStmt, ", ")
    updateStmt += " WHERE PK = ? AND SK = ?"

   
    managerParams := append(params,
        &types.AttributeValueMemberS{Value: pkManager},
        &types.AttributeValueMemberS{Value: skManager},
    )

    _, err = r.dynamoCliet.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
        Statement:  aws.String(updateStmt),
        Parameters: managerParams,
    })
    if err != nil {
        return err
    }

    
    pkEmployee := "USER#" + oldAssignedTo
    skEmployee := "TASK#PROJECT#" + projectId + "TASK#" + taskId

    // Check if employee copy exists before updating
    employeeSelectStmt := "SELECT * FROM " + r.tableName + " WHERE PK = ? AND SK = ?"
    employeeOut, err := r.dynamoCliet.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
        Statement: aws.String(employeeSelectStmt),
        Parameters: []types.AttributeValue{
            &types.AttributeValueMemberS{Value: pkEmployee},
            &types.AttributeValueMemberS{Value: skEmployee},
        },
    })
    if err != nil {
        return err
    }

    // Only update employee copy if it exists
    if len(employeeOut.Items) > 0 {
        employeeParams := append(params,
            &types.AttributeValueMemberS{Value: pkEmployee},
            &types.AttributeValueMemberS{Value: skEmployee},
        )

        _, err = r.dynamoCliet.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
            Statement:  aws.String(updateStmt),
            Parameters: employeeParams,
        })
        if err != nil {
            return err
        }
    }

   
    newAssignedTo, assignedChanged := updates["AssignedTo"].(string)

    if assignedChanged && newAssignedTo != oldAssignedTo {

        // DELETE OLD EMPLOYEE COPY
        deleteStmt := "DELETE FROM " + r.tableName + " WHERE PK = ? AND SK = ?"

        _, err = r.dynamoCliet.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
            Statement: aws.String(deleteStmt),
            Parameters: []types.AttributeValue{
                &types.AttributeValueMemberS{Value: pkEmployee},
                &types.AttributeValueMemberS{Value: skEmployee},
            },
        })
        if err != nil {
            return err
        }

        updatedTask := existing

        for k, v := range updates {
            switch k {
            case "Title":
                updatedTask.Title = v.(string)
            case "Description":
                updatedTask.Description = v.(string)
            case "AcceptanceCriteria":
                updatedTask.AcceptanceCriteria = v.(string)
            case "Deadline":
                updatedTask.Deadline = v.(string)
            case "TaskPriority":
                updatedTask.TaskPriority = v.(string)
            case "AssignedTo":
                updatedTask.AssignedTo = v.(string)
			case "TaskStatus":
                updatedTask.TaskStatus = v.(string)
            }
        }


        newPK := "USER#" + newAssignedTo
        newSK := "TASK#PROJECT#" + projectId + "TASK#" + taskId

        updatedTask.PK = newPK
        updatedTask.SK = newSK

        item, err := attributevalue.MarshalMap(updatedTask)
        if err != nil {
            return err
        }

        // Use PutItem instead of INSERT for the new employee copy
        _, err = r.dynamoCliet.PutItem(context.TODO(), &dynamodb.PutItemInput{
            TableName: aws.String(r.tableName),
            Item:      item,
        })
        if err != nil {
            return err
        }
    }

    return nil
}



func(taskRepo *TaskRepo)UpdateTaskEmailId(taskId string,managerId string, updates map[string]interface{}) error{
   return nil
}



// repository/project_repo.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ProjectRepo struct {
	db *sql.DB
    dynamoDb dynamodb.Client
    tableName string
    taskRepo TaskRepo
}

func NewProjectRepo(db *dynamodb.Client,tablename string, taskRepo TaskRepo) *ProjectRepo {
	return &ProjectRepo{dynamoDb: *db,tableName: tablename ,taskRepo: taskRepo}
}

func (pr *ProjectRepo) AddProject(newProject project.Project) error {
	query:=`INSERT INTO projects (project_id, project_name, project_description, deadline, created_by, assigned_manager_id) VALUES (?, ?, ?, ?, ?, ?)`

	_,err:=pr.db.Exec(query,newProject.ProjectId,newProject.ProjectName,newProject.ProjectDescription,newProject.Deadline,newProject.CreatedBy,newProject.AssignedManager)

	if err!=nil{
		return err
	}
	return nil
}

func (pr *ProjectRepo) ViewAllProjects() ([]project.Project, error) {
    var projects []project.Project

    query := `SELECT project_id, project_name, project_description, deadline, created_by, assigned_manager_id FROM projects`

    rows, err := pr.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
		var p project.Project
    var deadlineBytes []byte

    err := rows.Scan(&p.ProjectId, &p.ProjectName, &p.ProjectDescription, &deadlineBytes, &p.CreatedBy, &p.AssignedManager)
    if err != nil {
        return nil, err
    }

    // Parse the date string
    if len(deadlineBytes) > 0 {
        p.Deadline, err = time.Parse("2006-01-02", string(deadlineBytes)) // if DATE type
        if err != nil {
            return nil, err
        }
    }

    projects = append(projects, p)
}

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return projects, nil
}

func (pr *ProjectRepo) DeleteProject(creatorId, managerId, projectID string) error {
	ctx := context.TODO()

	creatorPk := fmt.Sprintf("USER#%s", creatorId)
	creatorSk := fmt.Sprintf("PROJECT#%s", projectID)

	_, err := pr.dynamoDb.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(pr.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: creatorPk},
			"SK": &types.AttributeValueMemberS{Value: creatorSk},
		},
	})
	if err != nil {
		return fmt.Errorf("error deleting creator project: %v", err)
	}

	managerPk := fmt.Sprintf("USER#%s", managerId)
	managerSk := fmt.Sprintf("ASSIGNED#PROJECT#%s", projectID)

	_, err = pr.dynamoDb.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(pr.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: managerPk},
			"SK": &types.AttributeValueMemberS{Value: managerSk},
		},
	})
	if err != nil {
		return fmt.Errorf("error deleting manager project: %v", err)
	}

	taskPrefix := fmt.Sprintf("PROJECT#%s#TASK#", projectID)
	creatorPk = fmt.Sprintf("USER#%s", managerId)

	queryOutput, err := pr.dynamoDb.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(pr.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: creatorPk},
			":sk": &types.AttributeValueMemberS{Value: taskPrefix},
		},
	})
	if err != nil {
		return fmt.Errorf("error querying tasks: %v", err)
	}

	tasks := []task.DynamoTask{}
	err = attributevalue.UnmarshalListOfMaps(queryOutput.Items, &tasks)
	if err != nil {
		return fmt.Errorf("task unmarshal error: %v", err)
	}

	for _, t := range tasks {
		fmt.Println(t)
		err = pr.taskRepo.DeleteTask(projectID, t.TaskId, t.CreatedBy, t.AssignedTo)
		if err != nil {
			return err
		}
	}

	return nil
}





func (pr *ProjectRepo) ViewAssignedProject(userId string) ([]project.Project, error) {
    var assignedProjects []project.Project

    query := `SELECT project_id, project_name, project_description, deadline, created_by, assigned_manager_id 
              FROM projects 
              WHERE assigned_manager_id = ?`

    rows, err := pr.db.Query(query, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
    var p project.Project
    var deadlineBytes []byte

    err := rows.Scan(&p.ProjectId, &p.ProjectName, &p.ProjectDescription, &deadlineBytes, &p.CreatedBy, &p.AssignedManager)
    if err != nil {
        return nil, err
    }

    // Parse the date string
    if len(deadlineBytes) > 0 {
        p.Deadline, err = time.Parse("2006-01-02", string(deadlineBytes)) // if DATE type
        if err != nil {
            return nil, err
        }
    }

    assignedProjects = append(assignedProjects, p)
}


    return assignedProjects, nil
}









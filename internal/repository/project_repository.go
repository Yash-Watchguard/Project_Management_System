// repository/project_repo.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

func (pr *ProjectRepo)UpdateProject(projectId,creatorId,managerId string,updates map[string]any)error{

    pkCreator := "USER#" + creatorId
    skCreator := "PROJECT#" + projectId

    selectStmt := "SELECT * FROM " + pr.tableName + " WHERE PK = ? AND SK = ?"

    out, err := pr.dynamoDb.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
        Statement: aws.String(selectStmt),
        Parameters: []types.AttributeValue{
            &types.AttributeValueMemberS{Value: pkCreator},
            &types.AttributeValueMemberS{Value: skCreator},
        },
    })
    if err != nil {
        return err
    }
    if len(out.Items) == 0 {
        return errors.New("project not found under creator")
    }

    var existing project.DynamoProject
    err = attributevalue.UnmarshalMap(out.Items[0], &existing)
    if err != nil {
        return err
    }

    oldManager := existing.Assigned_manager

    // ---------- STEP 2: BUILD UPDATE STATEMENT ----------
    modifiableFields := []string{
        "Assigned_manager",
        "Project_deadline",
        "Project_description",
        "Project_name",
    }

    updateStmt := "UPDATE " + pr.tableName + " SET "
    params := []types.AttributeValue{}

    for _, field := range modifiableFields {
        if val, ok := updates[field]; ok {
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

    // update creator copy
    creatorParams := append(params,
        &types.AttributeValueMemberS{Value: pkCreator},
        &types.AttributeValueMemberS{Value: skCreator},
    )

    _, err = pr.dynamoDb.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
        Statement:  aws.String(updateStmt),
        Parameters: creatorParams,
    })
    if err != nil {
        return err
    }

    // ---------- STEP 3: UPDATE MANAGER COPY IF EXISTS ----------
    pkManagerOld := "USER#" + oldManager
    skManagerOld := "ASSIGNED#PROJECT#" + projectId

    managerSelect := "SELECT * FROM " + pr.tableName + " WHERE PK = ? AND SK = ?"
    managerOut, err := pr.dynamoDb.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
        Statement: aws.String(managerSelect),
        Parameters: []types.AttributeValue{
            &types.AttributeValueMemberS{Value: pkManagerOld},
            &types.AttributeValueMemberS{Value: skManagerOld},
        },
    })
    if err != nil {
        return err
    }

    if len(managerOut.Items) > 0 {
        // update manager copy
        managerParams := append(params,
            &types.AttributeValueMemberS{Value: pkManagerOld},
            &types.AttributeValueMemberS{Value: skManagerOld},
        )

        _, err = pr.dynamoDb.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
            Statement:  aws.String(updateStmt),
            Parameters: managerParams,
        })
        if err != nil {
            return err
        }
    }

    // ---------- STEP 4: HANDLE MANAGER CHANGE ----------
    newManager, changed := updates["Assigned_manager"].(string)

    if changed && newManager != oldManager {

        // DELETE OLD MANAGER COPY
        deleteStmt := "DELETE FROM " + pr.tableName + " WHERE PK = ? AND SK = ?"
        _, err = pr.dynamoDb.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
            Statement: aws.String(deleteStmt),
            Parameters: []types.AttributeValue{
                &types.AttributeValueMemberS{Value: pkManagerOld},
                &types.AttributeValueMemberS{Value: skManagerOld},
            },
        })
        if err != nil {
            return err
        }

        // Create new manager copy
        updated := existing
        for k, v := range updates {
            switch k {
            case "Assigned_manager":
                updated.Assigned_manager = v.(string)
            case "Project_deadline":
                updated.Project_deadline = v.(string)
            case "Project_description":
                updated.Project_description = v.(string)
            case "Project_name":
                updated.Project_name = v.(string)
            }
        }

        newPK := "USER#" + newManager
        newSK := "ASSIGNED#PROJECT#" + projectId

        updated.PK = newPK
        updated.SK = newSK

        item, err := attributevalue.MarshalMap(updated)
        if err != nil {
            return err
        }

        _, err = pr.dynamoDb.PutItem(context.TODO(), &dynamodb.PutItemInput{
            TableName: aws.String(pr.tableName),
            Item:      item,
        })
        if err != nil {
            return err
        }
    }
    return nil
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









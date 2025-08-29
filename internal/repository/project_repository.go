// repository/project_repo.go
package repository

import (
	"database/sql"
	"errors"
    "time"
    "strings"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
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

func (pr *ProjectRepo) DeleteProject(projectID string) error {
   
    projectID = strings.TrimSpace(projectID)

 
    var existsInt int
    checkQuery := `SELECT EXISTS(SELECT 1 FROM projects WHERE project_id = ?)`
    err := pr.db.QueryRow(checkQuery, projectID).Scan(&existsInt)
    if err != nil {
        return errors.New("failed to check project existence")
    }
    if existsInt == 0 {
        return errors.New("project not found")
    }

    
    _, err = pr.db.Exec(`DELETE FROM tasks WHERE projectid = ?`, projectID)
    if err != nil {
        return errors.New("failed to delete tasks for the project")
    }

    // 3️⃣ Delete the project itself
    result, err := pr.db.Exec(`DELETE FROM projects WHERE project_id = ?`, projectID)
    if err != nil {
        return errors.New("failed to delete project")
    }

    // 4️⃣ Verify that the deletion affected a row
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return errors.New("failed to get rows affected")
    }
    if rowsAffected == 0 {
        return errors.New("no project deleted")
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









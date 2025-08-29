package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
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
	query := `INSERT INTO tasks 
        (task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := taskRepo.db.Exec(
		query,
		newTask.TaskId,
		newTask.Title,
		newTask.Description,
		newTask.AcceptanceCriteria,
		newTask.Deadline,
		newTask.TaskPriority,
		newTask.TaskStatus,
		newTask.AssignedTo,
		newTask.ProjectId,
		newTask.CreatedBy,
	)

	if err != nil {
		return err
	}
	return nil
}

func (taskRepo *TaskRepo) DeleteTask(taskId string) error {
	//Check  task exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM tasks WHERE task_id = ?)`
	err := taskRepo.db.QueryRow(checkQuery, taskId).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("task not found")
	}
	deleteQuery := `DELETE FROM tasks WHERE task_id = ?`
	_, err = taskRepo.db.Exec(deleteQuery, taskId)
	if err != nil {
		return err
	}

	return nil
}

func (taskRepo *TaskRepo) ViewAssignedTask(empId string) ([]task.Task, error) {
	var assignedTasks []task.Task

	query := `SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby 
              FROM tasks 
              WHERE assignesto = ?`

	rows, err := taskRepo.db.Query(query, empId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			t.Deadline, err = time.Parse("2006-01-02", string(deadlineBytes)) // if DATE type
			if err != nil {
				return nil, err
			}
		}
		assignedTasks = append(assignedTasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assignedTasks, nil
}

func (taskRepo *TaskRepo) UpdateTaskStatus(empId string, taskId string, updatedStatus status.TaskStatus) error {
	query := `UPDATE tasks 
              SET taskstatus = ? 
              WHERE task_id = ? AND assignesto = ?`

	res, err := taskRepo.db.Exec(query, updatedStatus, taskId, empId)
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

func (taskRepo *TaskRepo)ViewAllAssignedTasksInProject(projectId string, empId string) ([]task.Task, error) {
	var assignedTasks []task.Task

	query := `SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby 
              FROM tasks 
              WHERE projectid = ? AND assignesto = ?`

	rows, err := taskRepo.db.Query(query, projectId, empId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

		// Parse deadline safely
		if len(deadlineBytes) > 0 {
			t.Deadline, err = time.Parse("2006-01-02", string(deadlineBytes))
			if err != nil {
				return nil, err
			}
		}

		assignedTasks = append(assignedTasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assignedTasks, nil
}



package repository

import (
	"database/sql"
	"errors"
	
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	Priority "github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
)


func TestViewAllTask(t *testing.T) {
	query := regexp.QuoteMeta(
		`SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby
              FROM tasks WHERE projectid = ?`,
	)

	tests := []struct {
		name        string
		projectId   string
		setupMock   func(mock sqlmock.Sqlmock)
		expectedErr string
		expectTasks bool
	}{
		{
			name:      "success with valid deadline",
			projectId: "p1",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"task_id", "title", "description", "acceptance_criteria", "deadline",
					"taskpriority", "taskstatus", "assignesto", "projectid", "createdby",
				}).AddRow(
					"t1", "Task 1", "Desc", "Criteria", []byte("2025-08-20"),
					Priority.High, status.Done, "u1", "p1", "admin",
				)
				mock.ExpectQuery(query).
					WithArgs("p1").
					WillReturnRows(rows)
			},
			expectedErr: "",
			expectTasks: true,
		},
		{
			name:      "success with empty deadline",
			projectId: "p2",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"task_id", "title", "description", "acceptance_criteria", "deadline",
					"taskpriority", "taskstatus", "assignesto", "projectid", "createdby",
				}).AddRow(
					"t2", "Task 2", "Desc2", "Criteria2", []byte("2025-08-20"),
					Priority.Low, status.Done, "u2", "p2", "admin",
				)
				mock.ExpectQuery(query).
					WithArgs("p2").
					WillReturnRows(rows)
			},
			expectedErr: "",
			expectTasks: true,
		},
		{
			name:      "query error",
			projectId: "p3",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs("p3").
					WillReturnError(errors.New("db error"))
			},
			expectedErr: "db error",
			expectTasks: false,
		},
		{
			name:      "scan error",
			projectId: "p4",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"task_id", "title", "description", "acceptance_criteria", "deadline",
					"taskpriority", "taskstatus", "assignesto", "projectid", "createdby",
				}).AddRow(
					nil, "Bad", "Desc", "Criteria", []byte("2025-08-20"),
					Priority.High, status.Done, "u1", "p4", "admin",
				) // nil task_id → scan fails
				mock.ExpectQuery(query).
					WithArgs("p4").
					WillReturnRows(rows)
			},
			expectedErr: "converting NULL to string",
			expectTasks: false,
		},
		{
			name:      "deadline parse error",
			projectId: "p5",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"task_id", "title", "description", "acceptance_criteria", "deadline",
					"taskpriority", "taskstatus", "assignesto", "projectid", "createdby",
				}).AddRow(
					"t5", "Task 5", "Desc5", "Criteria5", []byte("invalid-date"),
					Priority.High, status.Done, "u5", "p5", "admin",
				)
				mock.ExpectQuery(query).
					WithArgs("p5").
					WillReturnRows(rows)
			},
			expectedErr: "parsing time",
			expectTasks: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open sqlmock db: %v", err)
			}
			defer db.Close()

			repo := NewTaskRepo(db)
			tt.setupMock(mock)

			tasks, err := repo.ViewAllTask(tt.projectId)

			if tt.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if !tt.expectTasks || len(tasks) == 0 {
					t.Fatalf("expected tasks, got none")
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("expected error containing %q, got %v", tt.expectedErr, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}


func TestSaveTask(t *testing.T) {
    query := regexp.QuoteMeta(`INSERT INTO tasks 
     (task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby)
     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

    tests := []struct {
        name        string
        setupMock   func(mock sqlmock.Sqlmock, newTask task.Task)
        expectedErr string
    }{
        {
            name: "success insert",
            setupMock: func(mock sqlmock.Sqlmock, newTask task.Task) {
                mock.ExpectExec(query).
                    WithArgs(
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
                    ).WillReturnResult(sqlmock.NewResult(1, 1))
            },
            expectedErr: "",
        },
        {
            name: "db error on insert",
            setupMock: func(mock sqlmock.Sqlmock, newTask task.Task) {
                mock.ExpectExec(query).
                    WithArgs(
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
                    ).WillReturnError(errors.New("insert failed"))
            },
            expectedErr: "insert failed",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db, mock, _ := sqlmock.New()
            defer db.Close()
            repo := NewTaskRepo(db)

            newTask := task.Task{
                TaskId:       "t1",
                Title:        "Task",
                Description:  "Some work",
                Deadline:     time.Now(),
                TaskPriority: Priority.High,
                TaskStatus:   status.Done,
                AssignedTo:   "u1",
                ProjectId:    "p1",
                CreatedBy:    "admin",
            }

            tt.setupMock(mock, newTask)

            err := repo.SaveTask(newTask)

            if tt.expectedErr == "" && err != nil {
                t.Fatalf("expected no error, got %v", err)
            }
            if tt.expectedErr != "" && (err == nil || !strings.Contains(err.Error(), tt.expectedErr)) {
                t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
            }

            if err := mock.ExpectationsWereMet(); err != nil {
                t.Fatalf("there were unfulfilled expectations: %v", err)
            }
        })
    }
}



func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name        string
		taskID      string
		mockSetup   func(mock sqlmock.Sqlmock, taskID string)
		expectedErr string
	}{
		{
			name:   "task exists and deleted successfully",
			taskID: "t1",
			mockSetup: func(mock sqlmock.Sqlmock, taskID string) {
				// SELECT EXISTS returns true
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM tasks WHERE task_id = ?)`)).
					WithArgs(taskID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				// DELETE succeeds
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM tasks WHERE task_id = ?`)).
					WithArgs(taskID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: "",
		},
		{
			name:   "task does not exist",
			taskID: "t2",
			mockSetup: func(mock sqlmock.Sqlmock, taskID string) {
				// SELECT EXISTS returns false
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM tasks WHERE task_id = ?)`)).
					WithArgs(taskID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			expectedErr: "task not found",
		},
		{
			name:   "db error on exists check",
			taskID: "t3",
			mockSetup: func(mock sqlmock.Sqlmock, taskID string) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM tasks WHERE task_id = ?)`)).
					WithArgs(taskID).
					WillReturnError(errors.New("db error"))
			},
			expectedErr: "db error",
		},
		{
			name:   "db error on delete",
			taskID: "t4",
			mockSetup: func(mock sqlmock.Sqlmock, taskID string) {
				// SELECT EXISTS returns true
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM tasks WHERE task_id = ?)`)).
					WithArgs(taskID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				// DELETE fails
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM tasks WHERE task_id = ?`)).
					WithArgs(taskID).
					WillReturnError(errors.New("delete failed"))
			},
			expectedErr: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo := NewTaskRepo(db)

			tt.mockSetup(mock, tt.taskID)

			err := repo.DeleteTask(tt.taskID)

			if tt.expectedErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.expectedErr != "" && (err == nil || !strings.Contains(err.Error(), tt.expectedErr)) {
				t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}


func TestViewAssignedTask(t *testing.T) {
	tests := []struct {
		name        string
		empID       string
		mockSetup   func(mock sqlmock.Sqlmock, empID string)
		expectError string
		expectCount int
	}{
		{
			name:  "success - tasks found",
			empID: "emp1",
			mockSetup: func(mock sqlmock.Sqlmock, empID string) {
				rows := sqlmock.NewRows([]string{
					"task_id", "title", "description", "acceptance_criteria", "deadline",
					"taskpriority", "taskstatus", "assignesto", "projectid", "createdby",
				}).AddRow(
					"t1", "Task 1", "Desc", "Criteria", "2025-08-20",
					Priority.High, status.Done, empID, "p1", "admin",
				).AddRow(
					"t2", "Task 2", "Another", "Criteria2", "2025-09-01",
					Priority.High, status.Done, empID, "p2", "manager",
				)

				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby 
              		FROM tasks 
              		WHERE assignesto = ?`)).
					WithArgs(empID).
					WillReturnRows(rows)
			},
			expectError: "",
			expectCount: 2,
		},
		{
			name:  "query error",
			empID: "emp2",
			mockSetup: func(mock sqlmock.Sqlmock, empID string) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby 
              		FROM tasks 
              		WHERE assignesto = ?`)).
					WithArgs(empID).
					WillReturnError(errors.New("db error"))
			},
			expectError: "db error",
			expectCount: 0,
		},
		{
			name:  "scan error - invalid deadline",
			empID: "emp3",
			mockSetup: func(mock sqlmock.Sqlmock, empID string) {
				rows := sqlmock.NewRows([]string{
					"task_id", "title", "description", "acceptance_criteria", "deadline",
					"taskpriority", "taskstatus", "assignesto", "projectid", "createdby",
				}).AddRow(
					"t3", "Task 3", "Bad date", "Criteria", "invalid-date",
					Priority.High, status.Done, empID, "p3", "admin",
				)

				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby 
              		FROM tasks 
              		WHERE assignesto = ?`)).
					WithArgs(empID).
					WillReturnRows(rows)
			},
			expectError: "cannot parse",
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo :=NewTaskRepo(db)

			tt.mockSetup(mock, tt.empID)

			result, err := repo.ViewAssignedTask(tt.empID)

			if tt.expectError == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.expectError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectError) {
					t.Fatalf("expected error containing %q, got %v", tt.expectError, err)
				}
			}
			if len(result) != tt.expectCount {
				t.Fatalf("expected %d tasks, got %d", tt.expectCount, len(result))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}

			// Extra check for success case
			if tt.expectCount > 0 {
				// ensure deadline parsed correctly
				if result[0].Deadline.IsZero() {
					t.Fatalf("expected deadline to be parsed, got zero value")
				}
				if result[0].Deadline.Format("2006-01-02") != "2025-08-20" {
					t.Fatalf("expected deadline 2025-08-20, got %s", result[0].Deadline)
				}
			}
		})
	}
}



func TestUpdateTaskStatus(t *testing.T) {
	tests := []struct {
		name        string
		empID       string
		taskID      string
		status      status.TaskStatus
		mockSetup   func(mock sqlmock.Sqlmock, empID, taskID string, st status.TaskStatus)
		expectError string
	}{
		{
			name:   "success - status updated",
			empID:  "emp1",
			taskID: "task1",
			status: status.InProgress,
			mockSetup: func(mock sqlmock.Sqlmock, empID, taskID string, st status.TaskStatus) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE tasks 
              		 SET taskstatus = ? 
              		 WHERE task_id = ? AND assignesto = ?`)).
					WithArgs(st, taskID, empID).
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
			},
			expectError: "",
		},
		{
			name:   "query error",
			empID:  "emp2",
			taskID: "task2",
			status: status.InProgress,
			mockSetup: func(mock sqlmock.Sqlmock, empID, taskID string, st status.TaskStatus) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE tasks 
              		 SET taskstatus = ? 
              		 WHERE task_id = ? AND assignesto = ?`)).
					WithArgs(st, taskID, empID).
					WillReturnError(errors.New("db exec error"))
			},
			expectError: "db exec error",
		},
		{
			name:   "rows affected error",
			empID:  "emp3",
			taskID: "task3",
			status: status.Done,
			mockSetup: func(mock sqlmock.Sqlmock, empID, taskID string, st status.TaskStatus) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE tasks 
              		 SET taskstatus = ? 
              		 WHERE task_id = ? AND assignesto = ?`)).
					WithArgs(st, taskID, empID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
			expectError: "rows affected error",
		},
		{
			name:   "no rows affected",
			empID:  "emp4",
			taskID: "task4",
			status: status.Done,
			mockSetup: func(mock sqlmock.Sqlmock, empID, taskID string, st status.TaskStatus) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE tasks 
              		 SET taskstatus = ? 
              		 WHERE task_id = ? AND assignesto = ?`)).
					WithArgs(st, taskID, empID).
					WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
			},
			expectError: "task not assigned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo := NewTaskRepo(db)

			tt.mockSetup(mock, tt.empID, tt.taskID, tt.status)

			err := repo.UpdateTaskStatus(tt.empID, tt.taskID, tt.status)

			if tt.expectError == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.expectError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectError) {
					t.Fatalf("expected error containing %q, got %v", tt.expectError, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

func TestViewAllAssignedTasksInProject(t *testing.T) {
	tests := []struct {
		name        string
		projectID   string
		empID       string
		mockSetup   func(mock sqlmock.Sqlmock, projectID, empID string)
		expectTasks int
		expectError string
	}{
		{
			name:      "success - returns tasks",
			projectID: "proj1",
			empID:     "emp1",
			mockSetup: func(mock sqlmock.Sqlmock, projectID, empID string) {
				rows := sqlmock.NewRows([]string{
					"task_id", "title", "description", "acceptance_criteria", "deadline",
					"taskpriority", "taskstatus", "assignesto", "projectid", "createdby",
				}).AddRow(
					"t1", "Task 1", "Desc 1", "Crit 1", "2025-08-30",
					Priority.High, status.Done, empID, projectID, "admin",
				).AddRow(
					"t2", "Task 2", "Desc 2", "Crit 2", "2025-09-01",
					Priority.Low, status.InProgress, empID, projectID, "manager",
				)

				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby 
              		 FROM tasks 
              		 WHERE projectid = ? AND assignesto = ?`)).
					WithArgs(projectID, empID).
					WillReturnRows(rows)
			},
			expectTasks: 2,
			expectError: "",
		},
		{
			name:      "query error",
			projectID: "proj2",
			empID:     "emp2",
			mockSetup: func(mock sqlmock.Sqlmock, projectID, empID string) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby 
              		 FROM tasks 
              		 WHERE projectid = ? AND assignesto = ?`)).
					WithArgs(projectID, empID).
					WillReturnError(sql.ErrConnDone)
			},
			expectTasks: 0,
			expectError: "conn",
		},
		// {
		// 	name:      "scan error",
		// 	projectID: "proj3",
		// 	empID:     "emp3",
		// 	mockSetup: func(mock sqlmock.Sqlmock, projectID, empID string) {
		// 		rows := sqlmock.NewRows([]string{
		// 			"task_id", "title", "description", "acceptance_criteria", "deadline",
		// 			"taskpriority", "taskstatus", "assignesto", "projectid", "createdby",
		// 		}).AddRow(
		// 			1, // wrong type, should be string
		// 			"Task X", "Desc X", "Crit X", "2025-09-10",
		// 			Priority.Low, status.Done, empID, projectID, "admin",
		// 		)

		// 		mock.ExpectQuery(regexp.QuoteMeta(
		// 			`SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby 
        //       		 FROM tasks 
        //       		 WHERE projectid = ? AND assignesto = ?`)).
		// 			WithArgs(projectID, empID).
		// 			WillReturnRows(rows)
		// 	},
		// 	expectTasks: 0,
		// 	expectError: "converting", // type conversion fails
		// },
	// 	{
	// 		name:      "rows iteration error",
	// 		projectID: "proj4",
	// 		empID:     "emp4",
	// 		mockSetup: func(mock sqlmock.Sqlmock, projectID, empID string) {
	// 			rows := sqlmock.NewRows([]string{
	// 				"task_id", "title", "description", "acceptance_criteria", "deadline",
	// 				"taskpriority", "taskstatus", "assignesto", "projectid", "createdby",
	// 			}).RowError(0, sql.ErrNoRows)

	// 			mock.ExpectQuery(regexp.QuoteMeta(
	// 				`SELECT task_id, title, description, acceptance_criteria, deadline, taskpriority, taskstatus, assignesto, projectid, createdby 
    //           		 FROM tasks 
    //           		 WHERE projectid = ? AND assignesto = ?`)).
	// 				WithArgs(projectID, empID).
	// 				WillReturnRows(rows)
	// 		},
	// 		expectTasks: 0,
	// 		expectError: "no rows",
	// 	},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo :=NewTaskRepo(db)

			tt.mockSetup(mock, tt.projectID, tt.empID)

			result, err := repo.ViewAllAssignedTasksInProject(tt.projectID, tt.empID)

			if tt.expectError == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.expectError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectError) {
					t.Fatalf("expected error containing %q, got %v", tt.expectError, err)
				}
			}

			if len(result) != tt.expectTasks {
				t.Fatalf("expected %d tasks, got %d", tt.expectTasks, len(result))
			}

			if tt.expectTasks > 0 {
				// ensure deadline parsed properly
				if result[0].Deadline.IsZero() {
					t.Fatalf("expected parsed deadline, got zero value")
				}
				// ensure first row values are correct
				if result[0].TaskId == "" || result[0].Title == "" {
					t.Fatalf("expected valid task data, got %+v", result[0])
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}




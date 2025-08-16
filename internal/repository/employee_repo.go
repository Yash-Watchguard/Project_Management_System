 package repository

// import (
// 	"encoding/json"
// 	"os"
//      "errors"
// 	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
// 	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
// )

// type EmployeeRepo struct {
// 	filepath string
// }

// func NewEmployeeRepo() *EmployeeRepo {
// 	return &EmployeeRepo{filepath: "C:/Users/ygoyal/Desktop/PMS_Project/Pms/internal/data/task.json"}
// }

// func (er *EmployeeRepo) ViewAssignedTask(empId string) ([]task.Task, error) {
// 	var tasks []task.Task
// 	var assignedTasks []task.Task

// 	data, err := os.ReadFile(er.filepath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(data) > 0 {
// 		err = json.Unmarshal(data, &tasks)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	for _, task := range tasks {
// 		if task.AssignedTo == empId {
// 			assignedTasks = append(assignedTasks, task)
// 		}
// 	}

// 	return assignedTasks, nil
// }
// func (er *EmployeeRepo) UpdateTaskStatus(empId string, taskId string, updatedStatus status.TaskStatus) error {
// 	var tasks []task.Task
// 	data, err := os.ReadFile(er.filepath)
// 	if err != nil {
// 		return err
// 	}
// 	if len(data) > 0 {
// 		err = json.Unmarshal(data, &tasks)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	updated := false
//     for i := range tasks {
// 	if tasks[i].TaskId == taskId && tasks[i].AssignedTo == empId {
// 		tasks[i].TaskStatus = updatedStatus
// 		updated = true
// 		break
// 	}
//     }


// 	if !updated {
// 		return errors.New("task not assigned to employee")
// 	}

// 	newData, err := json.MarshalIndent(tasks, "", "  ")
// 	if err != nil {
// 		return err
// 	}

// 	err = os.WriteFile(er.filepath, newData, 0644)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }


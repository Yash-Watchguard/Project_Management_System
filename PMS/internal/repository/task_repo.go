package repository

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
)

type TaskRepo struct {
	filepath string
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{filepath: "C:/Users/ygoyal/Desktop/PMS_Project/Pms/internal/data/task.json"}
}

func (taskRepo *TaskRepo) ViewAllTask(projectId string) ([]task.Task, error) {
	data, err := os.ReadFile(taskRepo.filepath)
	if err != nil {
		return nil, err
	}

	var allTasks []task.Task
	err = json.Unmarshal(data, &allTasks)
	if err != nil {
		return nil, err
	}

	var projectTasks []task.Task
	for _, t := range allTasks {
		if t.ProjectId == projectId {
			projectTasks = append(projectTasks, t)
		}
	}

	return projectTasks, nil
}
func (taskRepo *TaskRepo) SaveTask(newTask task.Task) error {
	var tasks []task.Task

	data, err := os.ReadFile(taskRepo.filepath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			return err
		}
	} else {
		tasks = []task.Task{} 
	}

	tasks = append(tasks, newTask)

	newData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(taskRepo.filepath, newData, 0644)
	if err != nil {
		return err
	}

	return nil
}


func (taskRepo *TaskRepo) DeleteTask(taskId string) error {
	var tasks []task.Task

	data, err := os.ReadFile(taskRepo.filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return err
	}


	found := false
	var updatedTasks []task.Task
	for _, t := range tasks {
		if t.TaskId == taskId {
			found = true
			continue 
		}
		updatedTasks = append(updatedTasks, t)
	}

	if !found {
		return errors.New("task not found")
	}

	newData, err := json.MarshalIndent(updatedTasks, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(taskRepo.filepath, newData, 0644)
	if err != nil {
		return err
	}

	return nil
}


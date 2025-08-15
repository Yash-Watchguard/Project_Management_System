package handler

import (
	"context"

	"github.com/Yash-Watchguard/Tasknest/internal/service1"

	"errors"
	"fmt"
	"time"

	"github.com/Yash-Watchguard/Tasknest/internal/model/priority"

	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"

	// "github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/fatih/color"
)

type TaskHandler struct {
	taskService *service1.TaskService
}
func NewTaskHandler(taskService *service1.TaskService)*TaskHandler{
	return &TaskHandler{taskService: taskService}
}

func (th *TaskHandler) ViewAllTask(ctx context.Context, projectId string) error {
	if projectId == "" {
		return errors.New("projectId cannot be empty")
	}

	// Fetch tasks for the project
	tasks, err := th.taskService.ViewAllTask(ctx, projectId)
	if err != nil {
		color.Red("Failed to fetch tasks: %v", err)
		return err
	}

	// If no tasks exist
	if len(tasks) == 0 {
		color.Yellow("No tasks found for this project.")
		return nil
	}

	// Display all tasks
	for i, task := range tasks {
		color.Cyan("------------ Task %d ------------", i+1)
		color.Cyan("Task ID        : %v", task.TaskId)
		color.Cyan("Title          : %v", task.Tile)
		color.Cyan("Description    : %v", task.Description)
		color.Cyan("Priority       : %v", task.TaskPriority)
		color.Cyan("Assigned To    : %v", task.AssignedTo)
		color.Cyan("Status         : %v", task.TaskStatus)
		color.Cyan("Deadline       : %v", task.Deadline)
		fmt.Println()
	}

	return nil
}
func (th *TaskHandler) CreateTask(ctx context.Context, projectId string) error {
	managerId := ctx.Value(ContextKey.UserId).(string)
	taskId := GenerateUUID()

	title, err := GetInput("Enter Task Title: ")
	if err != nil {
		return err
	}

	description, err := GetInput("Enter Task Description: ")
	if err != nil {
		return err
	}

	var deadline time.Time
	for {
		deadlineStr, err := GetInput("Enter Deadline in YYYY-MM-DD: ")
		if err != nil {
			return err
		}

		deadline, err = TimeParser(deadlineStr)
		if err != nil {
			color.Red("Invalid date format")
		} else {
			break
		}
	}

	var priority Priority.Priority
	for {
		priorityStr, err := GetInput("Enter Priority => Low/Medium/High: ")
		if err != nil {
			return err
		}

		priority, err = Priority.PriorityParser(priorityStr)
		if err != nil {
			color.Red("Invalid priority. Choose Low, Medium, or High.")
		} else {
			break
		}
	}

	assignedTo, err := GetInput("Enter Employee ID to assign this task to: ")
	if err != nil {
		return err
	}

	newTask := task.Task{
		TaskId:       taskId,
		Tile:         title,
		Description:  description,
		Deadline:     deadline,
		TaskPriority: priority,
		TaskStatus:   status.Pending,
		AssignedTo:   assignedTo,
		ProjectId:    projectId,
		CreatesBy:    managerId,
	}

	// Call service method
	if err := th.taskService.CreateTask(ctx, managerId, newTask); err != nil {
		return err
	}

	color.Green("Task created successfully!")

	color.Blue("Press Enter to go back...")
	fmt.Scanln()
	return nil
}
func (th *TaskHandler) DeleteTask(ctx context.Context, projectId string) error {
	managerId := ctx.Value(ContextKey.UserId).(string)

	// Fetch all tasks for the given project
	projectTasks, err := th.taskService.ViewAllTask(ctx, projectId)
	if err != nil {
		return err
	}

	if len(projectTasks) == 0 {
		return errors.New("no task created for this project")
	}

	// Display tasks
	for i, task := range projectTasks {
		color.Yellow("%d. Name: %s  ID: %s", i+1, task.Tile, task.TaskId)
	}

	// Get task ID to delete
	taskId, err := GetInput("Enter Task Id to delete: ")
	if err != nil {
		return err
	}

	// Delete task using service
	if err := th.taskService.DeleteTask(ctx, managerId, taskId); err != nil {
		return err
	}

	color.Green("Task deleted successfully!")

	color.Blue("Press Enter to go back...")
	fmt.Scanln()
	return nil
}
func (th *TaskHandler) GetAssignedTask(ctx context.Context, userId string) error {
    tasks, err := th.taskService.GetAssigenedTask(ctx, userId)
    if err != nil {
        return fmt.Errorf("failed to fetch assigned tasks: %v", err)
    }

    if len(tasks) == 0 {
        color.Yellow("No tasks assigned.")
        return nil
    }

    for i, task := range tasks {
        color.Cyan("------------ Task %d ------------", i+1)
        color.Cyan("Task ID       : %v", task.TaskId)
        color.Cyan("Title         : %v", task.Tile)
        color.Cyan("Description   : %v", task.Description)
        color.Cyan("Priority      : %v", task.TaskPriority)
        color.Cyan("Status        : %v", task.TaskStatus)
        color.Cyan("Created By    : %v", task.CreatesBy)
        color.Cyan("Deadline      : %v", task.Deadline.Format("2006-01-02 15:04:05"))
        fmt.Println()
    }

    
    return nil
}

func(th *TaskHandler)UpdateTaskStatus(ctx context.Context,taskId string)error{

updatedStatus,err:=GetInput("Enter Updated status (pending/in progress/done) : ")
   if err!=nil{
	color.Red("error in getting input")
   }
   newStatus:=status.GetStatusFromString(updatedStatus)
   err=th.taskService.UpdateTaskStatus(ctx,taskId,newStatus)
   if err!=nil{
	return err
   }
   color.Green("Task Status Updated Successfully")
   color.Blue("press enter to going back")
   fmt.Scanln()
   return nil
} 






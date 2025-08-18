package main

import (
	"context"

	"strconv"

	"github.com/Yash-Watchguard/Tasknest/handler"

	"github.com/Yash-Watchguard/Tasknest/internal/constants"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"

	"github.com/fatih/color"
)

func ManagerDashboard(
	ctx context.Context,
	user *user.User,
	userHandler *handler.UserHandler,
	projectHandler *handler.ProjectHandler,
	taskHandler *handler.TaskHandler,
	commentHandler *handler.CommentHandler,
) {
	for {
		color.Cyan(constants.ManagerDashbEntry)
		color.Cyan("1. View Profile")
		color.Cyan("2. View Assigned Projects")
		color.Cyan("3. View All Employees")
		color.Cyan("4. Promote Employee")
		color.Cyan("5. Logout")

		choiceStr, _ := handler.GetInput("\nEnter your choice: ")
        choice, _ := strconv.Atoi(choiceStr)

		switch choice {
		case 1:
			exist, err := userHandler.ViewUserProfile(ctx,user)
			if err != nil {
				color.Red("%v", err)
			}
			if exist {
				return
			}
		case 2:
			assignedProjectsMenu(ctx, user, projectHandler, taskHandler,commentHandler)
		case 3:
			userHandler.ViewAllEmployees(ctx)
			handler.Pause()
		case 4:
			userHandler.PromoteEmployee(ctx)
		case 5:
			color.Green("Logging out...")
			return
		default:
			color.Red("Invalid choice. Please select a valid option.")
		}
	}
}
func assignedProjectsMenu(ctx context.Context, user *user.User, ph *handler.ProjectHandler, th *handler.TaskHandler,ch *handler.CommentHandler) {
	err := ph.ViewAssignedProjects(ctx, user)
	if err != nil {
		color.Red("%v", err)
		return
	}

	color.Blue("1. manage a project")
	color.Blue("2. go back")
	choiceStr, _ := handler.GetInput("\nEnter your choice: ")
    choice, _ := strconv.Atoi(choiceStr)
	if choice == 1 {
		var projectId string
		projectId,_=handler.GetInput("Enter Project Id : ")
		projectTaskMenu(ctx, projectId, th,ph, ch)
	}
}
func projectTaskMenu(ctx context.Context, projectId string, th *handler.TaskHandler, ph *handler.ProjectHandler,ch *handler.CommentHandler) {
	for {
		color.Cyan("----------- Task Menu -----------")
		color.Cyan("1. View Tasks")
		color.Cyan("2. Create Task")
		color.Cyan("3. Delete Task")
		color.Cyan("4. Show Project Status")
		color.Cyan("5. Back")
		color.Cyan("---------------------------------")
		choiceStr, _ := handler.GetInput("\nEnter your choice: ")
        choice, _ := strconv.Atoi(choiceStr)

		switch choice {
		case 1:
			err:=th.ViewAllTask(ctx,projectId)
			if err != nil {
				color.Red("%v", err)
				handler.Pause()
			    continue
			}
			taskId, err := handler.GetInput("Enter Task ID For managing comments(or press Enter to go back): ")
            if err != nil || taskId == "" {
            break
            }
			commentMenu(ctx, ch, taskId)
		case 2:
			err:=th.CreateTask(ctx, projectId)
			if err != nil {
				color.Red("%v", err)
			}
		case 3:
			err:=th.DeleteTask(ctx, projectId)
			if err != nil {
				color.Red("%v", err)
			}
		case 4:
			err:=ph.ShowProjectStatus(ctx, projectId)
			if err != nil {
				color.Red("%v", err)
			}
			handler.Pause()
		case 5:
			return
		default:
			color.Red("Invalid choice.")
		}
	}
}


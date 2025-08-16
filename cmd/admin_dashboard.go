package main

import (
	"context"

	"github.com/Yash-Watchguard/Tasknest/handler"
	"github.com/Yash-Watchguard/Tasknest/internal/constants"
	"strconv"

	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	
	"github.com/fatih/color"
)

func AdminDashboard(ctx context.Context, user *user.User,userHandler *handler.UserHandler,taskHandler *handler.TaskHandler,projectHandler *handler.ProjectHandler,commentHandler *handler.CommentHandler) {
	for {
		color.Blue(constants.AdminDashbEntry)
		color.Blue("1. View Profile")
		color.Blue("2. View All Users")
		color.Blue("3. Delete User")
		color.Blue("4. Promote Employee to Manager")
		color.Blue("5. Manage Projects (Add, View, Delete)")
		// color.Blue("5. Add New Project")
		// color.Blue("6. View All Projects")
		// color.Blue("7. Delete Project")
		color.Blue("6. Logout")

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
			err := userHandler.ViewallUsers(ctx)
			if err != nil {
				color.Red("%v", err)
			}

		case 3:
			err :=userHandler.DeleteUser(ctx)
			if err != nil {
				color.Red("Error: %v", err)
			}
		case 4:
			err := userHandler.PromoteEmployee(ctx)
			if err != nil {
				color.Red("Error: %v", err)
			}
		case 5:
			 projectMenu(ctx, projectHandler, taskHandler, commentHandler)

		case 6:
			color.Green("Logging out...")
			return

		default:
			color.Red("Invalid choice. Please try again.")
		}
	}
}
func projectMenu(ctx context.Context, projectHandler *handler.ProjectHandler, taskHandler *handler.TaskHandler, commentHandler *handler.CommentHandler) {
    for {
        color.Blue("1. Add New Project")
        color.Blue("2. View All Projects")
        color.Blue("3. Delete Project")
        color.Blue("4. Back")

        choiceStr, _ := handler.GetInput("\nEnter your choice: ")
        choice, _ := strconv.Atoi(choiceStr)

        switch choice {
        case 1:
            err:= projectHandler.AddNewProject(ctx)
			if err != nil {
			 		color.Red("%v", err)
			 	}

        case 2:
            projectId, err := projectHandler.SelectAndReturnProjectId(ctx)
            if err != nil || projectId == "" {
                continue
            }
            taskMenu(ctx, taskHandler, commentHandler, projectId)

        case 3:
            err:= projectHandler.DeleteProject(ctx)
			if err != nil {
			 		color.Red("%v", err)
			 	}

        case 4:
            return

        default:
            color.Red("Invalid choice. Please try again.")
        }
    }
}

func taskMenu(ctx context.Context, taskHandler *handler.TaskHandler, commentHandler *handler.CommentHandler, projectId string) {
    err := taskHandler.ViewAllTask(ctx, projectId)
    if err != nil { return }

    taskId, err := handler.GetInput("Enter Task ID (or press Enter to go back): ")
    if err != nil || taskId == "" {
        return
    }
    commentMenu(ctx, commentHandler, taskId)
}

func commentMenu(ctx context.Context, commentHandler *handler.CommentHandler, taskId string) {
    for {
        color.Cyan("1. View All Comments")
        color.Cyan("2. Add Comment")
        color.Cyan("3. Update Comment")
        color.Cyan("4. Delete Comment")
        color.Cyan("5. Back")

       choiceStr, _ := handler.GetInput("\nEnter your choice: ")
    choice, _ := strconv.Atoi(choiceStr)

        switch choice {
        case 1:
            err:=commentHandler.ViewAllComment(ctx, taskId)
			if err!=nil{
				color.Red("%s",err)
			}
        case 2:
            err:= commentHandler.AddNewComment(ctx, taskId)
			if err!=nil{
				color.Red("%s",err)
			}
        case 3:
            _ = commentHandler.UpdateComment(ctx, taskId)
        case 4:
            _ = commentHandler.DeleteComment(ctx, taskId)
        case 5:
            return
        default:
            color.Red("Invalid choice")
        }
    }
}

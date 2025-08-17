package main

import (
	"context"
	"strconv"
	"github.com/Yash-Watchguard/Tasknest/handler"
	"github.com/Yash-Watchguard/Tasknest/internal/constants"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/fatih/color"
)

func employeeDashboard(ctx context.Context, user *user.User,userHandler *handler.UserHandler ,taskHandler *handler.TaskHandler, commentHandler *handler.CommentHandler) {
    // use taskHandler and commentHandler wherever needed

	for{
        color.Blue(constants.EmployDashbEntry)
		color.Blue("1. View Profile")
		color.Blue("2. View All Assigned Task")
		color.Blue("3  Update Task Status")
		color.Blue("4. Logout")

		choiceStr, _ := handler.GetInput("\nEnter your choice: ")
        choice, _ := strconv.Atoi(choiceStr)

		switch choice{
		case 1:
			exist, err := userHandler.ViewUserProfile(ctx,user)
			if err != nil {
				color.Red("%v", err)
			}
			if exist {
				return
			}
		case 2:
			err:=taskHandler.GetAssignedTask(ctx,user.Id)
			if err != nil {
				color.Red("%v", err)
			}
			taskId, err := handler.GetInput("Enter Task ID For managing comments(or press Enter to go back): ")
            if err != nil || taskId == "" {
            break
            }
			commentMenu(ctx, commentHandler, taskId)
			handler.Pause()
		case 3:
			taskId,err:=handler.GetInput("Enter Task Id : ")
			if err!=nil{
				color.Red("error in getting input")
			}
			err=taskHandler.UpdateTaskStatus(ctx,taskId)
			if err != nil {
				color.Red("%v", err)
			}
		case 4:
			color.Green("Logging Out......")
			return 
		}
	}
}


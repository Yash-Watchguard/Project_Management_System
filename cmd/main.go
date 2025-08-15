package main

import (
	// "fmt"

	"log"

	"context"
	// "errors"
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/handler"
	"github.com/Yash-Watchguard/Tasknest/internal/constants"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/fatih/color"
	// "github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	// "github.com/Yash-Watchguard/Tasknest/internal/model/project"
	// "github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	// "github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	
)

func main() {
	err := RunApp()
	if err != nil {
		log.Fatal(nil)
	}
}
func RunApp() error {
	var ctx context.Context
	userRepo := repository.NewUserRepo()
	projectRepo := repository.NewProjectRepo()
	
	taskRepo := repository.NewTaskRepo()
	commentRepo := repository.NewCommentRepo()
	
	empRepo:=repository.NewEmployeeRepo()

	authService:=service1.NewAuthService(userRepo)
    commentService:=service1.NewCommentService(commentRepo)
	projectService:=service1.NewProjectService(projectRepo)
	taskService:=service1.NewTaskService(taskRepo,empRepo)
	userService:=service1.NewUserService(userRepo)

    authHandler:=handler.NewAuthHandler(authService)
	userHandler:=handler.NewUserHandler(userService)
	projectHandler:=handler.NewProjectHandler(projectService,userService,taskService)
	taskHandler:=handler.NewTaskHandler(taskService)
	commentHandler:=handler.NewCommentHandler(commentService,userService)

	for {
		color.Red(constants.WelcomeMSG)
		color.Blue(constants.SignupChoice)

		color.Blue(constants.LoginChoice)
		color.Blue(constants.ExitChoice)

		fmt.Print(color.CyanString("Enter Your Choice : "))

		var Choice int
		fmt.Scanln(&Choice)

		switch Choice {
		case 1:
			err := authHandler.Signup()
			if err!=nil{
				color.Red("%v",err)
			}

		case 2:
			user,err:= authHandler.Login()
			if err!=nil{
				color.Red("%v",err)
			}
			ctx=context.Background()
	        ctx=context.WithValue(ctx,ContextKey.UserId,user.Id)
	        ctx=context.WithValue(ctx,ContextKey.UserPassword,user.Password)
	        ctx=context.WithValue(ctx,ContextKey.UserRole,user.Role)

	        color.Green("Welcom back %s in Worknest☺️", user.Name)
	
	        switch user.Role{
	        case 0:
		       AdminDashboard(ctx,user,userHandler,taskHandler,projectHandler,commentHandler)
            case 1:
		       ManagerDashboard(ctx,user,userHandler,projectHandler,taskHandler,commentHandler)
	        case 2:
	           employeeDashboard(ctx,user,userHandler,taskHandler,commentHandler)
			}
		case 3:
			   fmt.Println(constants.GoodByeMsg)
			   return nil
		default:
			color.Red(constants.InvalidChoice)
		}
	}

}

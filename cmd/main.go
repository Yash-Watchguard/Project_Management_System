// package main

// import (
// 	"log"

// 	"context"

// 	"fmt"

// 	"github.com/Yash-Watchguard/Tasknest/internal/config"
// 	"github.com/Yash-Watchguard/Tasknest/internal/constants"
// 	"github.com/Yash-Watchguard/Tasknest/internal/handler"
// 	"github.com/Yash-Watchguard/Tasknest/internal/logger"
// 	"github.com/Yash-Watchguard/Tasknest/internal/service1"
// 	"github.com/fatih/color"

// 	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"

// 	"github.com/Yash-Watchguard/Tasknest/internal/repository"
// )

// func main() {
// 	err := RunApp()
// 	if err != nil {
// 		log.Fatal(nil)
// 	}
// }
// func RunApp() error {
// 	logger.InitLogger()
// 	defer logger.CloseLogger()
// 	logger.Info("application started")
// 	db,err:=config.GetDbInstance()
// 	if err!=nil{
// 		log.Fatal(err)
// 	}

// 	defer func() {
// 		db.Close()
// 	}()

// 	var ctx context.Context
// 	userRepo := repository.NewUserRepo(db)
// 	projectRepo := repository.NewProjectRepo(db)

// 	taskRepo := repository.NewTaskRepo(db)
// 	commentRepo := repository.NewCommentRepo(db)

// 	authService:=service1.NewAuthService(userRepo)
//  commentService:=service1.NewCommentService(commentRepo)
// 	projectService:=service1.NewProjectService(projectRepo)
// 	taskService:=service1.NewTaskService(taskRepo)
// 	userService:=service1.NewUserService(userRepo)

//  authHandler:=handler.NewAuthHandler(authService)
// 	userHandler:=handler.NewUserHandler(userService)
// 	projectHandler:=handler.NewProjectHandler(projectService,userService,taskService)
// 	taskHandler:=handler.NewTaskHandler(taskService)
// 	commentHandler:=handler.NewCommentHandler(commentService,userService)

// 	for {
// 		color.Red(constants.WelcomeMSG)
// 		color.Blue(constants.SignupChoice)

// 		color.Blue(constants.LoginChoice)
// 		color.Blue(constants.ExitChoice)

// 		fmt.Print(color.CyanString("Enter Your Choice : "))

// 		var Choice int
// 		fmt.Scanln(&Choice)

// 		switch Choice {
// 		case 1:
// 			err := authHandler.Signup()
// 			if err!=nil{
// 				color.Red("%v",err)
// 			}

// 		case 2:
// 			user,err:= authHandler.Login()
// 			if err!=nil{
// 				color.Red("%v",err)
// 				continue
// 			}
// 			ctx=context.Background()
// 	        ctx=context.WithValue(ctx,ContextKey.UserId,user.Id)
// 	        ctx=context.WithValue(ctx,ContextKey.UserPassword,user.Password)
// 	        ctx=context.WithValue(ctx,ContextKey.UserRole,user.Role)

// 	        color.Green("Welcom back, %s, to Worknest☺️", user.Name)

// 	        switch user.Role{
// 	        case 0:
// 		       AdminDashboard(ctx,user,userHandler,taskHandler,projectHandler,commentHandler)
//             case 1:
// 		       ManagerDashboard(ctx,user,userHandler,projectHandler,taskHandler,commentHandler)
// 	        case 2:
// 	           employeeDashboard(ctx,user,userHandler,taskHandler,commentHandler)
// 			}
// 		case 3:
// 			   fmt.Println(constants.GoodByeMsg)
// 			   return nil
// 		default:
// 			color.Red(constants.InvalidChoice)
// 		}
// 	}

// }
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"fmt"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	
	"github.com/Yash-Watchguard/Tasknest/internal/logger"

	
)
func main(){
	logger.InitLogger()
    defer logger.CloseLogger()

	db,err:=config.GetDbInstance()
	if err!=nil{
		log.Fatal(err)
	}

	defer func() {
		db.Close()
	}()

	c:=make(chan os.Signal,1)

	signal.Notify(c,os.Interrupt,syscall.SIGTERM)
	go func(){
		<-c
		fmt.Println("\nDisconnecting from mysql...")
		db.Close()
		os.Exit(1)
	}()
	logger.Info("Start the application....")

	runApp(db)
}

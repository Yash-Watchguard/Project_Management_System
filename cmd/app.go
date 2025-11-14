
package main

import (
	"database/sql"
	// "fmt"
	// "log"
	// "net/http"
	// "os"

	// "github.com/Yash-Watchguard/Tasknest/internal/config"

	// "github.com/Yash-Watchguard/Tasknest/internal/middleware"
	// "github.com/Yash-Watchguard/Tasknest/internal/routers"
	// "github.com/Yash-Watchguard/Tasknest/internal/service1"

	// "github.com/Yash-Watchguard/Tasknest/internal/repository"
)

func runApp(db *sql.DB) {
	// userRepo := repository.NewUserRepo(db)
	// projectRepo := repository.NewProjectRepo(db)

	// taskRepo := repository.NewTaskRepo(db)
	// commentRepo := repository.NewCommentRepo(db)

	// authService := service1.NewAuthService(userRepo)
	// commentService := service1.NewCommentService(commentRepo)
	// projectService := service1.NewProjectService(projectRepo)
	// taskService := service1.NewTaskService(taskRepo)
	// userService := service1.NewUserService(userRepo)

	// router := routers.SetupRouter(authService, userService, projectService, taskService, commentService)

	// handler := middleware.CorsMiddleWare(
	// 	middleware.LoggingMiddleWare(router),
	// )

	// log.Println("Server starting on the Port 8080...")
	// fmt.Println("Server starting on the Port 8080...")

	// err := http.ListenAndServe(config.Port, handler)
	// if err != nil {
	// 	log.Println(err)
	// 	fmt.Println("Error starting server:", err)
	// 	os.Exit(1)
	// }
}

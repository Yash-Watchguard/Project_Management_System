package routers

import (
	"net/http"

	"github.com/Yash-Watchguard/Tasknest/internal/handler"
	"github.com/Yash-Watchguard/Tasknest/internal/middleware"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
  
)

func SetupRouter(authService service1.AuthServiceInterface, userService service1.UserServiceInterface, projectService service1.ProjectServiceInterface,taskSeervice service1.TaskServiceInterface,commentService service1.CommentServiceInterface)*http.ServeMux{
   r:=http.NewServeMux()

//    handlers
   authHandler:=handler.NewAuthHandler(authService,userService)
   userHandler:=handler.NewUserHandler(userService)
   projectHandler:=handler.NewProjectHandler(projectService,userService,taskSeervice)
   taskHandler:=handler.NewTaskHandler(taskSeervice)
   commentHandler:=handler.NewCommentHandler(commentService,userService)

//    routes for the auth
r.Handle("/v1/signup",http.HandlerFunc(authHandler.Signup))
r.Handle("/v1/login",http.HandlerFunc(authHandler.Login))

//  routes for the userService
r.Handle("/v1/users/", middleware.AuthMiddleWare(http.HandlerFunc(userHandler.UsersHandler)))

// routes for project
r.Handle("/v1/projects/", middleware.AuthMiddleWare(http.HandlerFunc(projectHandler.ProjectHandler)))
r.Handle("POST /v1/projects", middleware.AuthMiddleWare(http.HandlerFunc(projectHandler.CreateProject)))
r.Handle("DELETE /v1/projects/{project_id}", middleware.AuthMiddleWare(http.HandlerFunc(projectHandler.DeleteProject)))


// routes for the tasks

r.Handle("GET /v1/projects/{project_id}/tasks/",middleware.AuthMiddleWare(http.HandlerFunc(taskHandler.GetTask)))
r.Handle("GET /v1/projects/tasks/{employee_id}",middleware.AuthMiddleWare(http.HandlerFunc(taskHandler.AssignedTasks)))
r.Handle("POST /v1/projects/{project_id}/tasks",middleware.AuthMiddleWare(http.HandlerFunc(taskHandler.CreateTask)))
r.Handle("DELETE /v1/projects/{project_id}/tasks/{task_id}",middleware.AuthMiddleWare(http.HandlerFunc(taskHandler.DeleteTask)))
r.Handle("PATCH /v1/projects/{project_id}/tasks/{task_id}",middleware.AuthMiddleWare(http.HandlerFunc(taskHandler.UpdateStatus)))

// rotes for the comments
r.Handle("GET /v1/projects/{project_id}/tasks/{task_id}/comments",middleware.AuthMiddleWare(http.HandlerFunc(commentHandler.ViewAllComment)))
r.Handle("POST /v1/projects/{project_id}/tasks/{task_id}/comments",middleware.AuthMiddleWare(http.HandlerFunc(commentHandler.AddComment)))
r.Handle("PATCH /v1/projects/{project_id}/tasks/{task_id}/comments/{comment_id}",middleware.AuthMiddleWare(http.HandlerFunc(commentHandler.UpdateComment)))
r.Handle("DELETE /v1/projects/{project_id}/tasks/{task_id}/comments/{comment_id}",middleware.AuthMiddleWare(http.HandlerFunc(commentHandler.DeleteComment)))
return r
}
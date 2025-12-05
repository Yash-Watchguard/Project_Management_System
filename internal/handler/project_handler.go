package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Yash-Watchguard/Tasknest/internal/logger"
	
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/response"

	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"

	"encoding/json"

	"github.com/Yash-Watchguard/Tasknest/internal/service1"
)

type ProjectHandler struct {
	userService service1.UserServiceInterface
	projectService service1.ProjectServiceInterface
	taskService service1.TaskServiceInterface
}

func NewProjectHandler(projectService service1.ProjectServiceInterface ,userService service1.UserServiceInterface, taskService service1.TaskServiceInterface)*ProjectHandler{
	return &ProjectHandler{projectService: projectService,userService: userService,taskService: taskService}
}

func(ph *ProjectHandler)ProjectStatus(w http.ResponseWriter,r *http.Request){
     projectId:=r.PathValue("project_id")

	 if len(projectId)!=36{
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid Project Id", 1000)
		return
	 }

	 role:=r.Context().Value(ContextKey.UserRole).(roles.Role)

	 if role==roles.Employee{
		logger.Error("unauthorized person wants to view all projects")
        response.ErrorResponse(w, http.StatusForbidden, "Access denied", 1008)
        return
	 }
     
	projectTasks,err:=ph.taskService.ViewAllTask(projectId)

	if err!=nil{
		logger.Error("error in fatching tasks")
		response.ErrorResponse(w,http.StatusInternalServerError,"Error fatching the tasks",1010)
		return
	}
    //   calculate the project status on the basis of the tasks and task  status
     
    
	// if len(projectTasks) == 0 {
	// 	logger.Error("No task found")
	// 	response.ErrorResponse(w,http.StatusNotFound,"No Task",1000)
	// 	return
	// }
    
	total := len(projectTasks)
	done := 0

    if(total>0){
	for _, t := range projectTasks {
		if t.TaskStatus == status.Done {
			done++
		}
	}
}
    percerntDone:=float64(0)
	if(total>0){
	percerntDone = (float64(done) / float64(total)) * 100
	}

	

	type ProjectStatusResponse struct {
        ProjectID            string  `json:"projectId"`
        CompletedTasks       int     `json:"completedTasks"`
        TotalTasks           int     `json:"totalTasks"`
        CompletionPercentage float64 `json:"completionPercentage"`
    }

    statusResponse := &ProjectStatusResponse{
        ProjectID:            projectId,
        CompletedTasks:       done,
        TotalTasks:           total,
        CompletionPercentage: percerntDone,
    }

	logger.Info("status get successfully")
	response.SuccessResponse(w,statusResponse,"Status Fatched",http.StatusOK)
}
func (ph *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    role, ok := ctx.Value(ContextKey.UserRole).(roles.Role)
    if !ok {
        logger.Error("user role not found")
        response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
        return
    }

    pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/") // remove leading/trailing slashes

    if len(pathSegments) < 2 || pathSegments[0] != "v1" || pathSegments[1] != "projects" {
        response.ErrorResponse(w, http.StatusBadRequest, "Invalid path", 1008)
        return
    }

 
    if len(pathSegments) == 2 {
		if role!=roles.Admin{
			logger.Error("unauthorized person wants to view all projects")
            response.ErrorResponse(w, http.StatusForbidden, "Access denie", 1008)
            return
		}
        projects, err := ph.projectService.ViewAllProjects()
        if err != nil {
			logger.Error("Failed to fetch projects")
            response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch projects", 1009)
            return
        }

		if len(projects)==0{
			logger.Error("No projects assigned ")
            response.ErrorResponse(w, http.StatusNotFound, "No projects assig", 404)
            return
		}
		logger.Info("Projects Retrived Successfully")
        response.SuccessResponse(w, projects,"Projects Retrived Successfully",http.StatusOK)
        return
    }

    
    assignedUserID := pathSegments[2]
    projects, err := ph.projectService.ViewAssignedProject(assignedUserID)
    if err != nil {
		logger.Error("Failed to fetch user's projects")
        response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch Assigned projects", 1010)
        return
    }

    if len(projects)==0{
			logger.Error("No projects assigned ")
            response.ErrorResponse(w, http.StatusNotFound, "No projects assigned ", 404)
            return
	}
    logger.Info("Projects retrived successfully")
    response.SuccessResponse(w,projects,"Projects retrived successfully",http.StatusOK)
	
}

func(ph *ProjectHandler)CreateProject(w http.ResponseWriter,r * http.Request){
    ctx := r.Context()
	createdBy:=ctx.Value(ContextKey.UserId).(string)
    role, ok := ctx.Value(ContextKey.UserRole).(roles.Role)
    if !ok {
        logger.Error("user role not found")
        response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
        return
    }

	if role != roles.Admin {
		logger.Error("unauthorized to create projects")
        response.ErrorResponse(w, http.StatusForbidden, "Only admin can create projects", 1011)
        return
    }

	var projectReq struct{
		ProjectName       string `json:"projectName"`
        ProjectDescription string `json:"projectDescription"`
        Deadline          string `json:"deadline"` 
        AssignedManagerID string `json:"assignedManagerId"`
	}
  

	if err := json.NewDecoder(r.Body).Decode(&projectReq); err != nil {
		logger.Error("Invalid input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	ManagerProfile,_:=ph.userService.ViewProfile(projectReq.AssignedManagerID)

	if ManagerProfile[0].Status==user.InActive{
		logger.Error("Manger is not available of for project assignement")
		response.ErrorResponse(w,http.StatusBadRequest,"Manager is InActive",1000)
	}

	err := validate.Struct(projectReq)
	if err != nil {
		logger.Error("Validation error")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}

	projectId := GenerateUUID()
	var actualdeadline time.Time

	actualdeadline,_ = TimeParser(projectReq.Deadline)

	project := &project.Project{
		ProjectId:          projectId,
		ProjectName:        projectReq.ProjectName,
		ProjectDescription: projectReq.ProjectDescription,
		Deadline:           actualdeadline,
		CreatedBy:          createdBy,
		AssignedManager:    projectReq.AssignedManagerID,
	}

	err =ph.projectService.AddProject(*project)
	if err != nil {
		fmt.Printf("%v",err)
		logger.Error("Error creating project")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error creating project", 1006)
		return
	}

	logger.Info("Project created sucessfully")
	response.SuccessResponse(w, project, "Project created successfully", http.StatusCreated)
}

func (ph *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    role, ok := ctx.Value(ContextKey.UserRole).(roles.Role)
    if !ok {
        logger.Error("user role not found")
        response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
        return
    }

    
    if role != roles.Admin {
        logger.Error("unauthorized to delete projects")
        response.ErrorResponse(w, http.StatusForbidden, "Only admin can delete projects", 1008)
        return
    }

    
    
    projectId := r.PathValue("project_id")
    if projectId == "" {
        response.ErrorResponse(w, http.StatusBadRequest, "Project ID is required", 400)
        return
    }

    
    err := ph.projectService.DeleteProject(projectId,"","")
    if err != nil {
        logger.Error("error deleting project ")
        response.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete project",500)
        return
    }

    logger.Info("Project deleted successfully")
    response.SuccessResponse(w, nil, "Project deleted successfully", http.StatusOK)
}














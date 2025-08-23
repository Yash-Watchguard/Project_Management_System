package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"

	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"

	// "github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/fatih/color"
	
)

type ProjectHandler struct {
	userService service1.UserService
	projectService service1.ProjectService
	taskService service1.TaskService
}

func NewProjectHandler(projectService *service1.ProjectService ,userService *service1.UserService, taskService *service1.TaskService)*ProjectHandler{
	return &ProjectHandler{projectService: *projectService,userService: *userService,taskService: *taskService}
}

func(ph *ProjectHandler)AddNewProject(ctx context.Context)error{
createdBy := ctx.Value(ContextKey.UserId).(string)
	projectId := GenerateUUID()
	projectName, err := GetInput("Enter Project Name :")
	if err != nil {
		return err
	}
	projectDescription, err := GetInput("Enter Project Description :")
	if err != nil {
		return err
	}
	// deadline, _ := GetInput("Enter deadline (YYYY-MM-DD) :")
	// actualdeadline, err := TimeParser(deadline)
	// if err != nil {
	// 	return err
	// }
	var actualdeadline time.Time
    for {
		deadline, _ := GetInput("Enter deadline (YYYY-MM-DD) :")
	    actualdeadline, err = TimeParser(deadline)
	    if err != nil {
		color.Red("%v",err)
		continue
	   }
	   break
	}
	managers,err := ph.userService.GetAllManager(ctx)
	if len(managers)==0{
		return errors.New("no manager found")
	}
	if err != nil  {
		return errors.New("project add faild")
	}

	color.Blue("select Manager from Given list 👇")
	c:=1
	for _,man:= range managers{
     color.Yellow("%d . ID : %s , Name : %s",c,man.Id,man.Name)
	}
	var managerId string
	color.Blue("Enter Manager Id:")
	fmt.Scanln(&managerId)
	project := project.Project{
		ProjectId:          projectId,
		ProjectName:        projectName,
		ProjectDescription: projectDescription,
		Deadline:           actualdeadline,
		CreatedBy:          createdBy,
		AssignedManager:    managerId,
	}
	err = ph.projectService.AddProject(project)
	if err != nil {
		color.Red("Error adding project: %v", err)
	} else {
		color.Green("✅ Project added successfully!")
	}
	
	return nil
}

func (ph *ProjectHandler) SelectAndReturnProjectId(ctx context.Context) (string, error) {
	userId := ctx.Value(ContextKey.UserId).(string)

	projects, err := ph.projectService.ViewAllProjects(ctx)
	if err != nil {
		return "", err
	}

	counter := 1
	projectMap := make(map[int]string) // maps counter to projectId

	for _, project := range projects {
		if project.CreatedBy == userId {
			color.Cyan("----------------%d----------------", counter)
			fmt.Printf("%-20s: %v\n","Project Name",project.ProjectName)
			fmt.Printf("%-20s: %v\n","Project Id",project.ProjectId)
			fmt.Printf("%-20s: %v\n","Project Description",project.ProjectDescription)
			fmt.Printf("%-20s: %v\n","Project Deadline",project.Deadline)
			fmt.Printf("%-20s: %v\n","Assigned To",project.AssignedManager)
			
			projectMap[counter] = project.ProjectId
			counter++
		}
	}

	if counter == 1 {
		return "", errors.New("no project found")
	}

	color.Cyan("Enter project number to select, or press Enter to go back:")
	projectChoice, err := GetInput("")
	if err != nil {
		return "", err
	}

	if projectChoice == "" {
		return "", nil 
	}

	choiceNum, err := strconv.Atoi(projectChoice)
	if err != nil || choiceNum < 1 || choiceNum >= counter {
		return "", errors.New("invalid project choice")
	}

	return projectMap[choiceNum], nil
}
func (ph *ProjectHandler) DeleteProject(ctx context.Context) error {
	// Show project list first
	color.Green("--------- Project List 👇 ----------")

	projects, err := ph.projectService.ViewAllProjects(ctx)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		color.Red("No projects found!")
		return nil
	}

	counter := 1
	for _, p := range projects {
		color.Yellow("%v. Project Name: %v (Project Id: %v)", counter, p.ProjectName, p.ProjectId)
		counter++
	}

	// Ask for project ID to delete
	projectId, err := GetInput("Enter Project ID to delete:")
	if err != nil || projectId == "" {
		return errors.New("project ID cannot be empty")
	}

	// Call service layer to delete
	err = ph.projectService.DeleteProject(ctx, projectId)
	if err != nil {
		return err
	}

	color.Green("✅ Project deleted successfully!")
	Pause()

	return nil
}

func (ph *ProjectHandler) ViewAssignedProjects(ctx context.Context, user *user.User) error {


	assignedProjects, err := ph.projectService.ViewAssignedProject(ctx)
	if err != nil {
		return err
	}

	if len(assignedProjects) == 0 {
		return errors.New("no project assigned")
	}

	for i, project := range assignedProjects {
    color.Cyan("----------- Project %d -----------", i+1)

    fmt.Printf("%-15s : %s\n", "Project ID", project.ProjectId)
    fmt.Printf("%-15s : %s\n", "Project Name", project.ProjectName)
    fmt.Printf("%-15s : %s\n", "Description", project.ProjectDescription)
    fmt.Printf("%-15s : %s\n", "Deadline", project.Deadline.Format("02 Jan 2006"))
    fmt.Printf("%-15s : %s\n", "Created By", project.CreatedBy)

    color.Cyan("----------------------------------")
}


	return nil
}
func (ph *ProjectHandler) ShowProjectStatus(ctx context.Context, projectId string) error {
	// Fetch all tasks under this project
	projectTasks, err := ph.taskService.ViewAllTask(ctx, projectId)
	if err != nil {
		return err
	}

	if len(projectTasks) == 0 {
		return errors.New("no task created for this project — status is 0%")
	}

	total := len(projectTasks)
	done := 0

	for _, t := range projectTasks {
		if t.TaskStatus == status.Done {
			done++
		}
	}

	percentDone := (float64(done) / float64(total)) * 100

	color.Cyan("Project Status for Project ID: %s", projectId)
	color.Yellow("Completed Tasks: %d", done)
	color.Green("Completion: %.2f%%", percentDone)

	return nil
}


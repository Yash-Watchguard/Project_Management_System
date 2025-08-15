package handler

import (
	"context"
	"errors"
	"fmt"
    "strconv"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	
	// "github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/fatih/color"
)

import "github.com/Yash-Watchguard/Tasknest/internal/service1"

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
	deadline, _ := GetInput("Enter deadline (YYYY-MM-DD) :")
	actualdeadline, err := TimeParser(deadline)
	if err != nil {
		return err
	}
	color.Blue("select Manager from Given list 👇")
	err = ph.userService.GetAllManager(ctx)
	if err != nil {
		return errors.New("project Add Faild")
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
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()
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
			color.Yellow("----------------%d----------------", counter)
			color.Yellow("Project Name: %v", project.ProjectName)
			color.Yellow("Project Id: %v", project.ProjectId)
			color.Yellow("Project Description: %v", project.ProjectDescription)
			color.Yellow("Project Deadline: %v", project.Deadline)
			color.Yellow("Assigned To: %v", project.AssignedManager)
			
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
	color.Blue("Press Enter to go back...")
	fmt.Scanln()

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
		color.Yellow("Project ID     : %s", project.ProjectId)
		color.Yellow("Project Name   : %s", project.ProjectName)
		color.Yellow("Description    : %s", project.ProjectDescription)
		color.Yellow("Deadline       : %s", project.Deadline.Format("02 Jan 2006"))
		color.Yellow("Created By     : %s", project.CreatedBy)
		color.Cyan("----------------------------------")
	}
	fmt.Println()
	

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



// func viewTaskofProject(ph *ProjectHandler, ctx context.Context, projectId string) error {
// 	tasks, err := .ViewAllTask(ctx, projectId)
// 	if err != nil {
// 		color.Red(" Failed to fetch tasks: %v", err)
// 		return err
// 	}

// 	if len(tasks) == 0 {
// 		color.Yellow("No tasks found for this project.")
// 		return nil
// 	}

// 	for i, task := range tasks {
// 		color.Cyan("------------ Task %d ------------", i+1)
// 		color.Cyan("Task ID        : %v", task.TaskId)
// 		color.Cyan("Titel          : %v", task.Tile)
// 		color.Cyan("Description    : %v", task.Description)
// 		color.Cyan("Task Priority  : %v", task.TaskPriority)
// 		color.Cyan("Assigned To    : %v", task.AssignedTo)
// 		color.Cyan("Status         : %v", task.TaskStatus)
// 		color.Cyan("Deadline       : %v", task.Deadline)

// 		fmt.Println()
// 	}

// 	color.Blue("1.Enter task Id for add/edit/delete/view a comment on a specific task")
// 	color.Blue("2.Enter to go back")
// 	taskId, err := GetInput("")

// 	if err != nil {
// 		return nil
// 	}
// 	if taskId == "" {
// 		return nil
// 	}
// 	for {
// 		color.Cyan("1.View All Comments")
// 		color.Cyan("2.Add Comment")
// 		color.Cyan("3.Update Comment")
// 		color.Cyan("4.Delete Comment")
// 		color.Cyan("5.Return Back")
// 		var choice int
// 		fmt.Println(color.BlueString("Enter choice:"))
// 		fmt.Scan(&choice)

// 		switch choice {
// 		case 1:
// 			err := viewAllComment(ad, taskId)
// 			if err != nil {
// 				color.Red("%v", err)
// 			}
// 		case 2:
// 			err := addComment(ad, ctx)
// 			if err != nil {
// 				color.Red("%v", err)
// 			}
// 		case 3:
// 			err := updateComment(ad, ctx)
// 			if err != nil {
// 				color.Red("%v", err)
// 			}
// 		case 4:
// 			err := deleteCommnt(ad, ctx)
// 			if err != nil {
// 				color.Red("%v", err)
// 			}
// 		case 5:
// 			return nil
// 		default:
// 			color.Red("Enter Valid choice")
// 		}
// 	}
// }
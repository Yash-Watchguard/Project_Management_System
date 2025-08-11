package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Yash-Watchguard/Tasknest/internal/constants"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/priority"
	// "github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/task"
	status "github.com/Yash-Watchguard/Tasknest/internal/model/task_status"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/service"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	"github.com/fatih/color"
	"strings"
)

func ManagerDashboard(ctx context.Context, users *user.User) {
	userRepo := repository.NewUserRepo()
	projectRepo := repository.NewProjectRepo()
	taskRepo := repository.NewTaskRepo()
	managerRepo := repository.NewManagerRepo()
	
	commentRepo :=repository.NewCommentRepo()
	managerService := service.NewManagerService(userRepo, projectRepo, taskRepo, managerRepo,commentRepo)

	for {
		color.Cyan(constants.ManagerDashbEntry)
		color.Cyan("1. View Profile")
		color.Cyan("2. View Assigned Projects")
		color.Cyan("3. View All Employees")
        color.Cyan("4. Promote Employee")
		color.Cyan("5. Create Task")
		color.Cyan("6. Get Project Status")
		color.Cyan("7. Logout")
        fmt.Print(color.CyanString("Enter your choice: "))
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			exist, err := viewProfile(managerService, ctx, users)
			if err != nil {
				color.Red("%v", err)
			}
			if exist {
				return
			}
		case 2:
			err := viewAssignedProject(managerService, ctx, users)
			if err != nil {
				color.Red("%v", err)
			}
		case 3:
			err := viewAllEmplpyee(managerService, ctx)
			if err != nil {
				color.Red("%v", err)
			}
		case 4:
			err := promoteEmploye(managerService, ctx)
			if err != nil {
				color.Red("Error: %v", err)
			}
		case 5:
			err := taskCreator(managerService, ctx)
			if err != nil {
				color.Red("%v", err)
			}
		case 6:
			err:=projectStatus(managerService,ctx)
			if err != nil {
				color.Red("%v", err)
			}
		case 7:
			color.Green("Logout......")
			return
		default:
			color.Red("Invalid choice. Please select a valid option.")
		}
	}
}

func projectStatus(manager *service.ManagerService, ctx context.Context) error {
	err := onlyViewAssignedProject(manager, ctx)
	if err != nil {
		return errors.New("no project assigned")
	}
	fmt.Println()

	for {
		projectId, err := GetInput("Enter Project Id to see status: ")
		if err != nil {
			return err
		}

		err = showProjectStatus(manager, projectId)
		if err != nil {
			color.Red("%v", err)
		}

		fmt.Println()
		choice, err := GetInput("Do you want to see another project's status? yes or no : ")
		if err != nil {
			return err
		}

		if strings.ToLower(strings.TrimSpace(choice)) != "yes" {
			break
		}
	}

	return nil
}

// task creation under the main dashboard
func taskCreator(manager *service.ManagerService, ctx context.Context) error {

	color.Blue("--------- Select a project to create a task --------")
	err:= onlyViewAssignedProject(manager,ctx)
	if err!=nil{
		return errors.New("no project Assigned")
	}
	fmt.Println()
	projectId, err := GetInput("Enter Project Id (or press Enter to go back): ")
	if err != nil {
		return err
	}

	if strings.TrimSpace(projectId) == "" {
		return nil
	}

	for {
		err = createTask(manager, ctx, projectId)
		if err != nil {
			color.Red("Failed to create task: %v", err)
		} else {
			color.Green("✅ Task successfully created!")
		}

		choice, err := GetInput("Do you want to add another task? (y/n): ")
		if err != nil {
			return err
		}
		if strings.ToLower(choice) != "y" {
			break
		}
	}

	return nil
}

func viewAllEmplpyee(manager *service.ManagerService, ctx context.Context) error {
	employees, err := manager.ViewAllEmplpyee(ctx)
	if err != nil {
		return err
	}

	if len(employees) == 0 {
		color.Yellow(" No employees found.")
		return nil
	}

	color.Cyan("List of Employees:")
	for i, employee := range employees {
		color.Yellow("%d. ID: %s | Name: %s | Email: %s", i+1, employee.Id, employee.Name, employee.Email)
	}
	return nil
}

func viewAssignedProject(manager *service.ManagerService, ctx context.Context, users *user.User) error {
	mangerid := users.Id
	assignedProjects, err := manager.ViewAssignedProject(ctx)
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

	color.Blue("Press 1 for managing project")
	color.Blue("Press 2 for going back")
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	case 1:
		var projectId string
		color.Blue("Enter Project Id : ")
		fmt.Scanln(&projectId)

		for {
			var ch int
			color.Cyan("----------- Task Menu -----------")
			color.Cyan("1. View Tasks")
			color.Cyan("2. Create Task")
			color.Cyan("3. Delete Task")
			color.Cyan("4. Show Project Status")
			color.Cyan("5. Go Back")
			color.Cyan("---------------------------------")
			color.Blue("Enter your choice: ")
			fmt.Scanln(&ch)

			switch ch {
			case 1:
				err := viewTasks(manager, ctx, projectId)
				if err != nil {
					color.Red("%v", err)
				}
			case 2:
				err := createTask(manager, ctx, projectId)
				if err != nil {
					color.Red("%v", err)
				} else {
					color.Green("✅ Task created successfully!")
					fmt.Println("Press ENTER to return to dashboard...")
					fmt.Scanln()
				}
			case 3:
				err := deleteTask(manager, ctx, mangerid, projectId)
				if err != nil {
					color.Red("%v", err)
				}
			case 4:
				err := showProjectStatus(manager, projectId)
				if err != nil {
					color.Red("%v", err)
				}
			case 5:
				return nil
			default:
				color.Red("Invalid choice. Please try again.")
			}

		}
	case 2:
		return nil
	}

	// now for managing the project like add the task view all task

	return nil
}
// this is inside the showall project
func showProjectStatus(manager *service.ManagerService, projectId string) error {
	projectTasks, err := manager.ViewAllTask(ctx, projectId)
	if err != nil {
		return err
	}

	if len(projectTasks) == 0 {
		return errors.New(" No task created for this project , Project status is 0 %")
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
	color.Green("Completion: %v", percentDone)

	return nil
}
func deleteTask(manager *service.ManagerService, ctx context.Context, managerId string, projectId string) error {
	projectTasks, err := manager.ViewAllTask(ctx, projectId)
	if err != nil {
		return err
	}

	if len(projectTasks) == 0 {
		return errors.New(" No task created for this project")
	}
	for key, task := range projectTasks {
		color.Yellow("%d. Name: %s  ID: %s", key+1, task.Tile, task.TaskId)
	}
	taskId, err := GetInput("Enter Task Id :")
	if err != nil {
		return nil
	}

	return manager.DeleteTask(ctx, managerId, taskId)
}

func createTask(manager *service.ManagerService, ctx context.Context, projectId string) error {
	managerId := ctx.Value(ContextKey.UserId).(string)
	taskId := GenerateUUID()

	title, err := GetInput("Enter Task Title: ")
	if err != nil {
		return err
	}

	description, err := GetInput("Enter Task Description: ")
	if err != nil {
		return err
	}
	var deadline time.Time
	for {
		deadlineStr, err := GetInput("Enter Deadline in YYYY-MM-DD: ")
		if err != nil {
			return err
		}

		deadline, err = TimeParser(deadlineStr)
		if err != nil {
			color.Red("invalid date format")
		} else {
			break
		}
	}

	var priority Priority.Priority
	for {
		priorityStr, err := GetInput("Enter Priority =>Low/Medium/High: ")
		if err != nil {
			return err
		}

		priority, err = Priority.PriorityParser(priorityStr)
		if err != nil {
		    color.Red("Invalid priority. Choose Low, Medium, or High.")
		} else {
			break
		}
	}

	assignedTo, err := GetInput("Enter Employee ID to assign this task to: ")
	if err != nil {
		return err
	}

	newTask := task.Task{
		TaskId:       taskId,
		Tile:         title,
		Description:  description,
		Deadline:     deadline,
		TaskPriority: priority,
		TaskStatus:   status.Pending,
		AssignedTo:   assignedTo,
		ProjectId:    projectId,
		CreatesBy:    managerId,
	}

	err = manager.CreateTask(ctx, managerId, newTask)
	if err != nil {
		return err
	}
	return nil
}

func viewTasks(manager *service.ManagerService, ctx context.Context, projectId string) error {
	projectTasks, err := manager.ViewAllTask(ctx, projectId)
	if err != nil {
		return err
	}

	if len(projectTasks) == 0 {
		return errors.New(" No task created for this project")
	}

	color.New(color.FgHiCyan, color.Bold).Println("\n----------- Project Tasks -----------")

	for i, task := range projectTasks {
		color.Cyan("--------------- Task %d ---------------", i+1)
		color.Yellow("Task ID       : %s", task.TaskId)
		color.Yellow("Title         : %s", task.Tile)
		color.Yellow("Description   : %s", task.Description)
		color.Yellow("Deadline      : %s", task.Deadline.Format("2006-01-02"))
		color.Yellow("Priority      : %s", task.TaskPriority)
		color.Yellow("Status        : %s", status.GetStatusString(task.TaskStatus))
		color.Yellow("Assigned To   : %s", task.AssignedTo)
		color.Cyan("----------------------------------------\n")
	}
    
	taskId,err:=GetInput("Enter task Id for add/edit/delete a comment on a specific task")
	color.Blue("2.Enter to go back")

	if err!=nil{
		return nil
	}
	if taskId == "" {
		return nil 
	}
    for{
	color.Cyan("1.View All Comments")
	color.Cyan("2.Add Comment")
	color.Cyan("3.Update Comment")
	color.Cyan("4.Delete Comment")
	color.Cyan("5.Return Back")
	var choice int
	fmt.Println(color.BlueString("Enter choice:"))
	fmt.Scan(&choice)

	switch choice{
	case 1:
		err:=ViewAllComment(manager,taskId)
		if err!=nil{
			color.Red("%v",err)
		}
	case 2:
		err:=AddComment(manager,ctx)
		if err!=nil{
			color.Red("%v",err)
		}
	case 3:
		err:=UpdateComment(manager,ctx)
		if err!=nil{
			color.Red("%v",err)
		}
	case 4:
		err:=DeleteCommnt(manager,ctx)
		if err!=nil{
			color.Red("%v",err)
		}
	case 5:
		return nil
	default:
		color.Red("Enter Valid choice")
	}
}
}

func UpdateComment(manager *service.ManagerService,ctx context.Context)error{
	commentId,err:=GetInput("Enter Comment id for update the comment:")
	if err!=nil{
		return nil
	}
	updatedComment,err:=GetInput("Enter Updated Comment:")
    if err!=nil{
		return nil
	}
	err=manager.UpdateComment(ctx,commentId,updatedComment)
    if err!=nil{
		return err
	}
	color.Green("Comment Updated Successfully......")
	color.Cyan("Press enter for going back")
	fmt.Scanln()
	return nil
}

func ViewAllComment(manager *service.ManagerService,taskId string)error{

    comments,err:=manager.ViewAllComment(taskId)
    if err!=nil{
		return err
	}

	if len(comments)==0{
		return errors.New("no comments found for this task")
	}

	for _,comment:=range comments{
		color.Blue("Comment for %v (Created by-%v)",comment.TaskId,comment.CreatedBy)
		color.Cyan("%v",comment.Content)
		color.Blue("----------------------------------------------------------------")
	}
    color.Green("Press Enter to return to the previous menu...")
	fmt.Scanln()
	return nil
}
func AddComment(manager *service.ManagerService,ctx context.Context)error{
	
	var taskId string
	for{
	taskId, err := GetInput("Enter Task ID to comment on:")
	if err != nil || taskId == "" {
		color.Red("Invalid task ID")
	}else{
		break
	}
    }

	content, err := GetInput("Enter your comment:")
	if err != nil {
		fmt.Print(err)
	}

	createdBy:=ctx.Value(ContextKey.UserId).(string)
    
	commentId := GenerateUUID()

	newComment := comment.Comment{
		CommentId: commentId,
		Content:   content,
		CreatedBy: createdBy,
		TaskId:    taskId,
	}

	err = manager.AddComment(newComment)
	if err != nil {
		return err
	}

	color.Green("✅ Comment added successfully!")
	color.Blue("press enter to going back")
	fmt.Scanln()
	return nil
}

func DeleteCommnt(manager *service.ManagerService,ctx context.Context)error{
	commentId,err:=GetInput("Enter Comment id for update the comment:")
	if err!=nil{
		return nil
	}
	err= manager.DeleteComment(ctx,commentId)
    
	if err!=nil{
		return err
		
	}
	color.Green("comment deleted succesfully")
	color.Blue("press enter to going back")
	fmt.Scanln()
    return nil
}
func viewProfile(managerService *service.ManagerService, ctx context.Context, manager *user.User) (bool, error) {
	userProfiles, err := managerService.ViewProfile(ctx, manager.Id)
	if err != nil {
		return false, err
	}

	if len(userProfiles) == 0 {
		return false, errors.New("no user profile")
	}

	user := userProfiles[0]
	color.Cyan("----------- %s Profile -----------", user.Name)
	color.Yellow("ID           : %s", user.Id)
	color.Yellow("Name         : %s", user.Name)
	color.Yellow("Email        : %s", user.Email)
	color.Yellow("Phone Number : %s", user.PhoneNumber)
	color.Yellow("Role         : %v", user.Role)
	color.Cyan("----------------------------------")

	color.Blue("1. Update Profile")
	color.Blue("2. Delete Profile")
	color.Blue("3. Go Back")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		err := updateProfile(managerService, ctx, &user)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return false, nil

	case 2:
		err := managerService.DeleteUser(ctx, user.Id)
		if err != nil {
			fmt.Printf("%v\n", err)
			return false, err
		}
		color.Red("Profile Deleted Successfully")
		return true, nil

	case 3:
		return false, nil

	default:
		color.Red("Invalid choice.")
		return false, nil
	}
}

func updateProfile(manager *service.ManagerService, ctx context.Context, user *user.User) error {
	for {
		color.Blue("1. Update Name")
		color.Blue("2. Update Email")
		color.Blue("3. Update Password")
		color.Blue("4. Update Contact")
		color.Blue("5. Go Back")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			name, err := GetInput("Enter Name: ")
			if err != nil {
				color.Red("Error reading name: %v", err)
				continue
			}
			user.Name = name

		case 2:
			email, err := GetValidEmail()
			if err != nil {
				color.Red("Invalid email: %v", err)
				continue
			}
			user.Email = email

		case 3:
			password, err := GetValidPassword()
			if err != nil {
				color.Red("Invalid password: %v", err)
				continue
			}
			user.Password = password

		case 4:
			contact, err := GetValidPhoneNumber()
			if err != nil {
				color.Red("Invalid phone number: %v", err)
				continue
			}
			user.PhoneNumber = contact

		case 5:
			return nil

		default:
			color.Red("Invalid choice.")
			continue
		}

		// Save the updated user
		err := manager.UpdateProfile(user.Id, ctx, user.Name, user.Email, user.Password, user.PhoneNumber)
		if err != nil {
			color.Red("Update failed: %v", err)
		} else {
			color.Green("User updated successfully!")
		}
	}
}

func onlyViewAssignedProject(manager *service.ManagerService,ctx context.Context)error{
	projects, err := manager.ViewAssignedProject(ctx)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		return errors.New("no project assigned")
	}

	for key, project := range projects {
		color.Cyan("%d. %s (Project Id - %s)", key+1, project.ProjectName, project.ProjectId)
	}
	return nil
}

func promoteEmploye(manager *service.ManagerService, ctx context.Context) error {
	// view all employees
	users, err := manager.ViewAllEmplpyee(ctx)
	if err != nil {
		return nil
	}
	for key, user := range users {
		color.Blue("---------------------user %v----------------", key+1)
		color.Yellow("Name- %v ", user.Name)
		color.Yellow("Id - %v", user.Id)
		color.Yellow("Role- %v", roles.RoleParser(user.Role))
		color.Blue("--------------------------------------------")
	}

	employeeId, err := GetInput("Enter Employee Id To promot as Manager")
	if err!=nil{
		return err
	}

	err = manager.PromoteEmployee(ctx, employeeId)
	if err != nil {
		return err
	}
	return nil
}

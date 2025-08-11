package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/internal/constants"
	"github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/service"
	"github.com/fatih/color"
)

func AdminDashboard(ctx context.Context, user *user.User) {
	// userId := ctx.Value(ContextKey.UserId).(string)
	userRepo := repository.NewUserRepo()
	projectRepo := repository.NewProjectRepo()
	adminRepo := repository.NewAdminRepo()
	taskRepo := repository.NewTaskRepo()
	commentRepo := repository.NewCommentRepo()
	adminService := service.NewAdminServices(userRepo, projectRepo, adminRepo, taskRepo, commentRepo)

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

		var choice int
		fmt.Print(color.GreenString("Enter your choice: "))
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			exist, err := viewAdminProfile(adminService, ctx, user)
			if err != nil {
				color.Red("%v", err)
			}
			if exist {
				return
			}
		case 2:
			err := viewAllUsers(adminService, ctx)
			if err != nil {
				color.Red("%v", err)
			}

		case 3:
			err := deleteUser(adminService, ctx)
			if err != nil {
				color.Red("Error: %v", err)
			}
		case 4:
			err := promoteEmployee(adminService, ctx)
			if err != nil {
				color.Red("Error: %v", err)
			}
		case 5:
			color.Blue("1. Add New Project")
			color.Blue("2. View All Projects")
			color.Blue("3. Delete Project")

			var choice int
			fmt.Print(color.GreenString("Enter your choice: "))
			fmt.Scanln(&choice)
			switch choice {
			case 1:
				err := addNewProject(adminService, ctx)
				if err != nil {
					color.Red("%v", err)
				}
			case 2:
				err := viewAllProjects(adminService, ctx)
				if err != nil {
					color.Red("Error: %v", err)
				}
				fmt.Println("Press ENTER to return to dashboard...")
				fmt.Scanln()
			case 3:
				err := deleteProject(adminService, ctx)
				if err != nil {
					color.Red("Error: %v", err)
				}
				fmt.Println("Press ENTER to return to dashboard...")
				fmt.Scanln()
			default:
				color.Red("Invalid choice. Please try again.")
			}

		case 6:
			color.Green("Logging out...")
			return

		default:
			color.Red("Invalid choice. Please try again.")
		}
	}
}

//  function for view profile

func viewAdminProfile(ad *service.AdminService, ctx context.Context, admin *user.User) (bool, error) {
	userProfiles, err := ad.ViewProfile(ctx, admin.Id)
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
		err := updateAdminProfile(ad, ctx, &user)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return false, nil

	case 2:
		err := ad.DeleteUser(ctx, user.Id)
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

func promoteEmployee(ad *service.AdminService, ctx context.Context) error {

	users, err := ad.ViewAllUsers(ctx)
	if err != nil {
		return nil
	}
	counter := 0
	for key, user := range users {
		if user.Role == 0 || user.Role == 1 {
			continue
		}
		color.Blue("---------------------user %v----------------", key+1)
		color.Yellow("Name- %v ", user.Name)
		color.Yellow("Id - %v", user.Id)
		color.Yellow("Role- %v", roles.RoleParser(user.Role))
		color.Blue("--------------------------------------------")
		counter++
	}
	if counter == 0 {
		return errors.New("no employees for Promotion")
	}

	employeeId, err := GetInput("Enter Employee Id To promot as Manager")
	if err != nil {
		return err
	}

	err = ad.PromoteEmployee(ctx, employeeId)
	if err != nil {
		return err
	}

	color.Green("💐 Promoted as Manbager .......")
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()
	return nil
}

func viewAllUsers(ad *service.AdminService, ctx context.Context) error {
	users, err := ad.ViewAllUsers(ctx)

	if err != nil {
		return nil
	}

	color.Cyan("----------------------------- All Users -----------------------------------")
	counter := 1
	for _, user := range users {
		color.Yellow("%d. ID: %s, Name: %s, Email: %s, Role: %d\n", counter, user.Id, user.Name, user.Email, user.Role)
		counter++
	}
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()
	return nil
}

func deleteUser(ad *service.AdminService, ctx context.Context) error {
	ad.ViewAllUsers(ctx)
	fmt.Println("Enter User Id of user :")
	var userId string
	_, err := fmt.Scanln(&userId)
	if err != nil {
		return err
	}
	err = ad.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}
	color.Green("User deleted successfully!")
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()
	return nil
}

func addNewProject(ad *service.AdminService, ctx context.Context) error {
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
	err = ad.GetAllManager(ctx)
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
	err = ad.AddProject(project)
	if err != nil {
		color.Red("Error adding project: %v", err)
	} else {
		color.Green("✅ Project added successfully!")
	}
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()
	return nil
}

func viewAllProjects(ad *service.AdminService, ctx context.Context) error {
	userId := ctx.Value(ContextKey.UserId).(string)

	projects, err := ad.ViewAllProjects(ctx)
	if err != nil {
		return err
	}

	counter := 1
	for _, project := range projects {
		if project.CreatedBy == userId {
			color.Yellow("----------------%d----------------", counter)
			color.Yellow("Project Name: %v", project.ProjectName)
			color.Yellow("Project Id: %v", project.ProjectId)
			color.Yellow("Project Description: %v", project.ProjectDescription)
			color.Yellow("Project Deadline: %v", project.Deadline)
			color.Yellow("Assigned To: %v", project.AssignedManager)
			counter++
		}
	}

	if counter == 1 {
		color.Red("No projects found")
		return nil
	}

	color.Cyan("1.Enter Project ID to view tasks")
	color.Blue("press Enter to go back")
	projectId, err := GetInput("")
	

	if err != nil {
		return err
	}
	if projectId == "" {
		return nil
	}

	err = viewTaskofProject(ad, ctx, projectId)
	if err != nil {
		return err
	}
	return nil
}

func viewTaskofProject(ad *service.AdminService, ctx context.Context, projectId string) error {
	tasks, err := ad.ViewAllTask(ctx, projectId)
	if err != nil {
		color.Red(" Failed to fetch tasks: %v", err)
		return err
	}

	if len(tasks) == 0 {
		color.Yellow("No tasks found for this project.")
		return nil
	}

	for i, task := range tasks {
		color.Cyan("------------ Task %d ------------", i+1)
		color.Cyan("Task ID        : %v", task.TaskId)
		color.Cyan("Titel          : %v", task.Tile)
		color.Cyan("Description    : %v", task.Description)
		color.Cyan("Task Priority  : %v", task.TaskPriority)
		color.Cyan("Assigned To    : %v", task.AssignedTo)
		color.Cyan("Status         : %v", task.TaskStatus)
		color.Cyan("Deadline       : %v", task.Deadline)

		fmt.Println()
	}

	color.Blue("1.Enter task Id for add/edit/delete/view a comment on a specific task")
	color.Blue("2.Enter to go back")
	taskId, err := GetInput("")

	if err != nil {
		return nil
	}
	if taskId == "" {
		return nil
	}
	for {
		color.Cyan("1.View All Comments")
		color.Cyan("2.Add Comment")
		color.Cyan("3.Update Comment")
		color.Cyan("4.Delete Comment")
		color.Cyan("5.Return Back")
		var choice int
		fmt.Println(color.BlueString("Enter choice:"))
		fmt.Scan(&choice)

		switch choice {
		case 1:
			err := viewAllComment(ad, taskId)
			if err != nil {
				color.Red("%v", err)
			}
		case 2:
			err := addComment(ad, ctx)
			if err != nil {
				color.Red("%v", err)
			}
		case 3:
			err := updateComment(ad, ctx)
			if err != nil {
				color.Red("%v", err)
			}
		case 4:
			err := deleteCommnt(ad, ctx)
			if err != nil {
				color.Red("%v", err)
			}
		case 5:
			return nil
		default:
			color.Red("Enter Valid choice")
		}
	}
}

func deleteCommnt(ad *service.AdminService, ctx context.Context) error {
	commentId, err := GetInput("Enter Comment id for delete the comment:")
	if err != nil {
		return nil
	}
	err = ad.DeleteComment(ctx, commentId)

	if err != nil {
		return err

	}
	color.Green("comment deleted succesfully")
	color.Blue("press enter to going back")
	fmt.Scanln()
	return nil
}
func addComment(ad *service.AdminService, ctx context.Context) error {

	var taskId string
	for {
		taskId, err := GetInput("Enter Task ID to comment on:")
		if err != nil || taskId == "" {
			color.Red("Invalid task ID")
		} else {
			break
		}
	}

	content, err := GetInput("Enter your comment:")
	if err != nil {
		fmt.Print(err)
	}

	createdBy := ctx.Value(ContextKey.UserId).(string)

	commentId := GenerateUUID()

	newComment := comment.Comment{
		CommentId: commentId,
		Content:   content,
		CreatedBy: createdBy,
		TaskId:    taskId,
	}

	err = ad.AddComment(newComment)
	if err != nil {
		return err
	}

	color.Green("✅ Comment added successfully!")
	color.Blue("press enter to going back")
	fmt.Scanln()
	return nil
}

func viewAllComment(ad *service.AdminService, taskId string) error {

	comments, err := ad.ViewAllComment(taskId)
	if err != nil {
		return err
	}

	if len(comments) == 0 {
		return errors.New("no comments found for this Task")
	}

	for _, comment := range comments {
		color.Blue("Comment for %v (Created by-%v)", comment.TaskId, comment.CreatedBy)
		color.Cyan("%v", comment.Content)
		color.Blue("----------------------------------------------------------------")
	}
	color.Green("Press Enter to return to the previous menu...")
	fmt.Scanln()
	return nil
}
func updateComment(ad *service.AdminService, ctx context.Context) error {
	commentId, err := GetInput("Enter Comment id for update the comment:")
	if err != nil {
		return nil
	}
	updatedComment, err := GetInput("Enter Updated Comment:")
	if err != nil {
		return nil
	}
	err = ad.UpdateComment(ctx, commentId, updatedComment)
	if err != nil {
		return err
	}
	color.Green("Comment Updated Successfully......")
	color.Cyan("Press enter for going back")
	fmt.Scanln()
	return nil
}

func deleteProject(ad *service.AdminService, ctx context.Context) error {
	color.Green("---------Project List 👇----------")

	var projects []project.Project
	projects, err := ad.ViewAllProjects(ctx)
	if err != nil {
		return err
	}
	counter := 1
	for _, project := range projects {
		color.Yellow("%v. Project Name: %v (Project Id : %v)", counter, project.ProjectName, project.ProjectId)
		counter++
	}

	projectId, err := GetInput("Enter project Id")
	if err != nil {
		return nil
	}

	err = ad.DeleteProject(ctx, projectId)
	if err != nil {
		return errors.New("problem in deleting the project")
	}
	color.Green("Project Deleted Succesfully ✅")
	return err
}

func updateAdminProfile(ad *service.AdminService, ctx context.Context, user *user.User) error {
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

		err := ad.UpdateProfile(user.Id, ctx, user.Name, user.Email, user.Password, user.PhoneNumber)
		if err != nil {
			color.Red("Update failed: %v", err)
		} else {
			color.Green("User updated successfully!")
		}

	}
}

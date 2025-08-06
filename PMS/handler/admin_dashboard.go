package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/Yash-Watchguard/Tasknest/internal/constants"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/service"
	"github.com/fatih/color"
)

func AdminDashboard(ctx context.Context) {
	userId := ctx.Value(ContextKey.UserId).(string)
	userRepo := repository.NewUserRepo()
	projectRepo := repository.NewProjectRepo()
	adminService := service.NewAdminServices(userRepo, projectRepo)

	for {
		color.Blue(constants.AdminDashbEntry)
		color.Blue("1. View Profile")
		color.Blue("2. View All Users")
		color.Blue("3. Delete User")
		color.Blue("4. Add New Project")
		color.Blue("5. View All Projects")
		color.Blue("6. Delete Project")
		color.Blue("7. Logout")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			userProfiles, err := adminService.ViewProfile(ctx, userId)
			if err != nil {
				color.Red("%v", err)
				break
			}
			if len(userProfiles) == 0 {
				color.Red("No user profile found.")
				break
			}
			user := userProfiles[0]
			color.Cyan("----------- %s Profile -----------", user.Name)
			color.Yellow("ID     : %s", user.Id)
			color.Yellow("Name   : %s", user.Name)
			color.Yellow("Email  : %s", user.Email)
			color.Yellow("Role   : %s", user.Role)
			color.Cyan("-------------------------------------")
			fmt.Println("Press ENTER to return to dashboard...")
			fmt.Scanln()

		case 2:
			err := adminService.ViewAllUsers(ctx)
			if err != nil {
				color.Red("%v", err)
			}

		case 3:
			err := deleteUser(adminService, ctx)
			if err != nil {
				color.Red("Error: %v", err)
			}

		case 4:
			err := addNewProject(adminService, ctx)
			if err != nil {
				color.Red("%v", err)
			}

		case 5:
			err := viewAllProjects(adminService, ctx)
			if err != nil {
				color.Red("Error: %v", err)
			}
			fmt.Println("Press ENTER to return to dashboard...")
			fmt.Scanln()

		case 6:
			err := deleteProject(adminService, ctx)
			if err != nil {
				color.Red("Error: %v", err)
			}
			fmt.Println("Press ENTER to return to dashboard...")
			fmt.Scanln()

		case 7:
			color.Green("Logging out...")
			return

		default:
			color.Red("Invalid choice. Please try again.")
		}
	}
}

func deleteUser(ad *service.AdminService, ctx context.Context) error {
	ad.ViewAllUsers(ctx)
	fmt.Println("Enter User Id of user:")
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
	projectName, err := GetInput("Enter Project Name:")
	if err != nil {
		return err
	}
	projectDescription, err := GetInput("Enter Project Description")
	if err != nil {
		return err
	}
	deadline, _ := GetInput("Enter deadline (YYYY-MM-DD):")
	actualdeadline, err := TimeParser(deadline)
	if err != nil {
		return err
	}
	color.Blue("select Manager from Given list 👇")
	err = ad.GetAllManager(ctx)
	if err != nil {
		color.Red("%s", err)
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
	var projects []project.Project
	projects, err := ad.ViewAllProjects(ctx)
	if err != nil {
		return err
	}
	counter := 1
	for _, project := range projects {
		if project.CreatedBy == userId {
			color.Cyan("----------------%d----------------", counter)
			color.Cyan("Project Name: %v", project.ProjectName)
			color.Cyan("Project Id: %v", project.ProjectId)
			color.Cyan("Project Description: %v", project.ProjectDescription)
			color.Cyan("Project Deadline: %v", project.Deadline)
			color.Cyan("Assigned To: %v", project.AssignedManager)
			counter++
		}
	}
	return nil
}

func deleteProject(ad *service.AdminService, ctx context.Context) error {
	color.Green("---------Project List 👇----------")
	err := viewAllProjects(ad, ctx)
	if err != nil {
		return errors.New("no project for assign")
	}
	var proId string
	color.Blue("Enter the Project id which you want to delete")
	fmt.Scanln(&proId)
	err = ad.DeleteProject(ctx, proId)
	if err != nil {
		return errors.New("problem in deleting the project")
	}
	color.Green("Project Deleted Succesfully ✅")
	return err
}

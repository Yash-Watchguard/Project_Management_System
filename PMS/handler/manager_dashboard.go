package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/internal/constants"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/service"
	"github.com/fatih/color"
)

func ManagerDashboard(ctx context.Context, users *user.User) {
	userRepo := repository.NewUserRepo()
	projectRepo := repository.NewProjectRepo()

	managerService := service.NewManagerService(userRepo, projectRepo)

	for {
		color.Blue(constants.ManagerDashbEntry)
		color.Blue("1. View Profile")
		color.Blue("2. View Assigned Projects")
		color.Blue("3. View All Employees")
		color.Blue("4. Create Task")
		color.Blue("5. Get Project Status")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			exit, err := viewProfile(managerService, ctx, users)
			if err != nil {
				color.Red("%v", err)
			}
			if exit {
				return // Exit dashboard after deletion
			}
		case 2:
			err := viewAssignedProject(managerService, ctx)
			if err != nil {
				color.Red("%v", err)
			}
		default:
			color.Red("Invalid choice. Please select a valid option.")
		}
	}
}


func viewAssignedProject(manager *service.ManagerService, ctx context.Context) error {
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


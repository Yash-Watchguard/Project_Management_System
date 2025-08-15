package handler

import (
	"context"
	"errors"
	"fmt"

	
	// "github.com/Yash-Watchguard/Tasknest/internal/constants"
	// "github.com/Yash-Watchguard/Tasknest/internal/model/comment"
	// ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	// "github.com/Yash-Watchguard/Tasknest/internal/model/project"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/fatih/color"
)

type UserHandler struct {
	userService *service1.UserService
}

func NewUserHandler(userService *service1.UserService)*UserHandler{
	return &UserHandler{userService: userService}
}

func(uh *UserHandler)ViewUserProfile( ctx context.Context, admin *user.User) (bool, error) {
	userProfiles, err := uh.userService.ViewProfile(ctx, admin.Id)
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
    color.Green("Enter your Choice : ")
	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		err :=updateProfile(uh, ctx, &user)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		return false, nil

	case 2:
		err := uh.userService.DeleteUser(ctx, user.Id)
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

func updateProfile(uh *UserHandler, ctx context.Context, user *user.User) error {
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

		err := uh.userService.UpdateProfile(user.Id, ctx, user.Name, user.Email, user.Password, user.PhoneNumber)
		if err != nil {
			color.Red("Update failed: %v", err)
		} else {
			color.Green("User updated successfully!")
		}

	}
}

func(uh *UserHandler)ViewallUsers(ctx context.Context)error{
	users, err :=uh.userService.ViewAllUsers(ctx)

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
func(uh *UserHandler)DeleteUser( ctx context.Context) error {
	users, err := uh.userService.ViewAllUsers(ctx)
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

	fmt.Println("Enter User Id :")
	var userId string
	_, err = fmt.Scanln(&userId)

	if err != nil {
		return err
	}

	err = uh.userService.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}
	
	color.Green("User deleted successfully!")
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()
	return nil
}
func(uh *UserHandler)PromoteEmployee(ctx context.Context)error{
	users, err := uh.userService.ViewAllUsers(ctx)
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

	err = uh.userService.PromoteEmployee(ctx,employeeId)
	if err != nil {
		return err
	}

	color.Green("💐 Promoted as Manbager .......")
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()
	return nil
}

func (uh *UserHandler) ViewAllEmployees(ctx context.Context) error {
	users, err := uh.userService.ViewAllEmplpyee(ctx) // fixed typo: ViewAllEmployee
	if err != nil {
		color.Red("Error fetching employees: %v", err)
		return err
	}

	if len(users) == 0 {
		color.Yellow("No employees found.")
		return nil
	}

	color.Cyan("===== List of Employees =====")
	for i, user := range users {
		color.Yellow("%d. ID: %s | Name: %s | Email: %s | Role: %s", i+1, user.Id, user.Name, user.Email, user.Role)
	}
	color.Cyan("=============================")
	
	fmt.Println("Press ENTER to return to dashboard...")
	fmt.Scanln()
	return nil
}

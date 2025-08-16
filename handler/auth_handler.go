package handler

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	

	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"

	
	"github.com/Yash-Watchguard/Tasknest/internal/util"
	"github.com/fatih/color"
	"golang.org/x/crypto/bcrypt"
)
var GenerateUUID = util.GenerateUniqueUUID
var ValidEmail = util.ValidateEmail
var ValidPhoneNumber = util.ValidateMobileNumber
var inputReader *bufio.Reader = bufio.NewReader(os.Stdin)
var TimeParser = util.ParseDate

func SetInputReader(r io.Reader) {
	inputReader = bufio.NewReader(r)
}

type authHandler struct{
	userService *service1.AuthService
}
func NewAuthHandler(userService *service1.AuthService)*authHandler{
	return &authHandler{userService: userService}
}
func(au *authHandler)Signup() error {
	
	 
		name ,err:=GetInput("Enter name : ")
	    if err != nil {
	    return err
	   }

	var email string
	for {
		email, err = GetValidEmail()
		if err != nil {
			color.Red("Please Enter Valid Email Address : ")
			continue
		} else {
			break
		}

	}

	var password string
	for {
		password, err = GetValidPassword()
		if err != nil {
			color.Red("Please Enter Valid password : ")
			continue
		} else {
			break
		}
	}

	var phoneNumber string
	for {
		phoneNumber, err = GetValidPhoneNumber()
		if err != nil {
			color.Red("Please Enter Valid Phone number : ")
			continue
		} else {
			break
		}

	}
	// var userRole roles.Role
	// for {
	// 	color.Magenta("Please Enter Your Role")
	// 	color.Blue("Press 1 for Admin")
	// 	color.Blue("Press 2 for Manager")
	// 	color.Blue("Press 3 for Employee")

	// 	var choice int
		
	// 	fmt.Scanln(&choice)

	// 	switch choice {
	// 	case 1:
	// 		userRole = roles.Admin
	// 	case 2:
	// 		userRole = roles.Manager
	// 	case 3:
	// 		userRole = roles.Employee
	// 	default:
	// 		color.Red("Enter Valid Choice")
	// 		continue
	// 	}
	// 	break
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := user.User{
		Id:          GenerateUUID(),
		Name:        name,
		Email:       email,
		Password:    string(hashedPassword),
		PhoneNumber: phoneNumber,
		Role:        2,
	}

	
	// authService := service.NewAuthService(repo)

	if err := au.userService.Signup(&user); err != nil {
		color.Red("Signup failed: %v", err)
		return err
	}

	color.Green("✅ User signed up successfully!")
	return nil
}

func GetInput(prompt string) (string, error) {
	fmt.Print(color.RedString(prompt))
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// Pause waits for user to press Enter
func Pause() {
	fmt.Print(color.BlueString("Press Enter to go back..."))
	reader.ReadString('\n') // ignore error intentionally
}

func GetValidEmail() (string, error) {
	email, err := GetInput("Enter Email : ")
	if err != nil {
		return "", err
	}
	if err := ValidEmail(email); err != nil {

		return "", err
	}
	return email, nil
}

func GetValidPhoneNumber() (string, error) {
	number, err := GetInput("Enter Phone Number : ")
	if err != nil {
		return "", err
	}
	number = strings.TrimSpace(number)
	if err := ValidPhoneNumber(number); err != nil {
		return "", err
	}
	return number, nil
}

func GetValidPassword() (string, error) {
	password, err := GetInput("Enter Password: ")
	if err != nil {
		return "", err
	}
	if err := util.ValidatePassword(password); err != nil {

		return "", err
	}
	return password, nil
}

func(au * authHandler)Login() (*user.User,error) {

	
	name, err := GetInput("Enter name : ")
	if err != nil {
		color.Red("Login Faild")
		return nil,err
	}
	mailId, err := GetValidEmail()
	if err != nil {
		color.Red("Login Faild")
		return nil,err
	}
	password, err := GetValidPassword()
	if err != nil {
		color.Red("Login Faild")
		return nil,err
	}

	

	user, err := au.userService.Login(name, mailId, password)
	if err != nil {
		color.Red("----------Invalid details,Login Faild----------------")
		return nil,err
	}
	// 
	return user,err
}



// func DashBoard(user *model.User) {
// 	if user.Role == "Admin" {
// 		AdminDashboard(user)
// 	} else if user.Role == "Manager" {
// 		//    ManagerDashboard(user)
// 	} else {
// 		//    EmployeeDashboard(user)
// 	}

// 
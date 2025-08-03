package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Yash-Watchguard/Tasknest/model"
	"github.com/Yash-Watchguard/Tasknest/repository"
	"github.com/Yash-Watchguard/Tasknest/service"
	"github.com/Yash-Watchguard/Tasknest/util"

	"github.com/fatih/color"
	"golang.org/x/crypto/bcrypt"
)

var GenerateUUID = util.GenerateUniqueUUID
var ValidEmail = util.ValidateEmail
var ValidPhoneNumber = util.ValidateMobileNumber
var inputReader *bufio.Reader = bufio.NewReader(os.Stdin)
var TimeParser=util.ParseDate


func SetInputReader(r io.Reader) {
	inputReader = bufio.NewReader(r)
}

func SignupCli() error {
	name, err := GetInput("Enter Your Name: ")
	if err != nil {
		return err
	}

	email, err := GetValidEmail()
	if err != nil {
		return err
	}

	password, err := GetValidPassword()
	if err != nil {
		return err
	}

	phoneNumber, err := GetValidPhoneNumber()
	if err != nil {
		return err
	}

	var role string
	for {
		color.Magenta("Please Enter Your Role")
		color.Blue("Press 1 for Admin")
		color.Blue("Press 2 for Manager")
		color.Blue("Press 3 for Employee")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			role = "Admin"
		case 2:
			role = "Manager"
		case 3:
			role = "Employee"
		default:
			color.Red("Enter Valid Choice")
			continue
		}
		break
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Id:          GenerateUUID(),
		Name:        name,
		Email:       email,
		Password:    string(hashedPassword),
		PhoneNumber: phoneNumber,
		Role:        role,
	}

	repo := repository.NewUserRepo()
    authService := service.NewAuthService(repo)

if err := authService.Signup(&user); err != nil {
	color.Red("Signup failed: %v", err)
	return err
}


	color.Green("✅ User signed up successfully!")
	return nil
}

func GetInput(prompt string) (string, error) {
	color.Red(prompt)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func GetValidEmail() (string, error) {
	email, err := GetInput("Enter Email: ")
	if err != nil {
		return "", err
	}
	if err := ValidEmail(email); err != nil {
		color.Red("Invalid email: %v", err)
		return "", err
	}
	return email, nil
}

func GetValidPhoneNumber() (string, error) {
	number, err := GetInput("Enter Phone Number: ")
	if err != nil {
		return "", err
	}
	number = strings.TrimSpace(number)
	if err := ValidPhoneNumber(number); err != nil {
		color.Red("Invalid phone number: %v", err)
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
		color.Red("Invalid password: %v", err)
		return "", err
	}
	return password, nil
}

func LoginCli()error {
     name,err:=GetInput("Enter name:")
	 if err!=nil{
		return err
	 }
     mailId,err:=GetValidEmail()
	 if err!=nil{
		return err
	 }
	 password,err:=GetValidPassword()
	 if err!=nil{
		return err
	 }

	 repo:=repository.NewUserRepo()
	 authService:=service.NewAuthService(repo)

	 user,err:=authService.Login(name,mailId,password)
     if err!=nil{
		color.Red("Login Faild:",err)
		return err
	 }
	 
	 color.Green("Welcom back %s in Worknest☺️",user.Name)
     DashBoard(user)
	 return nil
}

func DashBoard(user *model.User){
if user.Role=="Admin"{
   AdminDashboard(user)
}else if user.Role=="Manager"{
//    ManagerDashboard(user)
}else{
//    EmployeeDashboard(user)
}

}
